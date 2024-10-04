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


