package tokenizer

import (
	"fmt"
	"strings"

	"github.com/ikawaha/kagome/v2/dict"
	"github.com/ikawaha/kagome/v2/tokenizer/lattice"
)

// TokenClass represents the token class.
type TokenClass lattice.NodeClass

const (
	// DUMMY represents the dummy token.
	DUMMY = TokenClass(lattice.DUMMY)
	// KNOWN represents the token in the dictionary.
	KNOWN = TokenClass(lattice.KNOWN)
	// UNKNOWN represents the token which is not in the dictionary.
	UNKNOWN = TokenClass(lattice.UNKNOWN)
	// USER represents the token in the user dictionary.
	USER = TokenClass(lattice.USER)
)

// String returns string representation of a token class.
func (c TokenClass) String() string {
	ret := ""
	switch c {
	case DUMMY:
		ret = "DUMMY"
	case KNOWN:
		ret = "KNOWN"
	case UNKNOWN:
		ret = "UNKNOWN"
	case USER:
		ret = "USER"
	}
	return ret
}

// Token represents a morph of a sentence.
type Token struct {
	ID      int
	Class   TokenClass
	Start   int
	End     int
	Surface string
	dict    *dict.Dict
	udict   *dict.UserDict
}

// Features returns contents of a token.
func (t Token) Features() []string {
	switch t.Class {
	case KNOWN:
		var c int
		if t.dict.Contents != nil {
			c = len(t.dict.Contents[t.ID])
		}
		features := make([]string, 0, len(t.dict.POSTable.POSs[t.ID])+c)
		for _, id := range t.dict.POSTable.POSs[t.ID] {
			features = append(features, t.dict.POSTable.NameList[id])
		}
		if t.dict.Contents != nil {
			features = append(features, t.dict.Contents[t.ID]...)
		}
		return features
	case UNKNOWN:
		features := make([]string, len(t.dict.UnkDict.Contents[t.ID]))
		for i := range t.dict.UnkDict.Contents[t.ID] {
			features[i] = t.dict.UnkDict.Contents[t.ID][i]
		}
		return features
	case USER:
		pos := t.udict.Contents[t.ID].Pos
		tokens := strings.Join(t.udict.Contents[t.ID].Tokens, "/")
		yomi := strings.Join(t.udict.Contents[t.ID].Yomi, "/")
		return []string{pos, tokens, yomi}
	}
	return nil
}

// FeatureAt returns the i th feature if exists.
func (t Token) FeatureAt(i int) (string, bool) {
	if i < 0 {
		return "", false
	}
	switch t.Class {
	case KNOWN:
		pos := t.dict.POSTable.POSs[t.ID]
		if i < len(pos) {
			id := pos[i]
			if id < 0 || int(id) > len(t.dict.POSTable.NameList) {
				return "", false
			}
			return t.dict.POSTable.NameList[id], true
		}
		i -= len(pos)
		c := t.dict.Contents[t.ID]
		if i < 0 || i >= len(c) {
			return "", false
		}
		return c[i], true
	case UNKNOWN:
		c := t.dict.UnkDict.Contents[t.ID]
		if i < 0 || i >= len(c) {
			return "", false
		}
		return c[i], true
	case USER:
		switch i {
		case 0:
			return t.udict.Contents[t.ID].Pos, true
		case 1:
			return strings.Join(t.udict.Contents[t.ID].Tokens, "/"), true
		case 2:
			return strings.Join(t.udict.Contents[t.ID].Yomi, "/"), true
		}
	}
	return "", false
}

// POS returns POS elements of features.
func (t Token) POS() []string {
	f := t.Features()
	if len(f) == 0 {
		return nil
	}
	var meta dict.ContentsMeta
	switch t.Class {
	case KNOWN:
		meta = t.dict.ContentsMeta
	case UNKNOWN:
		meta = t.dict.UnkDict.ContentsMeta
	}
	start := meta[dict.POSStartIndex]
	if start < 0 || int(start) > len(f) {
		start = 0
	}
	end, ok := meta[dict.POSHierarchy]
	if !ok || start+end < 0 || int(start+end) > len(f) {
		end = 1
	}
	return f[start : start+end] // default [0:1]
}

// InflectionalType returns the inflectional type feature if exists.
func (t Token) InflectionalType() (string, bool) {
	return t.pickupFromFeatures(dict.InflectionalType)
}

// InflectionalForm returns the inflectional form feature if exists.
func (t Token) InflectionalForm() (string, bool) {
	return t.pickupFromFeatures(dict.InflectionalForm)
}

// BaseForm returns the base form features if exists.
func (t Token) BaseForm() (string, bool) {
	return t.pickupFromFeatures(dict.BaseFormIndex)
}

// Reading returns the reading feature if exists.
func (t Token) Reading() (string, bool) {
	return t.pickupFromFeatures(dict.ReadingIndex)
}

// Pronunciation returns the pronunciation feature if exists.
func (t Token) Pronunciation() (string, bool) {
	return t.pickupFromFeatures(dict.PronunciationIndex)
}

func (t Token) pickupFromFeatures(key string) (string, bool) {
	var meta dict.ContentsMeta
	switch t.Class {
	case KNOWN:
		meta = t.dict.ContentsMeta
	case UNKNOWN:
		meta = t.dict.UnkDict.ContentsMeta
	}
	i, ok := meta[key]
	if !ok {
		return "", false
	}
	return t.FeatureAt(int(i))
}

// String returns a string representation of a token.
func (t Token) String() string {
	return fmt.Sprintf("%v(%v, %v)%v[%v]", t.Surface, t.Start, t.End, t.Class, t.ID)
}
