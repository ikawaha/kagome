package server

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	ipa "github.com/ikawaha/kagome-dict-ipa"
	ko "github.com/ikawaha/kagome-dict-ko"
	uni "github.com/ikawaha/kagome-dict-uni"
	"github.com/ikawaha/kagome/v2/dict"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

// subcommand property
var (
	CommandName  = "server"
	Description  = `run tokenize server`
	usageMessage = "%s [-http=:6060] [-userdict userdic_file] [-dict (ipa|uni|ko)]\n"
	ErrorWriter  = os.Stderr
)

// options
type option struct {
	http    string
	dict    string
	udict   string
	flagSet *flag.FlagSet
}

// ContinueOnError ErrorHandling // Return a descriptive error.
// ExitOnError                   // Call os.Exit(2).
// PanicOnError                  // Call panic with a descriptive error.flag.ContinueOnError
func newOption(eh flag.ErrorHandling) (o *option) {
	o = &option{
		flagSet: flag.NewFlagSet(CommandName, eh),
	}
	// option settings
	o.flagSet.StringVar(&o.http, "http", ":6060", "HTTP service address")
	o.flagSet.StringVar(&o.udict, "userdict", "", "user dict")
	o.flagSet.StringVar(&o.dict, "dict", "ipa", "system dict type (ipa|uni|ko)")
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
	if o.dict != "" && o.dict != "ipa" && o.dict != "uni" && o.dict != "ko" {
		return fmt.Errorf("invalid argument: -dict %v", o.dict)
	}
	return nil
}

//OptionCheck receives a slice of args and returns an error if it was not successfully parsed
func OptionCheck(args []string) error {
	opt := newOption(flag.ContinueOnError)
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
	case "ko":
		return ko.Dict(), nil
	}
	return nil, fmt.Errorf("unknown name type, %v", name)
}

// command main
func command(opt *option) error {
	d, err := selectDict(opt.dict)
	if err != nil {
		return err
	}
	udict := tokenizer.Nop()
	if opt.udict != "" {
		d, err := dict.NewUserDict(opt.udict)
		if err != nil {
			return err
		}
		udict = tokenizer.UserDict(d)
	}
	t, err := tokenizer.New(d, udict)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", &TokenizeDemoHandler{tokenizer: t})
	mux.Handle("/tokenize", &TokenizeHandler{tokenizer: t})
	log.Println(opt.http)
	log.Fatal(http.ListenAndServe(opt.http, mux))
	return nil
}

// Run receives the slice of args and executes the server
func Run(args []string) error {
	opt := newOption(flag.ExitOnError)
	if err := opt.parse(args); err != nil {
		Usage()
		PrintDefaults(flag.ExitOnError)
		return fmt.Errorf("%v, %v", CommandName, err)
	}
	return command(opt)
}

// Usage provides information on the use of the server
func Usage() {
	_, _ = fmt.Fprintf(os.Stderr, usageMessage, CommandName)
}

// PrintDefaults prints out the default flags
func PrintDefaults(eh flag.ErrorHandling) {
	o := newOption(eh)
	o.flagSet.PrintDefaults()
}
