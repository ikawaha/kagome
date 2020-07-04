package dic

import (
	"bytes"
	"encoding/gob"
	"io"
)

type CharClass []string

func (c CharClass) WriteTo(w io.Writer) (n int64, err error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(c); err != nil {
		return 0, err
	}
	return b.WriteTo(w)
}

type CharCategory []byte

func (c CharCategory) WriteTo(w io.Writer) (n int64, err error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(c); err != nil {
		return 0, err
	}
	return b.WriteTo(w)
}

type InvokeList []bool

func (l InvokeList) WriteTo(w io.Writer) (n int64, err error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(l); err != nil {
		return 0, err
	}
	return b.WriteTo(w)
}

type GroupList []bool

func (l GroupList) WriteTo(w io.Writer) (n int64, err error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(l); err != nil {
		return 0, err
	}
	return b.WriteTo(w)
}
