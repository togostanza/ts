package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
	"time"
)

var cmdNew = &Command{
	Name:      "new",
	Short:     "create a new stanza",
	UsageLine: "new [-stanza-base-dir dir] [stanza name]",
	Long:      "Create a new stanza",
}

func init() {
	cmdNew.Run = runNew // break init loop
	addBuildFlags(cmdNew)
}

type NewStanzaParameters struct {
	Name    string
	Created string
	Updated string
}

func newExtractStanzaAsset(dir, name string, st *NewStanzaParameters) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	t := template.Must(template.New(name).Parse(string(data)))

	s := strings.SplitN(name, "/", 2)
	if len(s) != 2 {
		fmt.Errorf("unexpected name: %s", name)
	}
	destName := s[1]
	err = os.MkdirAll(_filePath(dir, path.Dir(destName)), os.FileMode(0755))
	if err != nil {
		return err
	}

	destPath := _filePath(dir, destName)
	w, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer w.Close()

	err = t.Execute(w, st)
	if err != nil {
		return err
	}

	log.Printf("wrote %s", destPath)

	return nil
}

func newExtractStanzaAssets(dir, name string, st *NewStanzaParameters) error {
	children, err := AssetDir(name)
	if err != nil { // File
		return newExtractStanzaAsset(dir, name, st)
	} else { // Dir
		for _, child := range children {
			err = newExtractStanzaAssets(dir, path.Join(name, child), st)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func runNew(cmd *Command, args []string) {
	if len(args) != 1 {
		cmdNew.Flag.Usage()
		os.Exit(2)
	}

	stanzaName := args[0]
	stanzaDir := path.Join(flagStanzaBaseDir, stanzaName)
	log.Printf("creating stanza directory %#q", stanzaDir)

	t := time.Now()
	st := NewStanzaParameters{
		Name:    stanzaName,
		Created: t.Format("2006-01-02"),
		Updated: t.Format("2006-01-02"),
	}
	err := newExtractStanzaAssets(stanzaDir, "stanza-template", &st)
	if err != nil {
		log.Fatal(err)
	}
}
