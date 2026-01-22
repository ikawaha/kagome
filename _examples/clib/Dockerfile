FROM golang:latest AS builder

WORKDIR /app

COPY ./go/main.go /app/go/main.go
COPY ./Makefile /app/Makefile
COPY ./c_wrapper/kagome_wrapper.c /app/c_wrapper/kagome_wrapper.c

ENV CGO_ENABLED=1 CC=cc

RUN \
    go mod init example.com/kagome/clib && \
    go get "github.com/ikawaha/kagome/v2" && \
    go get "github.com/ikawaha/kagome-dict/ipa" && \
    go mod tidy

CMD ["make", "build"]
