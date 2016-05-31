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

package tokenizer

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

var testFile = "../_sample/userdic.txt"

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
		{inp: "日本経済新聞", id: 0, ok: true},
		{inp: "朝青龍", id: 1, ok: true},
		{inp: "関西国際空港", id: 2, ok: true},
		{inp: "成田国際空港", id: 9, ok: false},
	}
	for _, cr := range callAndRespose {
		if ids := udic.dic.Index.Search(cr.inp); (len(ids) != 0) != cr.ok {
			t.Errorf("got %v, expected %v\n", ids, cr.ok)
		}
	}
}

func TestNewUserDicRecords01(t *testing.T) {
	r := UserDicRecords{
		{
			Text:   "日本経済新聞",
			Tokens: []string{"日本", "経済", "新聞"},
			Yomi:   []string{"ニホン", "ケイザイ", "シンブン"},
			Pos:    "カスタム名詞",
		},
		{
			Text:   "朝青龍",
			Tokens: []string{"朝青龍"},
			Yomi:   []string{"アサショウリュウ"},
			Pos:    "カスタム人名",
		},
	}
	udic, err := r.NewUserDic()
	if err != nil {
		t.Fatalf("user dic build error, %v", err)
	}
	if ids := udic.dic.Index.Search("日本経済新聞"); len(ids) != 1 {
		t.Errorf("user dic search failed")
	} else if !reflect.DeepEqual(udic.dic.Contents[ids[0]].Tokens, []string{"日本", "経済", "新聞"}) {
		t.Errorf("got %+v, expected %+v", udic.dic.Contents[ids[0]].Tokens, []string{"日本", "経済", "新聞"})
	}
	if ids := udic.dic.Index.Search("関西国際空港"); len(ids) != 0 {
		t.Errorf("user dic build failed")
	}
	if ids := udic.dic.Index.Search("朝青龍"); len(ids) == 0 {
		t.Errorf("user dic search failed")
	} else if !reflect.DeepEqual(udic.dic.Contents[ids[0]].Tokens, []string{"朝青龍"}) {
		t.Errorf("got %+v, expected %+v", udic.dic.Contents[ids[0]].Tokens, []string{"朝青龍"})
	}
}

func TestNewUserDicRecords02(t *testing.T) {
	s := `
日本経済新聞,日本 経済 新聞,ニホン ケイザイ シンブン,カスタム名詞
# 関西国際空港,関西 国際 空港,カンサイ コクサイ クウコウ,カスタム地名
朝青龍,朝青龍,アサショウリュウ,カスタム人名
`
	r := strings.NewReader(s)
	rec, err := NewUserDicRecords(r)
	if err != nil {
		t.Fatalf("user dic build error, %v", err)
	}
	udic, err := rec.NewUserDic()
	if err != nil {
		t.Fatalf("user dic build error, %v", err)
	}
	if ids := udic.dic.Index.Search("日本経済新聞"); len(ids) != 1 {
		t.Errorf("user dic search failed")
	} else if !reflect.DeepEqual(udic.dic.Contents[ids[0]].Tokens, []string{"日本", "経済", "新聞"}) {
		t.Errorf("got %+v, expected %+v", udic.dic.Contents[ids[0]].Tokens, []string{"日本", "経済", "新聞"})
	}
	if ids := udic.dic.Index.Search("関西国際空港"); len(ids) != 0 {
		t.Errorf("user dic build failed")
	}
	if ids := udic.dic.Index.Search("朝青龍"); len(ids) == 0 {
		t.Errorf("user dic search failed")
	} else if !reflect.DeepEqual(udic.dic.Contents[ids[0]].Tokens, []string{"朝青龍"}) {
		t.Errorf("got %+v, expected %+v", udic.dic.Contents[ids[0]].Tokens, []string{"朝青龍"})
	}

}

func TestNewUserDicRecords03(t *testing.T) {
	r := UserDicRecords{
		{
			Text:   "日本経済新聞",
			Tokens: []string{"日本", "経済", "新聞"},
			Yomi:   []string{"ニホン", "ケイザイ"},
			Pos:    "カスタム名詞",
		},
	}
	_, err := r.NewUserDic()
	if err == nil {
		t.Errorf("expected error, but nil")
	}
}

func TestNewUserDicRecords04(t *testing.T) {
	r := UserDicRecords{
		{
			Text:   "日本経済新聞",
			Tokens: []string{"日本", "経済", "新聞"},
			Yomi:   []string{"ニホン", "ケイザイ", "シンブン"},
			Pos:    "カスタム名詞",
		},
		{
			Text:   "日本経済新聞",
			Tokens: []string{"日本", "経済", "新聞"},
			Yomi:   []string{"ニホン", "ケイザイ", "シンブン"},
			Pos:    "カスタム名詞",
		},
	}
	_, err := r.NewUserDic()
	if err == nil {
		t.Errorf("expected error, but nil")
	}
}

func TestUserDicRecordsLoadFromJSON(t *testing.T) {
	var rec UserDicRecords
	_ = json.Unmarshal([]byte(`[
        {
            "text":"日本経済新聞",
            "tokens":["日本","経済","新聞"],
            "yomi":["ニホン","ケイザイ","シンブン"],
            "pos":"カスタム名詞"
        },
        {
            "text":"朝青龍",
            "tokens":["朝青龍"],
            "yomi":["アサショウリュウ"],
            "pos":"カスタム人名"
        }]`), &rec)
	expected := UserDicRecords{
		{
			Text:   "日本経済新聞",
			Tokens: []string{"日本", "経済", "新聞"},
			Yomi:   []string{"ニホン", "ケイザイ", "シンブン"},
			Pos:    "カスタム名詞",
		},
		{
			Text:   "朝青龍",
			Tokens: []string{"朝青龍"},
			Yomi:   []string{"アサショウリュウ"},
			Pos:    "カスタム人名",
		},
	}

	if !reflect.DeepEqual(rec, expected) {
		t.Errorf("got %v, expected %v", rec, expected)
	}
}
