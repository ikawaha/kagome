package kagome

import (
	"fmt"
	"sort"
)

const (
	daInitBufferSize = 51200
	daExpandRatio    = 2
	daTerminator     = '\x00'
	daRootId         = 0
)

// DoubleArray represents the TRIE data structure.
type DoubleArray []struct {
	Base, Check int
}

func (d *DoubleArray) init() {
	*d = make(DoubleArray, daInitBufferSize)

	(*d)[daRootId].Base = 1
	(*d)[daRootId].Check = -1

	bufLen := len(*d)
	for i := 1; i < bufLen; i++ {
		(*d)[i].Base = -(i - 1)
		(*d)[i].Check = -(i + 1)
	}

	(*d)[1].Base = -(bufLen - 1)
	(*d)[bufLen-1].Check = -1
}

func (d *DoubleArray) setBase(p, aBase int) {
	if p == daRootId {
		return
	}
	if (*d)[p].Check < 0 {
		if (*d)[p].Base == (*d)[p].Check {
			d.expand()
		}
		prev := -(*d)[p].Base
		next := -(*d)[p].Check
		if -p == (*d)[daRootId].Check {
			(*d)[daRootId].Check = (*d)[p].Check
		}
		(*d)[next].Base = (*d)[p].Base
		(*d)[prev].Check = (*d)[p].Check
	}
	(*d)[p].Base = aBase
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
	dstSize := srcSize * daExpandRatio
	*dst = make(DoubleArray, dstSize)
	copy(*dst, *d)

	for i := srcSize; i < dstSize; i++ {
		(*dst)[i].Base = -(i - 1)
		(*dst)[i].Check = -(i + 1)
	}

	start := -(*d)[0].Check
	end := -(*dst)[start].Base
	(*dst)[srcSize].Base = -end
	(*dst)[start].Base = -(dstSize - 1)
	(*dst)[end].Check = -srcSize
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

func (d *DoubleArray) searchBytes(input []byte) (p, q, i int, ok bool) {
	if len(input) == 0 {
		return
	}
	bufLen := len(*d)
	for i = range input {
		p = q
		q = (*d)[p].Base + int(input[i])
		if q >= bufLen || (*d)[q].Check != p {
			return
		}
	}
	return p, q, i, true
}

func (d *DoubleArray) searchString(input string) (p, q, i int, ok bool) {
	if len(input) == 0 {
		return
	}
	bufLen := len(*d)
	inpLen := len(input)
	for i = 0; i < inpLen; i++ {
		p = q
		q = (*d)[p].Base + int(input[i])
		if q >= bufLen || (*d)[q].Check != p {
			return
		}
	}
	return p, q, i, true
}

func (d *DoubleArray) setCheck(p, aCheck int) {
	if (*d)[p].Base == (*d)[p].Check {
		d.expand()
	}
	prev := -(*d)[p].Base
	next := -(*d)[p].Check
	if -p == (*d)[daRootId].Check {
		(*d)[daRootId].Check = (*d)[p].Check
	}

	(*d)[next].Base = (*d)[p].Base
	(*d)[prev].Check = (*d)[p].Check
	(*d)[p].Check = aCheck

}

func (d *DoubleArray) seekAndMark(p int, aChars []byte) { // chars != nil
	free := daRootId
	rep := int(aChars[0])
	var base int
	for {
	L_start:
		if free != daRootId && (*d)[free].Check == (*d)[daRootId].Check {
			d.expand()
		}
		free = -(*d)[free].Check
		base = free - rep
		if base <= 0 {
			continue
		}
		for _, ch := range aChars {
			q := base + int(ch)
			if q < len(*d) && (*d)[q].Check >= 0 {
				goto L_start
			}
		}
		break
	}
	d.setBase(p, base)
	for _, ch := range aChars {
		q := (*d)[p].Base + int(ch)
		if q >= len(*d) {
			d.expand()
		}
		d.setCheck(q, p)
	}
}

func (d *DoubleArray) add(p, i int, branches []int, keywords []string, ids []int) {
	var chars []byte
	subtree := make(map[byte][]int)
	for _, keyId := range branches {
		str := []byte(keywords[keyId])
		var ch byte
		if i >= len(str) {
			ch = daTerminator
		} else {
			ch = str[i]
		}
		if size := len(chars); size == 0 || chars[len(chars)-1] != ch {
			chars = append(chars, ch)
		}
		if ch != daTerminator {
			subtree[ch] = append(subtree[ch], keyId)
		}
	}
	d.seekAndMark(p, chars)
	for _, ch := range chars {
		q := (*d)[p].Base + int(ch)
		if len(subtree[ch]) == 0 {
			if len(ids) == 0 {
				(*d)[q].Base = -branches[0]
			} else {
				(*d)[q].Base = -ids[branches[0]]
			}
		} else {
			d.add(q, i+1, subtree[ch], keywords, ids)
		}
	}
}

// Build constructs a double array from given keywords.
func (d *DoubleArray) Build(keywords []string) (err error) {
	d.init()
	s := len(keywords)
	if s == 0 {
		return
	}
	ids := make([]int, s, s)
	for i := range ids {
		ids[i] = i + 1
	}
	return d.BuildWithIds(keywords, ids)
}

// BuildWithIds constructs a double array from given keywords and ids.
func (d *DoubleArray) BuildWithIds(keywords []string, ids []int) (err error) {
	d.init()
	if len(keywords) != len(ids) {
		err = fmt.Errorf("invalid arguments")
		return
	}
	if len(keywords) == 0 {
		return
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
	return
}

// FindBytes searches TRIE by a given keyword and returns the id if found.
func (d *DoubleArray) FindBytes(input []byte) (id int, ok bool) {
	p, q, _, ok := d.searchBytes(input)
	if !ok {
		return
	}
	p = q
	q = (*d)[p].Base + int(daTerminator)
	if q >= len(*d) || (*d)[q].Check != p || (*d)[q].Base > 0 {
		return
	}
	return -(*d)[q].Base, true
}

// FindString searches TRIE by a given keyword and returns the id if found.
func (d *DoubleArray) FindString(input string) (id int, ok bool) {
	p, q, _, ok := d.searchString(input)
	if !ok {
		return
	}
	p = q
	q = (*d)[p].Base + int(daTerminator)
	if q >= len(*d) || (*d)[q].Check != p || (*d)[q].Base > 0 {
		return
	}
	return -(*d)[q].Base, true
}

// CommonPrefixSearchBytes finds keywords sharing common prefix in an input
// and returns the ids and it's lengths if found.
func (d *DoubleArray) CommonPrefixSearchBytes(input []byte) (ids, lens []int) {
	var p, q int
	bufLen := len(*d)
	for i := range input {
		p = q
		q = (*d)[p].Base + int(input[i])
		if q >= bufLen || (*d)[q].Check != p {
			break
		}
		ahead := (*d)[q].Base + daTerminator
		if ahead < bufLen && (*d)[ahead].Check == q && (*d)[ahead].Base <= 0 {
			ids = append(ids, -(*d)[ahead].Base)
			lens = append(lens, i+1)
		}
	}
	return
}

// CommonPrefixSearchString finds keywords sharing common prefix in an input
// and returns the ids and it's lengths if found.
func (d *DoubleArray) CommonPrefixSearchString(input string) (ids, lens []int) {
	var p, q int
	bufLen := len(*d)
	for i, size := 0, len(input); i < size; i++ {
		p = q
		q = (*d)[p].Base + int(input[i])
		if q >= bufLen || (*d)[q].Check != p {
			break
		}
		ahead := (*d)[q].Base + daTerminator
		if ahead < bufLen && (*d)[ahead].Check == q && (*d)[ahead].Base <= 0 {
			ids = append(ids, -(*d)[ahead].Base)
			lens = append(lens, i+1)
		}
	}
	return
}

// PrefixSearchBytes returns the longest commom prefix keyword in an input if found.
func (d *DoubleArray) PrefixSearchBytes(input []byte) (id int, ok bool) {
	var p, q, i int
	bufLen := len(*d)
	for size := len(input); i < size; i++ {
		p = q
		q = (*d)[p].Base + int(input[i])
		if q >= bufLen || (*d)[q].Check != p {
			break
		}
		ahead := (*d)[q].Base + daTerminator
		if ahead < bufLen && (*d)[ahead].Check == q && (*d)[ahead].Base <= 0 {
			id = -(*d)[ahead].Base
			ok = true
		}
	}
	return
}

// PrefixSearchString returns the longest commom prefix keyword in an input if found.
func (d *DoubleArray) PrefixSearchString(input string) (id int, ok bool) {
	var p, q, i int
	bufLen := len(*d)
	for size := len(input); i < size; i++ {
		p = q
		q = (*d)[p].Base + int(input[i])
		if q >= bufLen || (*d)[q].Check != p {
			break
		}
		ahead := (*d)[q].Base + daTerminator
		if ahead < bufLen && (*d)[ahead].Check == q && (*d)[ahead].Base <= 0 {
			id = -(*d)[ahead].Base
			ok = true
		}
	}
	return
}
