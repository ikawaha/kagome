package splitter

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestDefaultSplit(t *testing.T) {
	testdata := []struct {
		input  string
		expect []string
	}{
		{
			input:  "",
			expect: []string{},
		},
		{
			input:  "あああ",
			expect: []string{"あああ"},
		},
		{
			input:  "。。．",
			expect: []string{"。。．"},
		},
		{
			input:  "こんにちはは「さようなら」。\U0001f363が好き",
			expect: []string{"こんにちはは「さようなら」。", "\U0001f363が好き"},
		},
		{
			input:  "こんにちは。さようなら。おはよう．おやすみ．",
			expect: []string{"こんにちは。", "さようなら。", "おはよう．", "おやすみ．"},
		},
		{
			input:  "　 こ んに  ち　は。　 ．」　 さ\vよ\rう\tな\nら\n",
			expect: []string{"こんにちは。．」", "さようなら"},
		},
		{
			input:  "\n\nこんにちは\n\nさようなら。\n\n",
			expect: []string{"こんにちは", "さようなら。"},
		},
		{
			input:  "  こんにちは。。」さようなら。』ごきげんよう",
			expect: []string{"こんにちは。。」", "さようなら。』", "ごきげんよう"},
		},
	}

	for _, d := range testdata {
		scanner := bufio.NewScanner(strings.NewReader(d.input))
		scanner.Split(splitter.ScanSentences)
		r := make([]string, 0, len(d.expect))
		for scanner.Scan() {
			r = append(r, scanner.Text())
		}
		if !reflect.DeepEqual(r, d.expect) {
			t.Errorf("input %v, got %+v, expected %+v", d.input, r, d.expect)
		}
	}

}
