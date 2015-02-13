package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
)

var cmdServer = &Command{
	Run:       runServer,
	Name:      "server",
	Short:     "run server",
	UsageLine: "server [-port port] [-stanza-base-dir dir]",
	Long:      "Run ts server for development",
}

func init() {
	cmdServer.Flag.IntVar(&flagPort, "port", 8080, "port to listen on")
	cmdServer.Flag.StringVar(&flagStanzaBaseDir, "stanza-base-dir", ".", "stanza base directory")
}

func runServer(cmd *Command, args []string) {
	mux := http.NewServeMux()
	assetsHandler := http.FileServer(http.Dir(flagStanzaBaseDir))

	sp, err := NewStanzaProvider(flagStanzaBaseDir)
	if err != nil {
		log.Fatal(err)
	}
	if err := sp.Build(); err != nil {
		log.Fatal(err)
	}

	assetsDir := "assets"
	log.Printf("generating assets under %s", path.Join(flagStanzaBaseDir, assetsDir))
	if err := RestoreAssets(flagStanzaBaseDir, assetsDir); err != nil {
		log.Fatal(err)
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		assetsHandler.ServeHTTP(w, req)
	})

	addr := fmt.Sprintf(":%d", flagPort)
	log.Println("listening on", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
