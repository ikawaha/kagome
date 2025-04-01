package server

import (
	"bytes"
	"context"
	"flag"
	"os"
	"testing"
	"time"
)

func TestOptionCheck(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "empty",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "non-flag args",
			args:    []string{"piyo"},
			wantErr: true,
		},
		{
			name:    "unknown dict",
			args:    []string{"-dict", "piyo"},
			wantErr: true,
		},
		{
			name: "all args",
			args: []string{
				"-userdict", "../../testdata/userdict.txt",
				"-http", ":8888",
				"-dict", "ipa",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := OptionCheck(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("OptionCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "help option",
			args:    []string{"-help"},
			wantErr: true,
		},
		{
			name:    "unknown option",
			args:    []string{"piyo"},
			wantErr: true,
		},
		{
			name:    "normal operation w/o options",
			args:    []string{""},
			wantErr: true,
		},
		{
			name: "normal operation w/ options",
			args: []string{
				"-userdict", "../../testdata/userdict.txt",
				"-http", ":0",
				"-dict", "ipa",
			},
			wantErr: false,
		},
	}
	Stderr = os.NewFile(0, os.DevNull)
	defer func() {
		Stderr = os.Stderr
	}()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()
			if err := Run(ctx, tt.args); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
			<-ctx.Done()
		})
	}
}

func TestUsage(t *testing.T) {
	var b bytes.Buffer
	Stderr = &b
	defer func() {
		Stderr = os.Stderr
	}()
	Usage()
	want := `server [-http=:6060] [-userdict userdic_file] [-dict (ipa|uni)]` + "\n"
	if got := b.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func Test_command(t *testing.T) {
	tests := []struct {
		name    string
		opt     *option
		wantErr bool
	}{
		{
			name: "no options",
			opt: &option{
				flagSet: flag.NewFlagSet(CommandName, flag.ContinueOnError),
			},
			wantErr: true,
		},
		{
			name: "set options",
			opt: &option{
				http:    ":0",
				dict:    "ipa",
				udict:   "../../testdata/userdict.txt",
				flagSet: flag.NewFlagSet(CommandName, flag.ContinueOnError),
			},
			wantErr: false,
		},
	}
	Stderr = os.NewFile(0, os.DevNull)
	defer func() {
		Stderr = os.Stderr
	}()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()
			if err := command(ctx, tt.opt); (err != nil) != tt.wantErr {
				t.Errorf("command() error = %v, wantErr %v", err, tt.wantErr)
			}
			<-ctx.Done()
		})
	}
}
