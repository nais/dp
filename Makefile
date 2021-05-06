.PHONY: test
DATE = $(shell date "+%Y-%m-%d")
LAST_COMMIT = $(shell git --no-pager log -1 --pretty=%h)
VERSION ?= $(DATE)-$(LAST_COMMIT)
LDFLAGS := -X github.com/nais/dp/pkg/version.Revision=$(shell git rev-parse --short HEAD) -X github.com/nais/dp/pkg/version.Version=$(VERSION)
PKGID = io.nais.dp
GOPATH ?= ~/go

test:
	go test ./... -count=1

local:
	go run ./... --development-mode=true --bind-address=127.0.0.1:8080