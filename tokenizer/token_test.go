package tokenizer

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ikawaha/kagome-dict/dict"
)

const (
	userDictSample = "../testdata/userdict.txt"
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

func Test_FeaturesUserExtra(t *testing.T) {
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

	got := tok.UserExtra()
	want := &UserExtra{
		Tokens:   []string{"日本", "経済", "新聞"},
		Readings: []string{"ニホン", "ケイザイ", "シンブン"},
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
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

func TestNewTokenData(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	t.Run("DUMMY/UNKNOWN", func(t *testing.T) {
		tokens := tnz.Tokenize("トトロ")
		if want, got := 3, len(tokens); want != got {
			t.Fatalf("token length: want %d, got %d", want, got)
		}
		// BOS
		data0 := NewTokenData(tokens[0])
		if want, got := []string{}, data0.POS; !reflect.DeepEqual(got, want) {
			t.Errorf("want empty, got %+v", got)
		}
		if want, got := []string{}, data0.Features; !reflect.DeepEqual(got, want) {
			t.Errorf("want empty, got %+v", got)
		}
		if want, got := DUMMY.String(), data0.Class; got != want {
			t.Errorf("want %v, got %v", want, got)
		}
		// UNKNOWN
		data1 := NewTokenData(tokens[1])
		if want, got := "トトロ", data1.Surface; want != got {
			t.Errorf("want %v, got %v", want, got)
		}
		if want, got := []string{"名詞", "固有名詞", "組織", "*"}, data1.POS; !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}
		if want, got := []string{"名詞", "固有名詞", "組織", "*", "*", "*", "*"}, data1.Features; !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}
		if want, got := UNKNOWN.String(), data1.Class; got != want {
			t.Errorf("want %v, got %v", want, got)
		}
	})

	t.Run("known token", func(t *testing.T) {
		tokens := tnz.Tokenize("行か")
		if want, got := 3, len(tokens); want != got {
			t.Fatalf("token length: want %d, got %d", want, got)
		}
		data := NewTokenData(tokens[1])
		tok := tokens[1] // 行か
		if want, got := "行か", tok.Surface; want != got {
			t.Errorf("want %v, got %v", want, got)
		}
		if got, want := data.POS, tok.POS(); !reflect.DeepEqual(got, want) {
			t.Errorf("want %v, got %v", want, got)
		}
		if got, want := data.Features, tok.Features(); !reflect.DeepEqual(got, want) {
			t.Errorf("want %v, got %v", want, got)
		}
		if got, want := data.BaseForm, "行く"; !reflect.DeepEqual(got, want) {
			t.Errorf("want %v, got %v", want, got)
		}
		if got, want := data.Reading, "イカ"; !reflect.DeepEqual(got, want) {
			t.Errorf("want %v, got %v", want, got)
		}
		if got, want := data.Pronunciation, "イカ"; !reflect.DeepEqual(got, want) {
			t.Errorf("want %v, got %v", want, got)
		}
	})
}

func TestEqual(t *testing.T) {
	testdata := []struct {
		ID      int
		Surface string
		Class   TokenClass
		Want    bool
	}{
		{
			ID:      1,
			Surface: "ねこ",
			Class:   KNOWN,
			Want:    true,
		},
		{
			ID:      1,
			Surface: "いぬ",
			Class:   KNOWN,
			Want:    false,
		},
	}
	for i, v := range testdata {
		if got := (Token{
			ID:       v.ID,
			Class:    v.Class,
			Surface:  v.Surface,
			Index:    111,
			Position: 456,
			Start:    3,
			End:      4,
		}).Equal(Token{
			ID:       1,
			Class:    KNOWN,
			Surface:  "ねこ",
			Index:    1234, // ↓ ignored fields
			Position: 5,
			Start:    56,
			End:      58,
		}); got != v.Want {
			t.Errorf("#%d: want %t, got %t", i, v.Want, got)
		}
	}
}

func Test_EqualFeatures(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	tokens1 := tnz.Tokenize("公園に行くトトロ") // BOS/公園/に/行く/トトロ/EOS
	if want, got := 6, len(tokens1); want != got {
		t.Fatalf("token length: want %d, got %d", want, got)
	}
	tokens2 := tnz.Tokenize("学校に行くトトロ")
	if want, got := 6, len(tokens2); want != got {
		t.Fatalf("token length: want %d, got %d", want, got)
	}

	testdata := []struct {
		name     string
		lhs, rhs Token
		want     bool
	}{
		{
			name: "BOS vs BOS",
			lhs:  tokens1[0],
			rhs:  tokens2[0],
			want: true,
		},
		{
			name: "公園 vs 学校",
			lhs:  tokens1[1],
			rhs:  tokens2[1],
			want: false,
		},
		{
			name: "に vs に",
			lhs:  tokens1[2],
			rhs:  tokens2[2],
			want: true,
		},
		{
			name: "行く vs 行く",
			lhs:  tokens1[3],
			rhs:  tokens2[3],
			want: true,
		},
		{
			name: "トトロ vs トトロ",
			lhs:  tokens1[4],
			rhs:  tokens2[4],
			want: true,
		},
		{
			name: "EOS vs EOS",
			lhs:  tokens1[5],
			rhs:  tokens2[5],
			want: true,
		},
		{
			name: "BOS vs EOS",
			lhs:  tokens1[0],
			rhs:  tokens2[5],
			want: true,
		},
		{
			name: "学校 vs トトロ",
			lhs:  tokens1[0],
			rhs:  tokens2[4],
			want: false,
		},
	}
	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			if got := EqualFeatures(tt.lhs.Features(), tt.rhs.Features()); tt.want != got {
				t.Errorf("want %t, got %t, %q%+v vs %q%+v", tt.want, got, tt.lhs.Surface, tt.lhs.Features(), tt.rhs.Surface, tt.rhs.Features())
			}
			if got := tt.lhs.EqualFeatures(tt.rhs); tt.want != got {
				t.Errorf("want %t, got %t, %q%+v vs %q%+v", tt.want, got, tt.lhs.Surface, tt.lhs.Features(), tt.rhs.Surface, tt.rhs.Features())
			}
		})
	}
}

func Test_EqualPOS(t *testing.T) {
	d, err := dict.LoadDictFile(testDictPath)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	tnz, err := New(d)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	tokens1 := tnz.Tokenize("公園に行くトトロ") // BOS/公園/に/行く/トトロ/EOS
	if want, got := 6, len(tokens1); want != got {
		t.Fatalf("token length: want %d, got %d", want, got)
	}
	tokens2 := tnz.Tokenize("学校に行くトトロ")
	if want, got := 6, len(tokens2); want != got {
		t.Fatalf("token length: want %d, got %d", want, got)
	}

	testdata := []struct {
		name     string
		lhs, rhs Token
		want     bool
	}{
		{
			name: "BOS vs BOS",
			lhs:  tokens1[0],
			rhs:  tokens2[0],
			want: true,
		},
		{
			name: "公園 vs 学校",
			lhs:  tokens1[1],
			rhs:  tokens2[1],
			want: true,
		},
		{
			name: "に vs に",
			lhs:  tokens1[2],
			rhs:  tokens2[2],
			want: true,
		},
		{
			name: "行く vs 行く",
			lhs:  tokens1[3],
			rhs:  tokens2[3],
			want: true,
		},
		{
			name: "トトロ vs トトロ",
			lhs:  tokens1[4],
			rhs:  tokens2[4],
			want: true,
		},
		{
			name: "EOS vs EOS",
			lhs:  tokens1[5],
			rhs:  tokens2[5],
			want: true,
		},
		{
			name: "BOS vs EOS",
			lhs:  tokens1[0],
			rhs:  tokens2[5],
			want: true,
		},
		{
			name: "学校 vs トトロ",
			lhs:  tokens1[1],
			rhs:  tokens2[4],
			want: true,
		},
		{
			name: "学校 vs 行く",
			lhs:  tokens1[1],
			rhs:  tokens2[3],
			want: false,
		},
	}
	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.lhs.EqualPOS(tt.rhs); tt.want != got {
				t.Errorf("want %t, got %t, %q%+v vs %q%+v", tt.want, got, tt.lhs.Surface, tt.lhs.POS(), tt.rhs.Surface, tt.rhs.POS())
			}
		})
	}
}
