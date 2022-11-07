ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

BIN_DIR := bin
TOOLS_DIR := hack/tools
TOOLS_BIN_DIR := $(abspath $(TOOLS_DIR)/$(BIN_DIR))

LINK_CHECK_BIN := link-check
LINK_CHECK := $(abspath $(TOOLS_BIN_DIR)/$(LINK_CHECK_BIN))

.PHONY: serve-book
serve-book: ## Build and serve the book (with live-reload)
	docker run --rm -v $(ROOT_DIR):/src -w=/src/docs/book -p 1313:1313 --platform linux/amd64 docsy/docsy-example server

.PHONY: verify-markdown-link
verify-markdown-link: ## Verify links in markdown files
	act --rm --job markdown-link-check --container-architecture linux/amd64


.PHONY: $(LINK_CHECK_BIN)
$(LINK_CHECK_BIN): $(LINK_CHECK) ## Build a local copy of tilt-prepare.

$(LINK_CHECK): $(TOOLS_DIR)/go.mod # Build link-check from tools folder.
	cd $(TOOLS_DIR); go build -tags=tools -o $(BIN_DIR)/$(LINK_CHECK_BIN) github.com/fabriziopandini/cluster-api-website/hack/tools/link-check
    #use this command to test iteratively rm -f hack/tools/bin/link-check && make link-check && hack/tools/bin/link-check --verbose --hugo-folder docs/book