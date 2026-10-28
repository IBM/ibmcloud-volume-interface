
GOPACKAGES=$(shell go list ./... | grep -v /vendor/ | grep -v /samples)
GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")
ARCH = $(shell uname -m)
LINT_VERSION="2.12.2"

GOPATH := $(shell go env GOPATH)
LINT_BIN := $(GOPATH)/bin/golangci-lint

.PHONY: all
all: deps dofmt vet test

.PHONY: deps
deps:
	go get github.com/pierrre/gotestcover
	@if [ ! -x "$(LINT_BIN)" ] || ! "$(LINT_BIN)" version 2>/dev/null | grep -q "version $(LINT_VERSION)"; then \
		go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v$(LINT_VERSION); \
	fi

.PHONY: fmt
fmt:
	$(LINT_BIN) fmt --no-config --enable=gofmt --diff

.PHONY: dofmt
dofmt:
	$(LINT_BIN) fmt --no-config --enable=gofmt

.PHONY: lint
lint:
	$(LINT_BIN) run

.PHONY: makefmt
makefmt:
	gofmt -l -w ${GOFILES}

.PHONY: test
test:
ifeq ($(ARCH), ppc64le)
	# POWER
	$(GOPATH)/bin/gotestcover -v -coverprofile=cover.out ${GOPACKAGES} -timeout 90m
else
	# x86_64
	$(GOPATH)/bin/gotestcover -v -race -coverprofile=cover.out ${GOPACKAGES} -timeout 90m
endif

.PHONY: coverage
coverage:
	go tool cover -html=cover.out -o=cover.html

.PHONY: vet
vet:
	go vet ${GOPACKAGES}
