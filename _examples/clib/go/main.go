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

var (
	mu        sync.Mutex
	instances         = make(map[uintptr]*tokenizer.Tokenizer)
	nextID    uintptr = 1
)

//export KagomeTokenizeStruct
func KagomeTokenizeStruct(handle C.uintptr_t, input *C.char) *C.TokenArray {
	mu.Lock()
	t := instances[uintptr(handle)]
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
		// Always 4 elements, but check length for safety
		if len(pos) > 0 {
			slice[i].pos1 = C.CString(pos[0])
		} else {
			slice[i].pos1 = C.CString("")
		}
		if len(pos) > 1 {
			slice[i].pos2 = C.CString(pos[1])
		} else {
			slice[i].pos2 = C.CString("")
		}
		if len(pos) > 2 {
			slice[i].pos3 = C.CString(pos[2])
		} else {
			slice[i].pos3 = C.CString("")
		}
		if len(pos) > 3 {
			slice[i].pos4 = C.CString(pos[3])
		} else {
			slice[i].pos4 = C.CString("")
		}
		slice[i].start = C.int(tok.Start)
		slice[i].end = C.int(tok.End)
	}

	arr.tokens = cTokens
	return arr
}

//export KagomeFreeTokenArray
func KagomeFreeTokenArray(arr *C.TokenArray) {
	if arr == nil || arr.tokens == nil || arr.length == 0 {
		return
	}
	slice := (*[1 << 30]C.Token)(unsafe.Pointer(arr.tokens))[:arr.length:arr.length]
	for i := 0; i < int(arr.length); i++ {
		C.free(unsafe.Pointer(slice[i].surface))
		C.free(unsafe.Pointer(slice[i].pos1))
		C.free(unsafe.Pointer(slice[i].pos2))
		C.free(unsafe.Pointer(slice[i].pos3))
		C.free(unsafe.Pointer(slice[i].pos4))
	}
	C.free(unsafe.Pointer(arr.tokens))
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
