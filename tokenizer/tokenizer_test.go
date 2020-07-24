package tokenizer

import (
	"fmt"
	"testing"

	"github.com/ikawaha/kagome/v2/dict/ipa"
)

const (
	testUserDicPath = "../_sample/userdic.txt"
)

func TestHogeIPA(t *testing.T) {
	tz := New(ipa.New())
	fmt.Println(tz.Tokenize("すもももももももものうち"))
}

/*
func TestAnalyze01(t *testing.T) {
	tnz := New()
	tokens := tnz.Analyze("", Normal)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: -1, Surface: "EOS"},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}

}

func TestAnalyze02(t *testing.T) {
	tnz := New()
	tokens := tnz.Analyze("関西国際空港", Normal)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 372978, Surface: "関西国際空港", Start: 0, End: 6, Class: TokenClass(lattice.KNOWN)},
		{ID: -1, Surface: "EOS", Start: 6, End: 6},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}

}

func TestAnalyze03(t *testing.T) {
	tnz := New()
	udict, err := dict.NewUserDict(testUserDicPath)
	if err != nil {
		t.Fatalf("new user dict: unexpected error")
	}
	tnz.SetUserDict(udict)
	tokens := tnz.Analyze("関西国際空港", Normal)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 2, Surface: "関西国際空港", Start: 0, End: 6, Class: TokenClass(lattice.USER)},
		{ID: -1, Surface: "EOS", Start: 6, End: 6},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}

}

func TestAnalyze04(t *testing.T) {
	tnz := New()
	tokens := tnz.Analyze("ポポピ", Normal)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 34, Surface: "ポポピ", Start: 0, End: 3, Class: TokenClass(lattice.UNKNOWN)},
		{ID: -1, Surface: "EOS", Start: 3, End: 3},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}

}

func TestTokenizeSpecialCase(t *testing.T) {
	inputs := []string{
		"\u10000",
	}
	tnz := New()
	for _, s := range inputs {
		tnz.Tokenize(s)
	}
}

func TestTokenize(t *testing.T) {
	const input = "すもももももももものうち"
	tnz := New()
	x := tnz.Tokenize(input)
	y := tnz.Analyze(input, Normal)
	if !reflect.DeepEqual(x, y) {
		t.Errorf("got %v, expected %v", x, y)
	}
}

func TestSearcModeAnalyze01(t *testing.T) {
	tnz := New()
	tokens := tnz.Analyze("", Search)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: -1, Surface: "EOS"},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}

}

func TestSearchModeAnalyze02(t *testing.T) {
	tnz := New()
	tokens := tnz.Analyze("関西国際空港", Search)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 372968, Surface: "関西", Start: 0, End: 2, Class: TokenClass(lattice.KNOWN)},
		{ID: 168541, Surface: "国際", Start: 2, End: 4, Class: TokenClass(lattice.KNOWN)},
		{ID: 307133, Surface: "空港", Start: 4, End: 6, Class: TokenClass(lattice.KNOWN)},
		{ID: -1, Surface: "EOS", Start: 6, End: 6},
	}

	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}

}

func TestSearchModeAnalyze03(t *testing.T) {
	tnz := New()
	udict, err := dict.NewUserDict(testUserDicPath)
	if err != nil {
		t.Fatalf("new user dict: unexpected error")
	}
	tnz.SetUserDict(udict)
	tokens := tnz.Analyze("関西国際空港", Search)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 2, Surface: "関西国際空港", Start: 0, End: 6, Class: TokenClass(lattice.USER)},
		{ID: -1, Surface: "EOS", Start: 6, End: 6},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}

}

func TestSearchModeAnalyze04(t *testing.T) {
	tnz := New()
	tokens := tnz.Analyze("ポポピ", Search)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 34, Surface: "ポポピ", Start: 0, End: 3, Class: TokenClass(lattice.UNKNOWN)},
		{ID: -1, Surface: "EOS", Start: 3, End: 3},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}

}

func TestExtendedModeAnalyze01(t *testing.T) {
	tnz := New()
	tokens := tnz.Analyze("", Extended)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: -1, Surface: "EOS"},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}

}

func TestExtendedModeAnalyze02(t *testing.T) {
	tnz := New()
	tokens := tnz.Analyze("関西国際空港", Extended)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 372968, Surface: "関西", Start: 0, End: 2, Class: TokenClass(lattice.KNOWN)},
		{ID: 168541, Surface: "国際", Start: 2, End: 4, Class: TokenClass(lattice.KNOWN)},
		{ID: 307133, Surface: "空港", Start: 4, End: 6, Class: TokenClass(lattice.KNOWN)},
		{ID: -1, Surface: "EOS", Start: 6, End: 6},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}

}

func TestExtendedModeAnalyze03(t *testing.T) {
	tnz := New()
	udict, err := dict.NewUserDict(testUserDicPath)
	if err != nil {
		t.Fatalf("new user dict: unexpected error")
	}
	tnz.SetUserDict(udict)
	tokens := tnz.Analyze("関西国際空港", Extended)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 2, Surface: "関西国際空港", Start: 0, End: 6, Class: TokenClass(lattice.USER)},
		{ID: -1, Surface: "EOS", Start: 6, End: 6},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}

}

func TestExtendedModeAnalyze04(t *testing.T) {
	tnz := New()
	tokens := tnz.Analyze("ポポピ", Extended)
	expected := []Token{
		{ID: -1, Surface: "BOS"},
		{ID: 34, Surface: "ポ", Start: 0, End: 1, Class: TokenClass(lattice.DUMMY)},
		{ID: 34, Surface: "ポ", Start: 1, End: 2, Class: TokenClass(lattice.DUMMY)},
		{ID: 34, Surface: "ピ", Start: 2, End: 3, Class: TokenClass(lattice.DUMMY)},
		{ID: -1, Surface: "EOS", Start: 3, End: 3},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if tok.ID != expected[i].ID ||
			tok.Class != expected[i].Class ||
			tok.Start != expected[i].Start ||
			tok.End != expected[i].End ||
			tok.Surface != expected[i].Surface {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}

}

func TestTokenizerSetDic(t *testing.T) {
	d := SysDict()
	tnz := NewWithDic(d)

	tnz.SetDic(d)
	if tnz.dict != d.dict {
		t.Errorf("got %v, expected %v", tnz.dict, d)
	}
}

func TestNewWithDicPath(t *testing.T) {
	result, err := NewWithDicPath("../_sample/ipa.dict")

	if err != nil {
		t.Errorf("got err %v loading dictionary", err)
	}
	exp := []string{"BOS", "こんにちは", "、", "元気", "です", "か", "？", "EOS"}
	tokens := result.Tokenize("こんにちは、元気ですか？")
	for i, token := range tokens {
		if token.Surface != exp[i] {
			t.Errorf("expected %v, got %v", exp[i], token)
		}
	}
}

func TestNewWithInvalidPath(t *testing.T) {
	_, err := NewWithDicPath("invalid.zip")
	if err == nil {
		t.Errorf("no dictionary should have been loaded")
	}
}

func TestTokenizerDot(t *testing.T) {
	tnz := New()

	// test empty case
	var b bytes.Buffer
	tnz.Dot(&b, "")
	if b.String() == "" {
		t.Errorf("got empty string")
	}

	// only idling
	b.Reset()
	tnz.Dot(&b, "わたしまけましたわ")
	if b.String() == "" {
		t.Errorf("got empty string")
	}
}

func TestTokenizerAnalyzeGraph(t *testing.T) {
	tnz := New()

	// test empty case
	for _, mode := range []TokenizeMode{Normal, Search, Extended} {
		var b bytes.Buffer
		tnz.AnalyzeGraph(&b, "", mode)
		if b.String() == "" {
			t.Errorf("got empty string")
		}

		// only idling
		b.Reset()
		tnz.Dot(&b, "わたしまけましたわ")
		if b.String() == "" {
			t.Errorf("got empty string")
		}
	}
}

var benchSampleText = "人魚は、南の方の海にばかり棲んでいるのではありません。北の海にも棲んでいたのであります。北方の海の色は、青うございました。ある時、岩の上に、女の人魚があがって、あたりの景色を眺めながら休んでいました。"

func BenchmarkAnalyzeNormal(b *testing.B) {
	tnz := New()
	for i := 0; i < b.N; i++ {
		tnz.Analyze(benchSampleText, Normal)
	}
}

func BenchmarkAnalyzeSearch(b *testing.B) {
	tnz := New()
	for i := 0; i < b.N; i++ {
		tnz.Analyze(benchSampleText, Search)
	}
}

func BenchmarkAnalyzeExtended(b *testing.B) {
	tnz := New()
	for i := 0; i < b.N; i++ {
		tnz.Analyze(benchSampleText, Extended)
	}
}
*/
