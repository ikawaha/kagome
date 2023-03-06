package lattice

import (
	"bytes"
	"context"
	"flag"
	"os"
	"strings"
	"testing"
)

func TestOptionCheck(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "empty args",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "input only",
			args:    []string{"ねこです。ねこはいます。"},
			wantErr: false,
		},
		{
			name:    "unknown option",
			args:    []string{"-flag", "piyo"},
			wantErr: true,
		},
		{
			name:    "invalid dict",
			args:    []string{"-dict", "piyo"},
			wantErr: true,
		},
		{
			name:    "invalid mode",
			args:    []string{"-mode", "piyo"},
			wantErr: true,
		},
		{
			name: "all options and input",
			args: []string{
				"-udict", "../../sample/dict/userdict.txt",
				"-dict", "ipa",
				"-mode", "search",
				"-output", "/dev/null",
				"-v",
				"私は鰻",
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
			name:    "no options",
			args:    []string{},
			wantErr: true,
		},
		{
			name: "invalid options",
			args: []string{
				"-piyo", "foo",
			},
			wantErr: true,
		},
		{
			name: "all options and input",
			args: []string{
				"-udict", "../../sample/dict/userdict.txt",
				"-dict", "ipa",
				"-mode", "search",
				"-output", "/dev/null",
				"-v",
				"ねこです。",
			},
			wantErr: false,
		},
	}
	Stdout = os.NewFile(0, os.DevNull)
	Stderr = Stdout
	defer func() {
		Stdout = os.Stdout
		Stderr = os.Stderr
	}()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Run(context.Background(), tt.args); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
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
	want := `lattice [-udict userdict_file] [-dict (ipa|uni)] [-mode (normal|search|extended)] [-output output_file] [-v] sentence` + "\n"
	if got := b.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func Test_command(t *testing.T) {
	tests := []struct {
		name        string
		args        *option
		verboseWant string
		wantErr     bool
	}{
		{
			name: "no options",
			args: &option{
				flagSet: flag.NewFlagSet(CommandName, flag.ContinueOnError),
			},
			verboseWant: "",
			wantErr:     true,
		},
		{
			name: "verbose",
			args: &option{
				udict:   "../../sample/dict/userdict.txt",
				dict:    "uni",
				mode:    "extended",
				output:  "",
				verbose: true,
				input:   "関西国際空港",
				flagSet: flag.NewFlagSet(CommandName, flag.ContinueOnError),
			},
			verboseWant: "関西国際空港\tテスト名詞,関西/国際/空港,カンサイ/コクサイ/クウコウ\nEOS\n",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b, berr bytes.Buffer
			Stdout = &b
			Stderr = &berr
			defer func() {
				Stdout = os.Stdout
				Stderr = os.Stderr
			}()
			if err := command(context.TODO(), tt.args); (err != nil) != tt.wantErr {
				t.Errorf("command() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got := b.String(); !strings.HasPrefix(got, "graph lattice {") {
				if len(got) > 50 {
					got = got[:50]
				}
				t.Errorf("invalid graphviz format, %s", got)
			}
			if got, want := berr.String(), tt.verboseWant; got != want {
				t.Errorf("stdout error, got %q, want %q", got, want)
			}
		})
	}
}
