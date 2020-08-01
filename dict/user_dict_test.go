package dict

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

var testFile = "./testdata/userdic.txt"

func Test_NewUserDict(t *testing.T) {
	if _, err := NewUserDict(""); err == nil {
		t.Error("expected error, but no occurred\n")
	}
}

func Test_NewUserDictIndex(t *testing.T) {
	udict, err := NewUserDict(testFile)
	if err != nil {
		t.Fatalf("unexpected error: %v\n", err)
	}
	testdata := []struct {
		Input string
		ID    int
		Ok    bool
	}{
		{Input: "日本経済新聞", ID: 0, Ok: true},
		{Input: "朝青龍", ID: 1, Ok: true},
		{Input: "関西国際空港", ID: 2, Ok: true},
		{Input: "成田国際空港", ID: 9, Ok: false},
	}
	for _, v := range testdata {
		ids := udict.Index.Search(v.Input)
		if got := (len(ids) != 0); got != v.Ok {
			t.Errorf("want %v, got %v, %+v", v.Ok, got, ids)
		}
	}
}

func Test_NewUserDictRecords(t *testing.T) {
	r := UserDictRecords{
		{
			Text:   "日本経済新聞",
			Tokens: []string{"日本", "経済", "新聞"},
			Yomi:   []string{"ニホン", "ケイザイ", "シンブン"},
			Pos:    "カスタム名詞",
		},
		{
			Text:   "朝青龍",
			Tokens: []string{"朝青龍"},
			Yomi:   []string{"アサショウリュウ"},
			Pos:    "カスタム人名",
		},
	}
	udict, err := r.NewUserDict()
	if err != nil {
		t.Fatalf("user dict build error, %v", err)
	}
	if ids := udict.Index.Search("日本経済新聞"); len(ids) != 1 {
		t.Errorf("user dict search failed")
	} else if !reflect.DeepEqual(udict.Contents[ids[0]].Tokens, []string{"日本", "経済", "新聞"}) {
		t.Errorf("got %+v, expected %+v", udict.Contents[ids[0]].Tokens, []string{"日本", "経済", "新聞"})
	}
	if ids := udict.Index.Search("関西国際空港"); len(ids) != 0 {
		t.Errorf("user dict build failed")
	}
	if ids := udict.Index.Search("朝青龍"); len(ids) == 0 {
		t.Errorf("user dict search failed")
	} else if !reflect.DeepEqual(udict.Contents[ids[0]].Tokens, []string{"朝青龍"}) {
		t.Errorf("got %+v, expected %+v", udict.Contents[ids[0]].Tokens, []string{"朝青龍"})
	}
}

func Test_NewUserDicRecords(t *testing.T) {
	t.Run("from string", func(t *testing.T) {
		s := `
日本経済新聞,日本 経済 新聞,ニホン ケイザイ シンブン,カスタム名詞
# 関西国際空港,関西 国際 空港,カンサイ コクサイ クウコウ,カスタム地名
朝青龍,朝青龍,アサショウリュウ,カスタム人名
`
		r := strings.NewReader(s)
		rec, err := NewUserDicRecords(r)
		if err != nil {
			t.Fatalf("user dict build error, %v", err)
		}
		udict, err := rec.NewUserDict()
		if err != nil {
			t.Fatalf("user dict build error, %v", err)
		}
		if ids := udict.Index.Search("日本経済新聞"); len(ids) != 1 {
			t.Errorf("user dict search failed")
		} else if !reflect.DeepEqual(udict.Contents[ids[0]].Tokens, []string{"日本", "経済", "新聞"}) {
			t.Errorf("got %+v, expected %+v", udict.Contents[ids[0]].Tokens, []string{"日本", "経済", "新聞"})
		}
		if ids := udict.Index.Search("関西国際空港"); len(ids) != 0 {
			t.Errorf("user dict build failed")
		}
		if ids := udict.Index.Search("朝青龍"); len(ids) == 0 {
			t.Errorf("user dict search failed")
		} else if !reflect.DeepEqual(udict.Contents[ids[0]].Tokens, []string{"朝青龍"}) {
			t.Errorf("got %+v, expected %+v", udict.Contents[ids[0]].Tokens, []string{"朝青龍"})
		}

	})
	t.Run("from struct", func(t *testing.T) {
		r := UserDictRecords{
			{
				Text:   "日本経済新聞",
				Tokens: []string{"日本", "経済", "新聞"},
				Yomi:   []string{"ニホン", "ケイザイ"},
				Pos:    "カスタム名詞",
			},
		}
		_, err := r.NewUserDict()
		if err == nil {
			t.Errorf("expected error, but nil")
		}
	})
	t.Run("from records", func(t *testing.T) {
		r := UserDictRecords{
			{
				Text:   "日本経済新聞",
				Tokens: []string{"日本", "経済", "新聞"},
				Yomi:   []string{"ニホン", "ケイザイ", "シンブン"},
				Pos:    "カスタム名詞",
			},
			{
				Text:   "日本経済新聞",
				Tokens: []string{"日本", "経済", "新聞"},
				Yomi:   []string{"ニホン", "ケイザイ", "シンブン"},
				Pos:    "カスタム名詞",
			},
		}
		_, err := r.NewUserDict()
		if err == nil {
			t.Errorf("expected error, but nil")
		}
	})
	t.Run("from JSON", func(t *testing.T) {
		var got UserDictRecords
		_ = json.Unmarshal([]byte(`[
        {
            "text":"日本経済新聞",
            "tokens":["日本","経済","新聞"],
            "yomi":["ニホン","ケイザイ","シンブン"],
            "pos":"カスタム名詞"
        },
        {
            "text":"朝青龍",
            "tokens":["朝青龍"],
            "yomi":["アサショウリュウ"],
            "pos":"カスタム人名"
        }]`), &got)
		want := UserDictRecords{
			{
				Text:   "日本経済新聞",
				Tokens: []string{"日本", "経済", "新聞"},
				Yomi:   []string{"ニホン", "ケイザイ", "シンブン"},
				Pos:    "カスタム名詞",
			},
			{
				Text:   "朝青龍",
				Tokens: []string{"朝青龍"},
				Yomi:   []string{"アサショウリュウ"},
				Pos:    "カスタム人名",
			},
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}
	})
}
