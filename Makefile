# Makefile for Kelly betting calculator

BINARY_NAME=kelly
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

.PHONY: build
build:
	go build ${LDFLAGS} -o ${BINARY_NAME} .

.PHONY: install
install:
	go install ${LDFLAGS} .

.PHONY: test
test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: test-short
test-short:
	go test -short ./...

.PHONY: lint
lint:
	go vet ./...
	gofmt -s -w .
	go mod tidy

.PHONY: clean
clean:
	go clean
	rm -f ${BINARY_NAME}
	rm -f coverage.out coverage.html

.PHONY: run
run:
	go run .

.PHONY: run-cli
run-cli:
	go run . -a 2.56 -b 3.85 -t 10000

.PHONY: deps
deps:
	go mod download
	go mod verify

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build       - Build the binary"
	@echo "  install     - Install to GOPATH/bin"
	@echo "  test        - Run tests with coverage"
	@echo "  test-short  - Run quick tests"
	@echo "  lint        - Run linters and formatters"
	@echo "  clean       - Remove build artifacts"
	@echo "  run         - Run in interactive mode"
	@echo "  run-cli     - Run with example CLI args"
	@echo "  deps        - Download dependencies"
