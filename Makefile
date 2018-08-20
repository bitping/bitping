BINARY := bitping
VERSION := $(shell cat VERSION)
COMMIT := $(shell git rev-parse HEAD)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
BIN_DIR := $(shell pwd)/build
CURR_DIR := $(shell pwd)

PKGS := $(shell go list ./... | grep -v vendor)

COMMIT = $(shell git rev-parse HEAD | cut -c 1-6)
BUILD_TIME = $(shell date -u '+%Y-%m-%dT%I:%M:%S%p')

PLATFORMS := linux darwin
os = $(word 1, $@)

COM_COLOR   = \033[0;34m
OBJ_COLOR   = \033[0;36m
OK_COLOR    = \033[0;32m
ERROR_COLOR = \033[0;31m
WARN_COLOR  = \033[0;33m
NO_COLOR    = \033[m

BUILD_COM = "Building..."

LDFLAGS =-ldflags "-X github.com/auser/bitping/cmd.AppName=$(BINARY) -X github.com/auser/bitping/cmd.Branch=$(BRANCH) -X github.com/auser/bitping/cmd.Version=$(VERSION) -X github.com/auser/bitping/cmd.Commit=$(COMMIT) -X github.com/auser/bitping/cmd.BuildTime=$(BUILD_TIME)"

# LDFLAGS=-ldflags='-X "cmd.AppName=${BINARY}" -X "cmd.Version=${VERSION}" -X "cmd.Commit=${COMMIT}" -X "cmd.Branch=${BRANCH}" -X "cmd.BuildTime=$(shell date +%FT%T%Z)"'

deps:
	dep ensure
	make deps_first_time

deps_first_time:
	go get -u github.com/ethereum/go-ethereum
	cp -r \
  "${GOPATH}/src/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1" \
  "vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/"

build:
	@printf "%b\n" "$(COM_COLOR)$(BUILD_COM)$(OBJ_COLOR)$(NO_COLOR)\n";
	@go build ${LDFLAGS} -o $(CURR_DIR)/build/bin/$(BINARY)

.PHONY: $(PLATFORMS) build

$(PLATFORMS):
	GOOS=$(os) GOARCH=amd64 go build ${LDFLAGS} -o $(BIN_DIR)/$(BINARY)-$(VERSION)-$(os)-amd64


geth:
	geth --datadir .ethereum \
		--ipcpath .ethereum/geth.ipc \
		--syncmode "fast" --cache 512

test:
	go test ./... -v -ginkgo.v -ginkgo.progress
