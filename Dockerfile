# This is a Docker image of Kagome on Alpine Linux.
# =================================================

# 1st stage: Build binary.
# ------------------------
FROM --platform=$BUILDPLATFORM golang:alpine AS build-app

ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG GOOS=linux
ARG GOARCH=amd64
ARG GOARM=

# Copy the current dir including .git dir for versioning.
COPY . /go/src/github.com/ikawaha/kagome
WORKDIR /go/src/github.com/ikawaha/kagome

# Shell script to build the image (with the tag as version).
RUN \
  apk --no-cache add git && \
  version_app=$(git describe --tag) && \
  echo "- Running on ${BUILDPLATFORM}, building for ${TARGETPLATFORM}" && \
  echo "- Current platform: $(uname -a)" && \
  echo "- Current git tag: ${version_app}" && \
  echo "- Current Go version: $(go version)" && \
  echo "- Kagome version to be build:${version_app}" && \
  CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} GOARM=${GOARM} go build \
    --ldflags "-w -s -extldflags \"-static\" -X 'main.version=${version_app}'" \
    -o /go/bin/kagome \
    /go/src/github.com/ikawaha/kagome && \
  echo "- Smoke test (run kagome version command) ... $(/go/bin/kagome version)"

# 2nd stage: Copy only the built binary to shrink the image size.
# ---------------------------------------------------------------
FROM --platform=$BUILDPLATFORM alpine:latest

COPY --from=build-app /go/bin/kagome /usr/local/bin/kagome

ENTRYPOINT [ "kagome" ]
