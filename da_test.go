package kagome

import (
	"math"
	"testing"
)

func TestDaBuildAndSearch01(t *testing.T) {
	d := &DoubleArray{}
	if e := d.Build(nil); e != nil {
		t.Fatalf("unexpected error: %v", e)
	}

	if _, ok := d.FindString(""); ok {
		t.Errorf("unexpected result: %v", ok)
	}
}

func TestDaBuildAndSearch02(t *testing.T) {
	d := &DoubleArray{}
	if e := d.Build([]string{}); e != nil {
		t.Fatalf("unexpected error: %v", e)
	}

	if _, ok := d.FindString(""); ok {
		t.Errorf("unexpected result: %v", ok)
	}
}

func TestDaBuildAndSearch03(t *testing.T) {

	keywords := []string{
		"12345",
		"2345",
		"１２３",
		"abc",
		"ABCD",
		"あいう",
		"Ａ",
	}

	d := &DoubleArray{}
	if e := d.Build(keywords); e != nil {
		t.Fatalf("unexpected error: %v", e)
	}

	for _, k := range keywords {
		if _, ok := d.FindString(k); !ok {
			t.Errorf("does not detected: %v\n", k)
		}
	}
}

func TestDaBuildAndCommonPrefixSearch01(t *testing.T) {
	d := &DoubleArray{}
	if e := d.Build(nil); e != nil {
		t.Fatalf("unexpected error: %v", e)
	}

	if ids, lens := d.CommonPrefixSearchString(""); len(ids) != 0 || len(lens) != 0 {
		t.Errorf("unexpected result: %v", ids)
	}
}

func TestDaBuildAndCommonPrefixSearch02(t *testing.T) {
	d := &DoubleArray{}
	if e := d.Build([]string{}); e != nil {
		t.Fatalf("unexpected error: %v", e)
	}

	if ids, lens := d.CommonPrefixSearchString(""); len(ids) != 0 || len(lens) != 0 {
		t.Errorf("unexpected result: %v", ids)
	}
}

func TestDaBuildAndCommonPrefixSearch03(t *testing.T) {
	keywords := []string{
		"電気通信",              //1
		"電気",                //2
		"電気通信大学",            //3
		"電気通信大学院大学",         //4
		"電気通信大学大学院",         //5
		"電気通信大学大学院電気通信学研究科", //6
		"電気通信大学電気通信学部",      //7
	}

	d := &DoubleArray{}
	d.Build(keywords)

	expectedIds := []int{
		2, //"電気",
		1, //"電気通信",
		3, //"電気通信大学",
		5, //"電気通信大学大学院",
		6, //"電気通信大学大学院電気通信学研究科",
	}
	// byte length
	expectedLens := []int{
		6,  //"電気", 2*3
		12, //"電気通信", 4*3
		18, //"電気通信大学", 6*3
		27, //"電気通信大学大学院",9*3
		51, //"電気通信大学大学院電気通信学研究科",17*3
	}

	ids, lens := d.CommonPrefixSearchString("電気通信大学大学院電気通信学研究科")
	if len(ids) != len(expectedIds) || len(lens) != len(expectedLens) {
		t.Fatalf("got %v, expected %v\n", ids, expectedIds)
	}
	for i := range expectedIds {
		if ids[i] != expectedIds[i] {
			t.Fatalf("id: got %v, expected %v\n", ids, expectedIds)
		}
		if lens[i] != expectedLens[i] {
			t.Fatalf("len: got %v, expected %v\n", lens, expectedLens)
		}
	}
}

func TestDaBuildWithIdsAndCommonPrefixSearch01(t *testing.T) {
	d := &DoubleArray{}
	if e := d.BuildWithIds(nil, nil); e != nil {
		t.Fatalf("unexpected error: %v", e)
	}

	if ids, lens := d.CommonPrefixSearchString(""); len(ids) != 0 || len(lens) != 0 {
		t.Errorf("unexpected result: %v, %v", ids, lens)
	}
}

func TestDaBuildWithIdsAndCommonPrefixSearch02(t *testing.T) {
	d := &DoubleArray{}
	if e := d.BuildWithIds([]string{}, []int{}); e != nil {
		t.Fatalf("unexpected error: %v", e)
	}

	if ids, lens := d.CommonPrefixSearchString(""); len(ids) != 0 || len(lens) != 0 {
		t.Errorf("unexpected result: %v, %v", ids, lens)
	}
}

func TestDaBuildWithIdsAndPrefixSearch03(t *testing.T) {
	keywords := []string{
		"電気通信大学電気通信学部",
		"電気",
		"電気通信",
		"電気通信大学",
		"電気通信大学院大学",
		"電気通信大学大学院",
		"電気通信大学大学院電気通信学研究科",
	}

	ids := []int{1, 2, 3, 4, 5, 6, 7, 8}

	h := make(map[string]int)
	for i := range keywords {
		h[keywords[i]] = ids[i]
	}

	d := &DoubleArray{}
	if e := d.BuildWithIds(keywords, ids); e == nil {
		t.Errorf("expected error: invalid argument error\n")
	}

	ids = ids[0 : len(ids)-1]
	if e := d.BuildWithIds(keywords, ids); e != nil {
		t.Errorf("unexpected error: %v\n", e)
	}
	for key, expectedId := range h {
		if id, ok := d.FindString(key); !ok || id != expectedId {
			t.Errorf("got ok:%v, id:%v, expected ok:true, id:%v (keyword:%v)", ok, id, expectedId, key)
		}
	}
}

func TestDaBuildAndPrefixSearch01(t *testing.T) {
	d := &DoubleArray{}
	if e := d.Build(nil); e != nil {
		t.Fatalf("unexpected error: %v", e)
	}

	if id, ok := d.PrefixSearchString(""); ok {
		t.Errorf("unexpected result: %v", id)
	}
}

func TestDaBuildAndPrefixSearch02(t *testing.T) {
	d := &DoubleArray{}
	if e := d.Build([]string{}); e != nil {
		t.Fatalf("unexpected error: %v", e)
	}

	if id, ok := d.PrefixSearchString(""); ok {
		t.Errorf("unexpected result: %v", id)
	}
}

func TestDaBuildAndPrefixSearch03(t *testing.T) {
	keywords := []string{
		"電気",                //1
		"電気通信",              //2
		"電気通信大学",            //3
		"電気通信大学院大学",         //4
		"電気通信大学大学院",         //5
		"電気通信大学電気通信学部",      //6
		"電気通信大学大学院電気通信学研究科", //7
	}

	d := &DoubleArray{}
	d.Build(keywords)

	expected := 7 //電気通信大学大学院電気通信学研究科
	id, ok := d.PrefixSearchString("電気通信大学大学院電気通信学研究科")
	if !ok {
		t.Fatalf("cannot search the prefix\n", id, expected)
	}
	if id != expected {
		t.Fatalf("got %v, expected %v\n", id, expected)
	}
}

func TestDaEfficiency01(t *testing.T) {
	d := &DoubleArray{}
	d.init()
	unspent, size, rate := d.efficiency()
	if unspent != size || size != daInitBufferSize {
		t.Errorf("got unspent:%v, size:%v, expected both %v\n", unspent, size, daInitBufferSize)
	}
	if rate != 0.0 {
		t.Errorf("got :%v, expected 0.0\n", rate)
	}

	d.truncate()
	unspent, size, rate = d.efficiency()
	if unspent != size || size != 0 {
		t.Errorf("got unspent:%v, size:%v, expected 0\n", unspent, size)
	}
	if !math.IsNaN(rate) {
		t.Errorf("got :%v, expected NaN\n", rate)
	}
}
