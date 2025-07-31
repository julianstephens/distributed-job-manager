.PHONY: help migrate-create

help:
	@echo "Available targets:"
	@grep -E '^[a-zA-Z0-9_-]+:.*##' $(MAKEFILE_LIST) | sort | while read -r line; do \
		target=$$(echo "$$line" | cut -d':' -f1); \
		description=$$(echo "$$line" | sed -e 's/.*## //'); \
		printf "  \033[36m%-20s\033[0m %s\n" "$$target" "$$description"; \
	done

migrate-up: ## Apply all migrations
	@migrate -path cassandra-migrations/ -database $$CASS_URL -verbose up

migrate-down: ## Rollback all migrations
	@migrate -path cassandra-migrations/ -database $$CASS_URL -verbose down
