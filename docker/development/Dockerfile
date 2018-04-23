FROM golang:1.10-alpine

RUN apk --no-cache add \
    bash \
    git \
    g++ \
    openssl

# Go dep installation
RUN go get -u github.com/golang/dep/cmd/dep \
    && mkdir -p /go/pkg/dep \
    && chmod 777 /go/pkg/dep

# Create user
ARG uid=1000
ARG gid=1000
RUN addgroup -g $gid builder && \
    adduser -D -u $uid -G builder builder && \
    chown builder:builder -R /go

USER builder
WORKDIR /go/src/github.com/jeffersongirao/source-ranges-controller