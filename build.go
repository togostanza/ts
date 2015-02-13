package main

import (
	"log"
	"path"
)

var cmdBuild = &Command{
	Run:       runBuild,
	Name:      "build",
	Short:     "build stanza provider",
	UsageLine: "server [-stanza-base-dir dir]",
	Long:      "Build stanza provider",
}

func init() {
	cmdBuild.Flag.StringVar(&flagStanzaBaseDir, "stanza-base-dir", ".", "stanza base directory")
}

func runBuild(cmd *Command, args []string) {
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
}
