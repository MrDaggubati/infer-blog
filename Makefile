SHELL := /bin/bash

CLI := go run ./cmd/blogbuild

BIN := infer-blog
PREFIX ?= /usr/local

VERSION ?= 0.1.0
LDFLAGS := -X main.version=$(VERSION)

PORT ?= 8080
HOST ?= localhost

.PHONY: help deps check-deps tidy fmt build rebuild clean clean-cache clean-all serve cli install uninstall version

help:
	@$(CLI) --help

deps:
	@$(CLI) check

check-deps:
	@$(CLI) check

tidy:
	@go mod tidy

fmt:
	@gofmt -w cmd/blogbuild/*.go

build:
	@$(CLI) build

rebuild:
	@$(CLI) rebuild

clean:
	@$(CLI) clean

clean-cache:
	@$(CLI) clean-cache

clean-all:
	@$(CLI) clean-all

serve:
	@$(CLI) serve --host "$(HOST)" --port "$(PORT)"

cli:
	@go build \
		-ldflags "$(LDFLAGS)" \
		-o "$(BIN)" \
		./cmd/blogbuild
	@printf 'built %s %s\n' "$(BIN)" "$(VERSION)"

install:
	@go build \
		-ldflags "$(LDFLAGS)" \
		-o "$(BIN)" \
		./cmd/blogbuild
	@install -m 0755 "$(BIN)" "$(PREFIX)/bin/$(BIN)"
	@printf 'installed %s %s -> %s/bin/%s\n' \
		"$(BIN)" \
		"$(VERSION)" \
		"$(PREFIX)" \
		"$(BIN)"

uninstall:
	@rm -f "$(PREFIX)/bin/$(BIN)"
	@printf 'removed %s/bin/%s\n' "$(PREFIX)" "$(BIN)"

version:
	@printf '%s %s\n' "$(BIN)" "$(VERSION)"