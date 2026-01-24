# Kagome C Library Example: Multi-Language FFI

This directory shows how to use **Kagome** (a Japanese Morphological Analyzer written in Go) from other languages via [FFI](https://en.wikipedia.org/wiki/Foreign_function_interface) (Foreign Function Interface).

## What is this?

- Builds Kagome as a C shared library (`.so`, `.dll`, or `.dylib`)
- Provides a C wrapper with four simple functions: `kagome_init()`, `kagome_destroy()`, `kagome_tokenize()`, `kagome_free_token_array()`
- Includes example scripts for Python (using `ctypes`) and PHP (using `FFI`)
- Each example loads the library, tokenizes Japanese text, and verifies the results

## Directory Structure

```shellsession
% tree -L 1 --dirsfirst -F
./
├── bin/               # Built shared libraries and archives (ignored in git)
├── c_wrapper/         # C wrapper source code for stable ABI exposure
├── go/                # Go source for the C ABI wrapper
├── php8/              # PHP 8 example using FFI
├── python3/           # Python 3 example using ctypes
│
├── .gitignore         # Git ignore file (ignores artifacts in 'bin/' and '__pycache__/')
├── docker-compose.yml # for Docker and Docker Compose
├── Dockerfile         # Dockerfile for building the shared library and running tests
├── go.mod             # Go module file
├── go.sum             # Go module checksum file
├── Makefile           # Makefile for building and testing
└── README.md          # This README file
```

## Requirements (for local run)

- Go 1.24 or later (required to build the archive and header files for `libkagome`)
- C compiler (required to build the `libkagome` shared library. E.g., `gcc`, `clang`)
- Python 3.10 or later for the Python example
- PHP 8 or later with the `FFI` extension for the PHP example

> **For Docker users**, you can run these examples without installing the above dependencies locally using Docker and Docker Compose.

## How to build and test locally

From this directory (`_examples/clib/`), run:

```sh
make clean
make build
make test
```

If the test passes, you will see `PASS` in the output.

## How to build and test with Docker and Docker Compose

From this directory (`_examples/clib/`), run:

```sh
make docker-build
make docker-test
```

If the test passes, you will see `PASS` in the output.

## Tested Environments

- macOS Tahoe 26.2 (M4 Mac mini, 16 GB RAM)
  - Go 1.25.6
  - C compiler: Apple clang 17.0.0
  - Python 3.14.2
  - PHP n/a (not tested locally)
  - GNU Make 3.81 (via Homebrew)
  - Docker 29.1.3 (via Docker Desktop)
  - Docker Compose 5.0.1
- Docker:
  - Go 1.25.6 (`golang:1.25.6` Docker image)
  - Python 3.14.2 (`python:3.14.2` Docker image)
    - C compiler: gcc (Debian 14.2.0-19) 14.2.0
  - PHP 8.5.2 with FFI enabled (`php:8.5.2-cli` Docker image as base)
    - C compiler: gcc (Debian 14.2.0-19) 14.2.0

## More Information

See the README in each subdirectory (`python3/README.md`, `php8/README.md`) for details and expected output.
