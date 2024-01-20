package main_test

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/ikawaha/kagome-dict/dict"
	"github.com/ikawaha/kagome/v2/filter"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

func TestKagomeTokenizerGolden(t *testing.T) {
	const (
		inputText  = "./testdata/bocchan.txt"
		goldenText = "./testdata/bocchan.golden"
		dumpText   = "./testdata/bocchan.dump"
	)
	d, err := dict.LoadDictFile("testdata/ipa.dict")
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	kagome, err := tokenizer.New(d)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	in, err := os.Open(inputText)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	defer in.Close()

	golden, err := os.Open(goldenText)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	defer golden.Close()

	dump, err := os.OpenFile(dumpText, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0o600)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	defer dump.Close()

	gs := bufio.NewScanner(golden)
	is := bufio.NewScanner(in)
	is.Split(filter.ScanSentences)
	for is.Scan() {
		tokens := kagome.Tokenize(is.Text())
		for _, v := range tokens {
			var want string
			if gs.Scan() {
				want = gs.Text()
			} else {
				t.Errorf("unexpected golden EOF")
			}
			got := fmt.Sprintf("%s %s", v, v.POS())
			if want != got {
				t.Errorf("got %v, want %v", got, want)
			}
			dump.WriteString(got)
			dump.WriteString("\n")
		}
	}
	if err := is.Err(); err != nil {
		t.Errorf("unexpected input file scanning error, %v", err)
	}
	if err := gs.Err(); err != nil {
		t.Errorf("unexpected golden file scanning error, %v", err)
	}
}
