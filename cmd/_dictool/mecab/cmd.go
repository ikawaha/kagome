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

package mecab

import (
	"flag"
	"fmt"
	"os"
)

// subcommand property
var (
	CommandName  = "mecab"
	Description  = `mecab dic build tool`
	usageMessage = "%s -mecab mecabdic-path [-z]\n"
	errorWriter  = os.Stderr
)

// options
type option struct {
	output   string
	mecab    string
	archive  bool
	encoding string

	flagSet *flag.FlagSet
}

func newOption() (o *option) {
	o = &option{
		flagSet: flag.NewFlagSet(CommandName, flag.ContinueOnError),
	}
	// option settings
	o.flagSet.BoolVar(&o.archive, "z", true, "build an archived dic")
	o.flagSet.StringVar(&o.output, "output", "mecab.dic", "set output path")
	o.flagSet.StringVar(&o.mecab, "mecab", "", "set mecab src path")
	o.flagSet.StringVar(&o.encoding, "encoding", "utf8", "set mecab src encoding [utf8|eucjp|sjis|jis]")
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
	d, err := buildMecabDic(opt.mecab, opt.encoding)
	if err != nil {
		return err
	}
	if saveMecabDic(d, opt.output, opt.archive); err != nil {
		return fmt.Errorf("build error: %v", err)
	}
	return nil
}

// Run mecab command
func Run(args []string) error {
	if len(args) == 0 {
		usage()
		printDefaults()
		return nil
	}
	opt := newOption()
	if e := opt.parse(args); e != nil {
		usage()
		printDefaults()
		return e
	}
	return command(opt)
}

func usage() {
	fmt.Fprintf(os.Stderr, usageMessage, CommandName)
}

func printDefaults() {
	o := newOption()
	o.flagSet.PrintDefaults()
}
