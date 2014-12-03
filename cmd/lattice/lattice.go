//  Copyright (c) 2014 ikawaha.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ikawaha/kagome"
)

var usageMessage = "usage: lattice [-u userdic_file] [-o output_file] [-v] sentence"

func usage() {
	fmt.Fprintln(os.Stderr, usageMessage)
	flag.PrintDefaults()
	os.Exit(0)
}

var (
	fUserDicFile = flag.String("u", "", "user dic")
	fOutputFile  = flag.String("o", "", "output file")
	fVerbose     = flag.Bool("v", false, "verbose mode")
)

func Main(input string) {
	if input == "" {
		usage()
	}
	var out = os.Stdout
	if *fOutputFile != "" {
		var err error
		out, err = os.OpenFile(*fOutputFile, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer out.Close()
	}

	t := kagome.NewTokenizer()
	if *fUserDicFile != "" {
		if udic, err := kagome.NewUserDic(*fUserDicFile); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		} else {
			t.SetUserDic(udic)
		}
	}

	tokens := t.Dot(input, out)
	if *fVerbose {
		for i, size := 1, len(tokens); i < size; i++ {
			tok := tokens[i]
			f := tok.Features()
			if tok.Class == kagome.DUMMY {
				fmt.Fprintf(os.Stderr, "%s\n", tok.Surface)
			} else {

				fmt.Fprintf(os.Stderr, "%s\t%v\n", tok.Surface, strings.Join(f, ","))
			}
		}
	}
}

func main() {
	flag.Usage = usage
	flag.Parse()
	Main(flag.Arg(0))
}
