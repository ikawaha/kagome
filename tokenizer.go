package kagome

import (
	"io"
	"unicode/utf8"
)

type Tokenizer struct {
	lattice *lattice
}

func NewTokenizer() (t *Tokenizer) {
	t = new(Tokenizer)
	t.lattice = newLattice()
	t.lattice.setDic(NewSysDic())
	return
}

func (t *Tokenizer) SetDic(dic *Dic) {
	t.lattice.setDic(dic)
}

func (t *Tokenizer) SetUserDic(udic *UserDic) {
	t.lattice.setUserDic(udic)
}

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
		if tok.Id == BOSEOS_ID {
			tok.Surface = "EOS"
		}

		//XXX
		//if tok.Class == USER {
		//      udic := t.lattice.getUdic()
		//	tok.ex = udic.Contents[n.id]
		//}
		tokens = append(tokens, tok)
	}
	return
}

func (t *Tokenizer) Dot(input string, w io.Writer) (tokens []Token) {
	tokens = t.Tokenize(input)
	t.lattice.dot(w)
	return
}
