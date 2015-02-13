package main

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
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

func (sp *StanzaProvider) Stanza(name string) *Stanza {
	return sp.stanzas[name]
}

func (sp *StanzaProvider) NumStanzas() int {
	return len(sp.stanzas)
}
