package ja

import (
	"reflect"
	"testing"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

func TestFilter(t *testing.T) {
	f, err := NewFilter()
	if err != nil {
		t.Fatal(err)
	}
	tz, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		t.Fatal(err)
	}
	tokens := tz.Tokenize("人魚は、南の方の海にばかり棲んでいるのではありません。")
	want := []string{"人魚", "南", "方", "海", "棲む"}
	got := f.Yield(tokens)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
