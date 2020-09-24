package lattice

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ikawaha/kagome-dict/dict"
	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome-dict/uni"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

// subcommand property
var (
	CommandName  = "lattice"
	Description  = `lattice viewer`
	UsageMessage = "%s [-userDict userdic_file] [-dict (ipa|uni)] [-mode (normal|search|extended)] [-output output_file] [-v] sentence\n"
	ErrorWriter  = os.Stderr
)

// options
type option struct {
	userDict string
	dict     string
	mode     string
	output   string
	verbose  bool
	input    string
	flagSet  *flag.FlagSet
}

// ContinueOnError ErrorHandling // Return a descriptive error.
// ExitOnError                   // Call os.Exit(2).
// PanicOnError                  // Call panic with a descriptive error.flag.ContinueOnError
func newOption(w io.Writer, eh flag.ErrorHandling) (o *option) {
	o = &option{
		flagSet: flag.NewFlagSet(CommandName, eh),
	}
	// option settings
	o.flagSet.StringVar(&o.userDict, "userDict", "", "user dict")
	o.flagSet.StringVar(&o.dict, "dict", "ipa", "dict type (ipa|uni)")
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
	if o.dict != "" && o.dict != "ipa" && o.dict != "uni" {
		return fmt.Errorf("invalid argument: -dict %v", o.dict)
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

func selectDict(name string) (*dict.Dict, error) {
	switch name {
	case "ipa":
		return ipa.Dict(), nil
	case "uni":
		return uni.Dict(), nil
	}
	return nil, fmt.Errorf("unknown name type, %v", name)
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
	d, err := selectDict(opt.dict)
	if err != nil {
		return err
	}
	udict := tokenizer.Nop()
	if opt.userDict != "" {
		d, err := dict.NewUserDict(opt.userDict)
		if err != nil {
			return err
		}
		udict = tokenizer.UserDict(d)
	}
	t, err := tokenizer.New(d, udict)
	if err != nil {
		return err
	}
	var out = os.Stdout
	if opt.output != "" {
		var err error
		out, err = os.OpenFile(opt.output, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer out.Close()
	}

	mode := selectMode(opt.mode)
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
