# Python FFI Example for Kagome Tokenizer

This directory contains a minimal example showing how to use `kagome` (a Japanese Morphological Analyzer/tokenizer written in Go) from Python via **C ABI and Python ctypes** including a simple test.

## What is this?

- A small Python script that calls Kagome through a shared library (`.so`, `.dll`, or `.dylib`) built from Go.
- Python accesses the library using [ctypes](https://docs.python.org/3/library/ctypes.html).
- The script tokenizes a Japanese sentence and prints each token.
- It also verifies the result against an expected output and prints `PASS` on success.

## Requirements

- Python 3.12 or later (tested with Python 3.12.9 and 3.14.2)
- Go and a C compiler (e.g., gcc) are required to build the `libkagome` shared library

> **For Docker users**, you can run this example without installing the above dependencies locally using Docker and Docker Compose.
> See the "How to run (with Docker and Docker Compose)" section below.

## How to run (locally)

1. Build the shared library from the example root (`_examples/clib/`):

    ```sh
    cd _examples/clib/
    make build
    ```

2. Run the script:

    ```sh
    python3 ./python3/main.py
    ```

## How to run (with Docker and Docker Compose)

1. Build the shared library builder via Docker from the example root (`_examples/clib/`):

    ```sh
    cd _examples/clib/
    make docker-build
    ```

2. Build and run the Python test container:

    ```sh
    docker compose run --rm --remove-orphans python-test
    ```

If the test passes, you will see `PASS` in the output.

## Expected Output

When you run the script, you should see output similar to the following:

```txt
Loading library: /app/bin/libkagome.so
surface=すもも, pos=['名詞', '一般', '*', '*'], base_form=すもも, conj_type=*, conj_form=*, reading=スモモ, pronunciation=スモモ, start=0, end=3
surface=も, pos=['助詞', '係助詞', '*', '*'], base_form=も, conj_type=*, conj_form=*, reading=モ, pronunciation=モ, start=3, end=4
surface=もも, pos=['名詞', '一般', '*', '*'], base_form=もも, conj_type=*, conj_form=*, reading=モモ, pronunciation=モモ, start=4, end=6
surface=も, pos=['助詞', '係助詞', '*', '*'], base_form=も, conj_type=*, conj_form=*, reading=モ, pronunciation=モ, start=6, end=7
surface=もも, pos=['名詞', '一般', '*', '*'], base_form=もも, conj_type=*, conj_form=*, reading=モモ, pronunciation=モモ, start=7, end=9
surface=の, pos=['助詞', '連体化', '*', '*'], base_form=の, conj_type=*, conj_form=*, reading=ノ, pronunciation=ノ, start=9, end=10
surface=うち, pos=['名詞', '非自立', '副詞可能', '*'], base_form=うち, conj_type=*, conj_form=*, reading=ウチ, pronunciation=ウチ, start=10, end=12
PASS
```
