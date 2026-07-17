VERSION=0.0.4
GITCOMMIT?=$(shell git describe --dirty --always)
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION} -X main.commit=${GITCOMMIT}"

all: mackerel-plugin-command-status

.PHONY: mackerel-plugin-command-status

mackerel-plugin-command-status: main.go
	go build $(LDFLAGS) -o mackerel-plugin-command-status main.go

linux: main.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o mackerel-plugin-command-status main.go

check:
	go test -v ./...
