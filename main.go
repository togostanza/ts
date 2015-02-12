package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

//go:generate go-bindata data/

var flagPort int
var flagStanzaBaseDir string

func init() {
	flag.IntVar(&flagPort, "port", 8080, "port to listen on")
	flag.StringVar(&flagStanzaBaseDir, "stanza-base-dir", ".", "stanza base directory")
}

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	assetsHandler := http.FileServer(http.Dir(flagStanzaBaseDir))

	sp := NewStanzaProvider(flagStanzaBaseDir)

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		stanzaName := strings.TrimSuffix(strings.TrimPrefix(req.URL.Path, "/"), "/")
		st, err := sp.Stanza(stanzaName)
		if err != nil {
			log.Println("ERROR", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		if st == nil {
			assetsHandler.ServeHTTP(w, req)
			return
		}
		if err := st.Generate(w); err != nil {
			log.Println("ERROR", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	})

	addr := fmt.Sprintf(":%d", flagPort)
	log.Println("listening on", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
