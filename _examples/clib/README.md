# Kagome C Library Example: Python & PHP FFI

This directory contains minimal examples showing how to use **`kagome` (the Japanese Morphological Analyzer written in Go) as a C library** from both Python and PHP using [FFI](https://en.wikipedia.org/wiki/Foreign_function_interface) (Foreign Function Interface).

## What is this?

- Provides C ABI wrappers for the Kagome tokenizer built from Go.
- Includes example scripts for Python (using `ctypes`) and PHP (using `FFI`).
- Each script loads the shared library (`.so`, `.dll`, or `.dylib`), tokenizes a Japanese sentence, and checks the result.

## Directory Structure

- [go/](./go) — Go source for the C ABI wrapper
- [c_wrapper/](./c_wrapper) — C wrapper source code for stable ABI exposure (see its README for details)
- [python3/](./python3) — Python 3 example using `ctypes` (see its README for details)
- [php8/](./php8) — PHP 8 example using `FFI` (see its README for details)
- `bin/` — Built shared libraries and archives (output directory of `libkagome`)

## Requirements (for local run)

- Go and a C compiler (e.g., gcc) to build the shared library
- Python 3.12 or later for the Python example
- PHP 8 or later with the `FFI` extension for the PHP example

> **For Docker users**, you can run these examples without installing the above dependencies locally using Docker and Docker Compose.

## How to build (locally)

From this directory (`_examples/clib/`), run:

```sh
make build
```

This will build the shared library for your platform in `./bin/`.

## How to test (locally)

From this directory (`_examples/clib/`), run:

```sh
make test
```

If the test passes, you will see `PASS` in the output.

## How to test (with Docker and Docker Compose)

1. Build Docker images and the shared library:

    ```sh
    make docker-build
    ```

2. Run the test containers:

    ```sh
    make docker-test
    ```

If the test passes, you will see `PASS` in the output.

## More Information

See the README in each subdirectory (`python3/README.md`, `php/README.md`) for details and expected output.
