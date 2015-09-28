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

package dic

import (
	"reflect"
	"testing"
)

var testFile = "../../_sample/userdic.txt"

func TestNewUserDic01(t *testing.T) {
	if _, e := NewUserDic(""); e == nil {
		t.Error("expected error, but no occured\n")
	}
}

func TestNewUserDicIndex01(t *testing.T) {
	udic, e := NewUserDic(testFile)
	if e != nil {
		t.Fatalf("unexpected error: %v\n", e)
	}
	type tuple struct {
		inp string
		id  int
		ok  bool
	}
	callAndRespose := []tuple{
		tuple{inp: "日本経済新聞", id: 0, ok: true},
		tuple{inp: "朝青龍", id: 1, ok: true},
		tuple{inp: "関西国際空港", id: 2, ok: true},
		tuple{inp: "成田国際空港", id: 9, ok: false},
	}
	for _, cr := range callAndRespose {
		if ids := udic.Index.Search(cr.inp); (len(ids) != 0) != cr.ok {
			t.Errorf("got %v, expected %v\n", ids, cr.ok)
		}
	}
}

func TestNewUserDicContents01(t *testing.T) {
	udic, e := NewUserDic(testFile)
	if e != nil {
		t.Fatalf("unexpected error: %v\n", e)
	}
	expectedLen := 3
	if len(udic.Contents) != expectedLen {
		t.Errorf("got %v, expected %v\n", len(udic.Contents), expectedLen)
	}

	type tuple struct {
		inp int
		out UserDicContent
	}
	callAndRespose := []tuple{
		tuple{
			inp: 0,
			out: UserDicContent{
				Tokens: []string{"日本", "経済", "新聞"},
				Yomi:   []string{"ニホン", "ケイザイ", "シンブン"},
				Pos:    "カスタム名詞",
			},
		},
		tuple{
			inp: 1,
			out: UserDicContent{
				Tokens: []string{"朝青龍"},
				Yomi:   []string{"アサショウリュウ"},
				Pos:    "カスタム人名",
			},
		},
		tuple{inp: 2,
			out: UserDicContent{
				Tokens: []string{"関西", "国際", "空港"},
				Yomi:   []string{"カンサイ", "コクサイ", "クウコウ"},
				Pos:    "テスト名詞",
			},
		},
	}
	for _, cr := range callAndRespose {
		c := udic.Contents[cr.inp]
		if !reflect.DeepEqual(c, cr.out) {
			t.Errorf("got %v, expected %v\n", c, cr.out)
		}
	}
}
