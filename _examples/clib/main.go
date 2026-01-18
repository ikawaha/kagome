package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"sync"
	"unsafe"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

var (
	mu        sync.Mutex
	instances         = make(map[uintptr]*tokenizer.Tokenizer)
	nextID    uintptr = 1
)

//export Echo
func Echo(input *C.char) *C.char {
	if input == nil {
		return nil
	}
	// Just return a copy of the input string
	return C.CString(C.GoString(input))
}

//export EchoFree
func EchoFree(p *C.char) {
	if p == nil {
		return
	}
	C.free(unsafe.Pointer(p))
}

//export KagomeInit
func KagomeInit(dictPath *C.char) C.uintptr_t {
	mu.Lock()
	defer mu.Unlock()

	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return 0
	}

	id := nextID
	nextID++
	instances[id] = t

	return C.uintptr_t(id)
}

//export KagomeTokenize
func KagomeTokenize(handle C.uintptr_t, input *C.char) *C.char {
	mu.Lock()
	t := instances[uintptr(handle)]
	mu.Unlock()

	if t == nil || input == nil {
		return nil
	}

	text := C.GoString(input)
	tokens := t.Tokenize(text)

	// Build output using a string builder for safety
	var builder []byte
	for _, tok := range tokens {
		builder = append(builder, []byte(tok.Surface)...)
		builder = append(builder, '\n')
	}
	if len(builder) == 0 {
		return nil
	}

	// Always allocate with C.CString and return that pointer
	return C.CString(string(builder))
}

//export KagomeFree
func KagomeFree(p *C.char) {
	if p == nil {
		return
	}
	C.free(unsafe.Pointer(p))
}

//export KagomeDestroy
func KagomeDestroy(handle C.uintptr_t) {
	mu.Lock()
	delete(instances, uintptr(handle))
	mu.Unlock()
}

func main() {}
