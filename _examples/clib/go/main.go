package main

/*
#include <stdint.h>
#include <stdlib.h>
typedef struct {
	char* surface;
	char* pos1;
	char* pos2;
	char* pos3;
	char* pos4;
	char* base_form;
	char* conj_type;
	char* conj_form;
	char* reading;
	char* pronunciation;
	int start;
	int end;
} Token;
typedef struct {
  Token* tokens;
  int length;
} TokenArray;
*/
import "C"

import (
	"sync"
	"unsafe"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

// Go type aliases for C struct export
type Token C.struct_Token
type TokenArray C.struct_TokenArray

// Opaque handle struct for C API
type KagomeHandle struct{}

var (
	mu        sync.Mutex
	instances = make(map[*KagomeHandle]*tokenizer.Tokenizer)
)

//export KagomeTokenizeStruct
func KagomeTokenizeStruct(handle unsafe.Pointer, input *C.char) *C.TokenArray {
	mu.Lock()
	t := instances[(*KagomeHandle)(handle)]
	mu.Unlock()
	if t == nil || input == nil {
		return nil
	}

	text := C.GoString(input)
	tokens := t.Tokenize(text)
	n := len(tokens)

	arr := (*C.TokenArray)(C.malloc(C.size_t(unsafe.Sizeof(C.TokenArray{}))))
	arr.length = C.int(n)

	if n == 0 {
		arr.tokens = nil
		return arr
	}

	cTokens := (*C.Token)(C.malloc(C.size_t(n) * C.size_t(unsafe.Sizeof(C.Token{}))))
	slice := (*[1 << 30]C.Token)(unsafe.Pointer(cTokens))[:n:n]

	for i, tok := range tokens {
		slice[i].surface = C.CString(tok.Surface)
		pos := tok.POS()
		features := tok.Features()
		// POS
		slice[i].pos1 = C.CString(getOrEmpty(pos, 0))
		slice[i].pos2 = C.CString(getOrEmpty(pos, 1))
		slice[i].pos3 = C.CString(getOrEmpty(pos, 2))
		slice[i].pos4 = C.CString(getOrEmpty(pos, 3))
		// Features (IPA: 0-品詞,1-細1,2-細2,3-細3,4-活用型,5-活用形,6-原形,7-読み,8-発音)
		slice[i].conj_type = C.CString(getOrEmpty(features, 4))
		slice[i].conj_form = C.CString(getOrEmpty(features, 5))
		slice[i].base_form = C.CString(getOrEmpty(features, 6))
		slice[i].reading = C.CString(getOrEmpty(features, 7))
		slice[i].pronunciation = C.CString(getOrEmpty(features, 8))
		slice[i].start = C.int(tok.Start)
		slice[i].end = C.int(tok.End)
	}
	arr.tokens = cTokens
	return arr
}

// getOrEmpty returns the element at idx or "" if out of range
func getOrEmpty(arr []string, idx int) string {
	if idx < len(arr) {
		return arr[idx]
	}
	return ""
}

// Always free arr if it was malloc'ed, even if arr->tokens is nil.
// This ensures no memory leak occurs for empty results or allocation failures.
//
//export KagomeFreeTokenArray
func KagomeFreeTokenArray(arr *C.TokenArray) {
	if arr == nil {
		return
	}
	// Free token memory if present
	if arr.tokens != nil {
		slice := (*[1 << 30]C.Token)(unsafe.Pointer(arr.tokens))[:arr.length:arr.length]
		for i := 0; i < int(arr.length); i++ {
			C.free(unsafe.Pointer(slice[i].surface))
			C.free(unsafe.Pointer(slice[i].pos1))
			C.free(unsafe.Pointer(slice[i].pos2))
			C.free(unsafe.Pointer(slice[i].pos3))
			C.free(unsafe.Pointer(slice[i].pos4))
			C.free(unsafe.Pointer(slice[i].base_form))
			C.free(unsafe.Pointer(slice[i].conj_type))
			C.free(unsafe.Pointer(slice[i].conj_form))
			C.free(unsafe.Pointer(slice[i].reading))
			C.free(unsafe.Pointer(slice[i].pronunciation))
		}
		C.free(unsafe.Pointer(arr.tokens))
	}
	// Always free arr itself
	C.free(unsafe.Pointer(arr))
}

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
func KagomeInit() unsafe.Pointer {
	mu.Lock()
	defer mu.Unlock()

	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return nil
	}
	handle := &KagomeHandle{}
	instances[handle] = t
	return unsafe.Pointer(handle)
}

//export KagomeTokenize
func KagomeTokenize(handle unsafe.Pointer, input *C.char) *C.char {
	mu.Lock()
	t := instances[(*KagomeHandle)(handle)]
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
func KagomeDestroy(handle unsafe.Pointer) {
	mu.Lock()
	delete(instances, (*KagomeHandle)(handle))
	mu.Unlock()
}

func main() {}
