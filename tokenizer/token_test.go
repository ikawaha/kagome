package tokenizer

import (
	"fmt"
	"reflect"
	"testing"

	ipa "github.com/ikawaha/kagome-dict-ipa"
	"github.com/ikawaha/kagome/v2/dict"
	"github.com/ikawaha/kagome/v2/tokenizer/lattice"
)

const userDictSample = "../_sample/userdict.txt"

func Test_TokenClassString(t *testing.T) {
	testdata := []struct {
		inp TokenClass
		out string
	}{
		{inp: DUMMY, out: "DUMMY"},
		{inp: KNOWN, out: "KNOWN"},
		{inp: UNKNOWN, out: "UNKNOWN"},
		{inp: USER, out: "USER"},
	}

	for _, v := range testdata {
		if v.inp.String() != v.out {
			t.Errorf("got %v, expected %v", v.inp.String(), v.out)
		}
	}
}

func Test_FeaturesKnown(t *testing.T) {
	tok := Token{
		ID:      0,
		Class:   TokenClass(lattice.KNOWN),
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dict = ipa.Dict()

	f := tok.Features()
	expected := []string{"名詞", "一般", "*", "*", "*", "*", "Tシャツ", "ティーシャツ", "ティーシャツ"}
	if !reflect.DeepEqual(f, expected) {
		t.Errorf("got %v, expected %v", f, expected)
	}
}

func Test_FeaturesUnknown(t *testing.T) {
	tok := Token{
		ID:      0,
		Class:   TokenClass(lattice.UNKNOWN),
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dict = ipa.Dict()

	f := tok.Features()
	expected := []string{"名詞", "固有名詞", "地域", "一般", "*", "*", "*"}
	if !reflect.DeepEqual(f, expected) {
		t.Errorf("got %v, expected %v", f, expected)
	}
}

func Test_FeaturesUser(t *testing.T) {
	tok := Token{
		ID:      0,
		Class:   TokenClass(lattice.USER),
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dict = ipa.Dict()
	if udic, err := dict.NewUserDict(userDictSample); err != nil {
		t.Fatalf("build user dict error: %v", err)
	} else {
		tok.udict = udic
	}

	f := tok.Features()
	expected := []string{"カスタム名詞", "日本/経済/新聞", "ニホン/ケイザイ/シンブン"}
	if !reflect.DeepEqual(f, expected) {
		t.Errorf("got %v, expected %v", f, expected)
	}
}

func Test_FeaturesDummy(t *testing.T) {
	tok := Token{
		ID:      0,
		Class:   DUMMY,
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dict = ipa.Dict()
	if udic, err := dict.NewUserDict(userDictSample); err != nil {
		t.Fatalf("build user dict error: %v", err)
	} else {
		tok.udict = udic
	}

	f := tok.Features()
	if len(f) != 0 {
		t.Errorf("got %v, expected empty", f)
	}
}

func Test_FeaturesAndPosKnown(t *testing.T) {
	tok := Token{
		ID:      0,
		Class:   TokenClass(lattice.KNOWN),
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dict = ipa.Dict()

	f := tok.Features()
	expected := []string{"名詞", "一般", "*", "*", "*", "*", "Tシャツ", "ティーシャツ", "ティーシャツ"}
	if !reflect.DeepEqual(f, expected) {
		t.Errorf("got %v, expected %v", f, expected)
	}
	if p := tok.Pos(); p != f[0] {
		t.Errorf("got %v, expected %v", p, f[0])
	}
}

func Test_FeaturesAndPosUnknown(t *testing.T) {
	tok := Token{
		ID:      0,
		Class:   TokenClass(lattice.UNKNOWN),
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dict = ipa.Dict()

	f := tok.Features()
	expected := []string{"名詞", "固有名詞", "地域", "一般", "*", "*", "*"}
	if !reflect.DeepEqual(f, expected) {
		t.Errorf("got %v, expected %v", f, expected)
	}
	if p := tok.Pos(); p != f[0] {
		t.Errorf("got %v, expected %v", p, f[0])
	}
}

func Test_FeaturesAndPosUser(t *testing.T) {
	tok := Token{
		ID:      0,
		Class:   TokenClass(lattice.USER),
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dict = ipa.Dict()
	if udic, err := dict.NewUserDict(userDictSample); err != nil {
		t.Fatalf("build user dict error: %v", err)
	} else {
		tok.udict = udic
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

func Test_FeaturesAndPosUserDict(t *testing.T) {
	tok := Token{
		ID:      0,
		Class:   DUMMY,
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dict = ipa.Dict()
	if udic, err := dict.NewUserDict(userDictSample); err != nil {
		t.Fatalf("build user dict error: %v", err)
	} else {
		tok.udict = udic
	}

	f := tok.Features()
	if len(f) != 0 {
		t.Errorf("got %v, expected empty", f)
	}
	if p := tok.Pos(); p != "" {
		t.Errorf("got %v, expected empty", p)
	}
}

func Test_TokenString(t *testing.T) {
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
