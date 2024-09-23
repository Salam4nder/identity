# Identity: A lightweight identity service.

## Description

*Identity* is a simple yet easily extendable identity service that handles user registration, authentication and token management.

It produces stateless tokens and operates on given authentication strategies.

A `strategy` must implement the following interface:

```go
type Strategy interface {
	// Strategies require different inputs and outputs, 
    // so we store them in contexts.
	Register(context.Context) (context.Context, error)
	Authenticate(context.Context) error
}
```
One of the implemented strategies is authentication by a `personal number`, which is simply a 16-digit number.

It is a simple yet super convenient way for users to start using your prouducts without giving you their personal information.

This strategy is inspired by *Mullvad VPN*.

## Usage

See the `examples` package.

## Config

The application expects a `config.yaml` file in the root of the project.

## Run

Run `make api` to build the api image and `make up` to compose up the application and all its dependencies.


## Cleanup

Run `make down` to nuke everything down.


## Test

Run all tests with `make test`. This will spin up a required Postgres container defined in `internal/database/compose.yaml`. 

Running `go test./...` will exclude tests that require a db connection.


Running `make test-db` will only run tests that require a db connection.


## Lint

Run `make lint` to run the linter engine. Linters are described in the `.golangci.yaml` file.


## Endpoints

```proto
service Identity {
    // Exchange a valid refresh token for a new access token.
    rpc Refresh (TokenRequest) returns (RefreshResponse){}
    // Validate an access token.
    rpc Validate(TokenRequest) returns (google.protobuf.Empty){}
    // Register a user with the given strategy.
    rpc Register (RegisterRequest) returns (RegisterResponse){}
    // Verify a user that registered with the credentials strategy.
    rpc VerifyEmail (TokenRequest) returns (google.protobuf.Empty){}
    // Authenticate a user with the given strategy.
    rpc Authenticate (AuthenticateRequest) returns (AuthenticateResponse){}
}
```

## TODO
* TLS setup.
* Email sending implementation. For now a no-op stdout logger is used.
* Complete remaining unit tests.
* More auth strategies.
