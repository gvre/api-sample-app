# REST API Sample Application

## Requirements
- Docker
- Docker Compose

## Usage
- `make help` Display help, and the most useful Makefile targets
- `make start` Build all containers and start the services
- `make start_debug` Build container with the application and a debugger
- `make stop` Stop all containers
- `make app` Build application container
- `make app_debug` Build application and debugger container
- `make clean` Remove any existing containers and volumes of the application.
- `make openapi_validate`  Validate the OpenAPI specification file
- `make openapi_editor` Start the Swagger Editor (http://localhost:18081/)

## API endpoints
- `curl http://localhost:18080/users` Return all users
- `curl http://localhost:18080/users/1` Return the user with ID `1`
- `curl -H 'Content-type: application/json' http://localhost:18080/users -d '{"name":"new user"}'` Add a new user

## Testing
- `make test` Run all tests. Use the TAGS argument to pass specific tags (e.g. `make TAGS=api test`)

## Migrations
All database migrations should be backward compatible, so only the `up` ones are needed.
The migrations library supports the `down` ones though, in case you want to play with fire :)

See [MIGRATIONS.md](https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md) for more details.

## Architecture
<pre>┌─────────────┐          ┌─────────────┐          ┌─────────────┐          ┌─────────────┐
│             ├─────────▶│             │─────────▶│    User     │─────────▶│             │
│HTTP Handlers│          │User Service │          │ Repository  │          │  Database   │
│             │◀─────────┤             │◀─────────│             │◀─────────┤             │
└─────────────┘          └─────────────┘          └─────────────┘          └─────────────┘</pre>

## License
- MIT
