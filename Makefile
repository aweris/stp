GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

STP_MAIN	?=cmd/stp/main.go
GOFMT_FILES	?=$$(find . -name '*.go' | grep -v vendor)

DEV_DB = $(CURDIR)/development.store
V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

all: test

clean:  ; $(info $(M) Cleaning ...)
	$(Q) rm -rf $(DEV_DB)
test: ; $(info $(M) Running tests ...)
	$(Q) $(GOTEST) -v ./...

run: ; $(info $(M) Starting Application ...)
	$(Q) $(GOCMD) run $(STP_MAIN)

fmt: ; $(info $(M) gofmt ...)
	$(Q) $(GOFMT) -w $(GOFMT_FILES)

.PHONY: run fmt