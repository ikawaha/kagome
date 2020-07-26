package tokenize

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	ipa "github.com/ikawaha/kagome-dict-ipa"
	ko "github.com/ikawaha/kagome-dict-ko"
	uni "github.com/ikawaha/kagome-dict-uni"
	"github.com/ikawaha/kagome/v2/dict"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

// subcommand property
const (
	CommandName  = "tokenize"
	Description  = `command line tokenize`
	usageMessage = "%s [-file input_file] [-dict dic_file] [-userdict userdic_file] [-sysdict (ipa|uni|ko)] [-simple false] [-mode (normal|search|extended)]\n"
)

// ErrorWriter writes to stderr
var (
	ErrorWriter = os.Stderr
)

// options
type option struct {
	file    string
	dict    string
	udict   string
	sysdict string
	simple  bool
	mode    string
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
	o.flagSet.StringVar(&o.dict, "dict", "", "dict")
	o.flagSet.StringVar(&o.udict, "udict", "", "user dict")
	o.flagSet.StringVar(&o.sysdict, "sysdict", "ipa", "system dict type (ipa|uni|ko)")
	o.flagSet.BoolVar(&o.simple, "simple", false, "display abbreviated dictionary contents")
	o.flagSet.StringVar(&o.mode, "mode", "normal", "tokenize mode (normal|search|extended)")

	return
}

func (o *option) parse(args []string) error {
	if err := o.flagSet.Parse(args); err != nil {
		return err
	}
	// validations
	if nonFlag := o.flagSet.Args(); len(nonFlag) != 0 {
		return fmt.Errorf("invalid argument: %v", nonFlag)
	}
	if o.mode != "" && o.mode != "normal" && o.mode != "search" && o.mode != "extended" {
		return fmt.Errorf("invalid argument: -mode %v", o.mode)
	}
	if o.sysdict != "" && o.sysdict != "ipa" && o.sysdict != "uni" && o.sysdict != "ko" {
		return fmt.Errorf("invalid argument: -sysdict %v", o.sysdict)
	}
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

func selectDict(path, sysdict string, shurink bool) (*dict.Dict, error) {
	if path != "" {
		if shurink {
			return dict.LoadShurink(path)
		}
		return dict.LoadDictFile(path)
	}
	switch sysdict {
	case "ipa":
		if shurink {
			return ipa.DictShrink(), nil
		}
		return ipa.Dict(), nil
	case "uni":
		if shurink {
			return uni.DictShrink(), nil
		}
		return uni.Dict(), nil
	case "ko":
		if shurink {
			return ko.DictShrink(), nil
		}
		return ko.Dict(), nil
	}
	return nil, fmt.Errorf("unknown dict type, %v", sysdict)
}

func selectMode(mode string) tokenizer.TokenizeMode {
	switch mode {
	case "normal":
		return tokenizer.Normal
	case "search":
		return tokenizer.Search
	case "extended":
		return tokenizer.Extended
	}
	return tokenizer.Normal
}

// command main
func command(opt *option) error {
	d, err := selectDict(opt.dict, opt.sysdict, opt.simple)
	if err != nil {
		return err
	}
	t := tokenizer.New(d)
	if opt.udict != "" {
		udict, err := dict.NewUserDict(opt.udict)
		if err != nil {
			return err
		}
		t.SetUserDict(udict)
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
	mode := selectMode(opt.mode)
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

// Run receives the slice of args and executes the tokenize tool
func Run(args []string) error {
	opt := newOption(ErrorWriter, flag.ExitOnError)
	if err := opt.parse(args); err != nil {
		Usage()
		PrintDefaults(flag.ExitOnError)
		return fmt.Errorf("%v, %v", CommandName, err)
	}
	return command(opt)
}

// Usage provides information on the use of the tokenize tool
func Usage() {
	fmt.Fprintf(ErrorWriter, usageMessage, CommandName)
}

// PrintDefaults prints out the default flags
func PrintDefaults(eh flag.ErrorHandling) {
	o := newOption(ErrorWriter, eh)
	o.flagSet.PrintDefaults()
}
