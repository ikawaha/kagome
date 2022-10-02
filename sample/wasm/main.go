//go:build ignore
// +build ignore

package main

import (
	"strings"
	"syscall/js"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

func igOK(s string, _ bool) string {
	return s
}

func tokenize(_ js.Value, args []js.Value) interface{} {
	if len(args) == 0 {
		return nil
	}
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return nil
	}
	var ret []interface{}
	tokens := t.Tokenize(args[0].String())
	for _, v := range tokens {
		//fmt.Printf("%s\t%+v%v\n", v.Surface, v.POS(), strings.Join(v.Features(), ","))
		ret = append(ret, map[string]interface{}{
			"word_id":       v.ID,
			"word_type":     v.Class.String(),
			"word_position": v.Start,
			"surface_form":  v.Surface,
			"pos":           strings.Join(v.POS(), ","),
			"base_form":     igOK(v.BaseForm()),
			"reading":       igOK(v.Reading()),
			"pronunciation": igOK(v.Pronunciation()),
		})
	}
	return ret
}

func registerCallbacks() {
	_ = ipa.Dict()
	js.Global().Set("kagome_tokenize", js.FuncOf(tokenize))
}

func main() {
	c := make(chan struct{}, 0)
	registerCallbacks()
	println("Kagome Web Assembly Ready")
	<-c
}
