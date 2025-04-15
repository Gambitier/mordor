.PHONY: help gen-keys run examples

# Default target
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## ' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

RSA_KEYS_DIR ?= keys/rsa

gen-keys: ## Generate new RSA key pair (use RSA_KEYS_DIR=custom/path to override default path)
	@echo "Generating new RSA key pair in $(RSA_KEYS_DIR)..."
	@mkdir -p $(RSA_KEYS_DIR)
	@go run cmd/RSAKeyGenerator/main.go -output $(RSA_KEYS_DIR)

# Example usage help
examples: ## Show example commands
	@echo "Example commands:"
	@echo "  make gen-keys                         # Generate keys in default directory (keys/rsa)"
	@echo "  make gen-keys RSA_KEYS_DIR=./keys/rsa # Generate keys in custom directory"