#!/usr/bin/env bash
NAME=awesomeBot
source token
mkdir -p out
go get ./...
BUILDCONFIG="-X main.token=$TOKEN"
CGO_ENABLED=0 go build -ldflags "${BUILDCONFIG}" -v -a -o out/${NAME}
