#!/bin/bash

echo "Building shortenertest ---------------------------"
go build -o cmd/shortener/shortener cmd/shortener/main.go

echo "Running static tests -----------------------------"
go vet -vettool=statictest-darwin-arm64 ./...

echo "Running Iteration 1 tests ------------------------"
shortenertest -test.v -test.run=^TestIteration1$ -binary-path=cmd/shortener/shortener

echo "Running Iteration 2 tests ------------------------"
# shortenertest -test.v -test.run=^TestIteration2$ -binary-path=cmd/shortener/shortener -source-path=./internal
shortenertest -test.v -test.run=^TestIteration2$ -source-path=.

# go build -buildvcs=false -o cmd/shortener/shortener cmd/shortener/main.go
# shortenertestbeta-darwin-arm64 -test.v -test.run=^TestIteration1$ -binary-path=cmd/shortener/shortener
# shortenertestbeta-darwin-arm64 -test.v -test.run=^TestIteration2$ -source-path=.