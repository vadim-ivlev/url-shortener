#!/bin/bash

docker compose up -d

export DATABASE_DSN="postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable"
go run cmd/shortener/main.go $@
