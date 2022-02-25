.PHONY: bench
bench: ## Run benchmarks
	go test -benchmem -run=^$ github.com/demdxx/rpool -bench . -v -race

.PHONY: test
test: ## Run unit tests
	go test -v -race ./...

.PHONY: tidy
tidy: ## Run tidy util
	go mod tidy

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help