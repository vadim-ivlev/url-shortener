#!/bin/bash

docker compose up -d
go run cmd/shortener/main.go $@
