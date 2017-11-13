// Copyright 2015 ikawaha
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// 	You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ikawaha/kagome/cmd/kagome/lattice"
	"github.com/ikawaha/kagome/cmd/kagome/server"
	"github.com/ikawaha/kagome/cmd/kagome/tokenize"
)

var (
	errorWriter = os.Stderr

	subcommands = []struct {
		Name          string
		Description   string
		Run           func([]string) error
		Usage         func()
		OptionCheck   func([]string) error
		PrintDefaults func(flag.ErrorHandling)
	}{
		{
			tokenize.CommandName, tokenize.Description,
			tokenize.Run,
			tokenize.Usage, tokenize.OptionCheck, tokenize.PrintDefaults,
		},
		{
			server.CommandName, server.Description,
			server.Run,
			server.Usage, server.OptionCheck, server.PrintDefaults,
		},
		{
			lattice.CommandName, lattice.Description,
			lattice.Run,
			lattice.Usage, lattice.OptionCheck, lattice.PrintDefaults,
		},
	}

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
		if e := defaultSubcommand.OptionCheck(options); e != nil {
			Usage()
			PrintDefaults()
			fmt.Fprintln(errorWriter)
			defaultSubcommand.Usage()
			defaultSubcommand.PrintDefaults(flag.ExitOnError)
			os.Exit(1)
		}
		cmd = defaultSubcommand.Run
	}
	if e := cmd(options); e != nil {
		fmt.Fprintf(errorWriter, "%v\n", e)
		os.Exit(1)
	}
}
