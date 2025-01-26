.PHONY: build test lint clean run-test test-integration test-all

VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT  ?= $(shell git rev-parse --short HEAD)
DATE    ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -X github.com/ensync-cli/pkg/version.version=$(VERSION) \
           -X github.com/ensync-cli/pkg/version.commit=$(COMMIT) \
           -X github.com/ensync-cli/pkg/version.buildDate=$(DATE)

build:
	go build -ldflags "$(LDFLAGS)" -o bin/ensync

test:
	go test -v -race ./...

test-integration:
	go test -v -race ./test/integration/...

test-all: test test-integration

lint:
	golangci-lint run

clean:
	rm -rf bin/