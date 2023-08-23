ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
PROJECT_NAME := $(shell basename "${ROOT}")
COMPOSE_FILE := docker-compose.yml
COMPOSE_ARGS := -f "${COMPOSE_FILE}" -p "${PROJECT_NAME}"

.PHONY: run
run:
	go run .

.PHONY: up
up: ## Run docker compose up in the background
	docker compose ${COMPOSE_ARGS} up -d

PHNEY: down
down: ## Run docker compose down
	docker compose ${COMPOSE_ARGS} down

.PHONY: logs
logs: ## Tail docker compose logs
	docker compose ${COMPOSE_ARGS} logs -f

.PHONY: db-exec
db-exec:
	docker compose ${COMPOSE_ARGS} exec db bash -c \
	'mysql -u webapp --password=webapp_no_password'

generate: ## Generate codes
	go generate ./...

.PHONY: test
test: ## Execute tests
	${MAKE} up
	go test -race -shuffle=on ./...

.PHONY: lint
lint:
	golangci-lint run --config=./.golangci.yml ./...

.PHONY: clean
clean:
	${MAKE} down

help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
