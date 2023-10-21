# user

# Start:

1. Copy contents from the `.env` file to a `dev.env` file.
2. Fill in your `dev.env`.
3. Run `make docker` to build the API docker image.
4. Run `make up` to compose up the API and all its dependencies.

# Test
Unit test can be can with `make test` or simply `go test ./...`. 
This will however exclude tests that require a database connection.

For testing the database layer, please run `make test-db`
For viewing coverage, run `make test-cover`


## TODO
* Email verification task handler.
* Integration tests.
* TLS setup.
