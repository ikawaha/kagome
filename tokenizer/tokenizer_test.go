package tokenizer

import (
	"testing"
)

func TestTokenize(t *testing.T) {
	tokenizer := NewTokenizer()
	input := "すもももももももものうち"
	morphs, err := tokenizer.Tokenize(input)
	if err != nil {
		t.Errorf("TestTokenize: %v", err)
	}
	expectedLen := 8
	if len(morphs) != expectedLen {
		t.Errorf("TestTokenize: got %v, want %v\n%v", len(morphs), expectedLen, morphs)
	}
}

func BenchmarkTokenize(b *testing.B) {
	tokenizer := NewTokenizer()
	input := "村山富市首相は年頭にあたり首相官邸で内閣記者会と二十八日会見し、社会党の新民主連合所属議員の離党問題について「政権に影響を及ぼすことにはならない。離党者がいても、その範囲にとどまると思う」と述べ、大量離党には至らないとの見通しを示した。"
	for i := 0; i < b.N; i++ {
		tokenizer.Tokenize(input)
	}
}
