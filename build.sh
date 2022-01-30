#!/bin/bash
rm -rf dist/
mkdir dist/

# Create a universal binary that runs on M1 and older Macs.
# Create a binary that runs on 64-bit Linux, including Windows+WSL.
export CGO_ENABLED=0
GOOS=linux  GOARCH=amd64 go build -o dist/snk-linux
GOOS=darwin GOARCH=amd64 go build -o dist/snk-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o dist/snk-darwin-arm64
lipo -create -output dist/snk-mac dist/snk-darwin*
rm dist/snk-darwin*
sh -xc 'ls dist'
