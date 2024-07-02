#!/bin/bash

go build -o cmd/shortener/shortener cmd/shortener/main.go
shortenertest -test.v -test.run=^TestIteration1$ -binary-path=cmd/shortener/shortener
# shortenertest -test.v -test.run=^TestIteration2$ -binary-path=cmd/shortener/shortener -source-path=./internal
shortenertest -test.v -test.run=^TestIteration2$ -source-path=.

# go build -buildvcs=false -o cmd/shortener/shortener cmd/shortener/main.go
# shortenertestbeta-darwin-arm64 -test.v -test.run=^TestIteration1$ -binary-path=cmd/shortener/shortener
shortenertestbeta-darwin-arm64 -test.v -test.run=^TestIteration2$ -source-path=.