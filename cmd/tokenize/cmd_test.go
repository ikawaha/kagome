package tokenize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/ikawaha/kagome/v2/tokenizer"
)

func TestCommand_NormalOutput(t *testing.T) {
	// input
	{
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("unexpected pipe error, %v", err)
		}
		stdin := os.Stdin
		os.Stdin = r
		defer func() {
			os.Stdin = stdin
		}()
		go func() {
			fmt.Fprintf(w, "ねこです")
			w.Close()
		}()
	}
	// output
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("unexpected pipe error, %v", err)
	}
	stdout := os.Stdout
	os.Stdout = w
	defer func() {
		os.Stdout = stdout
	}()

	if err := command(&option{
		dict: "../../testdata/ipa.dict",
	}); err != nil {
		t.Errorf("unexpected error, command failed, %v", err)
	}
	w.Close()

	var b bytes.Buffer
	io.Copy(&b, r)
	want := `ねこ	名詞,一般,*,*,*,*,ねこ,ネコ,ネコ
です	助動詞,*,*,*,特殊・デス,基本形,です,デス,デス
EOS
`
	if got := b.String(); got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestCommand_JSONOutput(t *testing.T) {
	// input
	{
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("unexpected pipe error, %v", err)
		}
		stdin := os.Stdin
		os.Stdin = r
		defer func() {
			os.Stdin = stdin
		}()
		go func() {
			fmt.Fprintf(w, "ねこです")
			w.Close()
		}()
	}
	// output
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("unexpected pipe error, %v", err)
	}
	stdout := os.Stdout
	os.Stdout = w
	defer func() {
		os.Stdout = stdout
	}()

	// test
	if err := command(&option{
		dict: "../../testdata/ipa.dict",
		json: true,
	}); err != nil {
		t.Errorf("unexpected error, command failed, %v", err)
	}
	w.Close()

	// check output
	var b bytes.Buffer
	io.Copy(&b, r)
	want := `[
{"id":54873,"start":0,"end":2,"surface":"ねこ","class":"KNOWN","pos":["名詞","一般","*","*"],"base_form":"ねこ","reading":"ネコ","pronunciation":"ネコ","features":["名詞","一般","*","*","*","*","ねこ","ネコ","ネコ"]},
{"id":47492,"start":2,"end":4,"surface":"です","class":"KNOWN","pos":["助動詞","*","*","*"],"base_form":"です","reading":"デス","pronunciation":"デス","features":["助動詞","*","*","*","特殊・デス","基本形","です","デス","デス"]}
]
`
	if got := b.String(); got != want {
		t.Errorf("got %s, want %s", got, want)
	}
	var data []tokenizer.TokenData
	if err := json.Unmarshal(b.Bytes(), &data); err != nil {
		t.Errorf("json array to go array translation failed, %v", err)
	}
	if got, want := len(data), 2; got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func TestCommand_JSONOutput_issue249(t *testing.T) {
	// input
	{
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("unexpected pipe error, %v", err)
		}
		stdin := os.Stdin
		os.Stdin = r
		defer func() {
			os.Stdin = stdin
		}()
		go func() {
			fmt.Fprintf(w, "すもももももももものうち\n私は鰻\n猫\n")
			w.Close()
		}()
	}
	// output
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("unexpected pipe error, %v", err)
	}
	stdout := os.Stdout
	os.Stdout = w
	defer func() {
		os.Stdout = stdout
	}()

	// test
	if err := command(&option{
		dict: "../../testdata/ipa.dict",
		json: true,
	}); err != nil {
		t.Errorf("unexpected error, command failed, %v", err)
	}
	w.Close()

	// check output
	var b bytes.Buffer
	io.Copy(&b, r)
	want := `[
{"id":36163,"start":0,"end":3,"surface":"すもも","class":"KNOWN","pos":["名詞","一般","*","*"],"base_form":"すもも","reading":"スモモ","pronunciation":"スモモ","features":["名詞","一般","*","*","*","*","すもも","スモモ","スモモ"]},
{"id":73244,"start":3,"end":4,"surface":"も","class":"KNOWN","pos":["助詞","係助詞","*","*"],"base_form":"も","reading":"モ","pronunciation":"モ","features":["助詞","係助詞","*","*","*","*","も","モ","モ"]},
{"id":74988,"start":4,"end":6,"surface":"もも","class":"KNOWN","pos":["名詞","一般","*","*"],"base_form":"もも","reading":"モモ","pronunciation":"モモ","features":["名詞","一般","*","*","*","*","もも","モモ","モモ"]},
{"id":73244,"start":6,"end":7,"surface":"も","class":"KNOWN","pos":["助詞","係助詞","*","*"],"base_form":"も","reading":"モ","pronunciation":"モ","features":["助詞","係助詞","*","*","*","*","も","モ","モ"]},
{"id":74988,"start":7,"end":9,"surface":"もも","class":"KNOWN","pos":["名詞","一般","*","*"],"base_form":"もも","reading":"モモ","pronunciation":"モモ","features":["名詞","一般","*","*","*","*","もも","モモ","モモ"]},
{"id":55829,"start":9,"end":10,"surface":"の","class":"KNOWN","pos":["助詞","連体化","*","*"],"base_form":"の","reading":"ノ","pronunciation":"ノ","features":["助詞","連体化","*","*","*","*","の","ノ","ノ"]},
{"id":8027,"start":10,"end":12,"surface":"うち","class":"KNOWN","pos":["名詞","非自立","副詞可能","*"],"base_form":"うち","reading":"ウチ","pronunciation":"ウチ","features":["名詞","非自立","副詞可能","*","*","*","うち","ウチ","ウチ"]}
]
[
{"id":304999,"start":0,"end":1,"surface":"私","class":"KNOWN","pos":["名詞","代名詞","一般","*"],"base_form":"私","reading":"ワタシ","pronunciation":"ワタシ","features":["名詞","代名詞","一般","*","*","*","私","ワタシ","ワタシ"]},
{"id":57061,"start":1,"end":2,"surface":"は","class":"KNOWN","pos":["助詞","係助詞","*","*"],"base_form":"は","reading":"ハ","pronunciation":"ワ","features":["助詞","係助詞","*","*","*","*","は","ハ","ワ"]},
{"id":387420,"start":2,"end":3,"surface":"鰻","class":"KNOWN","pos":["名詞","一般","*","*"],"base_form":"鰻","reading":"ウナギ","pronunciation":"ウナギ","features":["名詞","一般","*","*","*","*","鰻","ウナギ","ウナギ"]}
]
[
{"id":286994,"start":0,"end":1,"surface":"猫","class":"KNOWN","pos":["名詞","一般","*","*"],"base_form":"猫","reading":"ネコ","pronunciation":"ネコ","features":["名詞","一般","*","*","*","*","猫","ネコ","ネコ"]}
]
`
	if got := b.String(); got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}