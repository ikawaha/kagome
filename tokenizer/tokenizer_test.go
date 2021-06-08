package tokenizer

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/ikawaha/kagome-dict/dict"
	"github.com/ikawaha/kagome/v2/tokenizer/lattice"
)

const (
	testDictPath = "../testdata/ipa.dict"
)

func equalTokens(lhs, rhs Token) bool {
	return lhs.Index == rhs.Index &&
		lhs.ID == rhs.ID &&
		lhs.Class == rhs.Class &&
		lhs.Position == rhs.Position &&
		lhs.Start == rhs.Start &&
		lhs.End == rhs.End &&
		lhs.Surface == rhs.Surface
}

func Example_tokenize_mode() {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		panic(err)
	}
	t, err := New(d)
	if err != nil {
		panic(err)
	}
	for _, mode := range []TokenizeMode{Normal, Search, Extended} {
		tokens := t.Analyze("関西国際空港", Normal)
		fmt.Printf("---%s---", mode)
		for _, token := range tokens {
			if token.Class == DUMMY {
				// BOS: Begin Of Sentence, EOS: End Of Sentence.
				fmt.Printf("%s\n", token.Surface)
				continue
			}
			features := strings.Join(token.Features(), ",")
			fmt.Printf("%s\t%v\n", token.Surface, features)
		}
	}
}

func Test_Tokenizer_New(t *testing.T) {
	if _, err := New(nil); err == nil {
		t.Errorf("expected empty dictionary error")
	} else if err.Error() != "empty dictionary" {
		t.Errorf("expected empty dictionary, got %v", err)
	}
}

func Test_TokenizeModeString(t *testing.T) {
	testdata := []struct {
		mode       TokenizeMode
		expression string
	}{
		{mode: Normal, expression: "normal"},
		{mode: Search, expression: "search"},
		{mode: Extended, expression: "extend"},
		{mode: TokenizeMode(999), expression: "unknown tokenize mode (999)"},
	}
	for _, v := range testdata {
		if want, got := v.expression, v.mode.String(); want != got {
			t.Errorf("want %q, got %q", want, got)
		}
	}
}

func Test_AnalyzeEmptyInput(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tokens := tnz.Analyze("", Normal)
	expected := []Token{
		{Index: 0, ID: -1, Surface: "BOS"},
		{Index: 1, ID: -1, Surface: "EOS"},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if !equalTokens(tok, expected[i]) {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}
}

func Test_Analyze(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	input := "関西国際空港"
	tokens := tnz.Analyze(input, Normal)
	want := []Token{
		{Index: 0, ID: -1, Surface: "BOS"},
		{Index: 1, ID: 372978, Surface: input, Position: 0, Start: 0, End: 6, Class: TokenClass(lattice.KNOWN)},
		{Index: 2, ID: -1, Surface: "EOS", Position: len(input), Start: 6, End: 6},
	}
	if len(tokens) != len(want) {
		t.Fatalf("got %v, want %v", tokens, want)
	}
	for i, tok := range tokens {
		if !equalTokens(tok, want[i]) {
			t.Errorf("got %+v, want %+v", tok, want[i])
		}
	}
}

func Test_AnalyzeUnknown(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tokens := tnz.Analyze("ポポピ", Normal)
	expected := []Token{
		{Index: 0, ID: -1, Surface: "BOS"},
		{Index: 1, ID: 34, Surface: "ポポピ", Start: 0, End: 3, Class: TokenClass(lattice.UNKNOWN)},
		{Index: 2, ID: -1, Surface: "EOS", Start: 3, End: 3, Position: 9},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if !equalTokens(tok, expected[i]) {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}
}

func Test_TokenizeSpecialCase(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	inputs := []string{
		"\u10000",
	}
	for _, s := range inputs {
		tnz.Tokenize(s) // does not panic.
	}
}

func Test_Tokenize(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	inputs := []string{
		"すもももももももものうち",
		"人魚は、南の方の海にばかり棲んでいるのではありません。",
	}
	for _, input := range inputs {
		x := tnz.Tokenize(input)
		y := tnz.Analyze(input, Normal)
		if !reflect.DeepEqual(x, y) {
			t.Errorf("got %v, expected %v", x, y)
		}
	}
}

func Test_AnalyzeWithSearchModeEmptyInput(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tokens := tnz.Analyze("", Search)
	expected := []Token{
		{Index: 0, ID: -1, Surface: "BOS"},
		{Index: 1, ID: -1, Surface: "EOS"},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if !equalTokens(tok, expected[i]) {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}
}

func Test_AnalyzeWithSearchMode(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tokens := tnz.Analyze("関西国際空港", Search)
	expected := []Token{
		{Index: 0, ID: -1, Surface: "BOS"},
		{Index: 1, ID: 372967, Surface: "関西", Start: 0, End: 2, Position: 0, Class: TokenClass(lattice.KNOWN)},
		{Index: 2, ID: 168542, Surface: "国際", Start: 2, End: 4, Position: 6, Class: TokenClass(lattice.KNOWN)},
		{Index: 3, ID: 307134, Surface: "空港", Start: 4, End: 6, Position: 12, Class: TokenClass(lattice.KNOWN)},
		{Index: 4, ID: -1, Surface: "EOS", Start: 6, End: 6, Position: 18},
	}

	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if !equalTokens(tok, expected[i]) {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}
}

func Test_AnalyzeWithSearchModeUnknown(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	tokens := tnz.Analyze("ポポピ", Search)
	expected := []Token{
		{Index: 0, ID: -1, Surface: "BOS"},
		{Index: 1, ID: 34, Surface: "ポポピ", Start: 0, End: 3, Class: TokenClass(lattice.UNKNOWN)},
		{Index: 2, ID: -1, Surface: "EOS", Start: 3, End: 3, Position: 9},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if !equalTokens(tok, expected[i]) {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}
}

func Test_AnalyzeWithExtendedModeEmpty(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	tokens := tnz.Analyze("", Extended)
	expected := []Token{
		{Index: 0, ID: -1, Surface: "BOS"},
		{Index: 1, ID: -1, Surface: "EOS"},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if !equalTokens(tok, expected[i]) {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}
}

func Test_AnalyzeWithExtendedMode(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	tokens := tnz.Analyze("関西国際空港", Extended)
	expected := []Token{
		{Index: 0, ID: -1, Surface: "BOS"},
		{Index: 1, ID: 372967, Surface: "関西", Start: 0, End: 2, Position: 0, Class: TokenClass(lattice.KNOWN)},
		{Index: 2, ID: 168542, Surface: "国際", Start: 2, End: 4, Position: 6, Class: TokenClass(lattice.KNOWN)},
		{Index: 3, ID: 307134, Surface: "空港", Start: 4, End: 6, Position: 12, Class: TokenClass(lattice.KNOWN)},
		{Index: 4, ID: -1, Surface: "EOS", Start: 6, End: 6, Position: 18},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if !equalTokens(tok, expected[i]) {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}
}

func Test_AnalyzeWithExtendedModeUnknown(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	tokens := tnz.Analyze("ポポピ", Extended)
	expected := []Token{
		{Index: 0, ID: -1, Surface: "BOS"},
		{Index: 1, ID: 34, Surface: "ポ", Start: 0, End: 1, Position: 0, Class: TokenClass(lattice.DUMMY)},
		{Index: 2, ID: 34, Surface: "ポ", Start: 1, End: 2, Position: 3, Class: TokenClass(lattice.DUMMY)},
		{Index: 3, ID: 34, Surface: "ピ", Start: 2, End: 3, Position: 6, Class: TokenClass(lattice.DUMMY)},
		{Index: 4, ID: -1, Surface: "EOS", Start: 3, End: 3, Position: 9},
	}
	if len(tokens) != len(expected) {
		t.Fatalf("got %v, expected %v", tokens, expected)
	}
	for i, tok := range tokens {
		if !equalTokens(tok, expected[i]) {
			t.Errorf("got %v, expected %v", tok, expected[i])
		}
	}
}

func Test_TokenizerDot(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

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

func Test_TokenizerAnalyzeGraph(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

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

func Test_Wakati(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	testdata := []struct {
		Input  string
		Output []string
	}{
		{
			Input:  "すもももももももものうち",
			Output: []string{"すもも", "も", "もも", "も", "もも", "の", "うち"},
		},
		{
			Input:  "寿司が食べたい。",
			Output: []string{"寿司", "が", "食べ", "たい", "。"},
		},
	}
	for _, v := range testdata {
		got := tnz.Wakati(v.Input)
		if want := v.Output; !reflect.DeepEqual(want, got) {
			t.Errorf("want %+v, got %+v", want, got)
		}
	}
}

var benchSampleText = "人魚は、南の方の海にばかり棲んでいるのではありません。北の海にも棲んでいたのであります。北方の海の色は、青うございました。ある時、岩の上に、女の人魚があがって、あたりの景色を眺めながら休んでいました。"

func BenchmarkAnalyzeNormal(b *testing.B) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		b.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		b.Fatalf("unexpected error, %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tnz.Analyze(benchSampleText, Normal)
	}
}

func BenchmarkAnalyzeSearch(b *testing.B) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		b.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		b.Fatalf("unexpected error, %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tnz.Analyze(benchSampleText, Search)
	}
}

func BenchmarkAnalyzeExtended(b *testing.B) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		b.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		b.Fatalf("unexpected error, %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tnz.Analyze(benchSampleText, Extended)
	}
}

func BenchmarkTooLongUnknownToken(b *testing.B) {
	input := `GOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO`
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		b.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		b.Fatalf("unexpected error, %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tnz.Tokenize(input)
	}
}
