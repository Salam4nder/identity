# user

# Start:

1. Copy contents from the `.env` file to a `dev.env` file.
2. Fill in the values in your `dev.env`.
3. Run `make docker` to build the API docker image.
4. Run `make up` to compose up the API and all its dependencies.

# Test
Unit test can be run with `make test` or simply with `go test ./...`  
This will however exclude tests that require a database connection.  

For testing the database layer, please run `make test-db`  
For viewing coverage, run `make test-cover`


## TODO
* TLS setup.
* Email sending implementation.
