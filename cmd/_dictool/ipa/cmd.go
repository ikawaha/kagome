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

package ipa

import (
	"flag"
	"fmt"
	"os"
)

// subcommand property
var (
	CommandName  = "ipa"
	Description  = `ipa dic build tool`
	usageMessage = "%s -mecab mecabdic-path [-neologd neologd-path] [-z]\n"
	errorWriter  = os.Stderr
)

// options
type option struct {
	output  string
	mecab   string
	neologd string
	archive bool

	flagSet *flag.FlagSet
}

func newOption() (o *option) {
	o = &option{
		flagSet: flag.NewFlagSet(CommandName, flag.ContinueOnError),
	}
	// option settings
	o.flagSet.BoolVar(&o.archive, "z", true, "build an archived dic")
	o.flagSet.StringVar(&o.output, "output", ".", "set output path")
	o.flagSet.StringVar(&o.mecab, "mecab", "", "set mecab src path")
	o.flagSet.StringVar(&o.neologd, "neologd", "", "set neologd src path")
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
	if o.mecab == "" {
		return fmt.Errorf("invalid argument: -mecab")
	}
	return
}

// command main
func command(opt *option) error {
	d, err := buildIpaDic(opt.mecab, opt.neologd)
	if err != nil {
		return err
	}
	if saveIpaDic(d, opt.output, opt.archive); err != nil {
		return fmt.Errorf("build error: %v\n", err)
	}
	return nil
}

func Run(args []string) error {
	if len(args) == 0 {
		Usage()
		PrintDefaults()
		return nil
	}
	opt := newOption()
	if e := opt.parse(args); e != nil {
		Usage()
		PrintDefaults()
		return e
	}
	return command(opt)
}

func Usage() {
	fmt.Fprintf(os.Stderr, usageMessage, CommandName)
}

func PrintDefaults() {
	o := newOption()
	o.flagSet.PrintDefaults()
}
