
.PHONY: serve-book
serve-book: ## Build and serve the book (with live-reload)
	docker compose -f docs/book/docker-compose.yaml up