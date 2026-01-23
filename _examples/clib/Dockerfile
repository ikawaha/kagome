FROM golang:latest AS build_base

RUN apt update && apt upgrade -y

WORKDIR /app

# This stage (layer) is intended to cache the Go module downloads.
#
# We do not `COPY` the local `go.mod` and `go.sum` here. The local `go.mod`
# contains `replace` directives which are not valid in Docker. So, we use
# the published kagome modules and create a temporary `go.mod` here.
RUN \
    go mod init kagome/examples/clib && \
    go get "github.com/ikawaha/kagome/v2" && \
    go get "github.com/ikawaha/kagome-dict/ipa" && \
    go mod download

FROM build_base AS build-and-test

COPY ./go/main.go /app/go/main.go
COPY ./go/main_test.go /app/go/main_test.go
COPY ./Makefile /app/Makefile
COPY ./c_wrapper/kagome_wrapper.h /app/c_wrapper/kagome_wrapper.h
COPY ./c_wrapper/kagome_wrapper.c /app/c_wrapper/kagome_wrapper.c

ENV CGO_ENABLED=1 CC=cc

RUN go mod tidy
RUN go clean -testcache && go test -race -v ./...

CMD ["make", "build"]
