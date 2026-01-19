FROM golang:latest AS builder

WORKDIR /app

COPY ./go/main.go /app/go/main.go
COPY ./Makefile /app/Makefile
COPY ./c/kagome_wrapper.c /app/c/kagome_wrapper.c

ENV CGO_ENABLED=1 CC=cc

RUN \
    go mod init example.com/kagome/clib && \
    go get "github.com/ikawaha/kagome/v2" && \
    go get "github.com/ikawaha/kagome-dict/ipa" && \
    go mod tidy

RUN make build
RUN ls /app/bin

FROM python:latest

WORKDIR /app

COPY --from=builder /app/bin /app/bin
COPY ./python3 /app/python3
COPY ./Makefile /app/Makefile

ENTRYPOINT ["make", "test"]
