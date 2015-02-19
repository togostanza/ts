package main

import (
	"log"
	"path"
	"path/filepath"
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
	sp, err := NewStanzaProvider(flagStanzaBaseDir)
	if err != nil {
		log.Fatal(err)
	}
	distPath := path.Join(flagStanzaBaseDir, "dist")
	if err := sp.Build(distPath); err != nil {
		log.Fatal(err)
	}

	assetsDir := "assets"
	log.Printf("generating assets under %s", path.Join(distPath, assetsDir))
	if err := RestoreAssets(distPath, assetsDir); err != nil {
		log.Fatal(err)
	}
}
