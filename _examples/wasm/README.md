# WebAssembly Example of Kagome

- Build

```sh
GOOS=js GOARCH=wasm go build -o kagome.wasm main.go
```

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

- Online demo: [https://ikawaha.github.io/kagome/](https://ikawaha.github.io/kagome/)
