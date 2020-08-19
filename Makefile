.DEFAULT_GOAL := help

TAGS = "api"
COVERAGE=

.PHONY: help
help: ## Display help. `make start` or `make start_debug` are what you will need most of the times. Use the `clean` target to remove any existing containers and volumes.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: build
build: ## Build images
	docker-compose build db app

.PHONY: start
start: stop build db migrate ## Start all necessary containers
	docker-compose up -d app

.PHONY: start_debug
start_debug: stop build db migrate ##Start a container with a debugger and the application
	docker-compose up -d debug

.PHONY: stop
stop: ## Stop all containers
	docker-compose stop

.PHONY: db
db:
	docker-compose up -d db

.PHONY: app
app: ## Build application container
	docker-compose stop app
	docker-compose rm --force app
	docker-compose up -d --build app

.PHONY: app_debug
app_debug: ## Build application and debugger container
	docker-compose stop debug
	docker-compose rm --force debug
	docker-compose up -d --build debug

.PHONY: test
test: ## Test the application. Use the TAGS argument to pass specific tags (e.g. make TAGS=api test)
	docker-compose run test sh -c 'CGO_ENABLED=1 go test -race -count=1 $(COVERAGE) -tags $(TAGS) ./...'

.PHONY: coverage
coverage: COVERAGE = -coverprofile=coverage.out
coverage: test ## Run tests and display coverage
	docker-compose run test sh -c 'go tool cover -func=coverage.out'

.PHONY: migrate
migrate: ## Run "up" database migrations
	sleep 5
	docker-compose run migrate sh -c 'migrate -verbose -path=/migrations -database="postgres://@" up'

.PHONY: openapi_validate
openapi_validate:
	@docker run -v "${PWD}/cmd/rest/spec/openapi.yaml:/openapi.yaml" --rm p1c2u/openapi-spec-validator /openapi.yaml

.PHONY: logs
logs: ## Watch the logs of all containers
	docker-compose logs -f

.PHONY: clean
clean:  ## Remove all containers and volumes of the app
	docker-compose rm --stop --force
	docker volume rm api-sample-app_pgdata &>/dev/null || exit 0
