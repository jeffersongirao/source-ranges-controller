#!/bin/sh

export CGO_ENABLED=0
export GOOS=linux

go build -ldflags '-extldflags "-static"' -a -v -o bin/linux/source-ranges-controller ./cmd/