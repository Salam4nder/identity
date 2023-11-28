# user

# Start:

1. Copy contents from the `.env` file to a `dev.env` file.
2. Fill in the values in your `dev.env`.
3. Run `make docker` to build the API docker image.
4. Run `make up` to compose up the API and all its dependencies.
5. `make logs` to follow the application logs and its dependencies.
6. `make down` to compose down with cleanup.

# Test
Unit test can be run with `make test` or simply with `go test ./...`  
This will however exclude tests that require a database connection.  

For testing the database layer, please run `make test-db`  
For viewing coverage, run `make test-cover`

# Endpoints
Easiest way to call the endpoints is with the evans gRPC client by running `make evans`. 
This will display the available endpoints based on reflection.
An example is `call ReadUser`.
See https://github.com/ktr0731/evans.


## TODO
* TLS setup.
* Email sending implementation.
* Complete remaining unit tests.
* oauth2 support.
