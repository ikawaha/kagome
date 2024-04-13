package main

import (
	"fmt"
	"strings"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

func main() {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		panic(err)
	}
	// tokenize
	fmt.Println("---tokenize---")
	tokens := t.Tokenize("すもももももももものうち")
	for _, token := range tokens {
		features := strings.Join(token.Features(), ",")
		fmt.Printf("%s\t%v\n", token.Surface, features)
	}

	// Output:
	//---tokenize---
	//すもも	名詞,一般,*,*,*,*,すもも,スモモ,スモモ
	//も	助詞,係助詞,*,*,*,*,も,モ,モ
	//もも	名詞,一般,*,*,*,*,もも,モモ,モモ
	//も	助詞,係助詞,*,*,*,*,も,モ,モ
	//もも	名詞,一般,*,*,*,*,もも,モモ,モモ
	//の	助詞,連体化,*,*,*,*,の,ノ,ノ
	//うち	名詞,非自立,副詞可能,*,*,*,うち,ウチ,ウチ
}
