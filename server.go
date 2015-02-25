package main

import (
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/togostanza/ts/provider"
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
	addBuildFlags(cmdServer)
}

func runServer(cmd *Command, args []string) {
	sp, err := provider.New(flagStanzaBaseDir)
	if err != nil {
		log.Fatal(err)
	}

	distPath := path.Join(flagStanzaBaseDir, "dist")
	distStanzaPath := path.Join(distPath, "stanza")
	if err := sp.Build(distStanzaPath); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	assetsHandler := http.FileServer(http.Dir(distPath))

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/" {
			http.Redirect(w, req, "/stanza/", http.StatusFound)
			return
		}
		assetsHandler.ServeHTTP(w, req)
	})

	addr := fmt.Sprintf(":%d", flagPort)
	log.Println("listening on", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
