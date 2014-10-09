package kagome

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFeatures01(t *testing.T) {
	tok := Token{
		Id:      0,
		Class:   KNOWN,
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dic = NewSysDic()

	f := tok.Features()
	expected := []string{"名詞", "一般", "*", "*", "*", "*", "Tシャツ", "ティーシャツ", "ティーシャツ"}
	if !reflect.DeepEqual(f, expected) {
		t.Errorf("got %v, expected %v\n", f, expected)
	}
}

func TestFeatures02(t *testing.T) {
	tok := Token{
		Id:      0,
		Class:   UNKNOWN,
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dic = NewSysDic()

	f := tok.Features()
	expected := []string{"名詞", "一般", "*", "*", "*", "*", "*"}
	if !reflect.DeepEqual(f, expected) {
		t.Errorf("got %v, expected %v\n", f, expected)
	}
}

func TestFeatures03(t *testing.T) {
	tok := Token{
		Id:      0,
		Class:   USER,
		Start:   0,
		End:     1,
		Surface: "",
	}
	tok.dic = NewSysDic()
	if udic, e := NewUserDic("_sample/userdic.txt"); e != nil {
		t.Fatalf("build user dic error: %v\n", e)
	} else {
		tok.udic = udic
	}

	f := tok.Features()
	expected := []string{"カスタム名詞", "日本/経済/新聞", "ニホン/ケイザイ/シンブン"}
	if !reflect.DeepEqual(f, expected) {
		t.Errorf("got %v, expected %v\n", f, expected)
	}
}

func TestString01(t *testing.T) {
	tok := Token{
		Id:      123,
		Class:   DUMMY,
		Start:   0,
		End:     1,
		Surface: "テスト",
	}
	expected := "テスト(0, 1)DUMMY"
	str := fmt.Sprintf("%v", tok)
	if str != expected {
		t.Errorf("got %v, expected %v\n", str, expected)
	}
}
