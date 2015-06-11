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
	"fmt"
	"os"
	"path"

	"github.com/ikawaha/kagome/cmd/kagome/lattice"
	"github.com/ikawaha/kagome/cmd/kagome/server"
	"github.com/ikawaha/kagome/cmd/kagome/tokenize"
)

var errorWriter = os.Stderr

var subcommands = []struct {
	Name        string
	Description string
	Run         func([]string) error
}{
	// subcommands
	{tokenize.CommandName, tokenize.Description, tokenize.Run},
	{server.CommandName, server.Description, server.Run},
	{lattice.CommandName, lattice.Description, lattice.Run},
}

func Usage() {
	fmt.Fprintf(errorWriter, "Japanese Morphological Analyzer -- github.com/ikawaha/kagome\n")
	fmt.Fprintf(errorWriter, "usage: %s <command>\n", path.Base(os.Args[0]))
}

func PrintDefaults() {
	fmt.Fprintln(errorWriter, "The commands are:")
	for _, c := range subcommands {
		fmt.Fprintf(errorWriter, "   %s - %s\n", c.Name, c.Description)
	}
}

func main() {
	if len(os.Args) < 2 {
		Usage()
		PrintDefaults()
		os.Exit(1)
	}
	var cmd func([]string) error
	for i := range subcommands {
		if os.Args[1] == subcommands[i].Name {
			cmd = subcommands[i].Run
		}
	}
	if cmd == nil {
		Usage()
		PrintDefaults()
		os.Exit(1)
	}
	if e := cmd(os.Args[2:]); e != nil {
		fmt.Fprintf(errorWriter, "%v\n", e)
		os.Exit(1)
	}
}
