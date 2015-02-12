package main

import (
	"path"
)

type StanzaProvider struct {
	baseDir string
}

func NewStanzaProvider(baseDir string) *StanzaProvider {
	sp := StanzaProvider{
		baseDir: baseDir,
	}

	return &sp
}

func (sp *StanzaProvider) Stanza(name string) (*Stanza, error) {
	return NewStanza(path.Join(sp.baseDir, name), name)
}
