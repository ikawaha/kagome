//  Copyright (c) 2014 ikawaha.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package kagome

import (
	"testing"
)

func TestTokenize01(t *testing.T) {
	tnz := NewTokenizer()
	tokens := tnz.Tokenize("")
	expected := []Token{
		{Id: -1, Surface: "BOS"},
		{Id: -1, Surface: "EOS"},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.Id != expected[i].Id ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestTokenize02(t *testing.T) {
	tnz := NewTokenizer()
	tokens := tnz.Tokenize("関西国際空港")
	expected := []Token{
		{Id: -1, Surface: "BOS"},
		{Id: 372977, Surface: "関西国際空港", Start: 0, End: 6, Class: KNOWN},
		{Id: -1, Surface: "EOS", Start: 6, End: 6},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.Id != expected[i].Id ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestTokenize03(t *testing.T) {
	tnz := NewTokenizer()
	udic, e := NewUserDic("_sample/userdic.txt")
	if e != nil {
		t.Fatalf("new user dic: unexpected error\n")
	}
	tnz.SetUserDic(udic)
	tokens := tnz.Tokenize("関西国際空港")
	expected := []Token{
		{Id: -1, Surface: "BOS"},
		{Id: 2, Surface: "関西国際空港", Start: 0, End: 6, Class: USER},
		{Id: -1, Surface: "EOS", Start: 6, End: 6},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.Id != expected[i].Id ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestTokenize04(t *testing.T) {
	tnz := NewTokenizer()
	tokens := tnz.Tokenize("ポポピ")
	expected := []Token{
		{Id: -1, Surface: "BOS"},
		{Id: 34, Surface: "ポポピ", Start: 0, End: 3, Class: UNKNOWN},
		{Id: -1, Surface: "EOS", Start: 3, End: 3},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.Id != expected[i].Id ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestSearcModeTokenize01(t *testing.T) {
	tnz := NewTokenizer()
	tokens := tnz.SearchModeTokenize("")
	expected := []Token{
		{Id: -1, Surface: "BOS"},
		{Id: -1, Surface: "EOS"},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.Id != expected[i].Id ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestSearchModeTokenize02(t *testing.T) {
	tnz := NewTokenizer()
	tokens := tnz.SearchModeTokenize("関西国際空港")
	expected := []Token{
		{Id: -1, Surface: "BOS"},
		{Id: 372968, Surface: "関西", Start: 0, End: 2, Class: KNOWN},
		{Id: 168541, Surface: "国際", Start: 2, End: 4, Class: KNOWN},
		{Id: 307133, Surface: "空港", Start: 4, End: 6, Class: KNOWN},
		{Id: -1, Surface: "EOS", Start: 6, End: 6},
	}

	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.Id != expected[i].Id ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestSearchModeTokenize03(t *testing.T) {
	tnz := NewTokenizer()
	udic, e := NewUserDic("_sample/userdic.txt")
	if e != nil {
		t.Fatalf("new user dic: unexpected error\n")
	}
	tnz.SetUserDic(udic)
	tokens := tnz.SearchModeTokenize("関西国際空港")
	expected := []Token{
		{Id: -1, Surface: "BOS"},
		{Id: 2, Surface: "関西国際空港", Start: 0, End: 6, Class: USER},
		{Id: -1, Surface: "EOS", Start: 6, End: 6},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.Id != expected[i].Id ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestSearchModeTokenize04(t *testing.T) {
	tnz := NewTokenizer()
	tokens := tnz.SearchModeTokenize("ポポピ")
	expected := []Token{
		{Id: -1, Surface: "BOS"},
		{Id: 34, Surface: "ポポピ", Start: 0, End: 3, Class: UNKNOWN},
		{Id: -1, Surface: "EOS", Start: 3, End: 3},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.Id != expected[i].Id ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestExtendedModeTokenize01(t *testing.T) {
	tnz := NewTokenizer()
	tokens := tnz.ExtendedModeTokenize("")
	expected := []Token{
		{Id: -1, Surface: "BOS"},
		{Id: -1, Surface: "EOS"},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.Id != expected[i].Id ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestExtendedModeTokenize02(t *testing.T) {
	tnz := NewTokenizer()
	tokens := tnz.ExtendedModeTokenize("関西国際空港")
	expected := []Token{
		{Id: -1, Surface: "BOS"},
		{Id: 372968, Surface: "関西", Start: 0, End: 2, Class: KNOWN},
		{Id: 168541, Surface: "国際", Start: 2, End: 4, Class: KNOWN},
		{Id: 307133, Surface: "空港", Start: 4, End: 6, Class: KNOWN},
		{Id: -1, Surface: "EOS", Start: 6, End: 6},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.Id != expected[i].Id ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestExtendedModeTokenize03(t *testing.T) {
	tnz := NewTokenizer()
	udic, e := NewUserDic("_sample/userdic.txt")
	if e != nil {
		t.Fatalf("new user dic: unexpected error\n")
	}
	tnz.SetUserDic(udic)
	tokens := tnz.ExtendedModeTokenize("関西国際空港")
	expected := []Token{
		{Id: -1, Surface: "BOS"},
		{Id: 2, Surface: "関西国際空港", Start: 0, End: 6, Class: USER},
		{Id: -1, Surface: "EOS", Start: 6, End: 6},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.Id != expected[i].Id ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestExtendedModeTokenize04(t *testing.T) {
	tnz := NewTokenizer()
	tokens := tnz.ExtendedModeTokenize("ポポピ")
	expected := []Token{
		{Id: -1, Surface: "BOS"},
		{Id: 34, Surface: "ポ", Start: 0, End: 1, Class: DUMMY},
		{Id: 34, Surface: "ポ", Start: 1, End: 2, Class: DUMMY},
		{Id: 34, Surface: "ピ", Start: 2, End: 3, Class: DUMMY},
		{Id: -1, Surface: "EOS", Start: 3, End: 3},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.Id != expected[i].Id ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestThreadsafeTokenize01(t *testing.T) {
	tnz := NewThreadsafeTokenizer()
	tokens := tnz.Tokenize("")
	expected := []Token{
		{Id: -1, Surface: "BOS"},
		{Id: -1, Surface: "EOS"},
	}
	if len(tokens) != len(expected) {
		t.Errorf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.Id != expected[i].Id ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestThreadsafeTokenize02(t *testing.T) {
	tnz := NewThreadsafeTokenizer()
	tokens := tnz.Tokenize("関西国際空港")
	expected := []Token{
		{Id: -1, Surface: "BOS"},
		{Id: 372977, Surface: "関西国際空港", Start: 0, End: 6, Class: KNOWN},
		{Id: -1, Surface: "EOS", Start: 6, End: 6},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.Id != expected[i].Id ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestThreadsafeTokenize03(t *testing.T) {
	tnz := NewThreadsafeTokenizer()
	udic, e := NewUserDic("_sample/userdic.txt")
	if e != nil {
		t.Fatalf("new user dic: unexpected error\n")
	}
	tnz.SetUserDic(udic)
	tokens := tnz.Tokenize("関西国際空港")
	expected := []Token{
		{Id: -1, Surface: "BOS"},
		{Id: 2, Surface: "関西国際空港", Start: 0, End: 6, Class: USER},
		{Id: -1, Surface: "EOS", Start: 6, End: 6},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.Id != expected[i].Id ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}

func TestThreadsafeTokenize04(t *testing.T) {
	tnz := NewThreadsafeTokenizer()
	tokens := tnz.Tokenize("ポポピ")
	expected := []Token{
		{Id: -1, Surface: "BOS"},
		{Id: 34, Surface: "ポポピ", Start: 0, End: 3, Class: UNKNOWN},
		{Id: -1, Surface: "EOS", Start: 3, End: 3},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v\n", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.Id != expected[i].Id ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v\n", tok, expected[i])
		}
	}

}
