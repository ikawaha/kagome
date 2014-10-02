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
