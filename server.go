package main

import (
	"fmt"
	"log"
	"net/http"
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
	assetsHandler := http.FileServer(http.Dir(flagStanzaBaseDir))

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		assetsHandler.ServeHTTP(w, req)
	})

	addr := fmt.Sprintf(":%d", flagPort)
	log.Println("listening on", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
