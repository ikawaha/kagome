FROM ikawaha/kagome:latest AS base

FROM alpine:latest

COPY --from=base --chown=user:group /kagome /usr/local/bin/kagome

RUN apk add --no-cache graphviz

EXPOSE 6060

ENTRYPOINT ["/usr/local/bin/kagome"]
CMD ["server"]
