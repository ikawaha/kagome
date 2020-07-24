package dict

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
)

type CharClass []string
type CharCategory []byte
type InvokeList []bool
type GroupList []bool

type CharDef struct {
	CharClass    CharClass
	CharCategory CharCategory
	InvokeList   InvokeList
	GroupList    GroupList
}

func ReadCharDef(r io.Reader) (*CharDef, error) {
	var ret CharDef
	dec := gob.NewDecoder(r)
	if err := dec.Decode(&ret); err != nil {
		return nil, fmt.Errorf("char class read error, %v", err)
	}
	return &ret, nil
}

func (d CharDef) WriteTo(w io.Writer) (n int64, err error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(d); err != nil {
		return 0, err
	}
	return b.WriteTo(w)
}
