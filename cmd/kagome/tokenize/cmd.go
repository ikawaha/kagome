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
	"os"
	"strings"

	"github.com/ikawaha/kagome/tokenizer"
)

// subcommand property
var (
	CommandName  = "tokenize"
	Description  = `command line tokenize`
	usageMessage = "%s [-file input_file] [-dic dic_file] [-udic userdic_file] [-mode (normal|search|extended)]\n"
	errorWriter  = os.Stderr
)

// options
type option struct {
	file    string
	dic     string
	udic    string
	mode    string
	flagSet *flag.FlagSet
}

func newOption() (o *option) {
	o = &option{
		flagSet: flag.NewFlagSet(CommandName, flag.ExitOnError),
	}
	// option settings
	o.flagSet.SetOutput(errorWriter)
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
	if o.mode != "" && o.mode != "normal" && o.mode != "search" && o.mode != "extended" {
		err = fmt.Errorf("invalid argument: -mode %v\n", o.mode)
		return
	}
	return
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

	t := tokenizer.New(dic)
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
		tokens := t.Tokenize(line, mode)
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
	opt := newOption()
	if e := opt.parse(args); e != nil {
		Usage()
		PrintDefaults()
		fmt.Fprintf(errorWriter, "%v\n", e)
		os.Exit(1)
	}
	return command(opt)
}

func Usage() {
	fmt.Fprintf(errorWriter, usageMessage, CommandName)
}

func PrintDefaults() {
	o := newOption()
	o.flagSet.PrintDefaults()
}
