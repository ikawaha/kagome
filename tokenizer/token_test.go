package tokenizer

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ikawaha/kagome-dict/dict"
)

const (
	userDictSample = "../sample/userdict.txt"
)

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
		if got, want := v.out, v.inp.String(); want != got {
			t.Errorf("want %v, got %v", got, v.inp.String())
		}
	}
}

func Test_FeaturesKnown(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tok := Token{
		ID:      0,
		Class:   KNOWN,
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dict = d

	got := tok.Features()
	want := []string{"名詞", "一般", "*", "*", "*", "*", "Tシャツ", "ティーシャツ", "ティーシャツ"}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func Test_FeatureAtKnown(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tok := Token{
		ID:      0,
		Class:   KNOWN,
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dict = d

	fs := tok.Features()
	want := []string{"名詞", "一般", "*", "*", "*", "*", "Tシャツ", "ティーシャツ", "ティーシャツ"}
	if got := fs; !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
	t.Run("features at", func(t *testing.T) {
		for i, want := range want {
			if got, ok := tok.FeatureAt(i); !ok {
				t.Errorf("want ok, got !ok, %d", i)
			} else if got != want {
				t.Errorf("want %s, got %s", want, got)
			}
		}
	})

	t.Run("index out of bound", func(t *testing.T) {
		if f, ok := tok.FeatureAt(-1); f != "" || ok {
			t.Errorf("index < 0: expected empty feature and false, got %q, %v", f, ok)
		}
		if f, ok := tok.FeatureAt(len(want)); f != "" || ok {
			t.Errorf("index >= len(want): expected empty feature and false, got %q, %v", f, ok)
		}
	})
}

func Test_FeaturesUnknown(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tok := Token{
		ID:      0,
		Class:   UNKNOWN,
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dict = d

	got := tok.Features()
	want := []string{"名詞", "固有名詞", "地域", "一般", "*", "*", "*"}

	t.Run("features at", func(t *testing.T) {
		if len(want) != len(got) {
			t.Errorf("len: want %d, got %d", len(want), len(got))
		}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("pos: want %v, got %v", want, got)
		}
	})

	t.Run("index out of bound", func(t *testing.T) {
		if f, ok := tok.FeatureAt(-1); f != "" || ok {
			t.Errorf("index < 0: expected empty feature and false, got %q, %v", f, ok)
		}
		if f, ok := tok.FeatureAt(len(want)); f != "" || ok {
			t.Errorf("index >= len(want): expected empty feature and false, got %q, %v", f, ok)
		}
	})
}

func Test_FeatureAtUnknown(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tok := Token{
		ID:      0,
		Class:   UNKNOWN,
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dict = d

	fs := tok.Features()
	want := []string{"名詞", "固有名詞", "地域", "一般", "*", "*", "*"}
	if got := fs; !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
	for i, want := range want {
		if got, ok := tok.FeatureAt(i); !ok {
			t.Errorf("want ok, got !ok, %d", i)
		} else if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	}
}

func Test_FeaturesUser(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tok := Token{
		ID:      0,
		Class:   USER,
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dict = d
	if udic, err := dict.NewUserDict(userDictSample); err != nil {
		t.Fatalf("build user dict error: %v", err)
	} else {
		tok.udict = udic
	}

	got := tok.Features()
	want := []string{"カスタム名詞", "日本/経済/新聞", "ニホン/ケイザイ/シンブン"}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func Test_FeatureAtUsers(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tok := Token{
		ID:      0,
		Class:   USER,
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dict = d
	if udic, err := dict.NewUserDict(userDictSample); err != nil {
		t.Fatalf("build user dict error: %v", err)
	} else {
		tok.udict = udic
	}

	fs := tok.Features()
	want := []string{"カスタム名詞", "日本/経済/新聞", "ニホン/ケイザイ/シンブン"}
	if got := fs; !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
	for i, want := range want {
		if got, ok := tok.FeatureAt(i); !ok {
			t.Errorf("want ok, got !ok, %d", i)
		} else if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	}
}

func Test_FeaturesDummy(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tok := Token{
		ID:      0,
		Class:   DUMMY,
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dict = d
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
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tok := Token{
		ID:      0,
		Class:   KNOWN,
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dict = d

	f := tok.Features()
	want := []string{"名詞", "一般", "*", "*", "*", "*", "Tシャツ", "ティーシャツ", "ティーシャツ"}
	if !reflect.DeepEqual(want, f) {
		t.Errorf("want %v, got %v", want, f)
	}
	if got, want := tok.POS(), want[0:4]; !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func Test_FeaturesAndPosUnknown(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tok := Token{
		ID:      0,
		Class:   UNKNOWN,
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dict = d

	got := tok.Features()
	want := []string{"名詞", "固有名詞", "地域", "一般", "*", "*", "*"}
	if len(want) != len(got) {
		t.Errorf("len: want %d, got %d", len(want), len(got))
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
	if want, got := want[0:4], tok.POS(); !reflect.DeepEqual(want, got) {
		t.Errorf("want %+v, got %+v", want, got)
	}
}

func Test_FeaturesAndPosUser(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tok := Token{
		ID:      0,
		Class:   USER,
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dict = d
	if udic, err := dict.NewUserDict(userDictSample); err != nil {
		t.Fatalf("build user dict error: %v", err)
	} else {
		tok.udict = udic
	}

	f := tok.Features()
	want := []string{"カスタム名詞", "日本/経済/新聞", "ニホン/ケイザイ/シンブン"}
	if !reflect.DeepEqual(want, f) {
		t.Errorf("want %v, got %v", want, f)
	}
	if got := tok.POS(); !reflect.DeepEqual(want[0:1], got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func Test_FeaturesAndPosUserDict(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tok := Token{
		ID:      0,
		Class:   DUMMY,
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dict = d
	if udic, err := dict.NewUserDict(userDictSample); err != nil {
		t.Fatalf("build user dict error: %v", err)
	} else {
		tok.udict = udic
	}

	f := tok.Features()
	if len(f) != 0 {
		t.Errorf("got %v, expected empty", f)
	}
	if got := tok.POS(); len(got) > 0 {
		t.Errorf("want empty, got %v", got)
	}
}

func Test_TokenString(t *testing.T) {
	tok := Token{
		ID:      123,
		Class:   DUMMY,
		Start:   0,
		End:     1,
		Surface: "テスト",
	}
	want := `0:"テスト" (0: 0, 1) DUMMY [123]`
	got := fmt.Sprintf("%v", tok)
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}

func Test_InflectionalType(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	tokens := tnz.Tokenize("寿司を食べたい")
	if want, got := 6, len(tokens); want != got {
		t.Fatalf("token length: want %d, got %d", want, got)
	}
	// BOS
	if got, ok := tokens[0].InflectionalType(); ok {
		t.Errorf("want !ok, got %q", got)
	}
	// 食べ 動詞,自立,*,*,一段,連用形,食べる,タベ,タベ
	if got, ok := tokens[3].InflectionalType(); !ok {
		t.Error("want ok, but !ok")
	} else if want := "一段"; want != got {
		t.Fatalf("want %s, got %s", want, got)
	}
}

func Test_InflectionalForm(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	t.Run("known", func(t *testing.T) {
		tokens := tnz.Tokenize("寿司を食べたい")
		if want, got := 6, len(tokens); want != got {
			t.Fatalf("token length: want %d, got %d", want, got)
		}
		// BOS
		if got, ok := tokens[0].InflectionalForm(); ok {
			t.Errorf("want !ok, got %q", got)
		}
		// 食べ 動詞,自立,*,*,一段,連用形,食べる,タベ,タベ
		if got, ok := tokens[3].InflectionalForm(); !ok {
			t.Error("want ok, but !ok")
		} else if want := "連用形"; want != got {
			t.Fatalf("want %s, got %s", want, got)
		}
	})

	t.Run("unknown", func(t *testing.T) {
		tokens := tnz.Tokenize("トトロ")
		if want, got := 3, len(tokens); want != got {
			t.Fatalf("token length: want %d, got %d", want, got)
		}
		// UNKNOWN
		if _, ok := tokens[1].InflectionalForm(); ok {
			t.Error("want !ok, but ok")
		}
	})
}

func Test_BaseForm(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	tokens := tnz.Tokenize("寿司を食べたい")
	if want, got := 6, len(tokens); want != got {
		t.Fatalf("token length: want %d, got %d", want, got)
	}
	// BOS
	if got, ok := tokens[0].BaseForm(); ok {
		t.Errorf("want !ok, got %q", got)
	}
	// 食べ->食べる
	if got, ok := tokens[3].BaseForm(); !ok {
		t.Error("want ok, but !ok")
	} else if want := "食べる"; want != got {
		t.Fatalf("want %s, got %s", want, got)
	}
}

func Test_Reading(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	tokens := tnz.Tokenize("公園に行く")
	if want, got := 5, len(tokens); want != got {
		t.Fatalf("token length: want %d, got %d", want, got)
	}
	// BOS
	if got, ok := tokens[0].Reading(); ok {
		t.Errorf("want !ok, got %q", got)
	}
	// 公園->コウエン
	if got, ok := tokens[1].Reading(); !ok {
		t.Error("want ok, but !ok")
	} else if want := "コウエン"; want != got {
		t.Fatalf("want %s, got %s", want, got)
	}
}

func Test_Pronunciation(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	tokens := tnz.Tokenize("公園に行く")
	if want, got := 5, len(tokens); want != got {
		t.Fatalf("token length: want %d, got %d", want, got)
	}
	// BOS
	if got, ok := tokens[0].Pronunciation(); ok {
		t.Errorf("want !ok, got %q", got)
	}
	// 公園->コウエン
	if got, ok := tokens[1].Pronunciation(); !ok {
		t.Error("want ok, but !ok")
	} else if want := "コーエン"; want != got {
		t.Fatalf("want %s, got %s", want, got)
	}
}

func Test_POS(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	t.Run("known", func(t *testing.T) {
		tokens := tnz.Tokenize("公園に行く")
		if want, got := 5, len(tokens); want != got {
			t.Fatalf("token length: want %d, got %d", want, got)
		}
		// BOS
		if got := tokens[0].POS(); len(got) > 0 {
			t.Errorf("want empty, got %+v", got)
		}
		// 行く
		if want, got := []string{"動詞", "自立", "*", "*"}, tokens[3].POS(); !reflect.DeepEqual(want, got) {
			t.Fatalf("want %+v, got %+v", want, got)
		}
	})

	t.Run("unknown", func(t *testing.T) {
		tokens := tnz.Tokenize("トトロ")
		if want, got := 3, len(tokens); want != got {
			t.Fatalf("token length: want %d, got %d", want, got)
		}
		// BOS
		if got := tokens[0].POS(); len(got) > 0 {
			t.Errorf("want empty, got %+v", got)
		}
		// UNKNOWN
		if want, got := []string{"名詞", "固有名詞", "組織", "*"}, tokens[1].POS(); !reflect.DeepEqual(want, got) {
			t.Fatalf("want %+v, got %+v", want, got)
		}
	})
}
