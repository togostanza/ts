package main

import (
	"log"
	"os"

	"github.com/togostanza/ts/new"
)

var cmdNew = &Command{
	Name:      "new",
	Short:     "create a new stanza",
	UsageLine: "new [-stanza-base-dir dir] [stanza name]",
	Long:      "Create a new stanza",
}

func init() {
	cmdNew.Run = runNew // break init loop
	addBuildFlags(cmdNew)
}

func runNew(cmd *Command, args []string) {
	if len(args) != 1 {
		cmdNew.Flag.Usage()
		os.Exit(2)
	}

	stanzaName := args[0]
	err := new.Generate(stanzaName, flagStanzaBaseDir)
	if err != nil {
		log.Fatal(err)
	}
}
