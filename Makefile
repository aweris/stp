GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

STP_MAIN	?=cmd/stp/main.go
GOFMT_FILES	?=$$(find . -name '*.go' | grep -v vendor)

V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

run: ; $(info $(M) Starting Application ...)
	$(Q) $(GOCMD) run $(STP_MAIN)

fmt: ; $(info $(M) gofmt ...)
	$(Q) $(GOFMT) -w $(GOFMT_FILES)

.PHONY: run fmt