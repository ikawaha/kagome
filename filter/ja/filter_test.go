package ja

import (
	"reflect"
	"testing"

	"github.com/ikawaha/kagome-dict/dict"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

const testDictPath = "../../testdata/ipa.dict"

func TestFilter(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatal(err)
	}
	tz, err := tokenizer.New(d, tokenizer.OmitBosEos())
	if err != nil {
		t.Fatal(err)
	}
	tokens := tz.Tokenize("人魚は、南の方の海にばかり棲んでいるのではありません。")
	t.Run("yield string from tokens", func(t *testing.T) {
		f, err := NewFilter()
		if err != nil {
			t.Fatal(err)
		}
		want := []string{"人魚", "南", "方", "海", "棲む"}
		got := f.Yield(tokens)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})
	t.Run("drop tokens", func(t *testing.T) {
		f, err := NewFilter()
		if err != nil {
			t.Fatal(err)
		}
		f.Drop(&tokens)
		if len(tokens) != 5 {
			t.Errorf("got %+v, want %+v", len(tokens), 5)
		}
		if tokens[0].Surface != "人魚" {
			t.Errorf("got %+v, want %+v", tokens[0].Surface, "人魚")
		}
		if tokens[1].Surface != "南" {
			t.Errorf("got %+v, want %+v", tokens[1].Surface, "南")
		}
		if tokens[2].Surface != "方" {
			t.Errorf("got %+v, want %+v", tokens[2].Surface, "方")
		}
		if tokens[3].Surface != "海" {
			t.Errorf("got %+v, want %+v", tokens[3].Surface, "海")
		}
		if tokens[4].Surface != "棲ん" {
			t.Errorf("got %+v, want %+v", tokens[4].Surface, "棲ん")
		}
	})
}
