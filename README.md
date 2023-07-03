# Twitter Clone

## Architecture

### Backend

I have used an architecture following the same principles as the [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) by Robert C. Martin. I have included the repository pattern to abstract the data layer from the business logic. The most important packages are:

- `internal/service/repository/postgres/orm`: the ORM auto generated code from the database schema that [SQLBoiler](https://github.com/volatiletech/sqlboiler) generates.

There are 2 repositories:
- `internal/service/repository/postgres`: the repository that implements the interfaces defined in `internal/service/repository/interfaces.go`
- `internal/service/repository/memory`: a memory mock repository that implements the interfaces defined in `internal/service/repositoryinterfaces.go`. This is used for unit testing.

The business logic is implemented in the `internal/service` package. This represents the use cases of the application - the domain service. Note that the business logic always use data entities defined in the `internal/entities` package. This ensures business logic is decoupled from the data layer defined in the `internal/service/repository` package.

### Frontend API

The use cases are exposed via the API layer. The API layer is implemented using the [go-chi](https://github.com/go-chi/chi). 

The endpoints routes were generated from a [Open-API](https://www.openapis.org/) [spec](openapi.yaml) using [oapi-codegen](https://github.com/deepmap/oapi-codegen). The code generated is in `internal/api/v1/openapi`.

The following endpoints are exposed:
```bash
# Get the API spec
GET /api/v1/api.json
# Get all tweets
GET /api/v1/tweets
# Create a tweet
POST /api/v1/tweets
# Get a tweet by id
GET /api/v1/tweets/{id}
# Create a user
POST /api/v1/users
# Get all users
GET /api/v1/users
# Get a user by id
GET /api/v1/users/{id}
```

## Running the application

### Prerequisites
- docker
- docker-compose

### Running the application

Ideally, the application should be launched from the makefile. This makes sure `docker-compose` is run with the correct environment variables. Otherwise, set the environment variables in the `env-template` file. See `env-template` for the instructions.

```bash
make docker-up
```

Below are a few examples of how to use the API. Note that the API spec is available at `http://localhost:8888/api/v1/api.json`.

```bash
## Create a user
curl -v -X POST -H "Content-Type: application/json" \
-d '{ "username": "foo", "name": "John Doe", "email": "jd@mail.com" }' \
http://localhost:8888/api/v1/users

## Get all users
curl -v http://localhost:8888/api/v1/users

## Set user_id as an environment variable
user_id=$(curl http://localhost:8888/api/v1/users | jq -r '.[0].id')

## Get a user by id
curl -v http://localhost:8888/api/v1/users/$user_id

## Create a tweet
curl -v -X POST -H "Content-Type: application/json" \
-d '{"user_id":"'$user_id'", "content": "Hello World!" }' \
http://localhost:8888/api/v1/tweets

# Get all tweets
curl -v http://localhost:8888/api/v1/tweets

## Set twitter_id as an environment variable
tweet_id=$(curl http://localhost:8888/api/v1/tweets | jq -r '.[0].id')

## Get a tweet by id
curl -v http://localhost:8888/api/v1/tweets/$tweet_id
```

## Development

### Makefile targets

The `make help` command lists all the available targets.

```bash
# Start / stop the database running in docker
db-start                       Postgres start
db-stop                        Postgres stop
# Open a Postgres CLI
db-cli                         Start the Postgres CLI
# Manage the database schema
db-migrate-down                Run database downgrade the last migration
db-migrate-up                  Run database upgrade migrations
db-migrate-version             Print the current migration version
# Generate code
openapi-generate               Generate OpenAPI client
db-orm-models                  Generate Go database models
# Development targets
dev                            Run development server
lint                           Lint and format source code based on golangci configuration
# Runs tests
test                           Run unit tests
test_integration               Run integration tests
# Starts the API in docker (starts the database and runs the migrations if needed)
api-start                      Run docker API container
api-stop                       Stop docker API container
# Runs docker-compose up/down with the correct environment variables
docker-down                    Stop docker container
docker-up                      Run docker container
```
