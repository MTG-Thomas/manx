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
secret-scan: secret-gate-python
	"$(SECRET_GATE_PYTHON)" -I scripts/secret_gate.py

pre-push: secret-gate-python
	"$(SECRET_GATE_PYTHON)" -I scripts/secret_gate.py --pre-push

# Native Windows: set MTG_SECRET_GATE_PYTHON to an approved absolute interpreter.
PYTHON ?= python3

.PHONY: secret-gate-python
ifeq ($(OS),Windows_NT)
SECRET_GATE_PYTHON ?= $(MTG_SECRET_GATE_PYTHON)
else
SECRET_GATE_PYTHON ?= $(PYTHON)
endif

secret-gate-python:
ifeq ($(OS),Windows_NT)
	$(if $(or $(findstring :/,$(subst \,/,$(SECRET_GATE_PYTHON))),$(filter //%,$(subst \,/,$(SECRET_GATE_PYTHON)))),,$(error UNAVAILABLE: set MTG_SECRET_GATE_PYTHON to an approved absolute Python path outside all Git checkouts; publication blocked))
endif
