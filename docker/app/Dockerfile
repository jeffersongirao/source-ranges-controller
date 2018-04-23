FROM golang:1.10-alpine

WORKDIR /go/src/github.com/jeffersongirao/source-ranges-controller
COPY . .
RUN ./scripts/build.sh

FROM alpine:latest
RUN apk --no-cache add \
    ca-certificates
COPY --from=0 /go/src/github.com/jeffersongirao/source-ranges-controller/bin/linux/source-ranges-controller /usr/local/bin

ENTRYPOINT ["/usr/local/bin/source-ranges-controller"]