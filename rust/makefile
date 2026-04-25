.SUFFIXES:

.PHONY: help test lint format

help:  ## Show this help
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

test:  ## Run the test suite
	cargo test

lint:  ## clippy + rustfmt check (no writes)
	cargo clippy --all-targets -- -D warnings
	cargo fmt -- --check

format:  ## rustfmt the tree in place
	cargo fmt
