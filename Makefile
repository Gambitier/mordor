output ?= keys/rsa
KEEP_COUNT ?= 2
SHELL_TYPE ?= $(shell basename $$SHELL)

.PHONY: help gen-keys clean examples build gen-rsa-keys cleanup-rsa-keys install-completion

# Default target
help:
	@echo "Usage: make <target> [-- <args>]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## ' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	@echo "Building the binary..."
	@go build -o ./bin/RSAKeyGenerator ./cmd/RSAKeyGenerator/main.go

gen-rsa-keys: build ## Generate new RSA key pair (use output=custom/path to override default path)
	@mkdir -p $(output)
	@./bin/RSAKeyGenerator generate --output $(output) $(filter-out $@, $(MAKECMDGOALS))

cleanup-rsa-keys: build ## Clean up old RSA key pairs (use KEEP_COUNT=N to keep N most recent pairs, default 2)
	@./bin/RSAKeyGenerator cleanup --output $(output) --keep $(KEEP_COUNT) $(filter-out $@, $(MAKECMDGOALS))

clean: ## Clean the binary and key directory
	@rm -rf ./bin
	@rm -rf $(output)

# Example usage help
examples: ## Show example commands
	@echo "Example commands:"
	@echo "  make build                                      # Build the binary"
	@echo "  make gen-rsa-keys                               # Generate keys in default directory (keys/rsa)"
	@echo "  make gen-rsa-keys output=./keys/rsa             # Generate keys in custom directory"
	@echo "  make gen-rsa-keys -- --help                     # Show help for key generation"
	@echo "  make cleanup-rsa-keys                           # Clean up keys, keeping 2 most recent pairs"
	@echo "  make cleanup-rsa-keys KEEP_COUNT=3              # Clean up keys, keeping 3 most recent pairs"
	@echo "  make cleanup-rsa-keys -- --help                 # Show help for cleanup"
	@echo "  make clean                                      # Clean the binary and key directory"
	@echo ""

# Avoid "No rule to make target '--help'" error
%:
	@:
