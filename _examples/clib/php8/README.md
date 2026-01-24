# PHP FFI Example for Kagome Tokenizer

This directory shows how to use Kagome (a Japanese Morphological Analyzer written in Go) from PHP via FFI.

## What is this?

- PHP script that calls the Kagome shared library using [PHP FFI](https://www.php.net/manual/en/class.ffi.php)
- Uses the `kagome_*` C wrapper functions for a stable API
- Tokenizes Japanese text and prints each token
- Includes a test that verifies the output

## Requirements (for local run)

- PHP 8 or later (tested with PHP 8.5.2)
- PHP extension `FFI` enabled
- Go 1.24 or later (to build the archive and header files for `libkagome`)
- C compiler (required to build the `libkagome` shared library. E.g., `gcc`)

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
    php ./php8/main.php
    ```

## How to run (with Docker and Docker Compose)

1. Build the shared library builder via Docker from the example root (`_examples/clib/`):

    ```sh
    cd _examples/clib/
    make docker-build
    ```

2. Build and run the PHP test container:

    ```sh
    docker compose run --rm --remove-orphans php-test
    ```

If the test passes, you will see `PASS` in the output.

## Expected Output

When you run the script, you should see output similar to the following:

```txt
surface=すもも, pos=[名詞, 一般, *, *], base_form=すもも, conj_type=*, conj_form=*, reading=スモモ, pronunciation=スモモ, start=0, end=3
surface=も, pos=[助詞, 係助詞, *, *], base_form=も, conj_type=*, conj_form=*, reading=モ, pronunciation=モ, start=3, end=4
surface=もも, pos=[名詞, 一般, *, *], base_form=もも, conj_type=*, conj_form=*, reading=モモ, pronunciation=モモ, start=4, end=6
surface=も, pos=[助詞, 係助詞, *, *], base_form=も, conj_type=*, conj_form=*, reading=モ, pronunciation=モ, start=6, end=7
surface=もも, pos=[名詞, 一般, *, *], base_form=もも, conj_type=*, conj_form=*, reading=モモ, pronunciation=モモ, start=7, end=9
surface=の, pos=[助詞, 連体化, *, *], base_form=の, conj_type=*, conj_form=*, reading=ノ, pronunciation=ノ, start=9, end=10
surface=うち, pos=[名詞, 非自立, 副詞可能, *], base_form=うち, conj_type=*, conj_form=*, reading=ウチ, pronunciation=ウチ, start=10, end=12
PASS
```
