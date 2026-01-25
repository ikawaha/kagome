# Go FFI Bridge Layer

This directory contains the **Go-to-C FFI bridge** that exposes Kagome's tokenizer to other languages via C.

## What is this?

- **main.go** — Go implementation that exports C-compatible functions for FFI
- **main_test.go** — Comprehensive tests including concurrency and memory safety checks

## Architecture

```txt
Other Languages (Python, PHP, etc.)
          ↓
Shared Library (../bin/kagome.(so|dll|dylib))
          ↓
C Wrapper (../c_wrapper/)
          ↓
Go Bridge (this directory)
          ↓
Kagome Library (github.com/ikawaha/kagome/v2)
```

This layer sits between the C wrapper and the Kagome library, handling:

- Memory management across the FFI boundary
- Thread-safe tokenizer instance management
- Conversion between Go and C data structures

## Why is this needed?

Go's cgo system can export functions to C, but the exported interface is not always clean:

- Go exports many internal symbols (runtime, helpers, etc.)
- Go pointers cannot be passed to C (cgo rule)
- Memory allocated by Go needs special handling

This layer solves these issues by:

- Managing tokenizer instances with opaque C handles (using `C.malloc`)
- Allocating all returned memory with `C.malloc` (not Go's allocator)
- Providing explicit cleanup functions to prevent memory leaks
- Using mutex locks for thread-safe instance management

## Key Functions Exported

**Tokenizer Lifecycle:**

- `KagomeInit()` — Create a tokenizer, returns opaque handle
- `KagomeDestroy(handle)` — Free a tokenizer

**Tokenization:**

- `KagomeTokenizeStruct(handle, input)` — Tokenize text, returns C-allocated token array
- `KagomeFreeTokenArray(arr)` — Free tokenization results

**Testing Utilities:**

- `Echo(input)` — Echo a string (for testing FFI setup)
- `EchoFree(str)` — Free echoed string

## Design Decisions

### Memory Management

All memory returned to C is allocated with `C.malloc` (not Go's allocator):

- Token arrays
- String fields in tokens
- The TokenArray container itself

This ensures FFI callers can safely hold pointers without triggering Go's garbage collector.

### Thread Safety

- Uses `sync.Mutex` to protect the tokenizer instance map
- Lock is acquired **only** for map access, not during tokenization
- Allows concurrent tokenization across different handles (performance optimization)

### Error Handling

- All exported functions are nil-safe (return NULL on error)
- Memory allocation failures are handled gracefully
- Partial allocation failures trigger complete cleanup (no memory leaks)

## Running Tests

From the example root (`_examples/clib/`):

```sh
# Run Go tests only
make go-test

# Run all tests (Go, Python, PHP)
make test
```

From this directory:

```sh
# Run tests with race detector
go clean -testcache && go test -race -v .

# Run specific test
go test -v -run TestKagomeTokenizeConcurrent
```

## Test Coverage

The test suite includes:

- Tokenizer initialization and cleanup
- Concurrent tokenization (race condition detection)
- Thread-safe instance map operations
- Integer overflow protection
- Bounds checking helpers
- Nil-safe string deallocation

## Important cgo Rules

When modifying this code, remember:

1. **Never pass Go pointers to C** — Use `C.malloc` for handles
2. **Free what you malloc** — Every `C.malloc` needs `C.free`
3. **Use `C.CString` carefully** — It mallocs, must be freed
4. **Go strings are safe** — `C.GoString` copies, no free needed

## More Information

See the main README at `../README.md` for:

- How to build the shared library
- How to use it from Python and PHP
- Docker support
- Example usage
