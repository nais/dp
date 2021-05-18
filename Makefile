.PHONY: test
DATE = $(shell date "+%Y-%m-%d")
LAST_COMMIT = $(shell git --no-pager log -1 --pretty=%h)
VERSION ?= $(DATE)-$(LAST_COMMIT)
LDFLAGS := -X github.com/nais/dp/backend/version.Revision=$(shell git rev-parse --short HEAD) -X github.com/nais/dp/backend/version.Version=$(VERSION)
PKGID = io.nais.dp
GOPATH ?= ~/go
APP = dp

test:
	go test ./... -count=1

local-with-auth:
	go run cmd/backend/main.go \
	--teams-url=https://raw.githubusercontent.com/navikt/teams/main/teams.json \
	--oauth2-client-secret=$(shell gcloud secrets versions access --secret dp-oauth2-client-secret latest --project aura-dev-d9f5) \
	--teams-token=$(shell gcloud secrets versions access --secret github-read-token latest --project aura-dev-d9f5 | cut -d= -f2) \
	--oauth2-client-id=$(shell gcloud secrets versions access --secret dp-oauth2-client-id latest --project aura-dev-d9f5) \
	--oauth2-tenant-id=62366534-1ec3-4962-8869-9b5535279d0b \
	--bind-address=127.0.0.1:8080 \
	--firestore-google-project-id=aura-dev-d9f5 \
	--firestore-collection=dp \
	--hostname=localhost \
	--log-level=debug \
	--state=$(shell gcloud secrets versions access --secret dp-state latest --project aura-dev-d9f5 | cut -d= -f2)

local:
	go run cmd/backend/main.go \
	--teams-url=https://raw.githubusercontent.com/navikt/teams/main/teams.json \
	--teams-token=$(shell gcloud secrets versions access --secret github-read-token latest --project aura-dev-d9f5 | cut -d= -f2) \
	--development-mode=true \
	--bind-address=127.0.0.1:8080 \
	--firestore-google-project-id=aura-dev-d9f5 \
	--firestore-collection=dp \
	--hostname=localhost \
	--log-level=debug \
	--state=$(shell gcloud secrets versions access --secret dp-state latest --project aura-dev-d9f5 | cut -d= -f2)

linux-build:
	go build -a -installsuffix cgo -o $(APP) -ldflags "-s $(LDFLAGS)" cmd/backend/main.go

docker-build:
	docker image build -t ghcr.io/nais/$(APP):$(VERSION) -t ghcr.io/nais/$(APP):latest .

docker-push:
	docker image push ghcr.io/nais/$(APP):$(VERSION)
	docker image push ghcr.io/nais/$(APP):latest

