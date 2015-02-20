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
	addBuildFlags(cmdServer)
}

func runServer(cmd *Command, args []string) {
	runBuild(nil, nil)

	mux := http.NewServeMux()
	distPath := path.Join(flagStanzaBaseDir, "dist")
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
