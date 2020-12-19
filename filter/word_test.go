package filter_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ikawaha/kagome-dict/dict"
	"github.com/ikawaha/kagome/v2/filter"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

func TestWordFilter_Match(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		panic(err)
	}
	tnz, err := tokenizer.New(d)
	if err != nil {
		panic(err)
	}
	tokens := tnz.Tokenize(input)

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
			filter := filter.NewWordFilter(v.wordList)
			var got []string
			for _, token := range tokens {
				if filter.Match(token.Surface) {
					got = append(got, token.Surface)
				}
			}
			if !reflect.DeepEqual(v.want, got) {
				t.Errorf("want %+v, got %+v", v.want, got)
			}
		})
	}
}

func TestWordFilter_Keep(t *testing.T) {
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
			filter := filter.NewWordFilter(v.wordList)
			var got []string
			filter.Keep(&tokens)
			for _, token := range tokens {
				got = append(got, token.Surface)
			}
			if !reflect.DeepEqual(v.want, got) {
				t.Errorf("want %+v, got %+v", v.want, got)
			}
		})
	}

	t.Run("empty input test", func(t *testing.T) {
		filter := filter.NewWordFilter(nil)
		filter.Keep(nil)
	})
}

func TestWordFilter_Drop(t *testing.T) {
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
			filter := filter.NewWordFilter(v.wordList)
			var got []string
			filter.Drop(&tokens)
			for _, token := range tokens {
				got = append(got, token.Surface)
			}
			if !reflect.DeepEqual(v.want, got) {
				t.Errorf("want %+v, got %+v", v.want, got)
			}
		})
	}

	t.Run("empty input test", func(t *testing.T) {
		filter := filter.NewWordFilter(nil)
		filter.Drop(nil)
	})
}

func Example_WordFilter() {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		panic(err)
	}
	t, err := tokenizer.New(d, tokenizer.OmitBosEos())
	if err != nil {
		panic(err)
	}
	stopWords := filter.NewWordFilter([]string{"私", "は", "が", "の", "。"})
	tokens := t.Tokenize("私の猫の名前はアプロです。")
	stopWords.Drop(&tokens)
	for _, v := range tokens {
		fmt.Println(v.Surface)
	}
	// Output:
	// 猫
	// 名前
	// アプロ
	// です
}
