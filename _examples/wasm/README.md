# WebAssembly Example with Kagome

In this example we will demonstrate how to use Kagome in a WebAssembly application and show how responsive it can be.

- See: "[Kagome As a Server Side Tokenizer (Feeling Kagome Slow?)](https://github.com/ikawaha/kagome/wiki/Kagome-As-a-Server-Side-Tokenizer)" | Wiki | kagome @ GitHub

## How to Use

```sh
# Build the wasm binary
GOOS=js GOARCH=wasm go build -o kagome.wasm main.go

# Copy wasm_exec.js which maches to the compiled binary
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
**snip**
```

Now call the `wasm_exec.js` and `kagome.wasm` from the HTML file and run a web server.

- Online demo: [https://ikawaha.github.io/kagome/](https://ikawaha.github.io/kagome/)

```shellsession
├── docs                 ... gh-pages
│   ├── index.html
│   ├── kagome.wasm
│   └── wasm_exec.js
├── _examples
│   └── wasm
│       ├── README.md     ... this document
│       ├── kagome.html   ... html sample
│       ├── main.go       ... source code
│       ├── go.mod
│       └── go.sum
```
