package provider

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"text/template"
	"time"

	"github.com/togostanza/ts/stanza"
)

//go:generate go-bindata -pkg=provider data/... assets/...

type StanzaProvider struct {
	baseDir      string
	stanzas      map[string]*stanza.Stanza
	lastModified time.Time
}

func New(baseDir string) (*StanzaProvider, error) {
	sp := StanzaProvider{
		baseDir: baseDir,
	}

	return &sp, nil
}

func (sp *StanzaProvider) LastModified() (time.Time, error) {
	var t time.Time
	err := filepath.Walk(sp.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(sp.baseDir, path)
		if err != nil {
			return err
		}
		if rel == "dist" {
			return filepath.SkipDir
		}
		mt := info.ModTime()
		if mt.After(t) {
			t = mt
		}
		return nil
	})
	return t, err
}

func (sp *StanzaProvider) Load() error {
	stanzaMetadataPaths, err := filepath.Glob(path.Join(sp.baseDir, "*/metadata.json"))
	if err != nil {
		return err
	}

	stanzas := make(map[string]*stanza.Stanza)
	for _, stanzaMetadataPath := range stanzaMetadataPaths {
		stanzaPath := filepath.Dir(stanzaMetadataPath)
		stanzaName := filepath.Base(stanzaPath)
		log.Printf("loading stanza %s", stanzaPath)
		stanza, err := stanza.NewStanza(path.Join(stanzaPath), stanzaName)
		if err != nil {
			return err
		}
		stanzas[stanzaName] = stanza
	}
	sp.stanzas = stanzas

	return nil
}

func (sp *StanzaProvider) build(distDir string) error {
	t0 := time.Now()

	if err := sp.Load(); err != nil {
		return err
	}

	if sp.NumStanzas() == 0 {
		return fmt.Errorf("no stanzas available under %s", sp.baseDir)
	}

	if err := os.MkdirAll(distDir, os.FileMode(0755)); err != nil {
		return err
	}

	if err := sp.buildStanzas(distDir); err != nil {
		return err
	}
	if err := sp.extractAssets(distDir); err != nil {
		return err
	}
	if err := sp.buildList(distDir); err != nil {
		return err
	}

	log.Println("built in", time.Since(t0))
	return nil
}

func (sp *StanzaProvider) Build(distDir string) error {
	lm, err := sp.LastModified()
	if err != nil {
		return err
	}
	sp.lastModified = lm

	if err := sp.build(distDir); err != nil {
		return err
	}
	return nil
}

func (sp *StanzaProvider) RebuildIfRequired(distDir string) error {
	lm, err := sp.LastModified()
	if err != nil {
		return err
	}

	if lm.After(sp.lastModified) {
		sp.lastModified = lm
		log.Println("update detected; rebuilding ...")
		if err := sp.Build(distDir); err != nil {
			return err
		}
	}
	return nil
}

func (sp *StanzaProvider) buildStanzas(distDir string) error {
	log.Println("building stanzas")
	numBuilt := 0
	for name, stanza := range sp.stanzas {
		destStanzaBase := path.Join(distDir, name)
		if err := stanza.Build(destStanzaBase); err != nil {
			return err
		}
		numBuilt++
	}

	log.Printf("%d stanza(s) built", numBuilt)
	return nil
}

func (sp *StanzaProvider) buildList(distDir string) error {
	tmpl := MustTemplateAsset("data/list.html")

	destPath := path.Join(distDir, "index.html")
	w, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer w.Close()

	context := struct {
		Stanzas []*stanza.Stanza
	}{
		Stanzas: sp.Stanzas(),
	}

	if err := tmpl.Execute(w, context); err != nil {
		return err
	}

	log.Printf("generated %s", destPath)

	return nil
}

func (sp *StanzaProvider) extractAssets(distStanzaPath string) error {
	assetsDir := "assets"
	log.Printf("generating assets under %s", path.Join(distStanzaPath, assetsDir))
	return RestoreAssets(distStanzaPath, assetsDir)
}

func (sp *StanzaProvider) Stanzas() []*stanza.Stanza {
	stanzas := make([]*stanza.Stanza, len(sp.stanzas))
	i := 0
	for _, stanza := range sp.stanzas {
		stanzas[i] = stanza
		i++
	}
	return stanzas
}

func (sp *StanzaProvider) Stanza(name string) *stanza.Stanza {
	return sp.stanzas[name]
}

func (sp *StanzaProvider) NumStanzas() int {
	return len(sp.stanzas)
}

func MustTemplateAsset(path string) *template.Template {
	data, err := Asset(path)
	if err != nil {
		panic("asset " + path + "not found")
	}

	return template.Must(template.New(path).Parse(string(data)))
}
