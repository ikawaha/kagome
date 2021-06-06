package sentence

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/ikawaha/kagome/v2/filter"
)

// subcommand property
const (
	CommandName  = "sentence"
	Description  = `tiny sentence splitter`
	usageMessage = "%s\n"
)

// ErrorWriter writes to stderr
var (
	ErrorWriter = os.Stderr
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
	return nil
}

// OptionCheck receives a slice of args and returns an error if it was not successfully parsed
func OptionCheck(args []string) error {
	opt := newOption(ioutil.Discard, flag.ContinueOnError)
	if err := opt.parse(args); err != nil {
		return fmt.Errorf("%v, %v", CommandName, err)
	}
	return nil
}

// command main
func command(opt *option) error {
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
		fmt.Println(scanner.Text())
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
