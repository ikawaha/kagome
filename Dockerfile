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
RUN apk add git gcc g++ && \
    version_app=$(git describe --tag) && \
    echo "- Current git tag: ${version_app}" && \
    cd /go/src/github.com/ikawaha/kagome/cmd && \
    go build \
      -a \
      --ldflags "-w -extldflags \"-static\" -X 'main.versionKagome=${version_app}'" \
      -o /go/bin/kagome \
      ./kagome && \
    echo '- Running tests ...' && \
    /go/bin/kagome -v && \
    /go/bin/kagome --version && \
    /go/bin/kagome version

# 2nd stage: Copy only the built binary to shrink the image size.
# ---------------------------------------------------------------
FROM alpine:latest

COPY --from=build-app /go/bin/kagome /usr/local/bin/kagome

ENTRYPOINT [ "kagome" ]
