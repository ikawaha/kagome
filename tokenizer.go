package kagome

import (
	"io"
	"unicode/utf8"
)

const (
	initialNodePoolCapacity = 512
)

// Tokenizer represents morphological analyzer.
type Tokenizer struct {
	lattice *lattice
}

// NewTokenizer create a tokenizer.
func NewTokenizer() (t *Tokenizer) {
	t = new(Tokenizer)
	t.lattice = newLattice()
	t.lattice.setDic(NewSysDic())
	t.lattice.setNodePool(initialNodePoolCapacity)
	return
}

// NewThreadsafeTokenizer create a threadsafe tokenizer.
func NewThreadsafeTokenizer() (t *Tokenizer) {
	t = new(Tokenizer)
	t.lattice = newLattice()
	t.lattice.setDic(NewSysDic())
	return
}

// SetDic sets dictionary to dic.
func (t *Tokenizer) SetDic(dic *Dic) {
	t.lattice.setDic(dic)
}

// SetUserDic sets user dictionary to udic.
func (t *Tokenizer) SetUserDic(udic *UserDic) {
	t.lattice.setUserDic(udic)
}

// Tokenize returns morphs of a sentence.
func (t *Tokenizer) Tokenize(input string) (tokens []Token) {
	t.lattice.build(input)
	t.lattice.forward()
	t.lattice.backward()
	size := len(t.lattice.output)
	tokens = make([]Token, 0, size)
	for i := range t.lattice.output {
		n := t.lattice.output[size-1-i]
		tok := Token{
			Id:      n.id,
			Class:   n.class,
			Start:   n.start,
			End:     n.start + utf8.RuneCountInString(n.surface),
			Surface: n.surface,
			dic:     t.lattice.dic,
			udic:    t.lattice.udic,
		}
		if tok.Id == BosEosId {
			if i == 0 {
				tok.Surface = "BOS"
			} else {
				tok.Surface = "EOS"
			}
		}
		tokens = append(tokens, tok)
	}
	return
}

// Dot returns morphs of a sentense and exports a lattice graph to dot format.
func (t *Tokenizer) Dot(input string, w io.Writer) (tokens []Token) {
	tokens = t.Tokenize(input)
	t.lattice.dot(w)
	return
}
