package tokenizer

import (
	"fmt"
	"strings"

	"github.com/ikawaha/kagome/v2/dict"
	"github.com/ikawaha/kagome/v2/tokenizer/lattice"
)

// TokenClass represents the token type.
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
	switch lattice.NodeClass(t.Class) {
	case lattice.DUMMY:
		return nil
	case lattice.KNOWN:
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
	case lattice.UNKNOWN:
		features := make([]string, len(t.dict.UnkDict.Contents[t.ID]))
		for i := range t.dict.UnkDict.Contents[t.ID] {
			features[i] = t.dict.UnkDict.Contents[t.ID][i]
		}
		return features
	case lattice.USER:
		pos := t.udict.Contents[t.ID].Pos
		tokens := strings.Join(t.udict.Contents[t.ID].Tokens, "/")
		yomi := strings.Join(t.udict.Contents[t.ID].Yomi, "/")
		return []string{pos, tokens, yomi}
	}
	return nil
}

// Pos returns the first element of features.
func (t Token) Pos() string {
	f := t.Features()
	if len(f) < 1 {
		return ""
	}
	return f[0]
}

// String returns a string representation of a token.
func (t Token) String() string {
	return fmt.Sprintf("%v(%v, %v)%v[%v]", t.Surface, t.Start, t.End, t.Class, t.ID)
}
