.SUFFIXES:

SUBDIRS := go rust python

.PHONY: help test lint format $(SUBDIRS)

help:  ## Show this help
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

test:  ## Run tests in all subdirs
	@for d in $(SUBDIRS); do $(MAKE) -C $$d test || exit $$?; done

lint:  ## Run lint in all subdirs
	@for d in $(SUBDIRS); do $(MAKE) -C $$d lint || exit $$?; done

format:  ## Run format in all subdirs
	@for d in $(SUBDIRS); do $(MAKE) -C $$d format || exit $$?; done
