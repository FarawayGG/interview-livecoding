export SERVICE_NAME := wisdom
export PROJECT_NAME := wisdom

export GOBIN := $(PWD)/bin
export GOPRIVATE := github.com/farawaygg/*
export GONOSUMDB := github.com/farawaygg
export PATH := $(GOBIN):$(PATH)

SHELL := env PATH=$(PATH) /bin/sh
PROTOC_VERSION := 3.18.0
MIGRATE_VERSION := 4.15.2
PROTOC_VALIDATE_VERSION := 0.6.2
GOLANGLINT_VERSION := 1.51.2

GOFLAGS ?=

UNAME_OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
UNAME_ARCH := $(shell uname -m)

PROTOC_ARCH := "$(UNAME_OS)-$(UNAME_ARCH)"
ifeq ($(UNAME_OS),darwin)
  PROTOC_ARCH = osx-x86_64
endif

ifeq ($(UNAME_ARCH),x86_64)
  UNAME_ARCH = amd64
endif

PROTOC_ZIP := protoc-$(PROTOC_VERSION)-$(PROTOC_ARCH).zip
MIGRATE_BIN := migrate.$(UNAME_OS)-$(UNAME_ARCH)

.PHONY: default
default: all

.PHONY: all
all: generate lint test build run

.PHONY: clean
clean:
	rm -rf ./bin

./bin:
	mkdir -p ./bin

./bin/protoc.zip: | ./bin
	curl -L https://github.com/google/protobuf/releases/download/v$(PROTOC_VERSION)/$(PROTOC_ZIP) -o ./bin/protoc.zip

./bin/protoc-validate.zip: | ./bin
	curl -L https://github.com/envoyproxy/protoc-gen-validate/archive/refs/tags/v$(PROTOC_VALIDATE_VERSION).zip -o ./bin/protoc-validate.zip

./bin/protoc: ./bin/protoc.zip
	unzip -o ./bin/protoc.zip -d ./ bin/protoc

./proto/include: ./bin/protoc.zip ./bin/protoc-validate.zip
	unzip -o ./bin/protoc.zip -d ./proto include/*
	unzip -o ./bin/protoc-validate.zip -d ./proto protoc-gen-validate-$(PROTOC_VALIDATE_VERSION)/validate/validate.proto

./bin/golangci-lint: | ./bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v$(GOLANGLINT_VERSION)

./bin/goimports: | ./bin
	go install -modfile tools/go.mod golang.org/x/tools/cmd/goimports@latest

./bin/grpcui: | bin
	go install -modfile tools/go.mod github.com/fullstorydev/grpcui/cmd/grpcui

./bin/protoc-gen-go: | ./bin
	go install -modfile tools/go.mod google.golang.org/protobuf/cmd/protoc-gen-go

./bin/protoc-gen-go-grpc: | ./bin
	go install -modfile tools/go.mod google.golang.org/grpc/cmd/protoc-gen-go-grpc

./bin/protoc-gen-validate: | ./bin
	go install -modfile tools/go.mod github.com/envoyproxy/protoc-gen-validate

./bin/minimock: | ./bin
	go install -modfile tools/go.mod github.com/gojuno/minimock/v3/cmd/minimock

./bin/gowrap: | ./bin
	go install -modfile tools/go.mod github.com/hexdigest/gowrap/cmd/gowrap

./bin/migrate: | ./bin
	curl -L https://github.com/golang-migrate/migrate/releases/download/v$(MIGRATE_VERSION)/$(MIGRATE_BIN).tar.gz -o ./bin/migrate.tar.gz
	@tar -xvf ./bin/migrate.tar.gz -C ./bin

PROTOS = $(wildcard ./proto/*.proto)

.PHONY: ./pkg
./pkg: PROTOC_OPT ?= module=github.com/farawaygg/$(PROJECT_NAME)/pkg:./pkg
./pkg: ./bin/protoc ./proto/include ./bin/protoc-gen-go ./bin/protoc-gen-go-grpc ./bin/protoc-gen-validate $(PROTOS)
	mkdir -p $@
	protoc -I ./proto \
    -I ./proto/include \
    -I ./proto/protoc-gen-validate-$(PROTOC_VALIDATE_VERSION) \
    --go_out=$(PROTOC_OPT) \
    --go-grpc_out=require_unimplemented_servers=false,$(PROTOC_OPT) \
    --validate_out="lang=go,$(PROTOC_OPT)" \
    --descriptor_set_out=./proto/$(SERVICE_NAME).protoset  \
    --include_imports \
    ./proto/*.proto


.PHONY: generate
generate: ./pkg ./bin/gowrap ./bin/minimock ./bin/goimports
	go generate ./...
	goimports -w -local github.com/farawaygg .


.PHONY: tidy
tidy:
	go mod tidy
	cd tools && go mod tidy

.PHONY: build
build: ./pkg
	GOGC=off go build -v -o ./bin/$(SERVICE_NAME) ./cmd/$(SERVICE_NAME)

.PHONY: prepare
prepare:
	echo 'this can be your custom actions'

.PHONY: run
run: ./pkg
	go run ./cmd/$(SERVICE_NAME)/*.go -config app/config.yaml

.PHONY: grpcui
grpcui: ./bin/grpcui
	grpcui \
		-protoset proto/$(SERVICE_NAME).protoset \
		-plaintext \
		-v \
		-port 30001 \
		127.0.0.1:48080

.PHONY: test-docker-pg-up
test-docker-pg-up:
	docker-compose -f docker-compose.yml up --no-deps -d --force-recreate pg
	timeout 30 sh -c 'until nc -z localhost 5432; do sleep 1; done' # waiting for port to come alive

.PHONY: test-docker-up
test-docker-up: test-docker-pg-up

migrate-test: ./bin/migrate test-docker-pg-up
	for i in $$(seq 1 10); do \
		psql -U postgres -h localhost -p 5432 -c "DROP DATABASE IF EXISTS wisdoms_test" \
		&& break || sleep 1; \
	done
	psql -U postgres -h localhost -p 5432 -c "CREATE DATABASE wisdoms_test"
	migrate -path ./migrations -database "postgres://postgres@localhost:5432?sslmode=disable&dbname=wisdoms_test" up

.PHONY: test-integration
test-integration: migrate-test test-docker-up
	DSN="postgresql://postgres@localhost:5432/wisdoms_test?sslmode=disable" \
	GOGC=off MallocNanoZone=0 \
		 GOGC=off go test -race -cover $(GOFLAGS) -v ./... -count 1 -tags integration -timeout 30s

.PHONY: test-unit
test-unit:
	GOGC=off go test -race $(GOFLAGS) -v ./... -count 1

.PHONY: test
test: test-unit test-integration

.PHONY: lint
lint: ./bin/golangci-lint
	golangci-lint run -v ./...

.PHONY: test-all
test-all: tidy lint test