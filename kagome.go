package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/ikawaha/kagome/v2/cmd/lattice"
	"github.com/ikawaha/kagome/v2/cmd/sentence"
	"github.com/ikawaha/kagome/v2/cmd/server"
	"github.com/ikawaha/kagome/v2/cmd/tokenize"
)

type subcommand struct {
	Name          string
	Description   string
	Run           func(context.Context, []string) error
	Usage         func()
	OptionCheck   func([]string) error
	PrintDefaults func(flag.ErrorHandling)
}

var (
	version     string // eg. go build --ldflags "-X 'main.version=$(git describe --tag)'"
	errorWriter = os.Stderr
	subcommands = []subcommand{
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
			Name:          sentence.CommandName,
			Description:   sentence.Description,
			Run:           sentence.Run,
			Usage:         sentence.Usage,
			OptionCheck:   sentence.OptionCheck,
			PrintDefaults: sentence.PrintDefaults,
		},
		{
			Name:        "version",
			Description: "show version",
			Run: func(context.Context, []string) error {
				ShowVersion()
				return nil
			},
			Usage:         func() {},
			OptionCheck:   func([]string) error { return nil },
			PrintDefaults: func(flag.ErrorHandling) {},
		},
	}
	defaultSubcommand = subcommands[0]
)

// Usage prints information about the tool
func Usage() {
	fmt.Fprintf(errorWriter, "Japanese Morphological Analyzer -- github.com/ikawaha/kagome/v2\n")
	fmt.Fprintf(errorWriter, "usage: %s <command>\n", filepath.Base(os.Args[0]))
}

// ShowVersion prints the version about the tool.
func ShowVersion() {
	info, ok := debug.ReadBuildInfo()
	if ok && version == "" {
		version = info.Main.Version
	}
	if version == "" {
		version = "(devel)"
	}
	fmt.Fprintln(errorWriter, version)
	if !ok {
		return
	}
	const prefix = "github.com/ikawaha/kagome-dict/"
	for _, v := range info.Deps {
		if strings.HasPrefix(v.Path, prefix) {
			fmt.Fprintln(errorWriter, "  ", v.Path[len(prefix):], v.Version)
		}
	}
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
	var run     func(context.Context, []string) error
	var options []string
	if len(os.Args) >= 2 {
		options = os.Args[2:]
		for i := range subcommands {
			if os.Args[1] == subcommands[i].Name {
				run = subcommands[i].Run
			}
		}
	}
	if run == nil {
		options = os.Args[1:]
		if err := defaultSubcommand.OptionCheck(options); err != nil {
			Usage()
			PrintDefaults()
			fmt.Fprintln(errorWriter)
			defaultSubcommand.Usage()
			defaultSubcommand.PrintDefaults(flag.ExitOnError)
			os.Exit(1)
		}
		run = defaultSubcommand.Run
	}
	if err := run(context.Background(), options); err != nil {
		fmt.Fprintf(errorWriter, "%v\n", err)
		os.Exit(1)
	}
}
