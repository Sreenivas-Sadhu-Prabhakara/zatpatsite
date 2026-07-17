.PHONY: run build test

run:
	go run ./cmd/server

build:
	mkdir -p bin
	go build -o bin/zatpatsite ./cmd/server

test:
	go test ./...
