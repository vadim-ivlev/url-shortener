#!/bin/bash

echo "Looking for dead code ---------------------------"
# go install github.com/deadcode/deadcode@latest
deadcode ./...

echo "Test coverage -----------------------------------"
go test -cover ./...
# go test -coverprofile=cover.out ./...
# go tool cover -html=cover.out

echo "Running tests -----------------------------------"
go test -count=1 ./...


echo "Building shortenertest ---------------------------"
go build -o cmd/shortener/shortener cmd/shortener/main.go

echo "Running static tests -----------------------------"
go vet -vettool=statictest-darwin-arm64 ./...

echo "Running Iteration 1 tests ------------------------"
shortenertest -test.v -test.run=^TestIteration1$ -binary-path=cmd/shortener/shortener

echo "Running Iteration 2 tests ------------------------"
# shortenertest -test.v -test.run=^TestIteration2$ -binary-path=cmd/shortener/shortener -source-path=./internal
shortenertest -test.v -test.run=^TestIteration2$ -source-path=.

echo "Running Iteration 3 tests ------------------------"
shortenertest -test.v -test.run=^TestIteration3$ -source-path=.


# go build -buildvcs=false -o cmd/shortener/shortener cmd/shortener/main.go
# shortenertestbeta-darwin-arm64 -test.v -test.run=^TestIteration1$ -binary-path=cmd/shortener/shortener
# shortenertestbeta-darwin-arm64 -test.v -test.run=^TestIteration2$ -source-path=.