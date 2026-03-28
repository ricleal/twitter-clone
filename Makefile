# Makefile settings
SHELL := /bin/bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := help

# DB settings
POSTGRES_USER ?= postgres
POSTGRES_PASSWORD ?= Pass_1234
POSTGRES_DB ?= twitter

DB_URL_LOCAL ?= postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5433/$(POSTGRES_DB)?sslmode=disable
DB_URL_COMPOSE ?= postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres:5432/$(POSTGRES_DB)?sslmode=disable
DB_URL_API_DEV ?= postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@host.docker.internal:5433/$(POSTGRES_DB)?sslmode=disable
DB_URL_E2E ?= postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres:5432/$(POSTGRES_DB)?sslmode=disable
DB_URL ?= $(DB_URL_LOCAL)

MIGRATIONS_PATH ?= $(shell pwd)/migrations

LOG_LEVEL ?= debug
API_PORT ?= 8888

ENV_VARS = \
	POSTGRES_USER=$(POSTGRES_USER) \
	POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	POSTGRES_DB=$(POSTGRES_DB) \
	DB_URL=$(DB_URL) \
	DB_URL_LOCAL=$(DB_URL_LOCAL) \
	DB_URL_COMPOSE=$(DB_URL_COMPOSE) \
	DB_URL_API_DEV=$(DB_URL_API_DEV) \
	DB_URL_E2E=$(DB_URL_E2E) \
	MIGRATIONS_PATH=$(MIGRATIONS_PATH) \
	LOG_LEVEL=$(LOG_LEVEL) \
	API_PORT=$(API_PORT) \
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
	@$(ENV_VARS) go test -race ./... -tags=integration

.PHONY: test_e2e
test_e2e: ## Run end-to-end tests
	@$(ENV_VARS) docker-compose -f docker-compose-e2e.yaml -p e2e down --volumes --remove-orphans >/dev/null 2>&1 || true
	@status=0; \
	$(ENV_VARS) docker-compose -f docker-compose-e2e.yaml -p e2e up --build --abort-on-container-exit --exit-code-from curl curl || status=$$?; \
	$(ENV_VARS) docker-compose -f docker-compose-e2e.yaml -p e2e down --volumes; \
	exit $$status

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
	@pgcli $(DB_URL)

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
# Install: go install "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest"
.PHONY: openapi-generate
openapi-generate: ## Generate OpenAPI client
	mkdir -p internal/api/v1/openapi
	oapi-codegen \
		-generate types \
		-package openapi \
		-o internal/api/v1/openapi/types.go \
		openapi.yaml
	oapi-codegen \
		-generate std-http-server \
		-package openapi \
		-o internal/api/v1/openapi/server.go \
		openapi.yaml
	oapi-codegen \
		-generate spec \
		-package openapi \
		-o internal/api/v1/openapi/spec.go \
		openapi.yaml

## DB ORM targets
# go install github.com/stephenafamo/bob/gen/bobgen-psql@latest
.PHONY: db-orm-models
db-orm-models: ## Generate Go database models
	@command -v bobgen-psql || (echo "Please install bobgen-psql: go install github.com/stephenafamo/bob/gen/bobgen-psql@latest" && exit 1)
	@tmp_cfg=$$(mktemp); \
	sed 's|__DB_URL_LOCAL__|$(DB_URL_LOCAL)|g' bobgen.yaml > "$$tmp_cfg"; \
	bobgen-psql -c "$$tmp_cfg"; \
	rm -f "$$tmp_cfg"

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
