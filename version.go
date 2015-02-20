package main

import (
	"fmt"
	"runtime"
)

var cmdVersion = &Command{
	Name:      "version",
	Short:     "print ts version",
	UsageLine: "version",
	Long:      "Print ts version",
	Run:       runVersion,
}

func runVersion(cmd *Command, args []string) {
	fmt.Printf("ts version %s (%s %s/%s)\n", VERSION, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
