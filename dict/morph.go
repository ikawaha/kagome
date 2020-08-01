package dict

import (
	"encoding/binary"
	"io"
)

// Morph represents part of speeches and an occurrence cost.
type Morph struct {
	LeftID, RightID, Weight int16
}

// Morphs represents a slice of morphs.
type Morphs []Morph

// WriteTo implements the io.WriterTo interface.
func (m Morphs) WriteTo(w io.Writer) (n int64, err error) {
	l := int64(len(m))
	if err := binary.Write(w, binary.LittleEndian, l); err != nil {
		return n, err
	}
	n += int64(binary.Size(l))
	for i := range m {
		if err := binary.Write(w, binary.LittleEndian, m[i].LeftID); err != nil {
			return n, err
		}
		n += int64(binary.Size(m[i].LeftID))
		if err := binary.Write(w, binary.LittleEndian, m[i].RightID); err != nil {
			return n, err
		}
		n += int64(binary.Size(m[i].RightID))
		if err := binary.Write(w, binary.LittleEndian, m[i].Weight); err != nil {
			return n, err
		}
		n += int64(binary.Size(m[i].Weight))
	}
	return n, nil
}

// ReadMorphs loads morph data from io.Reader.
func ReadMorphs(r io.Reader) (Morphs, error) {
	var l int64
	if err := binary.Read(r, binary.LittleEndian, &l); err != nil {
		return nil, err
	}
	m := make([]Morph, l)
	for i := range m {
		if err := binary.Read(r, binary.LittleEndian, &m[i].LeftID); err != nil {
			return m, err
		}
		if err := binary.Read(r, binary.LittleEndian, &m[i].RightID); err != nil {
			return m, err
		}
		if err := binary.Read(r, binary.LittleEndian, &m[i].Weight); err != nil {
			return m, err
		}
	}
	return m, nil
}
