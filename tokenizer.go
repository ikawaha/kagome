//  Copyright (c) 2014 ikawaha.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package kagome

import (
	"io"
	"unicode/utf8"
)

type tokenizeMode int

const (
	initialNodePoolCapacity              = 512
	normalModeTokenize      tokenizeMode = iota + 1
	searchModeTokenize
	extendedModeTokenize
)

// Tokenizer represents morphological analyzer.
type Tokenizer struct {
	dic  *Dic
	udic *UserDic
}

// NewTokenizer create a tokenizer.
func NewTokenizer() (t *Tokenizer) {
	t = new(Tokenizer)
	t.dic = NewSysDic()
	return
}

// SetDic sets dictionary to dic.
func (t *Tokenizer) SetDic(dic *Dic) {
	if dic != nil {
		t.dic = dic
	}
}

// SetUserDic sets user dictionary to udic.
func (t *Tokenizer) SetUserDic(udic *UserDic) {
	t.udic = udic
}

// Tokenize returns morphs of a sentence.
func (t *Tokenizer) Tokenize(input string) (tokens []Token) {
	la := newLattice()
	defer la.free()
	la.dic, la.udic = t.dic, t.udic
	la.build(input)
	la.forward(normalModeTokenize)
	la.backward(normalModeTokenize)
	size := len(la.output)
	tokens = make([]Token, 0, size)
	for i := range la.output {
		n := la.output[size-1-i]
		tok := Token{
			Id:      n.id,
			Class:   n.class,
			Start:   n.start,
			End:     n.start + utf8.RuneCountInString(n.surface),
			Surface: n.surface,
			dic:     t.dic,
			udic:    t.udic,
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

// SearchModeTokenize returns morphs of a sentence.
func (t *Tokenizer) SearchModeTokenize(input string) (tokens []Token) {
	la := newLattice()
	defer la.free()
	la.dic, la.udic = t.dic, t.udic
	la.build(input)
	la.forward(searchModeTokenize)
	la.backward(searchModeTokenize)
	size := len(la.output)
	tokens = make([]Token, 0, size)
	for i := range la.output {
		n := la.output[size-1-i]
		tok := Token{
			Id:      n.id,
			Class:   n.class,
			Start:   n.start,
			End:     n.start + utf8.RuneCountInString(n.surface),
			Surface: n.surface,
			dic:     t.dic,
			udic:    t.udic,
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

// ExtendedModeTokenize returns morphs of a sentence.
func (t *Tokenizer) ExtendedModeTokenize(input string) (tokens []Token) {
	la := newLattice()
	defer la.free()
	la.dic, la.udic = t.dic, t.udic
	la.build(input)
	la.forward(extendedModeTokenize)
	la.backward(extendedModeTokenize)
	size := len(la.output)
	tokens = make([]Token, 0, size)
	for i := range la.output {
		n := la.output[size-1-i]
		tok := Token{
			Id:      n.id,
			Class:   n.class,
			Start:   n.start,
			End:     n.start + utf8.RuneCountInString(n.surface),
			Surface: n.surface,
			dic:     t.dic,
			udic:    t.udic,
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
	la := newLattice()
	defer la.free()
	la.dic, la.udic = t.dic, t.udic
	la.build(input)
	la.forward(normalModeTokenize)
	la.backward(normalModeTokenize)
	size := len(la.output)
	tokens = make([]Token, 0, size)
	for i := range la.output {
		n := la.output[size-1-i]
		tok := Token{
			Id:      n.id,
			Class:   n.class,
			Start:   n.start,
			End:     n.start + utf8.RuneCountInString(n.surface),
			Surface: n.surface,
			dic:     t.dic,
			udic:    t.udic,
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
	tokens = t.Tokenize(input)
	la.dot(w)
	return
}
