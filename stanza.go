package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

type Stanza struct {
	BaseDir string
	Name    string
	Metadata
}

type Parameter struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	Example     string `json:"example"`
	Required    bool   `json:"required"`
}

type Metadata struct {
	Id          string      `json:"@id"`
	Label       string      `json:"label"`
	Parameters  []Parameter `json:"parameters"`
	Description string      `json:"description"`
	Usage       string      `json:"usage"`
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

func (st *Stanza) IndexHtmlPath() string {
	return path.Join(st.BaseDir, "index.html")
}

func (st *Stanza) HelpHtmlPath() string {
	return path.Join(st.BaseDir, "help.html")
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

func (st *Stanza) Build() error {
	if err := st.buildIndexHtml(); err != nil {
		return err
	}
	if err := st.buildHelpHtml(); err != nil {
		return err
	}
	return nil
}

func (st *Stanza) buildIndexHtml() error {
	data, err := Asset("data/template.html")
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

	buffer, err := json.Marshal(templates)
	if err != nil {
		return err
	}

	indexJs, err := st.IndexJs()
	if err != nil {
		return err
	}

	stylesheet, err := Asset("data/stanza.css")
	if err != nil {
		return err
	}

	b := struct {
		TemplatesJson    string
		IndexJs          string
		ElementName      string
		AttributesString string
		Attributes       []string
		Stylesheet       string
	}{
		TemplatesJson:    string(buffer),
		IndexJs:          string(indexJs),
		ElementName:      "togostanza-" + st.Name,
		AttributesString: strings.Join(st.Metadata.ParameterKeys(), " "),
		Attributes:       st.Metadata.ParameterKeys(),
		Stylesheet:       string(stylesheet),
	}

	destPath := st.IndexHtmlPath()
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

func (st *Stanza) buildHelpHtml() error {
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

	destPath := st.HelpHtmlPath()
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

	return nil
}
