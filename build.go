package main

import (
	"log"
	"path"
	"path/filepath"

	"github.com/togostanza/ts/provider"
)

var cmdBuild = &Command{
	Run:       runBuild,
	Name:      "build",
	Short:     "build stanza provider",
	UsageLine: "server [-stanza-base-dir dir]",
	Long:      "Build stanza provider",
}

func addBuildFlags(cmd *Command) {
	path := "."
	if absolutePath, err := filepath.Abs(path); err == nil {
		path = absolutePath
	}
	cmd.Flag.StringVar(&flagStanzaBaseDir, "stanza-base-dir", path, "stanza base directory")
}

func init() {
	addBuildFlags(cmdBuild)
}

func runBuild(cmd *Command, args []string) {
	sp, err := provider.NewStanzaProvider(flagStanzaBaseDir)
	if err != nil {
		log.Fatal(err)
	}
	distPath := path.Join(flagStanzaBaseDir, "dist")
	distStanzaPath := path.Join(distPath, "stanza")
	if err := sp.Build(distStanzaPath); err != nil {
		log.Fatal(err)
	}

	assetsDir := "assets"
	log.Printf("generating assets under %s", path.Join(distStanzaPath, assetsDir))
	if err := RestoreAssets(distStanzaPath, assetsDir); err != nil {
		log.Fatal(err)
	}
}
