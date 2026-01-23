FROM golang:latest AS builder

RUN apt update && apt upgrade -y

WORKDIR /app

# Copy dependency files FIRST (for Docker layer caching)
COPY go.mod go.sum ./
RUN go mod download

FROM builder AS build-and-test

COPY ./go/main.go /app/go/main.go
COPY ./go/main_test.go /app/go/main_test.go
COPY ./Makefile /app/Makefile
COPY ./c_wrapper/kagome_wrapper.h /app/c_wrapper/kagome_wrapper.h
COPY ./c_wrapper/kagome_wrapper.c /app/c_wrapper/kagome_wrapper.c

ENV CGO_ENABLED=1 CC=cc

RUN go clean -testcache && go test -race -v ./...

CMD ["make", "build"]
