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

## Testing
- `make test` Run all tests. Use the TAGS argument to pass specific tags (e.g. `make TAGS=api test`)

## Architecture
<pre>┌─────────────┐          ┌─────────────┐          ┌─────────────┐          ┌─────────────┐
│             ├─────────▶│             │─────────▶│    User     │─────────▶│             │
│HTTP Handlers│          │User Service │          │ Repository  │          │  Database   │
│             │◀─────────┤             │◀─────────│             │◀─────────┤             │
└─────────────┘          └─────────────┘          └─────────────┘          └─────────────┘</pre>

## License
- MIT