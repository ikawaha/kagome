package tokenizer

import (
	"fmt"
	"unicode/utf8"
)

func Tokenize(a_input string) (morphs []Morph, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("tokenizer.Tokenize: %v", e.(error))
		}
	}()
	lattice := &Lattice{}
	if err = lattice.build(&a_input); err != nil {
		return
	}
	if err = lattice.forward(); err != nil {
		return
	}
	lattice.backward()
	morphs = make([]Morph, 0, len(lattice.output))
	for _, n := range lattice.output {
		m := Morph{
			id:      n.id,
			class:   n.class,
			Start:   n.start,
			End:     n.start + utf8.RuneCount(n.surface),
			Surface: string(n.surface),
		}
		if m.id == BOSEOS {
			m.Surface = "BOSEOS"
		}
		morphs = append(morphs, m)
	}
	return
}
