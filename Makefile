.PHONY: test
DATE = $(shell date "+%Y-%m-%d")
LAST_COMMIT = $(shell git --no-pager log -1 --pretty=%h)
VERSION ?= $(DATE)-$(LAST_COMMIT)
LDFLAGS := -X github.com/nais/dp/pkg/version.Revision=$(shell git rev-parse --short HEAD) -X github.com/nais/dp/pkg/version.Version=$(VERSION)
PKGID = io.nais.dp
GOPATH ?= ~/go

test:
	go test ./... -count=1

local-with-auth:
	go run cmd/backend/main.go --oauth2-client-secret=$(gcloud secrets versions access --secret dp-oauth2-client-secret latest --project aura-dev-d9f5) --oauth2-client-id=791e3efd-28d6-4150-9978-20a37c340e7f --oauth2-tenant-id=62366534-1ec3-4962-8869-9b5535279d0b --bind-address=127.0.0.1:8080 --firestore-google-project-id=aura-dev-d9f5 --firestore-collection=dp

local:
	go run cmd/backend/main.go --development-mode=true --bind-address=127.0.0.1:8080 --firestore-google-project-id=aura-dev-d9f5 --firestore-collection=dp
