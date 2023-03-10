SHELL=/usr/bin/env bash

build:
	go build -o kyc-checks ./cmd

run:
	go run ./cmd run

clean:
	go clean
	rm -f kyc-checks

fmt:
	go fmt ./...
	gofumpt -w .

lint:
	golangci-lint run

test:
	go test -p 4 -v ./...

.PHONY: build run clean test
