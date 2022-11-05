ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

.PHONY: serve-book
serve-book: ## Build and serve the book (with live-reload)
	docker run --rm -v $(ROOT_DIR):/src -w=/src/docs/book -p 1313:1313 --platform linux/amd64 docsy/docsy-example server

.PHONY: verify-markdown-link
verify-markdown-link: ## Verify links in markdown files
	act --rm --job markdown-link-check --container-architecture linux/amd64