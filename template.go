package main

import (
	"text/template"
)

func MustTemplateAsset(path string) *template.Template {
	data, err := Asset(path)
	if err != nil {
		panic("asset " + path + "not found")
	}

	return template.Must(template.New(path).Parse(string(data)))
}
