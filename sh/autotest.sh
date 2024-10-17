#!/bin/bash

# Emulate CI environment 
export CI="home"

echo "Удаляем файл хранилища *******************************"
rm -rf ./data


echo "Останавливаем базу данных, для эмуляции поведения GitHub CI ***********************"
docker compose down



echo ; echo ; echo "Building shortenertest ---------------------------"
# go build -buildvcs=false -o cmd/shortener/shortener cmd/shortener/main.go
go build -o cmd/shortener/shortener cmd/shortener/main.go

echo ; echo ; echo "Running static tests -----------------------------"
go vet -vettool=statictest-darwin-arm64 ./...

echo ; echo ; echo "Code Increment #1 tests ------------------------"
shortenertestbeta-darwin-arm64 -test.v -test.run=^TestIteration1$ -binary-path=cmd/shortener/shortener

echo ; echo ; echo "Code Increment #2 tests ------------------------"
# shortenertest -test.v -test.run=^TestIteration2$ -binary-path=cmd/shortener/shortener -source-path=./internal
shortenertestbeta-darwin-arm64 -test.v -test.run=^TestIteration2$ -source-path=.

echo ; echo ; echo "Code Increment #3 tests ------------------------"
shortenertestbeta-darwin-arm64 -test.v -test.run=^TestIteration3$ -source-path=.

echo ; echo ; echo "Code Increment #4 tests ------------------------"
# SERVER_PORT=$(random unused-port)
SERVER_PORT=8082
echo "SERVER_PORT=$SERVER_PORT"
shortenertestbeta-darwin-arm64 -test.v -test.run=^TestIteration4$ -binary-path=cmd/shortener/shortener -server-port=$SERVER_PORT 

echo ; echo ; echo "Code Increment #5 tests ------------------------"
SERVER_PORT=8082
echo "SERVER_PORT=$SERVER_PORT"
shortenertestbeta-darwin-arm64 -test.v -test.run=^TestIteration5$ -binary-path=cmd/shortener/shortener -server-port=$SERVER_PORT

echo ; echo ; echo "Code Increment #6 tests ------------------------"
shortenertestbeta-darwin-arm64 -test.v -test.run=^TestIteration6$ -source-path=.

echo ; echo ; echo "Code Increment #7 tests ------------------------"
shortenertestbeta-darwin-arm64 -test.v -test.run=^TestIteration7$ -binary-path=cmd/shortener/shortener -source-path=.

echo ; echo ; echo "Code Increment #8 tests ------------------------"
shortenertestbeta-darwin-arm64 -test.v -test.run=^TestIteration8$ -binary-path=cmd/shortener/shortener

echo ; echo ; echo "Code Increment #9 tests ------------------------"
# TEMP_FILE=$(random tempfile)
# tempfile is not available in macos
TEMP_FILE=$(mktemp -p /tmp)
echo "TEMP_FILE=$TEMP_FILE"
shortenertestbeta-darwin-arm64 -test.v -test.run=^TestIteration9$ -binary-path=cmd/shortener/shortener -source-path=. -file-storage-path=$TEMP_FILE


echo "Запускаем базу данных, для эмуляции поведения GitHub CI ***********************"
docker compose up -d
sleep 2


echo ; echo ; echo "Code Increment #10 tests ------------------------"
shortenertestbeta-darwin-arm64 -test.v -test.run=^TestIteration10$ \
    -binary-path=cmd/shortener/shortener \
    -source-path=. \
    -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable'

echo ; echo ; echo "Code Increment #11 tests ------------------------"
shortenertestbeta-darwin-arm64 -test.v -test.run=^TestIteration11$ \
    -binary-path=cmd/shortener/shortener \
    -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable'

echo ; echo ; echo "Code Increment #12 tests ------------------------"
shortenertestbeta-darwin-arm64 -test.v -test.run=^TestIteration12$ \
    -binary-path=cmd/shortener/shortener \
    -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable'

echo ; echo ; echo "Code Increment #13 tests ------------------------"
shortenertestbeta-darwin-arm64 -test.v -test.run=^TestIteration13$ \
    -binary-path=cmd/shortener/shortener \
    -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable'
