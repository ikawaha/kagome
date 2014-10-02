package kagome

import (
	"reflect"
	"testing"
)

func TestNewUserDic01(t *testing.T) {
	if _, e := NewUserDic(""); e == nil {
		t.Error("expected error, but no occured\n")
	}
}

func TestNewUserDicIndex01(t *testing.T) {
	udic, e := NewUserDic("_sample/userdic.txt")
	if e != nil {
		t.Fatalf("unexpected error: %v\n", e)
	}
	type tuple struct {
		inp string
		id  int
		ok  bool
	}
	callAndRespose := []tuple{
		tuple{inp: "日本経済新聞", id: 0, ok: true},
		tuple{inp: "朝青龍", id: 1, ok: true},
		tuple{inp: "関西国際空港", id: 2, ok: true},
		tuple{inp: "成田国際空港", id: 9, ok: false},
	}
	for _, cr := range callAndRespose {
		id, ok := udic.Index.FindString(cr.inp)
		if ok != cr.ok {
			t.Errorf("got %v, expected %v\n", ok, cr.ok)
		} else {
			if !ok {
				continue
			}
		}
		if id != cr.id {
			t.Errorf("got %v, expected %v\n", id, cr.id)
		}
	}
}

func TestNewUserDicContents01(t *testing.T) {
	udic, e := NewUserDic("_sample/userdic.txt")
	if e != nil {
		t.Fatalf("unexpected error: %v\n", e)
	}
	expectedLen := 3
	if len(udic.Contents) != expectedLen {
		t.Errorf("got %v, expected %v\n", len(udic.Contents), expectedLen)
	}

	type tuple struct {
		inp int
		out UserDicContent
	}
	callAndRespose := []tuple{
		tuple{
			inp: 0,
			out: UserDicContent{
				Tokens: []string{"日本", "経済", "新聞"},
				Yomi:   []string{"ニホン", "ケイザイ", "シンブン"},
				Pos:    "カスタム名詞",
			},
		},
		tuple{
			inp: 1,
			out: UserDicContent{
				Tokens: []string{"朝青龍"},
				Yomi:   []string{"アサショウリュウ"},
				Pos:    "カスタム人名",
			},
		},
		tuple{inp: 2,
			out: UserDicContent{
				Tokens: []string{"関西", "国際", "空港"},
				Yomi:   []string{"カンサイ", "コクサイ", "クウコウ"},
				Pos:    "テスト名詞",
			},
		},
	}
	for _, cr := range callAndRespose {
		c := udic.Contents[cr.inp]
		if !reflect.DeepEqual(c, cr.out) {
			t.Errorf("got %v, expected %v\n", c, cr.out)
		}
	}
}
