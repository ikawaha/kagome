package dict

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
)

// CharClass represents a character class.
type CharClass []string

// CharCategory represents categories for characters.
type CharCategory []byte

// InvokeList represents whether to invoke unknown word processing.
type InvokeList []bool

// GroupList represents whether to make a new word by grouping the same character category.
type GroupList []bool

// CharDef represents char.def.
type CharDef struct {
	CharClass    CharClass
	CharCategory CharCategory
	InvokeList   InvokeList
	GroupList    GroupList
}

// ReadCharDef reads char.def format.
func ReadCharDef(r io.Reader) (*CharDef, error) {
	var ret CharDef
	dec := gob.NewDecoder(r)
	if err := dec.Decode(&ret); err != nil {
		return nil, fmt.Errorf("char class read error, %v", err)
	}
	return &ret, nil
}

// WriteTo implements the io.WriteTo interface.
func (d CharDef) WriteTo(w io.Writer) (n int64, err error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(d); err != nil {
		return 0, err
	}
	return b.WriteTo(w)
}
