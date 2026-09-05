VERSION=0.0.7
GITCOMMIT?=$(shell git describe --dirty --always)
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION} -X main.commit=${GITCOMMIT}"

all: mackerel-plugin-command-status

.PHONY: mackerel-plugin-command-status linux check lint

mackerel-plugin-command-status: *.go
	go build $(LDFLAGS) -o mackerel-plugin-command-status

linux: *.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o mackerel-plugin-command-status

check:
	go test -v ./...

lint:
	golangci-lint run ./...