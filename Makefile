BINARY	=	senscritique

all:	help

## install: Install binary in $GOBIN. Make sure it's set in your $PATH to run it from any directory.
install:
		go install -v ./cmd/...

## build: Build binary in local directory.
build:
		go build -o $(BINARY) -v ./cmd/$(BINARY)/...

## test: Runs `go test -v ./...`.
test:
		go test -v ./...

## clean: Runs `go clean` and clean build files.
clean:
		go clean
		rm -f $(BINARY)

## deps: Runs `GO111MODULE=on go mod vendor`.
deps:
		GO111MODULE=on go mod vendor

help:	Makefile
		@echo "Usage:"
		@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

.PHONY:	all install build test clean deps help
