#!/bin/bash

# Emulate CI environment 
export CI="home"

echo "Удаляем файл хранилища *******************************"
rm -rf ./data

# echo "Останавливаем базу данных, для эмуляции поведения GitHub CI ***********************"
# docker compose down

echo "Запускаем базу данных, для эмуляции поведения GitHub CI ***********************"
docker compose up -d
sleep 1

echo ; echo ; echo "Building shortenertest ---------------------------"
# go build -buildvcs=false -o cmd/shortener/shortener cmd/shortener/main.go
go build -o cmd/shortener/shortener cmd/shortener/main.go


echo ; echo ; echo "Code Increment #14 tests ------------------------"
shortenertestbeta-darwin-arm64 -test.v -test.run='^TestIteration14$' \
    -binary-path=cmd/shortener/shortener \
    -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable'
