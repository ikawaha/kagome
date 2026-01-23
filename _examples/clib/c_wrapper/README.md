# C Wrapper Layer

The C wrapper provides a **stable, public API** for using Kagome from other languages via [FFI](https://en.wikipedia.org/wiki/Foreign_function_interface) (Foreign Function Interface).

## What does it do?

The wrapper exposes four simple functions:

- `kagome_init()` — Create a tokenizer
- `kagome_destroy()` — Free a tokenizer
- `kagome_tokenize()` — Tokenize text
- `kagome_free_token_array()` — Free tokenization results

All FFI users (Python, PHP, Rust, etc.) call these functions. They never call Go functions directly.

## Why is this needed?

When Go builds a C library, it exports many internal symbols (runtime code, helper functions, etc.). This creates problems:

- FFI users don't know which symbols are safe to use
- Internal symbols may change between Go versions
- Risk of symbol name conflicts

The wrapper solves this by:

- Defining a small, explicit public API
- Using `-fvisibility=hidden` to hide all symbols except `kagome_*` functions
- Providing consistent naming (snake_case) for all languages

## Benefits

- **Stable**: Go changes don't break FFI code
- **Clear**: Only four functions to learn
- **Safe**: No access to internal symbols
- **Consistent**: Same API for all languages
