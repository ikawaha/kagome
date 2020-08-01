package dict

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

const (
	rowDelimiter = "\n"
	colDelimiter = "\a"
)

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
