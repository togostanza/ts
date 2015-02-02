package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

type Stanza struct {
	Name string
}

func NewStanza(name string) *Stanza {
	return &Stanza{
		Name: name,
	}
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

	paths, err := filepath.Glob("gene-attributes/templates/*")

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

	f, err := os.Open("gene-attributes/index.js")
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
		ElementName:   "togostanza-gene-attributes",
	}

	return tmpl.Execute(w, b)
}

func main() {
	mux := http.NewServeMux()
	assetsHandler := http.FileServer(http.Dir("."))
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/gene-attributes/" {
			st := NewStanza("gene-attributes")
			err := st.Generate(w)
			if err != nil {
				log.Println("ERROR", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		} else {
			assetsHandler.ServeHTTP(w, req)
		}
	})

	port := 8080
	addr := fmt.Sprintf(":%d", port)
	log.Println("listening on", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
