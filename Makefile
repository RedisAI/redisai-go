# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

.PHONY: all test coverage
all: test coverage

get:
	$(GOGET) -t -v ./...

test: get
	$(GOTEST) -v -race -covermode=atomic ./...

coverage: get test
	$(GOTEST) -v -race -coverprofile=coverage.txt -covermode=atomic ./redisai

