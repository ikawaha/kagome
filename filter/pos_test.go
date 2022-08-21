package filter_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ikawaha/kagome-dict/dict"
	"github.com/ikawaha/kagome/v2/filter"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

const testDictPath = "../testdata/ipa.dict"

//BOS
//赤い	形容詞,自立,*,*,形容詞・アウオ段,基本形,赤い,アカイ,アカイ
//蝋燭	名詞,一般,*,*,*,*,蝋燭,ロウソク,ローソク
//と	助詞,並立助詞,*,*,*,*,と,ト,ト
//人魚	名詞,一般,*,*,*,*,人魚,ニンギョ,ニンギョ
//。	記号,句点,*,*,*,*,。,。,。
//人魚	名詞,一般,*,*,*,*,人魚,ニンギョ,ニンギョ
//は	助詞,係助詞,*,*,*,*,は,ハ,ワ
//、	記号,読点,*,*,*,*,、,、,、
//南	名詞,一般,*,*,*,*,南,ミナミ,ミナミ
//の	助詞,連体化,*,*,*,*,の,ノ,ノ
//方	名詞,非自立,一般,*,*,*,方,ホウ,ホー
//の	助詞,連体化,*,*,*,*,の,ノ,ノ
//海	名詞,一般,*,*,*,*,海,ウミ,ウミ
//に	助詞,格助詞,一般,*,*,*,に,ニ,ニ
//ばかり	助詞,副助詞,*,*,*,*,ばかり,バカリ,バカリ
//棲ん	動詞,自立,*,*,五段・マ行,連用タ接続,棲む,スン,スン
//で	助詞,接続助詞,*,*,*,*,で,デ,デ
//いる	動詞,非自立,*,*,一段,基本形,いる,イル,イル
//の	名詞,非自立,一般,*,*,*,の,ノ,ノ
//で	助動詞,*,*,*,特殊・ダ,連用形,だ,デ,デ
//は	助詞,係助詞,*,*,*,*,は,ハ,ワ
//あり	動詞,自立,*,*,五段・ラ行,連用形,ある,アリ,アリ
//ませ	助動詞,*,*,*,特殊・マス,未然形,ます,マセ,マセ
//ん	助動詞,*,*,*,不変化型,基本形,ん,ン,ン
//EOS
var input = "赤い蝋燭と人魚。人魚は、南の方の海にばかり棲んでいるのではありません"

func TestPOSFilter_Match(t *testing.T) {
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
		title        string
		featuresList []filter.Features
		want         []string
	}{
		{
			title:        "empty filter test",
			featuresList: []filter.Features{},
			want:         nil,
		},
		{
			title: "POS filter test",
			featuresList: []filter.Features{
				{"名詞", "一般"},
				{"動詞", "自立"},
				{"形容詞"},
			},
			want: []string{
				"赤い", "蝋燭", "人魚", "人魚", "南", "海", "棲ん", "あり",
			},
		},
	}

	for _, v := range testdata {
		t.Run(v.title, func(t *testing.T) {
			fl := filter.NewPOSFilter(v.featuresList...)
			var got []string
			for _, token := range tokens {
				if fl.Match(token.POS()) {
					got = append(got, token.Surface)
				}
			}
			if !reflect.DeepEqual(v.want, got) {
				t.Errorf("want %+v, got %+v", v.want, got)
			}
		})
	}
}

func TestPOSFilter_Keep(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		panic(err)
	}
	tnz, err := tokenizer.New(d)
	if err != nil {
		panic(err)
	}

	testdata := []struct {
		title        string
		featuresList []filter.Features
		want         []string
	}{
		{
			title:        "empty filter test",
			featuresList: []filter.Features{},
			want:         nil,
		},
		{
			title: "POS filter test",
			featuresList: []filter.Features{
				{"名詞", "一般"},
				{"動詞", "自立"},
				{"形容詞"},
			},
			want: []string{
				"赤い", "蝋燭", "人魚", "人魚", "南", "海", "棲ん", "あり",
			},
		},
	}

	for _, v := range testdata {
		t.Run(v.title, func(t *testing.T) {
			tokens := tnz.Tokenize(input)
			fl := filter.NewPOSFilter(v.featuresList...)
			var got []string
			fl.Keep(&tokens)
			for _, token := range tokens {
				got = append(got, token.Surface)
			}
			if !reflect.DeepEqual(v.want, got) {
				t.Errorf("want %+v, got %+v", v.want, got)
			}
		})
	}

	t.Run("empty input test", func(t *testing.T) {
		fl := filter.NewPOSFilter(nil)
		fl.Keep(nil)
	})
}

func TestPOSFilter_Drop(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		panic(err)
	}
	tnz, err := tokenizer.New(d)
	if err != nil {
		panic(err)
	}

	testdata := []struct {
		title        string
		featuresList []filter.Features
		want         []string
	}{
		{
			title:        "empty filter test",
			featuresList: []filter.Features{},
			want: []string{
				"BOS", "赤い", "蝋燭", "と", "人魚", "。", "人魚", "は", "、", "南", "の", "方", "の", "海", "に", "ばかり", "棲ん", "で", "いる", "の", "で", "は", "あり", "ませ", "ん", "EOS",
			},
		},
		{
			title: "POS filter test",
			featuresList: []filter.Features{
				{"名詞", "一般"},
				{"動詞", "自立"},
				{"形容詞"},
			},
			want: []string{
				"BOS", "と", "。", "は", "、", "の", "方", "の", "に", "ばかり", "で", "いる", "の", "で", "は", "ませ", "ん", "EOS",
			},
		},
	}

	for _, v := range testdata {
		t.Run(v.title, func(t *testing.T) {
			tokens := tnz.Tokenize(input)
			fl := filter.NewPOSFilter(v.featuresList...)
			var got []string
			fl.Drop(&tokens)
			for _, token := range tokens {
				got = append(got, token.Surface)
			}
			if !reflect.DeepEqual(v.want, got) {
				t.Errorf("want %+v, got %+v", v.want, got)
			}
		})
	}

	t.Run("empty input test", func(t *testing.T) {
		fl := filter.NewPOSFilter(nil)
		fl.Drop(nil)
	})
}

func ExamplePOSFilter() {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		panic(err)
	}
	t, err := tokenizer.New(d, tokenizer.OmitBosEos())
	if err != nil {
		panic(err)
	}
	posFilter := filter.NewPOSFilter([]filter.POS{
		{"名詞", filter.Any, "人名"},
		{"形容詞"},
	}...)
	tokens := t.Tokenize("赤い蝋燭と人魚。小川未明")
	posFilter.Keep(&tokens)
	for _, v := range tokens {
		fmt.Println(v.Surface, v.POS())
	}
	// Output:
	// 赤い [形容詞 自立 * *]
	// 小川 [名詞 固有名詞 人名 姓]
	// 未明 [名詞 固有名詞 人名 名]
}
