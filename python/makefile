.SUFFIXES:

.PHONY: help test lint format

help:  ## Show this help
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

test:  ## Run pytest + mypy
	uv run --group dev pytest
	uv run --group dev mypy

lint:  ## ruff check + format check (no writes)
	uv run --group dev ruff check .
	uv run --group dev ruff format --check .

format:  ## ruff format the tree in place
	uv run --group dev ruff format .
