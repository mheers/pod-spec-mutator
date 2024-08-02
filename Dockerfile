FROM golang:1.22-alpine AS builder

RUN apk add --update \
    make \
    git \
    && rm -rf /var/cache/apk/*
RUN mkdir -p /go/src/github.com/mheers/pod-spec-mutator
ADD . /go/src/github.com/mheers/pod-spec-mutator
WORKDIR /go/src/github.com/mheers/pod-spec-mutator
RUN make build

FROM alpine:3.19

COPY --from=builder /go/src/github.com/mheers/pod-spec-mutator/bin/pod-spec-mutator /pod-spec-mutator
ENTRYPOINT ["/pod-spec-mutator"]
