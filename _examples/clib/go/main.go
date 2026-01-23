package main

/*
#include <stdint.h>
#include <stdlib.h>

// ------------------------------------------------------------------
// C ABI structures
// ------------------------------------------------------------------

// Token represents one morphological token.
// All strings are UTF-8, null-terminated, and allocated with malloc.
typedef struct {
	char* surface;        // Surface form (表層形)
	char* pos1;           // Part-of-speech hierarchy 0 (major class, 大分類)
	char* pos2;           // Part-of-speech hierarchy 1 (middle class, 中分類)
	char* pos3;           // Part-of-speech hierarchy 2 (small class, 小分類)
	char* pos4;           // Part-of-speech hierarchy 3 (fine class, 再分類)
	char* base_form;      // Base form / dictionary form (原形・基本形)
	char* conj_type;      // Conjugation type (活用型)
	char* conj_form;      // Conjugation form (活用形)
	char* reading;        // Reading in katakana (読み)
	char* pronunciation;  // Pronunciation (発音)
	int   start;          // Start position (開始位置)
	int   end;            // End position (終了位置)
} Token;

// TokenArray is an owned array returned to foreign languages.
// Both the array itself and all nested strings must be freed
// by calling KagomeFreeTokenArray.
typedef struct {
	Token* tokens;
	int    length;
} TokenArray;
*/
import "C"

import (
	"sync"
	"unsafe"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

// ------------------------------------------------------------------
// Internal state management
// ------------------------------------------------------------------

// instances maps opaque C handles to Kagome tokenizers.
//
// IMPORTANT (cgo rule):
//
//	Go pointers must never be passed to C.
//	Therefore, we allocate a dummy pointer using C.malloc()
//	and use that pointer as the external handle.
var (
	mu        sync.Mutex
	instances = make(map[unsafe.Pointer]*tokenizer.Tokenizer)
)

// ------------------------------------------------------------------
// Helper utilities
// ------------------------------------------------------------------

// getOrEmpty returns arr[idx] if it exists, otherwise an empty string.
// This avoids bounds checks at every call site.
func getOrEmpty(arr []string, idx int) string {
	if idx < len(arr) {
		return arr[idx]
	}
	return ""
}

// ------------------------------------------------------------------
// Exported C API
// ------------------------------------------------------------------

//export KagomeInit
func KagomeInit() unsafe.Pointer {
	mu.Lock()
	defer mu.Unlock()

	t, err := tokenizer.New(
		ipa.Dict(),
		tokenizer.OmitBosEos(),
	)
	if err != nil {
		return nil
	}

	// Allocate an opaque handle in C memory.
	// This pointer is safe to pass across FFI boundaries.
	handle := C.malloc(1)
	if handle == nil {
		return nil
	}

	instances[handle] = t
	return handle
}

//export KagomeDestroy
func KagomeDestroy(handle unsafe.Pointer) {
	if handle == nil {
		return
	}

	mu.Lock()
	delete(instances, handle)
	mu.Unlock()

	// Free the dummy handle allocated in KagomeInit.
	C.free(handle)
}

//export KagomeTokenizeStruct
func KagomeTokenizeStruct(handle unsafe.Pointer, input *C.char) *C.TokenArray {
	if handle == nil || input == nil {
		return nil
	}

	mu.Lock()
	t := instances[handle]
	mu.Unlock()

	if t == nil {
		return nil
	}

	text := C.GoString(input)
	tokens := t.Tokenize(text)
	n := len(tokens)

	// Allocate TokenArray (always owned by caller).
	arr := (*C.TokenArray)(C.malloc(C.size_t(unsafe.Sizeof(C.TokenArray{}))))
	if arr == nil {
		return nil
	}

	arr.length = C.int(n)

	if n == 0 {
		arr.tokens = nil
		return arr
	}

	// Allocate contiguous Token array.
	cTokens := (*C.Token)(C.malloc(
		C.size_t(n) * C.size_t(unsafe.Sizeof(C.Token{})),
	))
	if cTokens == nil {
		C.free(unsafe.Pointer(arr))
		return nil
	}

	slice := unsafe.Slice(cTokens, n)

	for i, tok := range tokens {
		pos := tok.POS()
		features := tok.Features()

		slice[i] = C.Token{
			surface:       C.CString(tok.Surface),
			pos1:          C.CString(getOrEmpty(pos, 0)),
			pos2:          C.CString(getOrEmpty(pos, 1)),
			pos3:          C.CString(getOrEmpty(pos, 2)),
			pos4:          C.CString(getOrEmpty(pos, 3)),
			conj_type:     C.CString(getOrEmpty(features, 4)),
			conj_form:     C.CString(getOrEmpty(features, 5)),
			base_form:     C.CString(getOrEmpty(features, 6)),
			reading:       C.CString(getOrEmpty(features, 7)),
			pronunciation: C.CString(getOrEmpty(features, 8)),
			start:         C.int(tok.Start),
			end:           C.int(tok.End),
		}
	}

	arr.tokens = cTokens
	return arr
}

//export KagomeFreeTokenArray
func KagomeFreeTokenArray(arr *C.TokenArray) {
	if arr == nil {
		return
	}

	if arr.tokens != nil {
		slice := unsafe.Slice(arr.tokens, int(arr.length))
		for _, t := range slice {
			C.free(unsafe.Pointer(t.surface))
			C.free(unsafe.Pointer(t.pos1))
			C.free(unsafe.Pointer(t.pos2))
			C.free(unsafe.Pointer(t.pos3))
			C.free(unsafe.Pointer(t.pos4))
			C.free(unsafe.Pointer(t.base_form))
			C.free(unsafe.Pointer(t.conj_type))
			C.free(unsafe.Pointer(t.conj_form))
			C.free(unsafe.Pointer(t.reading))
			C.free(unsafe.Pointer(t.pronunciation))
		}
		C.free(unsafe.Pointer(arr.tokens))
	}

	// Always free the container itself.
	C.free(unsafe.Pointer(arr))
}

// ------------------------------------------------------------------
// Optional utility API (example / test)
// ------------------------------------------------------------------

//export Echo
func Echo(input *C.char) *C.char {
	if input == nil {
		return nil
	}
	return C.CString(C.GoString(input))
}

//export EchoFree
func EchoFree(p *C.char) {
	if p != nil {
		C.free(unsafe.Pointer(p))
	}
}

func main() {}
