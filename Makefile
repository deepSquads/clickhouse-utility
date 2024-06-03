BINARY=main
PROJECT_NAME=clickhouse-utility
ENV ?= dev
LOG_LEVEL ?= debug
IMAGE_NAME ?= clickhouse-utility
IMAGE_TAG ?= latest
GOLANGCI_VERSION ?= v1.58.0

all: clean lint build

TEST ?= ...

clean:
	@echo "--> Target directory clean up"
	rm -rf ./.build/target
	rm -f ${BINARY}

lint:
	@echo "--> Running linters"
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCI_VERSION} run -c .golangci.yml

test:
	gotestsum --junitfile test_reports/unit-tests.xml -- -race $(shell go list ./...) -count=1

build:
	go build -o ${BINARY} ./cmd

build-docker:
	./.build/build_docker.sh ./cmd ${BINARY} ${IMAGE_NAME} ${IMAGE_TAG}

run-compose:
	docker-compose -f docker-compose.yml up
