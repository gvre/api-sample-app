.DEFAULT_GOAL := help

# Export variables from .env
ifneq (,$(wildcard ./.env))
	include .env
    export
endif

TAGS = "api"
COVERAGE=

.PHONY: help
help: ## Display help. `make start` or `make start_debug` are what you will need most of the times. Use the `clean` target to remove any existing containers and volumes.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: build
build: ## Build images
	docker-compose build db

.PHONY: start
start: stop build db migrate ## Start all necessary containers

.PHONY: stop
stop: ## Stop all containers
	docker-compose stop

.PHONY: db
db:
	docker-compose up -d db

.PHONY: test
test: ## Test the application. Use the TAGS argument to pass specific tags (e.g. make TAGS=api test)
	CGO_ENABLED=1 go test -buildvcs=false -shuffle=on -race -count=1 $(COVERAGE) -tags $(TAGS) ./...

.PHONY: coverage
coverage: COVERAGE = -coverprofile=coverage.out
coverage: test ## Run tests and display coverage
	go tool cover -func=coverage.out

.PHONY: build-migrations-image
build-migrations-image: ## Build and tag a Docker image for running database migrations
	DOCKER_BUILDKIT=1 docker build --no-cache --force-rm --secret id=dotenv,src=.env --target=migrations -t api-sample-app:migrations .

.PHONY: migrate
migrate: build-migrations-image ## Run "up" database migrations
	$(call wait_for_db)
	@docker-compose run migrate -url=jdbc:postgresql://db/$${PGDATABASE}?sslmode=disable -user=$${PGUSER} -password=$${PGPASSWORD} -connectRetries=3 -locations="filesystem:/migrations" -outOfOrder=true -ignoreMigrationPatterns='*:missing,*:future' migrate

.PHONY: logs
logs: ## Watch the logs of all containers
	docker-compose logs -f

.PHONY: clean
clean:  ## Remove all containers and volumes of the app
	docker-compose rm --stop --force
	docker volume rm api-sample-app_pgdata &>/dev/null || exit 0

.PHONY: openapi_validate
openapi_validate: ## Validate the OpenAPI specification file
	@docker run -v "${PWD}/cmd/rest/spec/openapi.yaml:/openapi.yaml" --rm p1c2u/openapi-spec-validator /openapi.yaml

.PHONY: openapi_editor
openapi_editor: ## Start the Swagger Editor
	docker-compose up -d openapi_editor

define wait_for_db
	@for i in {1..30}; do \
		if docker run -t postgres:16-alpine pg_isready -h $${DB_HOST:-host.docker.internal} -p $${PGPORT:-5432}; then \
			break; \
		fi; \
		sleep 1; \
	done
endef