ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

BIN_DIR := bin
TOOLS_DIR := hack/tools
TOOLS_BIN_DIR := $(abspath $(TOOLS_DIR)/$(BIN_DIR))

LINKCHECK_BIN := linkcheck
LINKCHECK := $(abspath $(TOOLS_BIN_DIR)/$(LINKCHECK_BIN))

.PHONY: serve-book
serve-book: ## Build and serve the book (with live-reload)
	docker run --rm -v $(ROOT_DIR):/src -w=/src/docs/book -p 1313:1313 --platform linux/amd64 docsy/docsy-example server

.PHONY: verify-markdown-link
verify-markdown-link: $(LINKCHECK_BIN) ## Verify links in markdown files
	$(TOOLS_BIN_DIR)/$(LINKCHECK_BIN) --root $(ROOT_DIR) --hugo-folder docs/book

.PHONY: $(LINKCHECK_BIN)
$(LINKCHECK_BIN): $(LINKCHECK) ## Build a local copy of linkcheck.

$(LINKCHECK): $(TOOLS_DIR)/go.mod # Build link-check from tools folder.
	cd $(TOOLS_DIR); go build  -o $(BIN_DIR)/$(LINKCHECK_BIN) github.com/fabriziopandini/cluster-api-website/hack/tools/linkcheck
