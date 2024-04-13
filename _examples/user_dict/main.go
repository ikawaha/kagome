package main

import (
	"fmt"

	"github.com/ikawaha/kagome-dict/dict"
	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

func main() {
	// Use IPA dictionary as a system dictionary.
	sysDic := ipa.Dict()

	// Build a user dictionary from a file.
	userDic, err := dict.NewUserDict("userdict.txt")
	if err != nil {
		panic(err)
	}

	// Specify the user dictionary as an option.
	t, err := tokenizer.New(sysDic, tokenizer.UserDict(userDic), tokenizer.OmitBosEos())
	if err != nil {
		panic(err)
	}

	tokens := t.Analyze("関西国際空港限定トートバッグ", tokenizer.Search)
	for _, token := range tokens {
		fmt.Printf("%s\t%v\n", token.Surface, token.Features())
	}

	// Output:
	// 関西国際空港    [テスト名詞 関西/国際/空港 カンサイ/コクサイ/クウコウ]
	// 限定    [名詞 サ変接続 * * * * 限定 ゲンテイ ゲンテイ]
	// トートバッグ    [名詞 一般 * * * * *]
}
