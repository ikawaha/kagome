package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ikawaha/kagome/v2/cmd/kagome/lattice"
	"github.com/ikawaha/kagome/v2/cmd/kagome/server"
	"github.com/ikawaha/kagome/v2/cmd/kagome/tokenize"
)

type subcommand struct {
	Name          string
	Description   string
	Run           func([]string) error
	Usage         func()
	OptionCheck   func([]string) error
	PrintDefaults func(flag.ErrorHandling)
}

var subcommands = []subcommand{
	{
		Name:          tokenize.CommandName,
		Description:   tokenize.Description,
		Run:           tokenize.Run,
		Usage:         tokenize.Usage,
		OptionCheck:   tokenize.OptionCheck,
		PrintDefaults: tokenize.PrintDefaults,
	},
	{
		Name:          server.CommandName,
		Description:   server.Description,
		Run:           server.Run,
		Usage:         server.Usage,
		OptionCheck:   server.OptionCheck,
		PrintDefaults: server.PrintDefaults,
	},
	{
		Name:          lattice.CommandName,
		Description:   lattice.Description,
		Run:           lattice.Run,
		Usage:         lattice.Usage,
		OptionCheck:   lattice.OptionCheck,
		PrintDefaults: lattice.PrintDefaults,
	},
	{
		Name:        "version",
		Description: "show version",
		Run: func([]string) error {
			fmt.Fprintf(os.Stderr, "%s\n", version)
			return nil
		},
		Usage:         func() {},
		OptionCheck:   func([]string) error { return nil },
		PrintDefaults: func(flag.ErrorHandling) {},
	},
}

var (
	// version is the app version.
	version = `!!version undefined!!
This must be specified by -X option during the go build. Such like:
	$ go build --ldflags "-X 'main.version=$(git describe --tag)'"`

	errorWriter       = os.Stderr
	defaultSubcommand = subcommands[0]
)

//Usage prints to stdout information about the tool
func Usage() {
	fmt.Fprintf(errorWriter, "Japanese Morphological Analyzer -- github.com/ikawaha/kagome\n")
	fmt.Fprintf(errorWriter, "usage: %s <command>\n", filepath.Base(os.Args[0]))
}

// PrintDefaults prints out the default flags
func PrintDefaults() {
	fmt.Fprintln(errorWriter, "The commands are:")
	for _, c := range subcommands {
		if c.Name == defaultSubcommand.Name {
			fmt.Fprintf(errorWriter, "   [%s] - %s (*default)\n", c.Name, c.Description)
		} else {
			fmt.Fprintf(errorWriter, "   %s - %s\n", c.Name, c.Description)
		}
	}
}

func main() {
	var (
		cmd     func([]string) error
		options []string
	)
	if len(os.Args) >= 2 {
		options = os.Args[2:]
		for i := range subcommands {
			if os.Args[1] == subcommands[i].Name {
				cmd = subcommands[i].Run
			}
		}
	}
	if cmd == nil {
		options = os.Args[1:]
		if err := defaultSubcommand.OptionCheck(options); err != nil {
			Usage()
			PrintDefaults()
			fmt.Fprintln(errorWriter)
			defaultSubcommand.Usage()
			defaultSubcommand.PrintDefaults(flag.ExitOnError)
			os.Exit(1)
		}
		cmd = defaultSubcommand.Run
	}
	if err := cmd(options); err != nil {
		fmt.Fprintf(errorWriter, "%v\n", err)
		os.Exit(1)
	}
}
