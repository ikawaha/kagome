package trie

import (
	"encoding/binary"
	"fmt"
	"io"
	"sort"
)

const (
	initBufferSize = 51200
	expandRatio    = 2
	terminator     = '\x00'
	rootID         = 0
)

// DoubleArray represents the TRIE data structure.
type DoubleArray []struct {
	Base, Check int32
}

// Build constructs a double array from given keywords.
func Build(keywords []string) (DoubleArray, error) {
	s := len(keywords)
	if s == 0 {
		return DoubleArray{}, nil
	}
	ids := make([]int, s)
	for i := range ids {
		ids[i] = i + 1
	}
	return BuildWithIDs(keywords, ids)
}

// BuildWithIDs constructs a double array from given keywords and ids.
func BuildWithIDs(keywords []string, ids []int) (DoubleArray, error) {
	d := DoubleArray{}
	d.init()
	if len(keywords) != len(ids) {
		return d, fmt.Errorf("invalid arguments")
	}
	if len(keywords) == 0 {
		return d, nil
	}
	if !sort.StringsAreSorted(keywords) {
		h := make(map[string]int)
		for i, key := range keywords {
			h[key] = ids[i]
		}
		sort.Strings(keywords)
		ids = ids[:0]
		for _, key := range keywords {
			ids = append(ids, h[key])
		}
	}
	branches := make([]int, len(keywords))
	for i := range keywords {
		branches[i] = i
	}
	d.add(0, 0, branches, keywords, ids)
	d.truncate()
	return d, nil
}

// Find searches TRIE by a given keyword and returns the id if found.
func (d DoubleArray) Find(input string) (id int, ok bool) {
	_, q, _, ok := d.search(input)
	if !ok {
		return
	}
	p := q
	q = int(d[p].Base) + int(terminator)
	if q >= len(d) || int(d[q].Check) != p || d[q].Base > 0 {
		return
	}
	return int(-d[q].Base), true
}

// CommonPrefixSearch finds keywords sharing common prefix in an input
// and returns the ids and it's lengths if found.
func (d DoubleArray) CommonPrefixSearch(input string) (ids, lens []int) {
	var p, q int
	bufLen := len(d)
	for i, size := 0, len(input); i < size; i++ {
		if input[i] == terminator {
			return
		}
		p = q
		q = int(d[p].Base) + int(input[i])
		if q >= bufLen || int(d[q].Check) != p {
			break
		}
		ahead := int(d[q].Base) + int(terminator)
		if ahead < bufLen && int(d[ahead].Check) == q && int(d[ahead].Base) <= 0 {
			ids = append(ids, int(-d[ahead].Base))
			lens = append(lens, i+1)
		}
	}
	return
}

// CommonPrefixSearchCallback finds keywords sharing common prefix in an input
// and callback with id and length.
func (d DoubleArray) CommonPrefixSearchCallback(input string, callback func(id, l int)) {
	var p, q int
	bufLen := len(d)
	for i := 0; i < len(input); i++ {
		if input[i] == terminator {
			return
		}
		p = q
		q = int(d[p].Base) + int(input[i])
		if q >= bufLen || int(d[q].Check) != p {
			break
		}
		ahead := int(d[q].Base) + int(terminator)
		if ahead < bufLen && int(d[ahead].Check) == q && int(d[ahead].Base) <= 0 {
			callback(int(-d[ahead].Base), i+1)
		}
	}
}

// PrefixSearch returns the longest common prefix keyword in an input if found.
func (d DoubleArray) PrefixSearch(input string) (id int, ok bool) {
	var p, q, i int
	bufLen := len(d)
	for size := len(input); i < size; i++ {
		if input[i] == terminator {
			return
		}
		p = q
		q = int(d[p].Base) + int(input[i])
		if q >= bufLen || int(d[q].Check) != p {
			break
		}
		ahead := int(d[q].Base) + int(terminator)
		if ahead < bufLen && int(d[ahead].Check) == q && int(d[ahead].Base) <= 0 {
			id = int(-d[ahead].Base)
			ok = true
		}
	}
	return
}

// WriteTo saves a double array.
func (d DoubleArray) WriteTo(w io.Writer) (n int64, err error) {
	sz := int64(len(d))
	//fmt.Println("write data len:", sz)
	if err := binary.Write(w, binary.LittleEndian, sz); err != nil {
		return n, err
	}
	n += int64(binary.Size(sz))
	for _, v := range d {
		if err := binary.Write(w, binary.LittleEndian, v.Base); err != nil {
			return n, err
		}
		n += int64(binary.Size(v.Base))
		if err := binary.Write(w, binary.LittleEndian, v.Check); err != nil {
			return n, err
		}
		n += int64(binary.Size(v.Check))
	}
	return n, nil
}

// Read loads a double array.
func Read(r io.Reader) (DoubleArray, error) {
	var sz int64
	if err := binary.Read(r, binary.LittleEndian, &sz); err != nil {
		return DoubleArray{}, err
	}
	//fmt.Println("read data len:", sz)
	d := make(DoubleArray, sz)
	for i := range d {
		if err := binary.Read(r, binary.LittleEndian, &d[i].Base); err != nil {
			return d, err
		}
		if err := binary.Read(r, binary.LittleEndian, &d[i].Check); err != nil {
			return d, err
		}
	}
	return d, nil
}

func (d *DoubleArray) init() {
	*d = make(DoubleArray, initBufferSize)

	(*d)[rootID].Base = 1
	(*d)[rootID].Check = -1

	bufLen := len(*d)
	for i := 1; i < bufLen; i++ {
		(*d)[i].Base = int32(-(i - 1))
		(*d)[i].Check = int32(-(i + 1))
	}

	(*d)[1].Base = int32(-(bufLen - 1))
	(*d)[bufLen-1].Check = int32(-1)
}

func (d *DoubleArray) setBase(p, base int) {
	if p == rootID {
		return
	}
	if (*d)[p].Check < 0 {
		if (*d)[p].Base == (*d)[p].Check {
			d.expand()
		}
		prev := -(*d)[p].Base
		next := -(*d)[p].Check
		if -p == int((*d)[rootID].Check) {
			(*d)[rootID].Check = (*d)[p].Check
		}
		(*d)[next].Base = (*d)[p].Base
		(*d)[prev].Check = (*d)[p].Check
	}
	(*d)[p].Base = int32(base)
}

func (d *DoubleArray) efficiency() (unspent int, size int, usageRate float64) {
	for _, pair := range *d {
		if pair.Check < 0 {
			unspent++
		}
	}
	return unspent, len(*d), float64(len(*d)-unspent) / float64(len(*d)) * 100
}

func (d *DoubleArray) expand() {
	srcSize := len(*d)
	dst := new(DoubleArray)
	dstSize := srcSize * expandRatio
	*dst = make(DoubleArray, dstSize)
	copy(*dst, *d)

	for i := srcSize; i < dstSize; i++ {
		(*dst)[i].Base = int32(-(i - 1))
		(*dst)[i].Check = int32(-(i + 1))
	}

	start := -(*d)[0].Check
	end := -(*dst)[start].Base
	(*dst)[srcSize].Base = -end
	(*dst)[start].Base = int32(-(dstSize - 1))
	(*dst)[end].Check = int32(-srcSize)
	(*dst)[dstSize-1].Check = -start

	*d = *dst
}

func (d *DoubleArray) truncate() {
	srcSize := len(*d)
	for i, size := 0, srcSize; i < size; i++ {
		if (*d)[size-i-1].Check < 0 {
			srcSize--
		} else {
			break
		}
	}
	if srcSize == len(*d) {
		return
	}
	dst := new(DoubleArray)
	*dst = make(DoubleArray, srcSize)
	copy(*dst, (*d)[:srcSize])
	*d = *dst
}

func (d *DoubleArray) search(input string) (p, q, i int, ok bool) {
	if len(input) == 0 {
		return
	}
	bufLen := len(*d)
	inpLen := len(input)
	for i = 0; i < inpLen; i++ {
		if input[i] == terminator {
			return
		}
		p = q
		q = int((*d)[p].Base) + int(input[i])
		if q >= bufLen || int((*d)[q].Check) != p {
			return
		}
	}
	return p, q, i, true
}

func (d *DoubleArray) setCheck(p, check int) {
	if (*d)[p].Base == (*d)[p].Check {
		d.expand()
	}
	prev := -(*d)[p].Base
	next := -(*d)[p].Check
	if -p == int((*d)[rootID].Check) {
		(*d)[rootID].Check = (*d)[p].Check
	}

	(*d)[next].Base = (*d)[p].Base
	(*d)[prev].Check = (*d)[p].Check
	(*d)[p].Check = int32(check)

}

func (d *DoubleArray) seekAndMark(p int, chars []byte) { // chars != nil
	free := rootID
	rep := int(chars[0])
	var base int
loop:
	for {
		if free != rootID && (*d)[free].Check == (*d)[rootID].Check {
			d.expand()
		}
		free = int(-(*d)[free].Check)
		base = free - rep
		if base <= 0 {
			continue
		}
		for _, ch := range chars {
			q := base + int(ch)
			if q < len(*d) && (*d)[q].Check >= 0 {
				goto loop
			}
		}
		break
	}
	d.setBase(p, base)
	for _, ch := range chars {
		q := int((*d)[p].Base) + int(ch)
		if q >= len(*d) {
			d.expand()
		}
		d.setCheck(q, p)
	}
}

func (d *DoubleArray) add(p, i int, branches []int, keywords []string, ids []int) {
	var chars []byte
	subtree := make(map[byte][]int)
	for _, keyID := range branches {
		str := []byte(keywords[keyID])
		var ch byte
		if i >= len(str) {
			ch = terminator
		} else {
			ch = str[i]
		}
		if size := len(chars); size == 0 || chars[len(chars)-1] != ch {
			chars = append(chars, ch)
		}
		if ch != terminator {
			subtree[ch] = append(subtree[ch], keyID)
		}
	}
	d.seekAndMark(p, chars)
	for _, ch := range chars {
		q := int((*d)[p].Base) + int(ch)
		if len(subtree[ch]) == 0 {
			if len(ids) == 0 {
				(*d)[q].Base = int32(-branches[0])
			} else {
				(*d)[q].Base = int32(-ids[branches[0]])
			}
		} else {
			d.add(q, i+1, subtree[ch], keywords, ids)
		}
	}
}
