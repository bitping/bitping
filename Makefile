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

build:
	go build ${LDFLAGS} -o $(CURR_DIR)/build/bin/$(BINARY)

.PHONY: $(PLATFORMS) build

$(PLATFORMS):
	GOOS=$(os) GOARCH=amd64 go build ${LDFLAGS} -o $(BIN_DIR)/$(BINARY)-$(VERSION)-$(os)-amd64


geth:
	geth --datadir /root/.ethereum \
		--ipcpath /root/.ethereum/geth.ipc \
		--syncmode "fast" --cache 512

