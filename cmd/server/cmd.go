package server

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ikawaha/kagome-dict/dict"
	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome-dict/uni"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

// Stderr is the standard error writer.
var Stderr io.Writer = os.Stderr

// subcommand property
var (
	CommandName  = "server"
	Description  = `run tokenize server`
	usageMessage = "%s [-http=:6060] [-userdict userdic_file] [-dict (ipa|uni)]"
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
func newOption(w io.Writer, eh flag.ErrorHandling) *option {
	ret := &option{
		flagSet: flag.NewFlagSet(CommandName, eh),
	}
	// option settings
	ret.flagSet.SetOutput(w)
	ret.flagSet.StringVar(&ret.http, "http", ":6060", "HTTP service address")
	ret.flagSet.StringVar(&ret.udict, "userdict", "", "user dict")
	ret.flagSet.StringVar(&ret.dict, "dict", "ipa", "system dict type (ipa|uni)")
	return ret
}

func (o *option) parse(args []string) error {
	if err := o.flagSet.Parse(args); err != nil {
		return err
	}
	// validations
	if nonFlag := o.flagSet.Args(); len(nonFlag) != 0 {
		return fmt.Errorf("invalid argument: %v", nonFlag)
	}
	if o.dict != "" && o.dict != "ipa" && o.dict != "uni" {
		return fmt.Errorf("invalid argument: -dict %v", o.dict)
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

func selectDict(name string) (*dict.Dict, error) {
	switch name {
	case "ipa":
		return ipa.Dict(), nil
	case "uni":
		return uni.Dict(), nil
	}
	return nil, fmt.Errorf("unknown name type, %v", name)
}

func command(ctx context.Context, opt *option) error {
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
	srv := http.Server{
		Addr:    opt.http,
		Handler: mux,
	}
	ch := make(chan error)
	go func() {
		log.Println(opt.http)
		ch <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Printf("shutdown, %v", ctx.Err())
	case err := <-ch:
		return fmt.Errorf("server error, %w", err)
	}
	log.Printf("shutting down HTTP server at %q", opt.http)

	return nil
}

// Run receives the slice of args and executes the server
func Run(ctx context.Context, args []string) error {
	opt := newOption(Stderr, flag.ContinueOnError)
	if err := opt.parse(args); err != nil {
		Usage()
		PrintDefaults(flag.ContinueOnError)
		return errors.New("")
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-c:
			cancel()
		}
	}()
	return command(ctx, opt)
}

// Usage provides information on the use of the server
func Usage() {
	fmt.Fprintf(Stderr, usageMessage+"\n", CommandName)
}

// PrintDefaults prints out the default flags
func PrintDefaults(eh flag.ErrorHandling) {
	o := newOption(Stderr, eh)
	o.flagSet.PrintDefaults()
}
