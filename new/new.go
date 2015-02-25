package new

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
	"time"
)

//go:generate go-bindata -pkg new blueprint/...

type parameters struct {
	Name    string
	Created string
	Updated string
}

func extractBlueprintAsset(dir, name string, params *parameters) error {
	t := MustTemplateAsset(name)

	s := strings.SplitN(name, "/", 2)
	if len(s) != 2 {
		fmt.Errorf("unexpected name: %s", name)
	}
	destName := s[1]
	err := os.MkdirAll(_filePath(dir, path.Dir(destName)), os.FileMode(0755))
	if err != nil {
		return err
	}

	destPath := _filePath(dir, destName)
	w, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer w.Close()

	err = t.Execute(w, params)
	if err != nil {
		return err
	}

	log.Printf("wrote %s", destPath)

	return nil
}

func extractBlueprintAssets(dir, name string, params *parameters) error {
	children, err := AssetDir(name)
	if err != nil { // File
		return extractBlueprintAsset(dir, name, params)
	} else { // Dir
		for _, child := range children {
			err = extractBlueprintAssets(dir, path.Join(name, child), params)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Generate(stanzaName string, stanzaBaseDir string) error {
	stanzaDir := path.Join(stanzaBaseDir, stanzaName)
	log.Printf("creating stanza directory %#q", stanzaDir)

	t := time.Now()
	params := parameters{
		Name:    stanzaName,
		Created: t.Format("2006-01-02"),
		Updated: t.Format("2006-01-02"),
	}
	return extractBlueprintAssets(stanzaDir, "blueprint", &params)
}

func MustTemplateAsset(path string) *template.Template {
	data, err := Asset(path)
	if err != nil {
		panic("asset " + path + "not found")
	}

	return template.Must(template.New(path).Parse(string(data)))
}
