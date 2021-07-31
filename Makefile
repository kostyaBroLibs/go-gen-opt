HELP_FUN = \
	%help; \
	while(<>) { push @{$$help{$$2 // 'options'}}, [$$1, $$3] if /^(\w+)\s*:.*\#\#(?:@(\w+))?\s(.*)$$/ }; \
	print "usage: make [target]\n\n"; \
	for (keys %help) { \
	print "$$_:\n"; $$sep = " " x (20 - length $$_->[0]); \
	print "  $$_->[0]$$sep$$_->[1]\n" for @{$$help{$$_}}; \
	print "\n"; }

help: ##@miscellaneous Show this help.
	@perl -e '$(HELP_FUN)' $(MAKEFILE_LIST)

lint: linter ## Run linter.

linter: ## Run linter.
	docker run --rm -v `pwd`:/app -w /app golangci/golangci-lint:v1.42 golangci-lint run -v

test: ## Run tests.
	go test -mod=vendor -v -race -bench=. -benchmem ./...
