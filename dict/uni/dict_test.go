package uni

import (
	"testing"

	"github.com/ikawaha/kagome/v2/dict"
)

const (
	IPADictEntrySize = 392126 + 1
)

func TestNew(t *testing.T) {
	d := New()
	if want, got := IPADictEntrySize, len(d.Morphs); want != got {
		t.Errorf("want %d, got %d", want, got)
	}
	if want, got := IPADictEntrySize, len(d.Contents); want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestLoadDictFile(t *testing.T) {
	d, err := dict.LoadDictFile("./ipa.dict")
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	if want, got := IPADictEntrySize, len(d.Morphs); want != got {
		t.Errorf("want %d, got %d", want, got)
	}
	if want, got := IPADictEntrySize, len(d.Contents); want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

//func TestNewDicSimple(t *testing.T) {
//	d, err := NewDictSimple(testDic)
//	if err != nil {
//		t.Fatalf("unexpected error: %v", err)
//	}
//	if expected, c := IPADictEntrySize, len(d.dic.Morphs); c != expected {
//		t.Errorf("got %v, expected %v", c, expected)
//	}
//	if expected, c := 0, len(d.dic.Contents); c != expected {
//		t.Errorf("got %v, expected %v", c, expected)
//	}
//}

func TestSingleton(t *testing.T) {
	a := New()
	b := New()
	if a != b {
		t.Errorf("got %p and %p, expected singleton", a, b)
	}
}
