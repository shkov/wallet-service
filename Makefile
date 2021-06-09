.PHONY: help
help:
	@echo "usage: make [target]"
	@echo
	@echo "targets:"
	@echo "  help            show this message"
	@echo "  test            run all tests (requires docker)"
	@echo "  clean-docker    stop test docker containers"
	@echo "  run             run wallet-service with all dependencies"
	@echo
