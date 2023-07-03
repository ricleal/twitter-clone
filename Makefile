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
	$(NULL)

## OpenAPI targets
# Install: go install "github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest"
.PHONY: openapi-generate
openapi-generate: ## Generate OpenAPI client
	oapi-codegen \
		-generate types \
		-package openapi \
		-o internal/api/openapi/types.go \
		openapi.yaml
	oapi-codegen \
		-generate chi-server \
		-package openapi \
		-o internal/api/openapi/chi.go \
		openapi.yaml
	oapi-codegen \
		-generate spec \
		-package openapi \
		-o internal/api/openapi/spec.go \
		openapi.yaml
#
## DB targets

.PHONY: db-start
db-start: ## Postgres start
	@$(ENV_VARS) docker-compose -f docker-compose.yaml up --detach postgres

.PHONY: db-stop
db-stop: ## Postgres stop
	@$(ENV_VARS) docker-compose -f docker-compose.yaml stop postgres

.PHONY: db-cli
db-cli: ## Postgres CLI
	@command -v pgcli || (echo "Please install `pgcli`." && exit 1)
	@PGPASSWORD='$(DB_PASSWORD)' \
		pgcli -h $(DB_HOSTNAME) -u $(DB_USERNAME) -p $(DB_PORT) -d $(DB_NAME)

#
## Docker targets

docker-build: ## Build docker image
	@$(ENV_VARS) docker-compose -f docker-compose.yaml build api

.PHONY: api-start
api-start: docker-build ## Run docker container
	@$(ENV_VARS) docker-compose -f docker-compose.yaml up --detach api

.PHONY: api-stop
api-stop: ## Stop docker container
	@$(ENV_VARS) docker-compose -f docker-compose.yaml stop api

#
## Development targets

.PHONY: dev
dev: ## Run development server
	DB_URL=$(DB_URL) LOG_LEVEL=$(LOG_LEVEL) go run ./cmd/twitter -port $(API_PORT)

.PHONY: test
test: ## Run tests
	go test -v ./...

.PHONY: test_integration
test_integration: ## Run integration tests
	@$(ENV_VARS) MIGRATIONS_PATH=$(MIGRATIONS_PATH) go test ./... -tags=integration


# brew install golangci-lint
.PHONY: format
format: ## Format source code based on golangci and prettier configuration
	@command -v golangci-lint || (echo "Please install `golangci-lint`" && exit 1)
	golangci-lint run --fix -v ./...

# https://github.com/golang-migrate/migrate
# brew install golang-migrate
db-migrate-up: ## Run database upgrade migrations
	migrate -verbose -database $(DB_URL) -path migrations up

db-migrate-down:  ## Run database downgrade the last migration
	migrate -verbose -database $(DB_URL) -path migrations down 1

db-migrate-version:  ## print the current migration version
	migrate -verbose -database $(DB_URL) -path migrations version

#### Code generation ####
# go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest
# go install github.com/volatiletech/sqlboiler/v4@latest
.PHONY: db-orm-models
db-orm-models: ## Generate Go database models
	@command -v sqlboiler || (echo "Please install `sqlboiler`" && exit 1)
	PSQL_USER=$(DB_USERNAME) PSQL_PASS='$(DB_PASSWORD)' PSQL_HOST=$(DB_HOSTNAME) \
		sqlboiler --wipe --no-tests --add-soft-deletes -o internal/service/repository/postgres/orm --pkgname orm -c sqlboiler.toml psql

.PHONY: help
help:
	@grep -hE '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
