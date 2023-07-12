ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
PROJECT_NAME := $(shell basename "${ROOT}")
COMPOSE_FILE := _tools/docker-compose.yml
COMPOSE_ARGS := -f "${COMPOSE_FILE}" -p "${PROJECT_NAME}"

.PHONY: run
run:
	go run .

.PHONY: db-up
db-up: ## Run docker compose up in the background
	docker compose ${COMPOSE_ARGS} up -d db

PHNEY: db-down
db-down: ## Run docker compose down
	docker compose ${COMPOSE_ARGS} down

.PHONY: db-logs
db-logs: ## Tail docker compose logs
	docker compose ${COMPOSE_ARGS} logs -f db

.PHONY: db-exec
db-exec:
	docker compose ${COMPOSE_ARGS} exec db bash -c \
	'mysql -u webapp --password=webapp_no_password'

.PHONY: db-clean
db-clean:
	${MAKE} db-down
# when you create docker volume, you need to delete it and add the command.

generate: ## Generate codes
	go generate ./...

.PHONY: test
test: ## Execute tests
	go test -race -shuffle=on ./...

.PHONY: clean
clean:
	${MAKE} db-clean

help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
