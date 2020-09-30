package filter_test

import (
	"reflect"
	"testing"

	"github.com/ikawaha/kagome-dict/dict"
	"github.com/ikawaha/kagome/v2/filter"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

func TestPickup(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		panic(err)
	}
	tnz, err := tokenizer.New(d)
	if err != nil {
		panic(err)
	}

	testdata := []struct {
		title    string
		wordList []string
		want     []string
	}{
		{
			title:    "empty filter test",
			wordList: nil,
			want:     nil,
		},
		{
			title: "word filter test",
			wordList: []string{
				"人魚", "南", "の",
			},
			want: []string{
				"人魚", "人魚", "南", "の", "の", "の",
			},
		},
	}

	for _, v := range testdata {
		t.Run(v.title, func(t *testing.T) {
			tokens := tnz.Tokenize(input)
			var got []string
			filter.PickUp(&tokens, func(t tokenizer.Token) bool {
				for _, v := range v.wordList {
					if v == t.Surface {
						return true
					}
				}
				return false
			})
			for _, token := range tokens {
				got = append(got, token.Surface)
			}
			if !reflect.DeepEqual(v.want, got) {
				t.Errorf("want %+v, got %+v", v.want, got)
			}
		})
	}

}

func TestDrop(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		panic(err)
	}
	tnz, err := tokenizer.New(d)
	if err != nil {
		panic(err)
	}

	testdata := []struct {
		title    string
		wordList []string
		want     []string
	}{
		{
			title:    "empty filter test",
			wordList: nil,
			want: []string{
				"BOS", "赤い", "蝋燭", "と", "人魚", "。", "人魚", "は", "、", "南", "の", "方", "の", "海", "に", "ばかり", "棲ん", "で", "いる", "の", "で", "は", "あり", "ませ", "ん", "EOS",
			},
		},
		{
			title: "word filter test",
			wordList: []string{
				"人魚", "南", "の",
			},
			want: []string{
				"BOS", "赤い", "蝋燭", "と", "。", "は", "、", "方", "海", "に", "ばかり", "棲ん", "で", "いる", "で", "は", "あり", "ませ", "ん", "EOS",
			},
		},
	}

	for _, v := range testdata {
		t.Run(v.title, func(t *testing.T) {
			tokens := tnz.Tokenize(input)
			filter.Drop(&tokens, func(t tokenizer.Token) bool {
				for _, v := range v.wordList {
					if v == t.Surface {
						return true
					}
				}
				return false
			})
			var got []string
			for _, token := range tokens {
				got = append(got, token.Surface)
			}
			if !reflect.DeepEqual(v.want, got) {
				t.Errorf("want %+v, got %+v", v.want, got)
			}
		})
	}
}

func Benchmark_TokenFilter_WordFilter(b *testing.B) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		panic(err)
	}
	tnz, err := tokenizer.New(d)
	if err != nil {
		panic(err)
	}
	b.Run("package token filter", func(b *testing.B) {
		tokens := tnz.Tokenize(input)
		words := map[string]struct{}{
			"人魚": {},
			"南":  {},
			"の":  {},
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			filter.PickUp(&tokens, func(t tokenizer.Token) bool {
				_, ok := words[t.Surface]
				return ok
			})
		}
	})
	b.Run("specific token filter", func(b *testing.B) {
		tokens := tnz.Tokenize(input)
		words := []string{"人魚", "南", "の"}
		filter := filter.NewWordFilter(words)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			filter.PickUp(&tokens)
		}
	})
}
