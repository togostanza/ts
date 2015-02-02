package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

func main() {
	f, err := os.Open("gene-attributes/index.html.erb")
	if err != nil {
		log.Fatal(err)
	}

	t, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.New("index").Parse(string(t))
	if err != nil {
		log.Fatal(err)
	}

	log.Println(tmpl)

	mux := http.NewServeMux()

	mux.HandleFunc("/gene-attributes/", func(w http.ResponseWriter, req *http.Request) {
		templates := make(map[string]string)

		paths, err := filepath.Glob("gene-attributes/templates/*")

		for _, path := range paths {
			f, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}

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

		js, err := ioutil.ReadAll(f)

		b := struct {
			TemplatesJson string
			IndexJs       string
		}{
			TemplatesJson: string(buffer),
			IndexJs:       string(js),
		}
		err = tmpl.Execute(w, b)
		if err != nil {
			log.Fatal(err)
		}
	})

	mux.Handle("/", http.FileServer(http.Dir(".")))

	http.ListenAndServe(":8080", mux)
}
