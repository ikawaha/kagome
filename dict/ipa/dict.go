package ipa

import (
	"archive/zip"
	"bytes"
	"fmt"
	"sort"
	"sync"

	data "github.com/ikawaha/kagome-dict-ipa"
	"github.com/ikawaha/kagome/v2/dict"
)

const dictPath = "dict"

type systemDict struct {
	once sync.Once
	dict *dict.Dict
}

var (
	full    systemDict
	shurink systemDict
)

func New() *dict.Dict {
	full.once.Do(func() {
		full.dict = loadInternalDict(dictPath, true)
		shurink.once.Do(func() {
			shurink.dict = full.dict
		})
	})
	return full.dict
}

func NewShrink() *dict.Dict {
	shurink.once.Do(func() {
		shurink.dict = loadInternalDict(dictPath, false)
	})
	return shurink.dict
}

func loadInternalDict(path string, full bool) (d *dict.Dict) {
	buf := make([]byte, 0, 36*1024*1024) // 36MB
	defer func() { buf = nil }()

	pieces := data.AssetNames()
	sort.Strings(pieces)
	for _, v := range pieces {
		b, err := data.Asset(v)
		if err != nil {
			panic(fmt.Errorf("assert error, %q, %v", v, err))
		}
		buf = append(buf, b...)
	}
	r := bytes.NewReader(buf)
	zr, err := zip.NewReader(r, r.Size())
	if err != nil {
		panic(err)
	}
	d, err = dict.Load(zr, full)
	if err != nil {
		panic(err)
	}
	return d
}
