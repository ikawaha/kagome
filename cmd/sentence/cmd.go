package sentence

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/ikawaha/kagome/v2/filter"
)

// subcommand property
const (
	CommandName  = "sentence"
	Description  = `tiny sentence splitter`
	usageMessage = "%s [-file filename]"
)

// Stderr writes to stderr
var (
	Stdout io.Writer = os.Stdout
	Stderr io.Writer = os.Stderr
)

// options
type option struct {
	file    string
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

func command(_ context.Context, w io.Writer, opt *option) error {
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
	scanner := bufio.NewScanner(fp)
	scanner.Split(filter.ScanSentences)
	for scanner.Scan() {
		fmt.Fprintln(w, scanner.Text())
	}
	return scanner.Err()
}

// Run receives the slice of args and executes the tokenize tool
func Run(ctx context.Context, args []string) error {
	opt := newOption(io.Discard, flag.ContinueOnError)
	if err := opt.parse(args); err != nil {
		Usage()
		PrintDefaults(flag.ContinueOnError)
		return fmt.Errorf("%v, %w", CommandName, err)
	}
	return command(ctx, Stdout, opt)
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
