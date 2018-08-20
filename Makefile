BINARY := bitping
VERSION := 0.0.1
COMMIT := $(shell git rev-parse HEAD)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
BIN_DIR := $(shell pwd)/build
CURR_DIR := $(shell pwd)

PKGS := $(shell go list ./... | grep -v vendor)

PLATFORMS := linux darwin
os = $(word 1, $@)

LDFLAGS = -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.branch=${BRANCH}"

deps:
	dep ensure

deps_first_time:
	go get -u github.com/ethereum/go-ethereum
	cp -r \
  "${GOPATH}/src/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1" \
  "vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/"

build:
	go build ${LDFLAGS} -o $(CURR_DIR)/build/bin/$(BINARY)

.PHONY: $(PLATFORMS) build

$(PLATFORMS):
	GOOS=$(os) GOARCH=amd64 go build ${LDFLAGS} -o $(BIN_DIR)/$(BINARY)-$(VERSION)-$(os)-amd64


geth:
	geth --datadir .ethereum \
		--ipcpath .ethereum/geth.ipc \
		--syncmode "fast" --cache 512

test:
	go test ./... -v -ginkgo.v -ginkgo.progress
