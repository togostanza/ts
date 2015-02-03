package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

var flagPort int

func init() {
	flag.IntVar(&flagPort, "port", 8080, "port to listen on")
}

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	assetsHandler := http.FileServer(http.Dir(cwd))

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		stanzaName := strings.TrimSuffix(strings.TrimPrefix(req.URL.Path, "/"), "/")
		st, err := NewStanza(path.Join(cwd, stanzaName), stanzaName)
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
