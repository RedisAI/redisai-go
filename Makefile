# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

.PHONY: all test coverage
all: test coverage examples

get:
	$(GOGET) -t -v ./redisai/...

TLS_CERT ?= redis.crt
TLS_KEY ?= redis.key
TLS_CACERT ?= ca.crt
REDISAI_TEST_HOST ?= 127.0.0.1:6379

examples: get
	@echo " "
	@echo "Building the examples..."
	$(GOBUILD) ./examples/redisai_pipelined_client/.
	$(GOBUILD) ./examples/redisai_simple_client/.
	$(GOBUILD) ./examples/redisai_tls_client/.
	./redisai_tls_client --tls-cert-file $(TLS_CERT) \
						 --tls-key-file $(TLS_KEY) \
						 --tls-ca-cert-file $(TLS_CACERT) \
						 --host $(REDISAI_TEST_HOST)

test: get
	$(GOTEST) -race -covermode=atomic ./...

coverage: get test
	$(GOTEST) -race -coverprofile=coverage.txt -covermode=atomic ./redisai

