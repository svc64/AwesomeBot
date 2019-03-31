#!/usr/bin/env bash
rm -rf out
NAME=awesomeBot
source config
mkdir -p out
dep ensure
BUILDCONFIG="-X main.token=$TOKEN -X main.DSN=$SENTRY"
go build -ldflags "${BUILDCONFIG}" -v -a -o out/${NAME}
