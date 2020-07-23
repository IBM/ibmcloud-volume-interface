
GOPACKAGES=$(shell go list ./... | grep -v /vendor/ | grep -v /samples)
GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")
GOLINTPACKAGES=$(shell go list ./... | grep -v /vendor/)
ARCH = $(shell uname -m)

.PHONY: all
all: deps dofmt vet test

.PHONY: deps
deps:
	go get github.com/pierrre/gotestcover
	go get golang.org/x/lint/golint

.PHONY: fmt
fmt: lint
	gofmt -l ${GOFILES}
	@if [ -n "$$(gofmt -l ${GOFILES})" ]; then echo 'Above Files needs gofmt fixes. Please run gofmt -l -w on your code.' && exit 1; fi

.PHONY: dofmt
dofmt:
	go fmt ./...

.PHONY: lint
lint:
	$(GOPATH)/bin/golint --set_exit_status ${GOLINTPACKAGES}

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
