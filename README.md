# User: A lightweight identity service.

## Config
The application expects a `config.yaml` file in the root of the project.

## Run
Run `make up` to compose up all application dependencies.

Run `go run .` in the root of the application to start serving requests.

You can run `make down` to compose down all the application dependencies.

You can alternatively run the server in a contianer:
1. Build the API by running `make api`.
2. Uncomment the api service in `compose.yaml`.
3. Edit the hosts for all the dependencies `config.yaml` file. 
4. Run `make up`.

## Test
Run all tests with `make test`. This will spin up a required Postgres container defined in `internal/db/compose.yaml`. 

Running `go test./...` will exclude tests that require a db connection.


Running `make test-db` will only run tests that require a db connection.


For testing the database layer, please run `make test-db`.
For viewing coverage, run `make test-cover`.

## Lint
Run `make lint` to run the linter engine. Linters are described in the `.golangci.yaml` file.

## Endpoints
This service serves gRPC requests. You can use GUI tools like **Insomnia** to call the endopints.


## TODO
* TLS setup.
* Email sending implementation.
* Complete remaining unit tests.
* oauth2 support.
