package spliter

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
		{input: "あああ", expect: []string{"あああ"}},
		{input: "。。．", expect: []string{"。。．"}},
		{input: "あいう「かきく」えお。abcd", expect: []string{"あいう「かきく」えお。", "abcd"}},
		{input: "あああ。いいい", expect: []string{"あああ。", "いいい"}},
		{input: "あああ．いいい．", expect: []string{"あああ．", "いいい．"}},
		{input: "   あああ．    　いいい。", expect: []string{"あああ．", "いいい。"}},
		{input: "   あああ．\nいいい。", expect: []string{"あああ．", "いいい。"}},
		{input: "   あああ．\nいいい。", expect: []string{"あああ．", "いいい。"}},
		{input: "  あああ\n\nいいい。", expect: []string{"あああ", "いいい。"}},
		{input: "  あ あ  あ。\n\n。いいい。", expect: []string{"あああ。", "。", "いいい。"}},
	}

	for _, d := range testdata {
		scanner := bufio.NewScanner(strings.NewReader(d.input))
		scanner.Split(spliter.ScanSentences)
		r := make([]string, 0, len(d.expect))
		for scanner.Scan() {
			//sen := scanner.Text()
			//r = append(r, sen)
			r = append(r, scanner.Text())
		}
		if !reflect.DeepEqual(r, d.expect) {
			t.Errorf("input %v, got %+v, expected %+v", d.input, r, d.expect)
		}
	}

}
