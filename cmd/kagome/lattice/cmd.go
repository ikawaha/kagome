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

package lattice

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ikawaha/kagome/v2/tokenizer"
)

// subcommand property
var (
	CommandName  = "lattice"
	Description  = `lattice viewer`
	UsageMessage = "%s [-udic userdic_file] [-sysdic (ipa|uni)] [-mode (normal|search|extended)] [-output output_file] [-v] sentence\n"
	ErrorWriter  = os.Stderr
)

// options
type option struct {
	udic    string
	sysdic  string
	mode    string
	output  string
	verbose bool
	input   string
	flagSet *flag.FlagSet
}

// ContinueOnError ErrorHandling // Return a descriptive error.
// ExitOnError                   // Call os.Exit(2).
// PanicOnError                  // Call panic with a descriptive error.flag.ContinueOnError
func newOption(w io.Writer, eh flag.ErrorHandling) (o *option) {
	o = &option{
		flagSet: flag.NewFlagSet(CommandName, eh),
	}
	// option settings
	o.flagSet.StringVar(&o.udic, "udic", "", "user dic")
	o.flagSet.StringVar(&o.sysdic, "sysdic", "ipa", "system dic type (ipa|uni)")
	o.flagSet.StringVar(&o.mode, "mode", "normal", "tokenize mode (normal|search|extended)")
	o.flagSet.StringVar(&o.output, "output", "", "output file")
	o.flagSet.BoolVar(&o.verbose, "v", false, "verbose mode")

	return
}

func (o *option) parse(args []string) error {
	if err := o.flagSet.Parse(args); err != nil {
		return err
	}
	// validations
	if o.flagSet.NArg() == 0 {
		return fmt.Errorf("input is empty")
	}
	if o.sysdic != "" && o.sysdic != "ipa" && o.sysdic != "uni" {
		return fmt.Errorf("invalid argument: -sysdic %v", o.sysdic)
	}
	if o.mode != "" && o.mode != "normal" && o.mode != "search" && o.mode != "extended" {
		return fmt.Errorf("invalid argument: -mode %v", o.mode)
	}
	o.input = strings.Join(o.flagSet.Args(), " ")
	return nil
}

//OptionCheck receives a slice of args and returns an error if it was not successfully parsed
func OptionCheck(args []string) error {
	opt := newOption(ioutil.Discard, flag.ContinueOnError)
	if err := opt.parse(args); err != nil {
		return fmt.Errorf("%v, %v", CommandName, err)
	}
	return nil
}

// command main
func command(opt *option) error {
	var dic *tokenizer.Dic
	if opt.sysdic == "ipa" {
		dic = tokenizer.SysDicIPA()
	} else if opt.sysdic == "uni" {
		dic = tokenizer.SysDicUni()
	} else {
		dic = tokenizer.SysDic()
	}
	t, err := tokenizer.NewWithDic(dic)
	if err != nil {
		return fmt.Errorf("invalid dictionary")
	}
	var out = os.Stdout
	if opt.output != "" {
		var err error
		out, err = os.OpenFile(opt.output, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
		if err != nil {
			fmt.Fprintln(ErrorWriter, err)
			os.Exit(1)
		}
		defer out.Close()
	}
	var udic *tokenizer.UserDic
	if opt.udic != "" {
		var err error
		udic, err = tokenizer.NewUserDic(opt.udic)
		if err != nil {
			return err
		}
		if err := t.SetUserDic(udic); err != nil {
			return err
		}
	}
	if opt.udic != "" {
		udic, err := tokenizer.NewUserDic(opt.udic)
		if err != nil {
			return err
		}
		t.SetUserDic(udic)
	}
	mode := tokenizer.Normal
	switch opt.mode {
	case "normal":
		mode = tokenizer.Normal
	case "search":
		mode = tokenizer.Search
	case "extended":
		mode = tokenizer.Extended
	}
	tokens := t.AnalyzeGraph(out, opt.input, mode)
	if opt.verbose {
		for i, size := 1, len(tokens); i < size; i++ {
			tok := tokens[i]
			f := tok.Features()
			if tok.Class == tokenizer.DUMMY {
				fmt.Fprintf(ErrorWriter, "%s\n", tok.Surface)
			} else {

				fmt.Fprintf(ErrorWriter, "%s\t%v\n", tok.Surface, strings.Join(f, ","))
			}
		}
	}
	return nil
}

// Run receives the slice of args and executes the lattice tool
func Run(args []string) error {
	opt := newOption(ErrorWriter, flag.ExitOnError)
	if err := opt.parse(args); err != nil {
		Usage()
		PrintDefaults(flag.ExitOnError)
		return fmt.Errorf("%v, %v", CommandName, err)
	}
	return command(opt)
}

// Usage provides information on the use of the lattice tool
func Usage() {
	fmt.Fprintf(ErrorWriter, UsageMessage, CommandName)
}

// PrintDefaults prints out the default flags
func PrintDefaults(eh flag.ErrorHandling) {
	o := newOption(ErrorWriter, eh)
	o.flagSet.PrintDefaults()
}
