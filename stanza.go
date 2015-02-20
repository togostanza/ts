package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

type Stanza struct {
	BaseDir string
	Name    string
	Metadata
}

type Parameter struct {
	Key         string `json:"stanza:key"`
	Description string `json:"stanza:description"`
	Example     string `json:"stanza:example"`
	Required    bool   `json:"stanza:required"`
}

type Metadata struct {
	Id         string      `json:"@id"`
	Label      string      `json:"stanza:label"`
	Parameters []Parameter `json:"stanza:parameters"`
	Definition string      `json:"stanza:definition"`
	Usage      string      `json:"stanza:usage"`
}

func (meta *Metadata) ParameterKeys() []string {
	keys := make([]string, len(meta.Parameters))
	for i, parameter := range meta.Parameters {
		keys[i] = parameter.Key
	}
	return keys
}

func LoadMetadata(metadataPath string) (*Metadata, error) {
	f, err := os.Open(metadataPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	var meta Metadata
	if err := decoder.Decode(&meta); err != nil {
		return nil, err
	}

	return &meta, nil
}

func NewStanza(baseDir, name string) (*Stanza, error) {
	st := &Stanza{
		BaseDir: baseDir,
		Name:    name,
	}
	if !st.MetadataExists() {
		return nil, nil
	}
	meta, err := LoadMetadata(st.MetadataPath())
	if err != nil {
		return nil, err
	}
	st.Metadata = *meta

	return st, nil
}

func (st *Stanza) MetadataPath() string {
	return path.Join(st.BaseDir, "metadata.json")
}

func (st *Stanza) MetadataExists() bool {
	_, err := os.Stat(st.MetadataPath())
	if err != nil {
		return false
	}
	return true
}

func (st *Stanza) TemplateGlobPattern() string {
	return path.Join(st.BaseDir, "templates/*")
}

func (st *Stanza) IndexJsPath() string {
	return path.Join(st.BaseDir, "index.js")
}

func (st *Stanza) AssetsDir() string {
	return path.Join(st.BaseDir, "assets")
}

func (st *Stanza) DestMetadataPath(destStanzaBase string) string {
	return path.Join(destStanzaBase, "metadata.json")
}

func (st *Stanza) DestIndexHtmlPath(destStanzaBase string) string {
	return path.Join(destStanzaBase, "index.html")
}

func (st *Stanza) DestHelpHtmlPath(destStanzaBase string) string {
	return path.Join(destStanzaBase, "help.html")
}

func (st *Stanza) DestAssetsDir(destStanzaBase string) string {
	return path.Join(destStanzaBase, "assets")
}

func (st *Stanza) ElementName() string {
	return "togostanza-" + st.Name
}

func (st *Stanza) IndexJs() ([]byte, error) {
	f, err := os.Open(st.IndexJsPath())
	if err != nil {
		return nil, err
	}
	defer f.Close()

	js, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return js, nil
}

func (st *Stanza) Build(destStanzaBase string) error {
	if err := os.MkdirAll(destStanzaBase, os.FileMode(0755)); err != nil {
		return err
	}
	if err := st.buildIndexHtml(destStanzaBase); err != nil {
		return err
	}
	if err := st.buildHelpHtml(destStanzaBase); err != nil {
		return err
	}
	if err := st.copyMetadataJson(destStanzaBase); err != nil {
		return err
	}
	if err := st.copyAssets(destStanzaBase); err != nil {
		return err
	}
	return nil
}

func (st *Stanza) copyMetadataJson(destStanzaBase string) error {
	destPath := st.DestMetadataPath(destStanzaBase)

	if err := copyFile(destPath, st.MetadataPath()); err != nil {
		return err
	}

	log.Printf("copied to %s", destPath)

	return nil
}

func copyFile(dest, src string) error {
	fin, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fin.Close()
	fout, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer fout.Close()
	_, err = io.Copy(fout, fin)
	return err
}

func (st *Stanza) copyAssets(destStanzaBase string) error {
	if _, err := os.Stat(st.AssetsDir()); os.IsNotExist(err) {
		return nil
	}
	return filepath.Walk(st.AssetsDir(), func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(st.AssetsDir(), srcPath)
		if err != nil {
			return err
		}
		destPath := path.Join(st.DestAssetsDir(destStanzaBase), rel)
		if info.Mode().IsDir() {
			if err := os.MkdirAll(destPath, os.FileMode(0755)); err != nil {
				return err
			}
			log.Printf("created directory %s", destPath)
		} else {
			if err := copyFile(destPath, srcPath); err != nil {
				return err
			}
			log.Printf("copied to %s", destPath)
		}
		return nil
	})
}

func (st *Stanza) buildIndexHtml(destStanzaBase string) error {
	data, err := Asset("data/index.html")
	if err != nil {
		return fmt.Errorf("asset not found")
	}

	tmpl, err := template.New("index").Parse(string(data))
	if err != nil {
		return err
	}

	templates := make(map[string]string)

	paths, err := filepath.Glob(st.TemplateGlobPattern())

	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		t, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		templates[filepath.Base(path)] = string(t)
	}

	indexJs, err := st.IndexJs()
	if err != nil {
		return err
	}

	stylesheet, err := Asset("data/stanza.css")
	if err != nil {
		return err
	}

	descriptor := struct {
		Templates   map[string]string `json:"templates"`
		Parameters  []string          `json:"parameters"`
		ElementName string            `json:"elementName"`
		Stylesheet  string            `json:"stylesheet"`
	}{
		Templates:   templates,
		Parameters:  st.Metadata.ParameterKeys(),
		ElementName: st.ElementName(),
		Stylesheet:  string(stylesheet),
	}
	descriptorJson, err := json.Marshal(descriptor)
	if err != nil {
		return err
	}

	b := struct {
		IndexJs        string
		DescriptorJson string
	}{
		IndexJs:        string(indexJs),
		DescriptorJson: string(descriptorJson),
	}

	destPath := st.DestIndexHtmlPath(destStanzaBase)
	w, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer w.Close()

	if err := tmpl.Execute(w, b); err != nil {
		return err
	}

	log.Printf("generated %s", destPath)

	return nil
}

func (st *Stanza) buildHelpHtml(destStanzaBase string) error {
	data, err := Asset("data/help.html")
	if err != nil {
		return fmt.Errorf("asset not found")
	}

	stylesheet, err := Asset("data/stanza.css")
	if err != nil {
		return err
	}

	tmpl, err := template.New("help.html").Parse(string(data))
	if err != nil {
		return err
	}

	destPath := st.DestHelpHtmlPath(destStanzaBase)
	w, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer w.Close()

	context := struct {
		Name       string
		Metadata   Metadata
		Stylesheet string
	}{
		Name:       st.Name,
		Metadata:   st.Metadata,
		Stylesheet: string(stylesheet),
	}

	if err := tmpl.Execute(w, context); err != nil {
		return err
	}

	log.Printf("generated %s", destPath)

	return nil
}
