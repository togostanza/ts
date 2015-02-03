package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
	Key string `json:"key"`
}

type Metadata struct {
	Parameters []Parameter `json:"parameters"`
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

func (st *Stanza) Generate(w io.Writer) error {
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

	f, err := os.Open(st.IndexJsPath())
	if err != nil {
		return err
	}
	defer f.Close()

	js, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	b := struct {
		TemplatesJson string
		IndexJs       string
		ElementName   string
	}{
		TemplatesJson: string(buffer),
		IndexJs:       string(js),
		ElementName:   "togostanza-" + st.Name,
	}

	return tmpl.Execute(w, b)
}
