# Makefile settings
SHELL := /bin/bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := help

# DB settings
DB_HOSTNAME ?= localhost
DB_PORT ?= 5432
DB_NAME ?= twitter
DB_USERNAME ?= postgres
DB_PASSWORD ?= Pass1234!

DB_URL ?= "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOSTNAME):$(DB_PORT)/$(DB_NAME)?sslmode=disable"

MIGRATIONS_PATH ?= $(shell pwd)/migrations

LOG_LEVEL ?= debug
API_PORT ?= 8888

ENV_VARS = \
	DB_HOSTNAME=$(DB_HOSTNAME) \
	DB_PORT=$(DB_PORT) \
	DB_NAME=$(DB_NAME) \
	DB_USERNAME=$(DB_USERNAME) \
	DB_PASSWORD=$(DB_PASSWORD) \
	LOG_LEVEL=$(LOG_LEVEL) \
	API_PORT=$(API_PORT) \
	DB_URL=$(DB_URL) \
	$(NULL)


## Development targets

.PHONY: dev
dev: ## Run development server
	DB_URL=$(DB_URL) LOG_LEVEL=$(LOG_LEVEL) go run ./cmd/twitter -port $(API_PORT)

.PHONY: test
test: ## Run unit tests
	go test -v ./...

.PHONY: test_integration
test_integration: ## Run integration tests
	@$(ENV_VARS) MIGRATIONS_PATH=$(MIGRATIONS_PATH) go test -v ./... -tags=integration

.PHONY: test_e2e
test_e2e: ## Run end-to-end tests
	@$(ENV_VARS) docker-compose -f docker-compose-e2e.yaml -p e2e up --detach --build
	@$(ENV_VARS) docker-compose -f docker-compose-e2e.yaml -p e2e logs curl
	@OUT=$$( $(ENV_VARS) docker-compose -f docker-compose-e2e.yaml -p e2e logs curl | tail -n 1) ; \
	$(ENV_VARS) docker-compose -f docker-compose-e2e.yaml -p e2e down --volumes; \
	if [ "$${OUT:(-15)}" != "PASS: E2E Tests" ]; then \
		echo "E2E tests failed:"; \
		exit 1; \
	else \
		echo "E2E tests passed"; \
	fi

# Instalation: brew install golangci-lint
.PHONY: lint
lint: ## Lint and format source code based on golangci configuration
	@command -v golangci-lint || (echo "Please install `golangci-lint`" && exit 1)
	golangci-lint run --fix -v ./...

## DB targets

.PHONY: db-start
db-start: ## Postgres start
	@$(ENV_VARS) docker-compose -f docker-compose-db.yaml -p db up --detach postgres-dev

.PHONY: db-stop
db-stop: ## Postgres stop
	@$(ENV_VARS) docker-compose -f docker-compose-db.yaml -p db stop postgres-dev

.PHONY: db-cli
db-cli: ## Start the Postgres CLI
	@command -v pgcli || (echo "Please install `pgcli`." && exit 1)
	@PGPASSWORD='$(DB_PASSWORD)' \
		pgcli -h $(DB_HOSTNAME) -u $(DB_USERNAME) -p $(DB_PORT) -d $(DB_NAME)

#
## API targets


.PHONY: api-start
api-start: ## Run docker API container
	@$(ENV_VARS) docker-compose -f docker-compose-api.yaml -p api up --detach --build api-dev

.PHONY: api-stop
api-stop: ## Stop docker API container
	@$(ENV_VARS) docker-compose -f docker-compose-api.yaml -p api stop api-dev


### DB migration targets

# https://github.com/golang-migrate/migrate
# brew install golang-migrate
db-migrate-up: ## Run database upgrade migrations
	migrate -verbose -database $(DB_URL) -path migrations up

db-migrate-down:  ## Run database downgrade the last migration
	migrate -verbose -database $(DB_URL) -path migrations down 1

db-migrate-version:  ## Print the current migration version
	migrate -verbose -database $(DB_URL) -path migrations version

#### Code generation ####

## OpenAPI targets
# Install: go install "github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest"
.PHONY: openapi-generate
openapi-generate: ## Generate OpenAPI client
	mkdir -p internal/api/v1/openapi
	oapi-codegen \
		-generate types \
		-package openapi \
		-o internal/api/v1/openapi/types.go \
		openapi.yaml
	oapi-codegen \
		-generate chi-server \
		-package openapi \
		-o internal/api/v1/openapi/chi.go \
		openapi.yaml
	oapi-codegen \
		-generate spec \
		-package openapi \
		-o internal/api/v1/openapi/spec.go \
		openapi.yaml

## DB ORM targets
# go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest
# go install github.com/volatiletech/sqlboiler/v4@latest
.PHONY: db-orm-models
db-orm-models: ## Generate Go database models
	@command -v sqlboiler || (echo "Please install `sqlboiler`" && exit 1)
	PSQL_USER=$(DB_USERNAME) PSQL_PASS='$(DB_PASSWORD)' PSQL_HOST=$(DB_HOSTNAME) \
		sqlboiler --wipe --no-tests --add-soft-deletes \
		-o internal/service/repository/postgres/orm --pkgname orm \
		-c sqlboiler.toml psql

.PHONY: help
help:
	@grep -hE '^[a-zA-Z_-][0-9a-zA-Z_-]*:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

#### Docker targets ####

.PHONY: docker-up
docker-up: ## Run docker container
	@$(ENV_VARS) docker-compose -f docker-compose.yaml up --build

.PHONY: docker-down
docker-down: ## Stop docker container
	@$(ENV_VARS) docker-compose -f docker-compose.yaml down
