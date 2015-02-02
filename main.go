package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

func main() {
	data, err := Asset("data/template.html")
	if err != nil {
		log.Fatal("asset not found")
	}

	tmpl, err := template.New("index").Parse(string(data))
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/gene-attributes/", func(w http.ResponseWriter, req *http.Request) {
		templates := make(map[string]string)

		paths, err := filepath.Glob("gene-attributes/templates/*")

		for _, path := range paths {
			f, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			t, err := ioutil.ReadAll(f)
			if err != nil {
				log.Fatal(err)
			}

			templates[filepath.Base(path)] = string(t)
		}

		buffer, err := json.Marshal(templates)
		if err != nil {
			log.Fatal(err)
		}

		f, err := os.Open("gene-attributes/index.js")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		js, err := ioutil.ReadAll(f)

		b := struct {
			TemplatesJson string
			IndexJs       string
			ElementName   string
		}{
			TemplatesJson: string(buffer),
			IndexJs:       string(js),
			ElementName:   "togostanza-gene-attributes",
		}
		err = tmpl.Execute(w, b)
		if err != nil {
			log.Fatal(err)
		}
	})

	mux.Handle("/", http.FileServer(http.Dir(".")))

	port := 8080
	addr := fmt.Sprintf(":%d", port)
	log.Println("listening on", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
