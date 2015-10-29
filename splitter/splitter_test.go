// Copyright 2015 ikawaha
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// 	You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package splitter

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestDefaultSplitter(t *testing.T) {
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
		scanner.Split(ScanSentences)
		r := make([]string, 0, len(d.expect))
		for scanner.Scan() {
			r = append(r, scanner.Text())
		}
		if !reflect.DeepEqual(r, d.expect) {
			t.Errorf("input %v, got %+v, expected %+v", d.input, r, d.expect)
		}
	}

}

func TestScanSentences(t *testing.T) {
	testdata := []struct {
		atEnd   bool
		data    []byte
		advance int
		token   []byte
		err     error
	}{
		{atEnd: true, data: []byte{}, advance: 0, token: []byte{}, err: nil},
		{atEnd: false, data: []byte{}, advance: 0, token: []byte{}, err: nil},
	}
	for _, d := range testdata {
		advance, token, err := ScanSentences(d.data, d.atEnd)
		if err != nil {
			t.Errorf("got err=%+v, expected nil", d.err)
		}
		if advance != 0 {
			t.Errorf("got advance %v, expected 0", d.advance)
		}
		if reflect.DeepEqual(token, d.token) {
			t.Errorf("got token=%+v, expected []", d.token)
		}
	}
}
