# Identity: A lightweight identity service.

## Config

The application expects a `config.yaml` file in the root of the project.

## Run

Run `make up` to compose up the application and all its dependencies.


## Cleanup


Run `make down` to nuke everything down.


## Test

Run all tests with `make test`. This will spin up a required Postgres container defined in `internal/db/compose.yaml`. 

Running `go test./...` will exclude tests that require a db connection.


Running `make test-db` will only run tests that require a db connection.


## Lint
Run `make lint` to run the linter engine. Linters are described in the `.golangci.yaml` file.


## Endpoints
This service serves gRPC requests. You can use GUI tools like **Insomnia** to test the endpoints.


## TODO
* TLS setup.
* Email sending implementation.
* Complete remaining unit tests.
* oauth2 support.
