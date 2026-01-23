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

// wouldOverflowTokenAllocation checks if allocating n tokens would cause integer overflow.
// Returns true if the allocation would be unsafe.
func wouldOverflowTokenAllocation(n int) bool {
	if n < 0 {
		return true
	}
	tokenSize := C.size_t(unsafe.Sizeof(C.Token{}))
	maxSafeTokens := (^C.size_t(0)) / tokenSize
	return C.size_t(n) > maxSafeTokens
}

// freeStrings safely frees multiple C strings.
// Checks for nil before freeing (safe to pass nil pointers).
func freeStrings(strs ...*C.char) {
	for _, s := range strs {
		if s != nil {
			C.free(unsafe.Pointer(s))
		}
	}
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
	defer mu.Unlock()

	t := instances[handle]
	if t == nil {
		return nil
	}

	text := C.GoString(input)
	tokens := t.Tokenize(text)
	n := len(tokens)

	// Check for integer overflow in allocation
	if wouldOverflowTokenAllocation(n) {
		return nil
	}

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

		// Allocate all strings for this token
		surface := C.CString(tok.Surface)
		pos1 := C.CString(getOrEmpty(pos, 0))
		pos2 := C.CString(getOrEmpty(pos, 1))
		pos3 := C.CString(getOrEmpty(pos, 2))
		pos4 := C.CString(getOrEmpty(pos, 3))
		conjType := C.CString(getOrEmpty(features, 4))
		conjForm := C.CString(getOrEmpty(features, 5))
		baseForm := C.CString(getOrEmpty(features, 6))
		reading := C.CString(getOrEmpty(features, 7))
		pronunciation := C.CString(getOrEmpty(features, 8))

		// Check if any allocation failed
		if surface == nil || pos1 == nil || pos2 == nil || pos3 == nil || pos4 == nil ||
			conjType == nil || conjForm == nil || baseForm == nil || reading == nil || pronunciation == nil {

			// Free strings we just allocated for current token
			freeStrings(surface, pos1, pos2, pos3, pos4, conjType, conjForm, baseForm, reading, pronunciation)

			// Free all previously completed tokens
			for j := 0; j < i; j++ {
				freeStrings(
					slice[j].surface,
					slice[j].pos1,
					slice[j].pos2,
					slice[j].pos3,
					slice[j].pos4,
					slice[j].conj_type,
					slice[j].conj_form,
					slice[j].base_form,
					slice[j].reading,
					slice[j].pronunciation,
				)
			}

			C.free(unsafe.Pointer(cTokens))
			C.free(unsafe.Pointer(arr))
			return nil
		}

		slice[i] = C.Token{
			surface:       surface,
			pos1:          pos1,
			pos2:          pos2,
			pos3:          pos3,
			pos4:          pos4,
			conj_type:     conjType,
			conj_form:     conjForm,
			base_form:     baseForm,
			reading:       reading,
			pronunciation: pronunciation,
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
// Test utilities (wrapped as kagome_echo/kagome_echo_free)
// ------------------------------------------------------------------

// Echo copies a string and returns it.
// Used for testing FFI setup (string passing, memory allocation).
//
// FFI users should call kagome_echo() from the C wrapper, not this directly.
//
//export Echo
func Echo(input *C.char) *C.char {
	if input == nil {
		return nil
	}
	return C.CString(C.GoString(input))
}

// EchoFree frees a string returned by Echo.
// FFI users should call kagome_echo_free() from the C wrapper, not this directly.
//
//export EchoFree
func EchoFree(p *C.char) {
	if p != nil {
		C.free(unsafe.Pointer(p))
	}
}

func main() {}
