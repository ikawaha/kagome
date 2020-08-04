package dict

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"strings"
)

const (
	rowDelimiter = "\n"
	colDelimiter = "\a"
)

const (
	POSStartIndex      = "_pos_start"
	POSHierarchy       = "_pos_hierarchy"
	InflectionalType   = "_inflectional_type"
	InflectionalForm   = "_inflectional_form"
	BaseFormIndex      = "_base"
	ReadingIndex       = "_reading"
	PronunciationIndex = "_pronunciation"
)

// ContentsMeta represents the contents record information.
type ContentsMeta map[string]int8

func (c ContentsMeta) WriteTo(w io.Writer) (n int64, err error) {
	return writeMapStringInt8(w, c)
}

func writeMapStringInt8(w io.Writer, m map[string]int8) (n int64, err error) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	sz := int64(len(keys))
	if err := binary.Write(w, binary.LittleEndian, sz); err != nil {
		return n, err
	}
	n += int64(binary.Size(sz))
	for _, k := range keys {
		sz := int64(len(k))
		if err := binary.Write(w, binary.LittleEndian, sz); err != nil {
			return n, err
		}
		n += int64(binary.Size(sz))
		x, err := io.Copy(w, strings.NewReader(k))
		if err != nil {
			return n, err
		}
		n += x
		v := m[k]
		if err := binary.Write(w, binary.LittleEndian, v); err != nil {
			return n, err
		}
		n += int64(binary.Size(v))
	}
	return n, err
}

func ReadContentsMeta(r io.Reader) (ContentsMeta, error) {
	return readMapStringInt8(r)
}

func readMapStringInt8(r io.Reader) (map[string]int8, error) {
	var sz int64
	if err := binary.Read(r, binary.LittleEndian, &sz); err != nil {
		return nil, err
	}
	m := make(map[string]int8, sz)
	for i := int64(0); i < sz; i++ {
		var x int64
		if err := binary.Read(r, binary.LittleEndian, &x); err != nil {
			return nil, err
		}
		var buf strings.Builder
		_, err := io.CopyN(&buf, r, x)
		if err != nil {
			return nil, err
		}
		var v int8
		if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
			return nil, err
		}
		m[buf.String()] = v
	}
	return m, nil
}

// Contents represents dictionary contents.
type Contents [][]string

// WriteTo implements the io.WriterTo interface.
func (c Contents) WriteTo(w io.Writer) (n int64, err error) {
	for i := 0; i < len(c)-1; i++ {
		x, err := fmt.Fprintf(w, "%s%s", strings.Join(c[i], colDelimiter), rowDelimiter)
		if err != nil {
			return n, err
		}
		n += int64(x)
	}
	if i := len(c) - 1; i > 0 {
		x, err := fmt.Fprintf(w, "%s", strings.Join(c[i], colDelimiter))
		if err != nil {
			return n, err
		}
		n += int64(x)
	}
	return n, nil
}

// NewContents creates dictionary contents from byte slice.
func NewContents(b []byte) [][]string {
	str := string(b)
	rows := strings.Split(str, rowDelimiter)
	m := make([][]string, len(rows))
	for i, r := range rows {
		m[i] = strings.Split(r, colDelimiter)
	}
	return m
}

// ReadContents reads dictionary contents from io.Reader.
func ReadContents(r io.Reader) (Contents, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read contents error, %v", err)
	}
	return NewContents(buf), nil
}
