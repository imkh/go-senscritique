GO = GO111MODULE=on go

BINARY = senscritique
PACKAGE = ./cmd/$(BINARY)

.PHONY: all
all: help ## help

.PHONY: run
run: ## Run the program
	$(GO) run $(PACKAGE)

.PHONY: install
install: ## Install binary in $GOBIN (make sure it's set in your $PATH)
	$(GO) install -v $(PACKAGE)

.PHONY: build
build: ## Build binary in this directory
	$(GO) build -v -o $(BINARY) $(PACKAGE)

.PHONY: dep
dep: ## Get the dependencies
	$(GO) get -d -v ./...

.PHONY: test
test: ## Run unit tests
	$(GO) test -v ./...

.PHONY: lint
lint: ## Lint the files (requires golangci-lint installed)
	golangci-lint run

.PHONY: clean
clean: ## Remove previous build files
	rm -f $(BINARY)

.PHONY: help
help: ## Display this help screen
	@echo "Usage:"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf " \033[36m%-15s\033[0m %s\n", $$1, $$2}'
