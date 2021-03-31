WebAssembly Sample
---

```
GOOS=js GOARCH=wasm go build -o kagome.wasm main.go
```


```
├── docs                 ... gh-pages
│   ├── index.html
│   ├── kagome.wasm
│   └── wasm_exec.js
├── sample
│   └── wasm
│       ├── README.md     ... this document.
│       ├── go.mod
│       ├── kagome.html   ... html sample
│       └── main.go       ... 
```
demo. https://ikawaha.github.io/kagome/