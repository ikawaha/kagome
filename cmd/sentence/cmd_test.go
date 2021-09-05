package sentence

import (
	"bytes"
	"context"
	"flag"
	"os"
	"testing"
)

func TestOptionCheck(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no options",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "no flag options",
			args:    []string{"ねこです。ねこはいます。"},
			wantErr: true,
		},
		{
			name:    "unknown option",
			args:    []string{"-flag", "piyo"},
			wantErr: true,
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
			args:    nil,
			wantErr: false,
		},
		{
			name: "invalid options",
			args: []string{
				"-piyo", "foo",
			},
			wantErr: true,
		},
		{
			name: "file option",
			args: []string{
				"-file", "../../testdata/nekodearu.txt",
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
	want := `sentence [-file filename]` + "\n"
	if got := b.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func Test_command(t *testing.T) {
	tests := []struct {
		name    string
		args    *option
		want    string
		wantErr bool
	}{
		{
			name: "no options",
			args: &option{
				flagSet: flag.NewFlagSet(CommandName, flag.ContinueOnError),
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "file",
			args: &option{
				file:    "../../testdata/nekodearu.txt",
				flagSet: flag.NewFlagSet(CommandName, flag.ContinueOnError),
			},
			want:    "吾輩は猫である。\n名前はまだ無い。\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			if err := command(context.TODO(), &b, tt.args); (err != nil) != tt.wantErr {
				t.Errorf("command() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.want != "" {
				if got, want := b.String(), tt.want; got != want {
					t.Errorf("got %s, want %s", got, want)
				}
			}
		})
	}
}
