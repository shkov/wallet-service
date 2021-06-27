.PHONY: help
help:
	@echo "usage: make [target]"
	@echo
	@echo "targets:"
	@echo "  help            show this message"
	@echo "  test            run all tests (requires docker)"
	@echo "  clean           stop docker containers"
	@echo "  run             run wallet-service with all dependencies"
	@echo

.PHONY: test
test:
	@go test -race -cover ./...

.PHONY: run
run:
	@docker-compose -f deployments/docker-compose.yml up -d --build
	# OK

.PHONY: clean
clean:
	@docker-compose -f deployments/docker-compose.yml down
	# OK
