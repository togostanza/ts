package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"text/template"
	"time"
)

type StanzaProvider struct {
	baseDir string
	stanzas map[string]*Stanza
}

func NewStanzaProvider(baseDir string) (*StanzaProvider, error) {
	sp := StanzaProvider{
		baseDir: baseDir,
	}

	if err := sp.Load(); err != nil {
		return nil, err
	}

	return &sp, nil
}

func (sp *StanzaProvider) Load() error {
	stanzaMetadataPaths, err := filepath.Glob(path.Join(sp.baseDir, "*/metadata.json"))
	if err != nil {
		return err
	}

	stanzas := make(map[string]*Stanza)
	for _, stanzaMetadataPath := range stanzaMetadataPaths {
		stanzaPath := filepath.Dir(stanzaMetadataPath)
		stanzaName := filepath.Base(stanzaPath)
		log.Printf("loading stanza %s", stanzaPath)
		stanza, err := NewStanza(path.Join(stanzaPath), stanzaName)
		if err != nil {
			return err
		}
		stanzas[stanzaName] = stanza
	}
	sp.stanzas = stanzas

	return nil
}

func (sp *StanzaProvider) Build() error {
	if err := sp.buildStanzas(); err != nil {
		return err
	}
	if err := sp.buildList(); err != nil {
		return err
	}
	return nil
}

func (sp *StanzaProvider) buildStanzas() error {
	log.Println("building stanzas")
	t0 := time.Now()
	numBuilt := 0
	for _, stanza := range sp.stanzas {
		if err := stanza.Build(); err != nil {
			return err
		}
		numBuilt++
	}

	if numBuilt == 0 {
		return fmt.Errorf("no stanzas available under %s", sp.baseDir)
	}

	log.Printf("%d stanza(s) built in %s", numBuilt, time.Since(t0))
	return nil
}

func (sp *StanzaProvider) IndexPath() string {
	return path.Join(sp.baseDir, "index.html")
}

func (sp *StanzaProvider) buildList() error {
	data, err := Asset("data/list.html")
	if err != nil {
		return fmt.Errorf("asset list.html not found")
	}
	tmpl, err := template.New("list.html").Parse(string(data))
	if err != nil {
		return err
	}

	destPath := sp.IndexPath()
	w, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer w.Close()

	context := struct {
		Stanzas []*Stanza
	}{
		Stanzas: sp.Stanzas(),
	}

	if err := tmpl.Execute(w, context); err != nil {
		return err
	}

	log.Printf("generated %s", destPath)

	return nil
}

func (sp *StanzaProvider) Stanzas() []*Stanza {
	stanzas := make([]*Stanza, len(sp.stanzas))
	i := 0
	for _, stanza := range sp.stanzas {
		stanzas[i] = stanza
		i++
	}
	return stanzas
}

func (sp *StanzaProvider) Stanza(name string) *Stanza {
	return sp.stanzas[name]
}

func (sp *StanzaProvider) NumStanzas() int {
	return len(sp.stanzas)
}
