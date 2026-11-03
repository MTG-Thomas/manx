SHELL := /bin/bash
GOFLAGS ?= -mod=mod

all: check build

build:
	@mkdir -p bin
	GOFLAGS=$(GOFLAGS) go build -o bin/manx ./cmd/manx
	GOFLAGS=$(GOFLAGS) go build -o bin/manx-tui ./cmd/manx-tui

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w cmd internal

fmt-check:
	test -z "$$(gofmt -l cmd internal)" || (gofmt -l cmd internal && exit 1)

vulncheck:
	govulncheck ./...

check: fmt-check vet test

clean:
	rm -rf bin

.PHONY: secret-scan pre-push
# Local authenticated gate; keep CI/check/test targets unchanged.
secret-scan:
	$(PYTHON) -I scripts/secret_gate.py

pre-push:
	$(PYTHON) -I scripts/secret_gate.py --pre-push

# Override with PYTHON=python on native Windows.
PYTHON ?= python3
