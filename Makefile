# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

.PHONY: all test coverage
all: test coverage examples

checkfmt:
	@echo 'Checking gofmt';\
 	bash -c "diff -u <(echo -n) <(gofmt -d .)";\
	EXIT_CODE=$$?;\
	if [ "$$EXIT_CODE"  -ne 0 ]; then \
		echo '$@: Go files must be formatted with gofmt'; \
	fi && \
	exit $$EXIT_CODE

get:
	GO111MODULE=on $(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint
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
	$(GOFMT) ./...
	golangci-lint run
	$(GOTEST) -race -covermode=atomic ./...

coverage: get test
	$(GOTEST) -race -coverprofile=coverage.txt -covermode=atomic ./redisai

