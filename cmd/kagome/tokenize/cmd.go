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

package tokenize

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ikawaha/kagome/tokenizer"
)

// subcommand property
const (
	CommandName  = "tokenize"
	Description  = `command line tokenize`
	usageMessage = "%s [-file input_file] [-dic dic_file] [-udic userdic_file] [-mode (normal|search|extended)]\n"
)

var (
	ErrorWriter = os.Stderr
)

// options
type option struct {
	file    string
	dic     string
	udic    string
	mode    string
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
	o.flagSet.SetOutput(w)
	o.flagSet.StringVar(&o.file, "file", "", "input file")
	o.flagSet.StringVar(&o.dic, "dic", "", "dic")
	o.flagSet.StringVar(&o.udic, "udic", "", "user dic")
	o.flagSet.StringVar(&o.mode, "mode", "normal", "tokenize mode (normal|search|extended)")

	return
}

func (o *option) parse(args []string) (err error) {
	if err = o.flagSet.Parse(args); err != nil {
		return
	}
	// validations
	if nonFlag := o.flagSet.Args(); len(nonFlag) != 0 {
		return fmt.Errorf("invalid argument: %v", nonFlag)
	}
	if o.mode != "" && o.mode != "normal" && o.mode != "search" && o.mode != "extended" {
		return fmt.Errorf("invalid argument: -mode %v\n", o.mode)
	}
	return
}

func OptionCheck(args []string) (err error) {
	opt := newOption(ioutil.Discard, flag.ContinueOnError)
	if e := opt.parse(args); e != nil {
		return fmt.Errorf("%v, %v", CommandName, e)
	}
	return nil
}

// command main
func command(opt *option) error {
	var dic tokenizer.Dic
	if opt.dic == "" {
		dic = tokenizer.SysDic()
	} else {
		var err error
		dic, err = tokenizer.NewDic(opt.dic)
		if err != nil {
			return err
		}
	}
	var udic tokenizer.UserDic
	if opt.udic != "" {
		var err error
		udic, err = tokenizer.NewUserDic(opt.udic)
		if err != nil {
			return err
		}
	}
	var fp = os.Stdin
	if opt.file != "" {
		var err error
		fp, err = os.Open(opt.file)
		if err != nil {
			return err
		}
		defer fp.Close()
	}

	t := tokenizer.NewWithDic(dic)
	t.SetUserDic(udic)

	mode := tokenizer.Normal
	switch opt.mode {
	case "normal":
		mode = tokenizer.Normal
		break
	case "search":
		mode = tokenizer.Search
	case "extended":
		mode = tokenizer.Extended
	}

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := t.Analyze(line, mode)
		for i, size := 1, len(tokens); i < size; i++ {
			tok := tokens[i]
			c := tok.Features()
			if tok.Class == tokenizer.DUMMY {
				fmt.Printf("%s\n", tok.Surface)
			} else {
				fmt.Printf("%s\t%v\n", tok.Surface, strings.Join(c, ","))
			}
		}
	}
	return scanner.Err()
}

func Run(args []string) error {
	opt := newOption(ErrorWriter, flag.ExitOnError)
	if e := opt.parse(args); e != nil {
		Usage()
		PrintDefaults(flag.ExitOnError)
		return fmt.Errorf("%v, %v", CommandName, e)
	}
	return command(opt)
}

func Usage() {
	fmt.Fprintf(ErrorWriter, usageMessage, CommandName)
}

func PrintDefaults(eh flag.ErrorHandling) {
	o := newOption(ErrorWriter, eh)
	o.flagSet.PrintDefaults()
}
