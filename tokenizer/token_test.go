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
	"fmt"
	"reflect"
	"testing"

	"github.com/ikawaha/kagome/v2/internal/dic"
	"github.com/ikawaha/kagome/v2/internal/lattice"
)

func TestTokenClassString(t *testing.T) {
	pairs := []struct {
		inp TokenClass
		out string
	}{
		{DUMMY, "DUMMY"},
		{KNOWN, "KNOWN"},
		{UNKNOWN, "UNKNOWN"},
		{USER, "USER"},
	}

	for _, p := range pairs {
		if p.inp.String() != p.out {
			t.Errorf("got %v, expected %v", p.inp.String(), p.out)
		}
	}
}

func TestFeatures01(t *testing.T) {
	tok := Token{
		ID:      0,
		Class:   TokenClass(lattice.KNOWN),
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dic = dic.SysDic()

	f := tok.Features()
	expected := []string{"名詞", "一般", "*", "*", "*", "*", "Tシャツ", "ティーシャツ", "ティーシャツ"}
	if !reflect.DeepEqual(f, expected) {
		t.Errorf("got %v, expected %v", f, expected)
	}
}

func TestFeatures02(t *testing.T) {
	tok := Token{
		ID:      0,
		Class:   TokenClass(lattice.UNKNOWN),
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dic = dic.SysDic()

	f := tok.Features()
	expected := []string{"名詞", "固有名詞", "地域", "一般", "*", "*", "*"}
	if !reflect.DeepEqual(f, expected) {
		t.Errorf("got %v, expected %v", f, expected)
	}
}

func TestFeatures03(t *testing.T) {
	tok := Token{
		ID:      0,
		Class:   TokenClass(lattice.USER),
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dic = dic.SysDic()
	if udic, err := dic.NewUserDic("../_sample/userdic.txt"); err != nil {
		t.Fatalf("build user dic error: %v", err)
	} else {
		tok.udic = udic
	}

	f := tok.Features()
	expected := []string{"カスタム名詞", "日本/経済/新聞", "ニホン/ケイザイ/シンブン"}
	if !reflect.DeepEqual(f, expected) {
		t.Errorf("got %v, expected %v", f, expected)
	}
}

func TestFeatures04(t *testing.T) {
	tok := Token{
		ID:      0,
		Class:   DUMMY,
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dic = dic.SysDic()
	if udic, err := dic.NewUserDic("../_sample/userdic.txt"); err != nil {
		t.Fatalf("build user dic error: %v", err)
	} else {
		tok.udic = udic
	}

	f := tok.Features()
	if len(f) != 0 {
		t.Errorf("got %v, expected empty", f)
	}
}

func TestFeaturesAndPos01(t *testing.T) {
	tok := Token{
		ID:      0,
		Class:   TokenClass(lattice.KNOWN),
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dic = dic.SysDic()

	f := tok.Features()
	expected := []string{"名詞", "一般", "*", "*", "*", "*", "Tシャツ", "ティーシャツ", "ティーシャツ"}
	if !reflect.DeepEqual(f, expected) {
		t.Errorf("got %v, expected %v", f, expected)
	}
	if p := tok.Pos(); p != f[0] {
		t.Errorf("got %v, expected %v", p, f[0])
	}
}

func TestFeaturesAndPos02(t *testing.T) {
	tok := Token{
		ID:      0,
		Class:   TokenClass(lattice.UNKNOWN),
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dic = dic.SysDic()

	f := tok.Features()
	expected := []string{"名詞", "固有名詞", "地域", "一般", "*", "*", "*"}
	if !reflect.DeepEqual(f, expected) {
		t.Errorf("got %v, expected %v", f, expected)
	}
	if p := tok.Pos(); p != f[0] {
		t.Errorf("got %v, expected %v", p, f[0])
	}
}

func TestFeaturesAndPos03(t *testing.T) {
	tok := Token{
		ID:      0,
		Class:   TokenClass(lattice.USER),
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dic = dic.SysDic()
	if udic, err := dic.NewUserDic("../_sample/userdic.txt"); err != nil {
		t.Fatalf("build user dic error: %v", err)
	} else {
		tok.udic = udic
	}

	f := tok.Features()
	expected := []string{"カスタム名詞", "日本/経済/新聞", "ニホン/ケイザイ/シンブン"}
	if !reflect.DeepEqual(f, expected) {
		t.Errorf("got %v, expected %v", f, expected)
	}
	if p := tok.Pos(); p != f[0] {
		t.Errorf("got %v, expected %v", p, f[0])
	}
}

func TestFeaturesAndPos04(t *testing.T) {
	tok := Token{
		ID:      0,
		Class:   DUMMY,
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dic = dic.SysDic()
	if udic, err := dic.NewUserDic("../_sample/userdic.txt"); err != nil {
		t.Fatalf("build user dic error: %v", err)
	} else {
		tok.udic = udic
	}

	f := tok.Features()
	if len(f) != 0 {
		t.Errorf("got %v, expected empty", f)
	}
	if p := tok.Pos(); p != "" {
		t.Errorf("got %v, expected empty", p)
	}
}

func TestTokenString01(t *testing.T) {
	tok := Token{
		ID:      123,
		Class:   TokenClass(lattice.DUMMY),
		Start:   0,
		End:     1,
		Surface: "テスト",
	}
	expected := "テスト(0, 1)DUMMY[123]"
	str := fmt.Sprintf("%v", tok)
	if str != expected {
		t.Errorf("got %v, expected %v", str, expected)
	}
}
