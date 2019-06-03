# This is a Docker image of Kagome on Alpine Linux.
# =================================================
# Aimed to check the functionality of version command but the built image will
# work as a single 'kagome' binary.
# - USAGE:
#   - To build the image, run:
#     $ docker build --tag kagome:latest ./Dockerfile && docker image prune -f
#   - To run the container, run:
#     $ # This is equivalent to "kagome version"
#     $ docker run --rm kagome:latest version

# 1st stage: Build binary.
# ------------------------
FROM golang:alpine AS build-app

# Copy the current dir including .git dir for versioning.
COPY . /go/src/github.com/ikawaha/kagome
WORKDIR /go/src/github.com/ikawaha/kagome

# Shell script to build the image (with the tag as version).
RUN apk --no-cache add git && \
    version_app=$(git describe --tag) && \
    echo "- Current git tag: ${version_app}" && \
    go version && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
      -a \
      -installsuffix cgo \
      --ldflags "-w -s -extldflags \"-static\" -X 'main.version=${version_app}'" \
      -o /go/bin/kagome \
      ./cmd/kagome && \
    echo '- Running tests ...' && \
    /go/bin/kagome version

# 2nd stage: Copy only the built binary to shrink the image size.
# ---------------------------------------------------------------
FROM alpine:latest

COPY --from=build-app /go/bin/kagome /usr/local/bin/kagome

ENTRYPOINT [ "kagome" ]
