# Configuration
APP_NAME := zen
IMAGE ?= ghcr.io/senither/zen-lang:latest

# Build metadata
VERSION := $(shell git describe --tags --abbrev=0 2>NUL || echo dev)
GIT_COMMIT := $(shell git rev-parse --short=12 HEAD 2>NUL)
# Optional: override GOOS and GOARCH for cross-compilation
GOOS ?=
GOARCH ?=

# Platform-specific settings
ifeq ($(OS),Windows_NT)
    BIN_NAME := $(APP_NAME).exe
    BUILD_DATE := $(shell powershell -NoProfile -Command "Get-Date -Format 'MMM_dd_yyyy_HH:mm:ss'" 2>NUL)
else
    BIN_NAME := $(APP_NAME)
    BUILD_DATE := $(shell date +"%b_%d_%Y_%H:%M:%S" 2>/dev/null)
endif

# Linker flags for embedding build metadata
LDFLAGS := -s -w -X github.com/senither/zen-lang/cli.Version=$(VERSION) -X github.com/senither/zen-lang/cli.GitCommit=$(GIT_COMMIT) -X github.com/senither/zen-lang/cli.BuildDate=$(BUILD_DATE)

.PHONY: all install build docker test test-integration test-language bench clean

all: clean build

install:
	go mod download
	go mod verify

build:
	@if not defined GOOS (set GOOS=$(GOOS)) else (set GOOS=%GOOS%)
	@if not defined GOARCH (set GOARCH=$(GOARCH)) else (set GOARCH=%GOARCH%)
	go build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN_NAME) .

docker:
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg BUILD_DATE="$(BUILD_DATE)" \
		-t $(IMAGE) .

test: test-integration test-language

test-integration:
	go test -v ./...

test-language:
	go run main.go test

bench:
	go test ./... -bench=. -run=^$ -benchmem -benchtime=5s

ifeq ($(OS),Windows_NT)
    RM_BIN := @if exist $(BIN_NAME) del /q $(BIN_NAME)
else
    RM_BIN := @rm -f $(BIN_NAME)
endif

clean:
	$(RM_BIN)
