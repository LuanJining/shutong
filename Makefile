SHELL := /bin/bash
ROOT_DIR := $(shell pwd)
GOCACHE := $(ROOT_DIR)/.gocache

export GOCACHE

.PHONY: build test tidy docker-iam docker-workflow docker-compose

build:
	@echo "==> building binaries"
	@mkdir -p $(GOCACHE)
	@GOCACHE=$(GOCACHE) go build ./...

test:
	@echo "==> running tests"
	@mkdir -p $(GOCACHE)
	@GOCACHE=$(GOCACHE) go test ./...

tidy:
	@go mod tidy

docker-iam:
	docker build -f deploy/docker/Dockerfile.iam -t iam-service:dev .

docker-workflow:
	docker build -f deploy/docker/Dockerfile.workflow -t workflow-service:dev .

docker-compose:
	docker compose -f deploy/docker-compose.yml up --build
