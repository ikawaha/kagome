# Why is `kagome_wrapper.c` Needed?

> In short, `kagome_wrapper.c` acts as a **gatekeeper at the ABI boundary**, ensuring that only the intended API is exposed for [FFI](https://en.wikipedia.org/wiki/Foreign_function_interface) (Foreign Function Interface) use.

The C wrapper (`kagome_wrapper.c`) exists to provide a **clean, minimal, and stable C ABI** for using Kagome from other languages via FFI.

## Purpose

- When Go code is built as a C library (`c-archive` / `c-shared`), the resulting binary contains many symbols related to the Go runtime and internal glue code.
- While only `//export`-annotated functions are intended as the public API, it can be unclear to FFI users which symbols are safe and supported to use.
- Exposing too much at the ABI boundary increases the risk of misuse, symbol conflicts, and accidental dependency on implementation details.

## What does the wrapper do?

- The C wrapper defines a **small, explicit public API** for FFI users.
- Only the intended functions are exported using `__attribute__((visibility("default")))`.
- All other symbols are hidden by compiling the wrapper with `-fvisibility=hidden`.
- FFI users (Python, PHP, etc.) interact **only** with this wrapper, not directly with Go-generated symbols.

## Benefits

- Provides a clear and stable ABI for other languages.
- Prevents accidental use of internal or unsupported functions.
- Decouples the FFI-facing API from Go implementation details.
- Makes future refactoring and maintenance safer and easier.
