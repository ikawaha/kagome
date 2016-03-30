// Copyright 2015 ikawaha
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package da

import (
	"bufio"
	"bytes"
	"math"
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestDaBuildAndSearch01(t *testing.T) {
	d, err := Build(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := d.Find(""); ok {
		t.Errorf("unexpected result: %v", ok)
	}
}

func TestDaBuildAndSearch02(t *testing.T) {
	d, err := Build([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := d.Find(""); ok {
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
	d, err := Build(keywords)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, k := range keywords {
		if _, ok := d.Find(k); !ok {
			t.Errorf("does not detected: %v\n", k)
		}
	}
}

func TestDaBuildAndCommonPrefixSearch01(t *testing.T) {
	d, err := Build(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ids, lens := d.CommonPrefixSearch(""); len(ids) != 0 || len(lens) != 0 {
		t.Errorf("unexpected result: %v", ids)
	}
}

func TestDaBuildAndCommonPrefixSearch02(t *testing.T) {
	d, err := Build([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ids, lens := d.CommonPrefixSearch(""); len(ids) != 0 || len(lens) != 0 {
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
	d, err := Build(keywords)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedIDs := []int{
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

	ids, lens := d.CommonPrefixSearch("電気通信大学大学院電気通信学研究科")
	if len(ids) != len(expectedIDs) || len(lens) != len(expectedLens) {
		t.Fatalf("got %v, expected %v\n", ids, expectedIDs)
	}
	for i := range expectedIDs {
		if ids[i] != expectedIDs[i] {
			t.Fatalf("id: got %v, expected %v\n", ids, expectedIDs)
		}
		if lens[i] != expectedLens[i] {
			t.Fatalf("len: got %v, expected %v\n", lens, expectedLens)
		}
	}
}

func TestDaBuildAndCommonPrefixSearchCallback01(t *testing.T) {
	d, err := Build(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	d.CommonPrefixSearchCallback("", func(id, l int) {
		t.Errorf("unexpected callback, id:%v, l:%v", id, l)
	})
}

func TestDaBuildAndCommonPrefixSearchCallback02(t *testing.T) {
	d, err := Build([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	d.CommonPrefixSearchCallback("", func(id, l int) {
		t.Errorf("unexpected callback, id:%v, l:%v", id, l)
	})
}

func TestDaBuildAndCommonPrefixSearchCallback03(t *testing.T) {
	keywords := []string{
		"電気通信",              //1
		"電気",                //2
		"電気通信大学",            //3
		"電気通信大学院大学",         //4
		"電気通信大学大学院",         //5
		"電気通信大学大学院電気通信学研究科", //6
		"電気通信大学電気通信学部",      //7
	}
	d, err := Build(keywords)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []struct {
		id, l int
	}{
		{2, 3 * 2},  //"電気"
		{1, 3 * 4},  //"電気通信"
		{3, 3 * 6},  //"電気通信大学"
		{5, 3 * 9},  //"電気通信大学大学院"
		{6, 3 * 17}, //"電気通信大学大学院電気通信学研究科"
	}

	var i int
	d.CommonPrefixSearchCallback("電気通信大学大学院電気通信学研究科", func(id, l int) {
		if expected[i].id != id {
			t.Errorf("id: got %v, expected %v", id, expected[i].id)
		}
		if expected[i].l != l {
			t.Errorf("len: got %v, expected %v", l, expected[i].l)
		}
		i++
	})
}

func TestDaBuildWithIDsAndCommonPrefixSearch01(t *testing.T) {
	d, err := BuildWithIDs(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ids, lens := d.CommonPrefixSearch(""); len(ids) != 0 || len(lens) != 0 {
		t.Errorf("unexpected result: %v, %v", ids, lens)
	}
}

func TestDaBuildWithIDsAndCommonPrefixSearch02(t *testing.T) {
	d, err := BuildWithIDs([]string{}, []int{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ids, lens := d.CommonPrefixSearch(""); len(ids) != 0 || len(lens) != 0 {
		t.Errorf("unexpected result: %v, %v", ids, lens)
	}
}

func TestDaBuildWithIDsAndPrefixSearch03(t *testing.T) {
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

	d, err := BuildWithIDs(keywords, ids)
	if err == nil {
		t.Errorf("invalid argument error was expected")
	}

	ids = ids[0 : len(ids)-1]
	d, err = BuildWithIDs(keywords, ids)
	if err != nil {
		t.Errorf("unexpected error: %v\n", err)
	}
	for key, expectedID := range h {
		if id, ok := d.Find(key); !ok || id != expectedID {
			t.Errorf("got ok:%v, id:%v, expected ok:true, id:%v (keyword:%v)", ok, id, expectedID, key)
		}
	}
}

func TestDaBuildAndPrefixSearch01(t *testing.T) {
	d, err := Build(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id, ok := d.PrefixSearch(""); ok {
		t.Errorf("unexpected result: %v", id)
	}
}

func TestDaBuildAndPrefixSearch02(t *testing.T) {
	d, err := Build([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id, ok := d.PrefixSearch(""); ok {
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

	d, err := Build(keywords)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := 7 //電気通信大学大学院電気通信学研究科
	id, ok := d.PrefixSearch("電気通信大学大学院電気通信学研究科")
	if !ok {
		t.Fatalf("cannot search the prefix. id=%v, %v\n", id, expected)
	}
	if id != expected {
		t.Fatalf("got %v, expected %v\n", id, expected)
	}
}

func TestDaEfficiency01(t *testing.T) {
	d := &DoubleArray{}
	d.init()
	unspent, size, rate := d.efficiency()
	if unspent != size || size != initBufferSize {
		t.Errorf("got unspent:%v, size:%v, expected both %v", unspent, size, initBufferSize)
	}
	if rate != 0.0 {
		t.Errorf("got :%v, expected 0.0", rate)
	}

	d.truncate()
	unspent, size, rate = d.efficiency()
	if unspent != size || size != 0 {
		t.Errorf("got unspent:%v, size:%v, expected 0", unspent, size)
	}
	if !math.IsNaN(rate) {
		t.Errorf("got :%v, expected NaN", rate)
	}
}

func TestDaTorture(t *testing.T) {
	const testdata = "_testdata/words.txt"
	fp, err := os.Open(testdata)
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	defer fp.Close()

	var keys []string
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		keys = append(keys, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	sort.Strings(keys)

	d, err := Build(keys)
	if err != nil {
		t.Fatalf("unexpected build error, %v", err)
	}
	for _, k := range keys {
		id, ok := d.Find(k)
		if !ok || id < 1 {
			t.Errorf("input: %v, detect: %v, id: %v", k, ok, id)
			continue
		}
		if k != keys[id-1] {
			t.Errorf("got %v, expected %v", keys[id-1], k)
		}
	}
}

func TestReadAndWrite(t *testing.T) {
	keywords := []string{
		"電気",                //1
		"電気通信",              //2
		"電気通信大学",            //3
		"電気通信大学院大学",         //4
		"電気通信大学大学院",         //5
		"電気通信大学電気通信学部",      //6
		"電気通信大学大学院電気通信学研究科", //7
	}

	org, err := Build(keywords)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var b bytes.Buffer
	n, err := org.WriteTo(&b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != int64(b.Len()) {
		t.Errorf("write len: got %v, expected %v", n, b.Len())
	}

	cpy, err := Read(&b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(org, cpy) {
		t.Errorf("got %v, expected %v", cpy, org)
	}
}

func TestDaExpand(t *testing.T) {
	const initSize = 5
	d := make(DoubleArray, initSize)
	d[rootID].Base = 1
	d[rootID].Check = -1

	for i := 1; i < len(d); i++ {
		d[i].Base = int32(-(i - 1))
		d[i].Check = int32(-(i + 1))
	}
	d[1].Base = int32(-(len(d) - 1))
	d[len(d)-1].Check = int32(-1)

	d.expand()

	if len(d) != initSize*expandRatio {
		t.Errorf("expand error: got %v, expected %v", len(d), initSize*expandRatio)
	}
}
