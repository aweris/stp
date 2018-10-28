GOCMD	=go
GOBUILD	=$(GOCMD) build
GOCLEAN	=$(GOCMD) clean
GOTEST	=$(GOCMD) test
GOGET	=$(GOCMD) get

STP_MAIN	 ?=cmd/stp/main.go
GOFMT_FILES	 ?=$$(find . -name '*.go' | grep -v vendor)
GOTEST_FILES ?=$$(go list ./... | grep -v /vendor/) -cover

DEV_DB 		= $(CURDIR)/development.store
RELEASE_DIR = release

BINARY_NAME=stp

V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

.PHONY: clean test run fmt
all: clean test build

clean:  ; $(info $(M) Cleaning ...)
	$(Q) rm -rf $(DEV_DB) $(RELEASE_DIR)
	$(Q) $(GOCLEAN)

test: ; $(info $(M) Running tests ...)
	$(Q) $(GOTEST) $(GOTEST_FILES)

run: ; $(info $(M) Starting Application ...)
	$(Q) $(GOCMD) run $(STP_MAIN)

fmt: ; $(info $(M) gofmt ...)
	$(Q) $(GOFMT) -w $(GOFMT_FILES)

.PHONY: build build-linux build-darwin

build: build-darwin build-linux

build-linux: ; $(info $(M) Building - GOOS=linux GOARCH=amd64 ...)
	$(Q) mkdir -p $(RELEASE_DIR)
	$(Q) GOOS=linux GOARCH=amd64 $(GOBUILD) -o release/$(BINARY_NAME)-linux-amd64 $(STP_MAIN)

build-darwin: ; $(info $(M) Building - GOOS=darwin GOARCH=amd64 ...)
	$(Q) mkdir -p $(RELEASE_DIR)
	$(Q) GOOS=darwin GOARCH=amd64 $(GOBUILD) -o release/$(BINARY_NAME)-darwin-amd64 $(STP_MAIN)
