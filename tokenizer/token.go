package tokenizer

import (
	"fmt"
	"strings"

	"github.com/ikawaha/kagome-dict/dict"
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
	Index    int
	ID       int
	Class    TokenClass
	Position int // byte position
	Start    int
	End      int
	Surface  string
	dict     *dict.Dict
	udict    *dict.UserDict
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
			if int(id) > len(t.dict.POSTable.NameList) {
				return "", false
			}
			return t.dict.POSTable.NameList[id], true
		}
		i -= len(pos)
		if len(t.dict.Contents) <= t.ID {
			return "", false
		}
		c := t.dict.Contents[t.ID]
		if i >= len(c) {
			return "", false
		}
		return c[i], true
	case UNKNOWN:
		if len(t.dict.UnkDict.Contents) <= t.ID {
			return "", false
		}
		c := t.dict.UnkDict.Contents[t.ID]
		if i >= len(c) {
			return "", false
		}
		return c[i], true
	case USER:
		if len(t.udict.Contents) <= t.ID {
			return "", false
		}
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

// UserExtra represents custom segmentation and custom reading for user entries.
type UserExtra struct {
	Tokens   []string
	Readings []string
}

// UserExtra returns extra data if token comes from a user dict.
func (t Token) UserExtra() *UserExtra {
	if t.Class != USER {
		return nil
	}
	tokens := make([]string, len(t.udict.Contents[t.ID].Tokens))
	copy(tokens, t.udict.Contents[t.ID].Tokens)
	yomi := make([]string, len(t.udict.Contents[t.ID].Yomi))
	copy(yomi, t.udict.Contents[t.ID].Yomi)
	return &UserExtra{
		Tokens:   tokens,
		Readings: yomi,
	}
}

// POS returns POS elements of features.
func (t Token) POS() []string {
	switch t.Class {
	case KNOWN:
		ret := make([]string, 0, len(t.dict.POSTable.POSs[t.ID]))
		for _, id := range t.dict.POSTable.POSs[t.ID] {
			ret = append(ret, t.dict.POSTable.NameList[id])
		}
		return ret
	case UNKNOWN:
		start := 0
		if v, ok := t.dict.UnkDict.ContentsMeta[dict.POSStartIndex]; ok {
			start = int(v)
		}
		end := 1
		if v, ok := t.dict.UnkDict.ContentsMeta[dict.POSHierarchy]; ok {
			end = start + int(v)
		}
		feature := t.dict.UnkDict.Contents[t.ID]
		if start >= end || end > len(feature) {
			return nil
		}
		ret := make([]string, 0, end-start)
		for i := start; i < end; i++ {
			ret = append(ret, feature[i])
		}
		return ret
	case USER:
		pos := t.udict.Contents[t.ID].Pos
		return []string{pos}
	}
	return nil
}

// EqualFeatures returns true, if the features of tokens are equal.
func (t Token) EqualFeatures(tt Token) bool {
	return EqualFeatures(t.Features(), tt.Features())
}

// EqualPOS returns true, if the POSs of tokens are equal.
func (t Token) EqualPOS(tt Token) bool {
	return EqualFeatures(t.POS(), tt.POS())
}

// EqualFeatures returns true, if the features are equal.
func EqualFeatures(lhs, rhs []string) bool {
	if len(lhs) != len(rhs) {
		return false
	}
	for i := 0; i < len(lhs); i++ {
		if lhs[i] != rhs[i] {
			return false
		}
	}
	return true
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
	return fmt.Sprintf("%d:%q (%d: %d, %d) %v [%d]", t.Index, t.Surface, t.Position, t.Start, t.End, t.Class, t.ID)
}

// Equal returns true if tokens are equal.
func (t Token) Equal(v Token) bool {
	return t.ID == v.ID &&
		t.Class == v.Class &&
		t.Surface == v.Surface
}

// TokenData is a data format with all the contents of the token.
type TokenData struct {
	ID            int      `json:"id"`
	Start         int      `json:"start"`
	End           int      `json:"end"`
	Surface       string   `json:"surface"`
	Class         string   `json:"class"`
	POS           []string `json:"pos"`
	BaseForm      string   `json:"base_form"`
	Reading       string   `json:"reading"`
	Pronunciation string   `json:"pronunciation"`
	Features      []string `json:"features"`
}

// NewTokenData returns a data which has with all the contents of the token.
func NewTokenData(t Token) TokenData {
	ret := TokenData{
		ID:       t.ID,
		Start:    t.Start,
		End:      t.End,
		Surface:  t.Surface,
		Class:    t.Class.String(),
		POS:      t.POS(),
		Features: t.Features(),
	}
	if ret.POS == nil {
		ret.POS = []string{}
	}
	if ret.Features == nil {
		ret.Features = []string{}
	}
	ret.BaseForm, _ = t.BaseForm()
	ret.Reading, _ = t.Reading()
	ret.Pronunciation, _ = t.Pronunciation()
	return ret
}
