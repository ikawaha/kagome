package da

import (
	"sort"
)

const (
	_BUFSIZE     = 51200
	_EXPANDRATIO = 2
	_TERMINATOR  = '\x00'
	_ROOTID      = 0
)

type DoubleArray []struct {
	Base, Check int
}

func NewDoubleArray() *DoubleArray {
	da := new(DoubleArray)
	*da = make(DoubleArray, _BUFSIZE)

	(*da)[_ROOTID].Base = 1
	(*da)[_ROOTID].Check = -1

	size := len(*da)
	for i := 1; i < size; i++ {
		(*da)[i].Base = -(i - 1)
		(*da)[i].Check = -(i + 1)
	}

	(*da)[1].Base = -(size - 1)
	(*da)[size-1].Check = -1

	return da
}

func (this *DoubleArray) SearchBytes(a_keyword []byte) (id int, ok bool) {
	vec := append(a_keyword, _TERMINATOR)
	p, q, _ := this.search(vec)
	if (*this)[q].Check != p || (*this)[q].Base > 0 {
		return 0, false
	}
	return -(*this)[q].Base, true
}

func (this *DoubleArray) Search(a_keyword string) (id int, ok bool) {
	return this.SearchBytes([]byte(a_keyword))
}

func (this *DoubleArray) CommonPrefixSearchBytes(a_keyword []byte) (keywords [][]byte, ids []int) {
	keywords, ids = make([][]byte, 0), make([]int, 0)
	p, q, i := 0, 0, 0
	str := a_keyword
	buf_size := len(*this)
	for size := len(str); i < size; i++ {
		p = q
		ch := int(str[i])
		q = (*this)[p].Base + ch
		if q >= buf_size || (*this)[q].Check != p {
			break
		}
		ahead := (*this)[q].Base + _TERMINATOR
		if ahead < buf_size && (*this)[ahead].Check == q && (*this)[ahead].Base < 0 {
			keywords = append(keywords, str[0:i+1])
			ids = append(ids, -(*this)[ahead].Base)
		}
	}
	return
}

func (this *DoubleArray) CommonPrefixSearch(a_keyword string) (keywords []string, ids []int) {

	var commonPrefixs [][]byte
	commonPrefixs, ids = this.CommonPrefixSearchBytes([]byte(a_keyword))
	keywords = make([]string, 0, len(commonPrefixs))
	for _, prefix := range commonPrefixs {
		keywords = append(keywords, string(prefix))
	}
	return
}

func (this *DoubleArray) PrefixSearchBytes(a_keyword []byte) (keyword string, id int, ok bool) {
	p, q, i := 0, 0, 0
	buf_size := len(*this)
	for size := len(a_keyword); i < size; i++ {
		p = q
		ch := int(a_keyword[i])
		q = (*this)[p].Base + ch
		if q >= buf_size || (*this)[q].Check != p {
			break
		}
		ahead := (*this)[q].Base + _TERMINATOR
		if ahead < buf_size && (*this)[ahead].Check == q && (*this)[ahead].Base < 0 {
			keyword = string(a_keyword[0 : i+1])
			id = -(*this)[ahead].Base
			ok = true
		}
	}
	return
}

func (this *DoubleArray) PrefixSearch(a_keyword string) (keyword string, id int, ok bool) {
	return this.PrefixSearchBytes([]byte(a_keyword))
}

func (this *DoubleArray) Efficiency() (int, int, float64) {
	unspent := 0
	for _, pair := range *this {
		if pair.Check < 0 {
			unspent++
		}
	}
	return unspent, len(*this), float64(len(*this)-unspent) / float64(len(*this)) * 100
}

func (this *DoubleArray) expand() {
	srcSize := len(*this)
	dst := new(DoubleArray)
	dstSize := srcSize * _EXPANDRATIO
	*dst = make(DoubleArray, dstSize)
	copy(*dst, *this)

	for i := srcSize; i < dstSize; i++ {
		(*dst)[i].Base = -(i - 1)
		(*dst)[i].Check = -(i + 1)
	}

	start := -(*this)[0].Check
	end := -(*dst)[start].Base
	(*dst)[srcSize].Base = -end
	(*dst)[start].Base = -(dstSize - 1)
	(*dst)[end].Check = -srcSize
	(*dst)[dstSize-1].Check = -start

	*this = *dst
}

func (this *DoubleArray) shrink() {
	srcSize := len(*this)
	for i, size := 0, srcSize; i < size; i++ {
		if (*this)[size-i-1].Check < 0 {
			srcSize--
		} else {
			break
		}
	}
	if srcSize == len(*this) {
		return
	}
	var dst *DoubleArray = new(DoubleArray)
	*dst = make(DoubleArray, srcSize)
	copy(*dst, (*this)[:srcSize])
	*this = *dst
}

func (this *DoubleArray) search(a_str []byte) (p, q, i int) {
	p, q, i = 0, 0, 0
	bufSize := len(*this)
	for size := len(a_str); i < size; i++ {
		p = q
		ch := int(a_str[i])
		q = (*this)[p].Base + ch
		if q >= bufSize || (*this)[q].Check != p {
			return p, q, i
		}
	}
	return p, q, i
}

func (this *DoubleArray) Build(a_keywords []string) {
	list := a_keywords
	if len(list) == 0 {
		return
	}
	if !sort.StringsAreSorted(list) {
		sort.Strings(list)
	}
	branches := make([]int, len(a_keywords))
	for i, size := 0, len(a_keywords); i < size; i++ {
		branches[i] = i
	}
	this.append(0, 0, branches, a_keywords, nil)
	this.shrink()
}

func (this *DoubleArray) BuildWithIds(a_keywords []string, a_ids []int) {
	if len(a_keywords) == 0 || len(a_keywords) != len(a_ids) {
		return
	}
	list := a_keywords
	if !sort.StringsAreSorted(list) {
		sort.Strings(list)
	}
	branches := make([]int, len(a_keywords))
	for i, size := 0, len(a_keywords); i < size; i++ {
		branches[i] = i
	}
	this.append(0, 0, branches, a_keywords, a_ids)
	this.shrink()
}

func (this *DoubleArray) setBase(a_p, a_base int) {
	if a_p == _ROOTID {
		return
	}
	if (*this)[a_p].Check < 0 {
		if (*this)[a_p].Base == (*this)[a_p].Check {
			this.expand()
		}
		prev := -(*this)[a_p].Base
		next := -(*this)[a_p].Check
		if -a_p == (*this)[_ROOTID].Check {
			(*this)[_ROOTID].Check = (*this)[a_p].Check
		}
		(*this)[next].Base = (*this)[a_p].Base
		(*this)[prev].Check = (*this)[a_p].Check
	}
	(*this)[a_p].Base = a_base
}

func (this *DoubleArray) setCheck(a_p, a_check int) {
	if (*this)[a_p].Base == (*this)[a_p].Check {
		this.expand()
	}
	prev := -(*this)[a_p].Base
	next := -(*this)[a_p].Check
	if -a_p == (*this)[_ROOTID].Check {
		(*this)[_ROOTID].Check = (*this)[a_p].Check
	}

	(*this)[next].Base = (*this)[a_p].Base
	(*this)[prev].Check = (*this)[a_p].Check
	(*this)[a_p].Check = a_check

}

func (this *DoubleArray) seekAndMark(a_p int, a_chars []byte) { // chars != nil
	free := _ROOTID
	rep := int(a_chars[0])
	var base int
	for {
	L_start:
		if free != _ROOTID && (*this)[free].Check == (*this)[_ROOTID].Check {
			this.expand()
		}
		free = -(*this)[free].Check
		base = free - rep
		if base <= 0 {
			continue
		}
		for _, ch := range a_chars {
			q := base + int(ch)
			if q < len(*this) && (*this)[q].Check >= 0 {
				goto L_start
			}
		}
		break
	}
	this.setBase(a_p, base)
	for _, ch := range a_chars {
		q := (*this)[a_p].Base + int(ch)
		if q >= len(*this) {
			this.expand()
		}
		this.setCheck(q, a_p)
	}
}

func (this *DoubleArray) append(a_p, a_i int, a_branches []int, a_keywords []string, a_ids []int) {
	chars := make([]byte, 0)
	subtree := make(map[byte][]int)
	for _, keyId := range a_branches {
		str := []byte(a_keywords[keyId])
		var ch byte
		if a_i >= len(str) {
			ch = _TERMINATOR
		} else {
			ch = str[a_i]
		}
		if size := len(chars); size == 0 || chars[len(chars)-1] != ch {
			chars = append(chars, ch)
		}
		if ch != _TERMINATOR {
			subtree[ch] = append(subtree[ch], keyId)
		}
	}
	this.seekAndMark(a_p, chars)
	for _, ch := range chars {
		q := (*this)[a_p].Base + int(ch)
		if len(subtree[ch]) == 0 {
			if a_ids == nil {
				(*this)[q].Base = -a_branches[0]
			} else {
				(*this)[q].Base = -a_ids[a_branches[0]]
			}
		} else {
			this.append(q, a_i+1, subtree[ch], a_keywords, a_ids)
		}
	}
}
