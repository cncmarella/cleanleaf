.PHONY: run test lint build docker deploy

VERSION ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo dev)

run:
	go run ./cmd/server

test:
	go test ./... -race

lint:
	go vet ./...
	gofmt -l .

build:
	CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$(VERSION)" -o bin/server ./cmd/server

docker:
	docker build --build-arg VERSION=$(VERSION) -t cleanleaf-api:$(VERSION) .

deploy:
	fly deploy --build-arg VERSION=$(VERSION)
