# user

# Start:

1. Copy contents from the `.env` file to a `dev.env` file.
2. Fill in your `dev.env`.
2. `make up`

# Test
Unit test can be can with `make test` or simply `go test ./...`. 
This will however exclude tests that require a database connection.

For testing the database layer, please run `make test-db`
For viewing coverage, run `make test-cover`


## TODO
* Email verification task handler.
* Integration tests.
* TLS setup.
