#!/bin/bash
rm -rf dist/
mkdir dist/
export CGO_ENABLED=0
GOOS=linux  GOARCH=amd64 sh -c 'go build -o dist/snk-$GOOS-$GOARCH'
GOOS=linux  GOARCH=arm   sh -c 'go build -o dist/snk-$GOOS-$GOARCH'
GOOS=darwin GOARCH=amd64 sh -c 'go build -o dist/snk-$GOOS-$GOARCH'
GOOS=darwin GOARCH=arm64 sh -c 'go build -o dist/snk-$GOOS-$GOARCH'
