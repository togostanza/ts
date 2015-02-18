package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
)

//go:generate go-bindata data/ assets/...

var flagPort int
var flagStanzaBaseDir string

type Command struct {
	Run       func(cmd *Command, args []string)
	Name      string
	Short     string
	Long      string
	UsageLine string
	Flag      flag.FlagSet
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(c.Long))
	os.Exit(2)
}

var commands = []*Command{
	cmdBuild,
	cmdServer,
}

const usageTemplate = `Usage:

	ts command [arguments]

The commands are:
{{range .}}
	{{.Name | printf "%-10s"}} {{.Short}}{{end}}
`

var helpTemplate = `usage: ts {{.UsageLine}}

{{.Long}}
`

func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

func usage() {
	tmpl(os.Stderr, usageTemplate, commands)
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	if args[0] == "help" {
		if len(args) == 1 {
			usage()
		}
		arg := args[1]
		for _, cmd := range commands {
			if cmd.Name == arg {
				tmpl(os.Stdout, helpTemplate, cmd)
				os.Exit(2)
			}
		}
		fmt.Fprintf(os.Stderr, "Unknown command %#q. Run 'ts help'\n", arg)
		os.Exit(2)
	}

	for _, cmd := range commands {
		if cmd.Name == args[0] {
			cmd.Flag.Usage = func() { cmd.Usage() }
			cmd.Flag.Parse(args[1:])
			args = cmd.Flag.Args()
			cmd.Run(cmd, args)
			os.Exit(0)
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown command %#q. Run 'ts help'\n", args[0])
	os.Exit(2)
}
