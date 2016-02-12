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

var testDic = "../../_sample/ipa.dic"

func TestDicLoad(t *testing.T) {
	dic, err := Load(testDic)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if expected, c := 392126, len(dic.Morphs); c != expected {
		t.Errorf("got %v, expected %v\n", c, expected)
	}
	if expected, c := 392126, len(dic.Contents); c != expected {
		t.Errorf("got %v, expected %v\n", c, expected)
	}
}

func TestDicIndex01(t *testing.T) {
	dic, err := Load(testDic)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	type callAndResponse struct {
		input string
		ids   []int
	}
	testSet := []callAndResponse{
		{"すもも", []int{36163}},
	}
	for _, cr := range testSet {
		ids := dic.Index.Search(cr.input)
		if !reflect.DeepEqual(ids, cr.ids) {
			t.Errorf("input %v, got %v, expected %v\n", cr.input, ids, cr.ids)
		}
	}
}

func TestDicIndex02(t *testing.T) {
	dic, err := Load(testDic)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	type callAndRespose struct {
		input string
		lens  []int
		ids   [][]int
	}
	testSet := []callAndRespose{
		{input: "あい",
			lens: []int{3, 6},
			ids: [][]int{
				{122, 123, 124, 125},
				{141, 142, 143, 144, 145}}},
		{input: "あいあい",
			lens: []int{3, 6, 12},
			ids: [][]int{
				{122, 123, 124, 125},
				{141, 142, 143, 144, 145},
				{146}}},
		{input: "すもも",
			lens: []int{3, 6, 9},
			ids: [][]int{
				{34563, 34564, 34565, 34566, 34567, 34568, 34569},
				{36161},
				{36163}}},
	}
	for _, cr := range testSet {
		lens, ids := dic.Index.CommonPrefixSearch(cr.input)
		if !reflect.DeepEqual(lens, cr.lens) {
			t.Errorf("input %v, got lens %v,\n expected %v\n", cr.input, lens, cr.lens)
		}
		if len(ids) != len(cr.ids) {
			t.Errorf("input %v, got ids len %v, expected len %v\n", cr.input, len(ids), len(cr.ids))
		}
		if !reflect.DeepEqual(ids, cr.ids) {
			t.Errorf("input %v, got ids %v,\n expected %v\n", cr.input, ids, cr.ids)
		}
	}

}

func TestDicCharClass01(t *testing.T) {
	dic, err := Load(testDic)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{
		"DEFAULT",      // 0
		"SPACE",        // 1
		"KANJI",        // 2
		"SYMBOL",       // 3
		"NUMERIC",      // 4
		"ALPHA",        // 5
		"HIRAGANA",     // 6
		"KATAKANA",     // 7
		"KANJINUMERIC", // 8
		"GREEK",        // 9
		"CYRILLIC",     //10
	}
	if !reflect.DeepEqual(dic.CharClass, expected) {
		t.Errorf("got %v, expected %v\n", dic.CharClass, expected)
	}
}

func TestDicCharCategory01(t *testing.T) {
	const (
		DEFAULT      = 0
		SPACE        = 1
		KANJI        = 2
		SYMBOL       = 3
		NUMERIC      = 4
		ALPHA        = 5
		HIRAGANA     = 6
		KATAKANA     = 7
		KANJINUMERIC = 8
		GREEK        = 9
		CYRILLIC     = 10
	)
	dic, err := Load(testDic)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	type callAndRespose struct {
		input    int
		category byte
	}
	testSet := []callAndRespose{
		{input: 0x0020, category: SPACE}, // 0x0020 SPACE  # DO NOT REMOVE THIS LINE, 0x0020 is reserved for SPACE
		//{input: 0x00D0, category: SPACE},    // 0x00D0 SPACE
		{input: 0x000d, category: SPACE},    // 0x00D0 SPACE
		{input: 0x0009, category: SPACE},    // 0x0009 SPACE
		{input: 0x000B, category: SPACE},    // 0x000B SPACE
		{input: 0x000A, category: SPACE},    // 0x000A SPACE
		{input: 0x0021, category: SYMBOL},   // 0x0021..0x002F SYMBOL
		{input: 0x002F, category: SYMBOL},   // 0x0021..0x002F SYMBOL
		{input: 0x0030, category: NUMERIC},  // 0x0030..0x0039 NUMERIC
		{input: 0x0039, category: NUMERIC},  // 0x0030..0x0039 NUMERIC
		{input: 0x003A, category: SYMBOL},   // 0x003A..0x0040 SYMBOL
		{input: 0x0040, category: SYMBOL},   // 0x003A..0x0040 SYMBOL
		{input: 0x0041, category: ALPHA},    // 0x0041..0x005A ALPHA
		{input: 0x005A, category: ALPHA},    // 0x0041..0x005A ALPHA
		{input: 0x005B, category: SYMBOL},   // 0x005B..0x0060 SYMBOL
		{input: 0x0060, category: SYMBOL},   // 0x005B..0x0060 SYMBOL
		{input: 0x0061, category: ALPHA},    // 0x0061..0x007A ALPHA
		{input: 0x007A, category: ALPHA},    // 0x0061..0x007A ALPHA
		{input: 0x007B, category: SYMBOL},   // 0x007B..0x007E SYMBOL
		{input: 0x007E, category: SYMBOL},   // 0x007B..0x007E SYMBOL
		{input: 0x00A1, category: SYMBOL},   // 0x00A1..0x00BF SYMBOL # Latin 1
		{input: 0x00BF, category: SYMBOL},   // 0x00A1..0x00BF SYMBOL # Latin 1
		{input: 0x00C0, category: ALPHA},    // 0x00C0..0x00FF ALPHA  # Latin 1
		{input: 0x00FF, category: ALPHA},    // 0x00C0..0x00FF ALPHA  # Latin 1
		{input: 0x0100, category: ALPHA},    // 0x0100..0x017F ALPHA  # Latin Extended A
		{input: 0x017F, category: ALPHA},    // 0x0100..0x017F ALPHA  # Latin Extended A
		{input: 0x0180, category: ALPHA},    // 0x0180..0x0236 ALPHA  # Latin Extended B
		{input: 0x0236, category: ALPHA},    // 0x0180..0x0236 ALPHA  # Latin Extended B
		{input: 0x1E00, category: ALPHA},    // 0x1E00..0x1EF9 ALPHA  # Latin Extended Additional
		{input: 0x1EF9, category: ALPHA},    // 0x1E00..0x1EF9 ALPHA  # Latin Extended Additional
		{input: 0x0400, category: CYRILLIC}, // 0x0400..0x04F9 CYRILLIC
		{input: 0x04F9, category: CYRILLIC}, // 0x0400..0x04F9 CYRILLIC
		{input: 0x0500, category: CYRILLIC}, // 0x0500..0x050F CYRILLIC # Cyrillic supplementary
		{input: 0x050F, category: CYRILLIC}, // 0x0500..0x050F CYRILLIC # Cyrillic supplementary
		{input: 0x0374, category: GREEK},    // 0x0374..0x03FB GREEK # Greek and Coptic
		{input: 0x03FB, category: GREEK},    // 0x0374..0x03FB GREEK # Greek and Coptic
		{input: 0x3041, category: HIRAGANA}, // 0x3041..0x309F  HIRAGANA
		{input: 0x309F, category: HIRAGANA}, // 0x3041..0x309F  HIRAGANA
		{input: 0x30A1, category: KATAKANA}, // 0x30A1..0x30FF  KATAKANA
		{input: 0x30FF, category: KATAKANA}, // 0x30A1..0x30FF  KATAKANA
		{input: 0x31F0, category: KATAKANA}, // 0x31F0..0x31FF  KATAKANA  # Small KU .. Small RO
		{input: 0x31FF, category: KATAKANA}, // 0x31F0..0x31FF  KATAKANA  # Small KU .. Small RO
		{input: 0x30FC, category: KATAKANA}, // 0x30FC          KATAKANA
		{input: 0xFF66, category: KATAKANA}, // 0xFF66..0xFF9D  KATAKANA
		{input: 0xFF9D, category: KATAKANA}, // 0xFF66..0xFF9D  KATAKANA
		{input: 0xFF9E, category: KATAKANA}, // 0xFF9E..0xFF9F  KATAKANA
		{input: 0xFF9F, category: KATAKANA}, // 0xFF9E..0xFF9F  KATAKANA
		{input: 0x2E80, category: KANJI},    // 0x2E80..0x2EF3  KANJI # CJK Raidcals Supplement
		{input: 0x2EF3, category: KANJI},    // 0x2E80..0x2EF3  KANJI # CJK Raidcals Supplement
		{input: 0x2F00, category: KANJI},    // 0x2F00..0x2FD5  KANJI
		{input: 0x2FD5, category: KANJI},    // 0x2F00..0x2FD5  KANJI
		//{input: 0x3005, category: KANJI},        // 0x3005          KANJI
		//{input: 0x3007, category: KANJI},        // 0x3007          KANJI
		{input: 0x3400, category: KANJI}, // 0x3400..0x4DB5  KANJI # CJK Unified Ideographs Extension
		{input: 0x4DB5, category: KANJI}, // 0x3400..0x4DB5  KANJI # CJK Unified Ideographs Extension
		//{input: 0x4E00, category: KANJI},        // 0x4E00..0x9FA5  KANJI
		{input: 0x9FA5, category: KANJI},        // 0x4E00..0x9FA5  KANJI
		{input: 0xF900, category: KANJI},        // 0xF900..0xFA2D  KANJI
		{input: 0xFA2D, category: KANJI},        // 0xF900..0xFA2D  KANJI
		{input: 0xFA30, category: KANJI},        // 0xFA30..0xFA6A  KANJI
		{input: 0xFA6A, category: KANJI},        // 0xFA30..0xFA6A  KANJI
		{input: 0x4E00, category: KANJINUMERIC}, // 0x4E00 KANJINUMERIC KANJI
		{input: 0x4E8C, category: KANJINUMERIC}, // 0x4E8C KANJINUMERIC KANJI
		{input: 0x4E09, category: KANJINUMERIC}, // 0x4E09 KANJINUMERIC KANJI
		{input: 0x56DB, category: KANJINUMERIC}, // 0x56DB KANJINUMERIC KANJI
		{input: 0x4E94, category: KANJINUMERIC}, // 0x4E94 KANJINUMERIC KANJI
		{input: 0x516D, category: KANJINUMERIC}, // 0x516D KANJINUMERIC KANJI
		{input: 0x4E03, category: KANJINUMERIC}, // 0x4E03 KANJINUMERIC KANJI
		{input: 0x516B, category: KANJINUMERIC}, // 0x516B KANJINUMERIC KANJI
		{input: 0x4E5D, category: KANJINUMERIC}, // 0x4E5D KANJINUMERIC KANJI
		{input: 0x5341, category: KANJINUMERIC}, // 0x5341 KANJINUMERIC KANJI
		{input: 0x767E, category: KANJINUMERIC}, // 0x767E KANJINUMERIC KANJI
		{input: 0x5343, category: KANJINUMERIC}, // 0x5343 KANJINUMERIC KANJI
		{input: 0x4E07, category: KANJINUMERIC}, // 0x4E07 KANJINUMERIC KANJI
		{input: 0x5104, category: KANJINUMERIC}, // 0x5104 KANJINUMERIC KANJI
		{input: 0x5146, category: KANJINUMERIC}, // 0x5146 KANJINUMERIC KANJI
		{input: 0xFF10, category: NUMERIC},      // 0xFF10..0xFF19 NUMERIC
		{input: 0xFF19, category: NUMERIC},      // 0xFF10..0xFF19 NUMERIC
		{input: 0xFF21, category: ALPHA},        // 0xFF21..0xFF3A ALPHA
		{input: 0xFF3A, category: ALPHA},        // 0xFF21..0xFF3A ALPHA
		{input: 0xFF41, category: ALPHA},        // 0xFF41..0xFF5A ALPHA
		{input: 0xFF5A, category: ALPHA},        // 0xFF41..0xFF5A ALPHA
		{input: 0xFF01, category: SYMBOL},       // 0xFF01..0xFF0F SYMBOL
		{input: 0xFF0F, category: SYMBOL},       // 0xFF01..0xFF0F SYMBOL
		{input: 0xFF1A, category: SYMBOL},       // 0xFF1A..0xFF1F SYMBOL
		{input: 0xFF1F, category: SYMBOL},       // 0xFF1A..0xFF1F SYMBOL
		{input: 0xFF3B, category: SYMBOL},       // 0xFF3B..0xFF40 SYMBOL
		{input: 0xFF40, category: SYMBOL},       // 0xFF3B..0xFF40 SYMBOL
		{input: 0xFF5B, category: SYMBOL},       // 0xFF5B..0xFF65 SYMBOL
		{input: 0xFF65, category: SYMBOL},       // 0xFF5B..0xFF65 SYMBOL
		{input: 0xFFE0, category: SYMBOL},       // 0xFFE0..0xFFEF SYMBOL # HalfWidth and Full width Form
		{input: 0xFFEF, category: SYMBOL},       // 0xFFE0..0xFFEF SYMBOL # HalfWidth and Full width Form
		{input: 0x2000, category: SYMBOL},       // 0x2000..0x206F  SYMBOL # General Punctuation
		{input: 0x206F, category: SYMBOL},       // 0x2000..0x206F  SYMBOL # General Punctuation
		{input: 0x2070, category: NUMERIC},      // 0x2070..0x209F  NUMERIC # Superscripts and Subscripts
		{input: 0x209F, category: NUMERIC},      // 0x2070..0x209F  NUMERIC # Superscripts and Subscripts
		{input: 0x20A0, category: SYMBOL},       // 0x20A0..0x20CF  SYMBOL # Currency Symbols
		{input: 0x20CF, category: SYMBOL},       // 0x20A0..0x20CF  SYMBOL # Currency Symbols
		{input: 0x20D0, category: SYMBOL},       // 0x20D0..0x20FF  SYMBOL # Combining Diaritical Marks for Symbols
		{input: 0x20FF, category: SYMBOL},       // 0x20D0..0x20FF  SYMBOL # Combining Diaritical Marks for Symbols
		{input: 0x2100, category: SYMBOL},       // 0x2100..0x214F  SYMBOL # Letterlike Symbols
		{input: 0x214F, category: SYMBOL},       // 0x2100..0x214F  SYMBOL # Letterlike Symbols
		{input: 0x2150, category: NUMERIC},      // 0x2150..0x218F  NUMERIC # Number forms
		{input: 0x218F, category: NUMERIC},      // 0x2150..0x218F  NUMERIC # Number forms
		{input: 0x2100, category: SYMBOL},       // 0x2100..0x214B  SYMBOL # Letterlike Symbols
		{input: 0x214B, category: SYMBOL},       // 0x2100..0x214B  SYMBOL # Letterlike Symbols
		{input: 0x2190, category: SYMBOL},       // 0x2190..0x21FF  SYMBOL # Arrow
		{input: 0x21FF, category: SYMBOL},       // 0x2190..0x21FF  SYMBOL # Arrow
		{input: 0x2200, category: SYMBOL},       // 0x2200..0x22FF  SYMBOL # Mathematical Operators
		{input: 0x22FF, category: SYMBOL},       // 0x2200..0x22FF  SYMBOL # Mathematical Operators
		{input: 0x2300, category: SYMBOL},       // 0x2300..0x23FF  SYMBOL # Miscellaneuos Technical
		{input: 0x23FF, category: SYMBOL},       // 0x2300..0x23FF  SYMBOL # Miscellaneuos Technical
		{input: 0x2460, category: SYMBOL},       // 0x2460..0x24FF  SYMBOL # Enclosed NUMERICs
		{input: 0x24FF, category: SYMBOL},       // 0x2460..0x24FF  SYMBOL # Enclosed NUMERICs
		{input: 0x2501, category: SYMBOL},       // 0x2501..0x257F  SYMBOL # Box Drawing
		{input: 0x257F, category: SYMBOL},       // 0x2501..0x257F  SYMBOL # Box Drawing
		{input: 0x2580, category: SYMBOL},       // 0x2580..0x259F  SYMBOL # Block Elements
		{input: 0x259F, category: SYMBOL},       // 0x2580..0x259F  SYMBOL # Block Elements
		{input: 0x25A0, category: SYMBOL},       // 0x25A0..0x25FF  SYMBOL # Geometric Shapes
		{input: 0x25FF, category: SYMBOL},       // 0x25A0..0x25FF  SYMBOL # Geometric Shapes
		{input: 0x2600, category: SYMBOL},       // 0x2600..0x26FE  SYMBOL # Miscellaneous Symbols
		{input: 0x26FE, category: SYMBOL},       // 0x2600..0x26FE  SYMBOL # Miscellaneous Symbols
		{input: 0x2700, category: SYMBOL},       // 0x2700..0x27BF  SYMBOL # Dingbats
		{input: 0x27BF, category: SYMBOL},       // 0x2700..0x27BF  SYMBOL # Dingbats
		{input: 0x27F0, category: SYMBOL},       // 0x27F0..0x27FF  SYMBOL # Supplemental Arrows A
		{input: 0x27FF, category: SYMBOL},       // 0x27F0..0x27FF  SYMBOL # Supplemental Arrows A
		{input: 0x27C0, category: SYMBOL},       // 0x27C0..0x27EF  SYMBOL # Miscellaneous Mathematical Symbols-A
		{input: 0x27EF, category: SYMBOL},       // 0x27C0..0x27EF  SYMBOL # Miscellaneous Mathematical Symbols-A
		{input: 0x2800, category: SYMBOL},       // 0x2800..0x28FF  SYMBOL # Braille Patterns
		{input: 0x28FF, category: SYMBOL},       // 0x2800..0x28FF  SYMBOL # Braille Patterns
		{input: 0x2900, category: SYMBOL},       // 0x2900..0x297F  SYMBOL # Supplemental Arrows B
		{input: 0x297F, category: SYMBOL},       // 0x2900..0x297F  SYMBOL # Supplemental Arrows B
		{input: 0x2B00, category: SYMBOL},       // 0x2B00..0x2BFF  SYMBOL # Miscellaneous Symbols and Arrows
		{input: 0x2BFF, category: SYMBOL},       // 0x2B00..0x2BFF  SYMBOL # Miscellaneous Symbols and Arrows
		{input: 0x2A00, category: SYMBOL},       // 0x2A00..0x2AFF  SYMBOL # Supplemental Mathematical Operators
		{input: 0x2AFF, category: SYMBOL},       // 0x2A00..0x2AFF  SYMBOL # Supplemental Mathematical Operators
		{input: 0x3300, category: SYMBOL},       // 0x3300..0x33FF  SYMBOL
		{input: 0x33FF, category: SYMBOL},       // 0x3300..0x33FF  SYMBOL
		{input: 0x3200, category: SYMBOL},       // 0x3200..0x32FE  SYMBOL # ENclosed CJK Letters and Months
		{input: 0x32FE, category: SYMBOL},       // 0x3200..0x32FE  SYMBOL # ENclosed CJK Letters and Months
		{input: 0x3000, category: SYMBOL},       // 0x3000..0x303F  SYMBOL # CJK Symbol and Punctuation
		{input: 0x303F, category: SYMBOL},       // 0x3000..0x303F  SYMBOL # CJK Symbol and Punctuation
		{input: 0xFE30, category: SYMBOL},       // 0xFE30..0xFE4F  SYMBOL # CJK Compatibility Forms
		{input: 0xFE4F, category: SYMBOL},       // 0xFE30..0xFE4F  SYMBOL # CJK Compatibility Forms
		{input: 0xFE50, category: SYMBOL},       // 0xFE50..0xFE6B  SYMBOL # Small Form Variants
		{input: 0xFE6B, category: SYMBOL},       // 0xFE50..0xFE6B  SYMBOL # Small Form Variants
		{input: 0x3007, category: SYMBOL},       // 0x3007 SYMBOL KANJINUMERIC

	}
	for _, cr := range testSet {
		category := dic.CharCategory[cr.input]
		if category != cr.category {
			t.Errorf("input %04X, got %v, expected %v\n", cr.input, category, cr.category)
		}
		category = dic.CharacterCategory(rune(cr.input))
		if category != cr.category {
			t.Errorf("input %04X, got %v, expected %v\n", cr.input, category, cr.category)
		}
	}
}

func TestCharCategory02(t *testing.T) {
	dic, err := Load(testDic)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c := dic.CharacterCategory(rune(len(dic.CharCategory) + 1))
	expected := dic.CharCategory[0]
	if c != expected {
		t.Errorf("got %v, expected %v", c, expected)
	}

}

func TestDicInvokeList01(t *testing.T) {
	const (
		DEFAULT      = 0
		SPACE        = 1
		KANJI        = 2
		SYMBOL       = 3
		NUMERIC      = 4
		ALPHA        = 5
		HIRAGANA     = 6
		KATAKANA     = 7
		KANJINUMERIC = 8
		GREEK        = 9
		CYRILLIC     = 10
	)
	crs := []struct {
		class  int
		invoke bool
	}{
		{DEFAULT, false},
		{SPACE, false},
		{KANJI, false},
		{SYMBOL, true},
		{NUMERIC, true},
		{ALPHA, true},
		{HIRAGANA, false},
		{KATAKANA, true},
		{KANJINUMERIC, true},
		{GREEK, true},
		{CYRILLIC, true},
	}
	dic, err := Load(testDic)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, cr := range crs {
		if iv := dic.InvokeList[cr.class]; iv != cr.invoke {
			t.Errorf("input %v: got %v, expected %v\n", cr.class, iv, cr.invoke)
		}
	}
}

func TestDicGroupList01(t *testing.T) {
	const (
		DEFAULT      = 0
		SPACE        = 1
		KANJI        = 2
		SYMBOL       = 3
		NUMERIC      = 4
		ALPHA        = 5
		HIRAGANA     = 6
		KATAKANA     = 7
		KANJINUMERIC = 8
		GREEK        = 9
		CYRILLIC     = 10
	)
	crs := []struct {
		class  int
		invoke bool
	}{
		{DEFAULT, true},
		{SPACE, true},
		{KANJI, false},
		{SYMBOL, true},
		{NUMERIC, true},
		{ALPHA, true},
		{HIRAGANA, true},
		{KATAKANA, true},
		{KANJINUMERIC, true},
		{GREEK, true},
		{CYRILLIC, true},
	}
	dic, err := Load(testDic)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, cr := range crs {
		if iv := dic.GroupList[cr.class]; iv != cr.invoke {
			t.Errorf("input %v: got %v, expected %v\n", cr.class, iv, cr.invoke)
		}
	}
}
