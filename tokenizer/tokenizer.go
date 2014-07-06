package tokenizer

import (
	"fmt"
	"unicode/utf8"

	"github.com/ikawaha/kagome/dic"
)

type Tokenizer struct {
	lattice *lattice
	userdic *dic.UserDic
}

func NewTokenizer() *Tokenizer {
	ret := new(Tokenizer)
	ret.lattice = NewLattice()
	return ret
}

func (this *Tokenizer) SetUserDic(a_udic *dic.UserDic) {
	this.userdic = a_udic
	this.lattice.setUserDic(a_udic)
}

func (this *Tokenizer) Tokenize(a_input string) (morphs []Morph, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("tokenizer.Tokenize(): %v", e.(error))
		}
	}()

	if err = this.lattice.build(&a_input); err != nil {
		return
	}
	if err = this.lattice.forward(); err != nil {
		return
	}
	this.lattice.backward()

	size := len(this.lattice.output)
	morphs = make([]Morph, 0, size)
	for i := 1; i < size; i++ {
		n := this.lattice.output[size-1-i]
		m := Morph{
			Id:      n.id,
			Class:   n.class,
			Start:   n.start,
			End:     n.start + utf8.RuneCount(n.surface),
			Surface: string(n.surface),
		}
		if m.Id == BOSEOS {
			m.Surface = "EOS"
		}
		if m.Class == USER {
			m.ex = &this.userdic.Contents[n.id]
		}
		morphs = append(morphs, m)
	}
	return
}
