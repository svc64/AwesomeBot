#!/usr/bin/env bash
NAME=awesomeBot
source token
mkdir -p out
dep ensure
BUILDCONFIG="-X main.token=$TOKEN"
go build -ldflags "${BUILDCONFIG}" -v -a -o out/${NAME}
