package tokenize

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ikawaha/kagome-dict/dict"
	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome-dict/uni"
	"github.com/ikawaha/kagome/v2/filter"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

// subcommand property
const (
	CommandName  = "tokenize"
	Description  = `command line tokenize`
	usageMessage = "%s [-file input_file] [-dict dic_file] [-userdict user_dic_file]" +
		" [-sysdict (ipa|uni)] [-simple false] [-mode (normal|search|extended)] [-split] [-json]"
)

var (
	// Stdout is the standard writer.
	Stdout io.Writer = os.Stdout
	// Stderr is the standard error writer.
	Stderr io.Writer = os.Stderr
)

// options
type option struct {
	file    string
	dict    string
	udict   string
	sysdict string
	simple  bool
	mode    string
	split   bool
	json    bool
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
	o.flagSet.StringVar(&o.sysdict, "sysdict", "ipa", "system dict type (ipa|uni)")
	o.flagSet.BoolVar(&o.simple, "simple", false, "display abbreviated dictionary contents")
	o.flagSet.StringVar(&o.mode, "mode", "normal", "tokenize mode (normal|search|extended)")
	o.flagSet.BoolVar(&o.split, "split", false, "use tiny sentence splitter")
	o.flagSet.BoolVar(&o.json, "json", false, "outputs in JSON format")

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
	if o.sysdict != "" && o.sysdict != "ipa" && o.sysdict != "uni" {
		return fmt.Errorf("invalid argument: -sysdict %v", o.sysdict)
	}
	return nil
}

// OptionCheck receives a slice of args and returns an error if it was not successfully parsed
func OptionCheck(args []string) error {
	opt := newOption(io.Discard, flag.ContinueOnError)
	if err := opt.parse(args); err != nil {
		return fmt.Errorf("%v, %w", CommandName, err)
	}
	return nil
}

func selectDict(path, sysdict string, shrink bool) (*dict.Dict, error) {
	if path != "" {
		if shrink {
			return dict.LoadShrink(path)
		}
		return dict.LoadDictFile(path)
	}
	switch sysdict {
	case "ipa":
		if shrink {
			return ipa.DictShrink(), nil
		}
		return ipa.Dict(), nil
	case "uni":
		if shrink {
			return uni.DictShrink(), nil
		}
		return uni.Dict(), nil
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

func command(_ context.Context, opt *option) error {
	d, err := selectDict(opt.dict, opt.sysdict, opt.simple)
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

	fp := os.Stdin
	if opt.file != "" {
		var err error
		fp, err = os.Open(opt.file)
		if err != nil {
			return err
		}
		defer func() {
			_ = fp.Close()
		}()
	}
	mode := selectMode(opt.mode)
	s := bufio.NewScanner(fp)
	if opt.split {
		s.Split(filter.ScanSentences)
	}
	for s.Scan() {
		tokens := t.Analyze(s.Text(), mode)
		if !opt.json {
			printTokens(tokens)
			continue
		}
		if err := printTokensJSON(tokens); err != nil {
			return err
		}
	}
	return s.Err()
}

func printTokens(tokens []tokenizer.Token) {
	w := bufio.NewWriter(Stdout)
	defer w.Flush()
	for _, v := range tokens {
		if v.ID == tokenizer.BosEosID {
			continue
		}
		w.WriteString(v.Surface)
		if v.Class != tokenizer.DUMMY {
			w.WriteString("\t")
			w.WriteString(strings.Join(v.Features(), ","))
		}
		w.WriteString("\n")
	}
	w.WriteString("EOS\n")
}

func printTokensJSON(tokens []tokenizer.Token) error {
	w := bufio.NewWriter(Stdout)
	defer w.Flush()

	if len(tokens) > 0 {
		w.WriteString("[\n")
	}
	var array [][]byte
	for _, v := range tokens {
		if v.Class == tokenizer.DUMMY {
			continue
		}
		r := tokenizer.NewTokenData(v)
		obj, err := json.Marshal(r)
		if err != nil {
			return err
		}
		array = append(array, obj)
	}
	w.Write(bytes.Join(array, []byte(",\n")))
	if len(tokens) > 0 {
		w.WriteString("\n]\n")
	}
	return nil
}

// Run receives the slice of args and executes the tokenize tool
func Run(ctx context.Context, args []string) error {
	opt := newOption(Stderr, flag.ContinueOnError)
	if err := opt.parse(args); err != nil {
		Usage()
		PrintDefaults(flag.ContinueOnError)
		return errors.New("")
	}
	return command(ctx, opt)
}

// Usage provides information on the use of the tokenize tool
func Usage() {
	fmt.Fprintf(Stderr, usageMessage+"\n", CommandName)
}

// PrintDefaults prints out the default flags
func PrintDefaults(eh flag.ErrorHandling) {
	o := newOption(Stderr, eh)
	o.flagSet.PrintDefaults()
}
