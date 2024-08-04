# Identity: A lightweight identity service.

## Quickstart

*Identity* operates on a given authentication `strategy`.

A `strategy` must implement the following interface:

```go
// WIP.
type Strategy interface {
	// ConfiguredStrategy exposes the current configured strategy.
	ConfiguredStrategy() gen.Strategy

	// Renew will trade a valid refresh token for a new access token.
	Renew(context.Context) error
	// Revoke all active tokens in the configured hot-storage
    // for the user.
	Revoke(context.Context) error
	// Register an entry with the configured strategy.
	Register(context.Context) error
	// Authenticate the user with the configured strategy.
	Authenticate(context.Context) error
}
```


## Usage

```go
// gRPC client stubs.
func ExampleClient() {
	cfg := config.Server{
		GRPCHost: "0.0.0.0",
		GRPCPort: "50051",
	}
	conn, err := grpc.NewClient(cfg.Addr())
	if err != nil {
		// handle err
	}
	defer conn.Close()

	client := gen.NewIdentityClient(conn)

	// Register with the credentials strategy.
	_, err = client.Register(context.TODO(), &gen.Input{
		Strategy: gen.Strategy_Credentials,
		Data: &gen.Input_Credentials{
			Credentials: &gen.CredentialsInput{
				Email:    "email@email.com",
				Password: "securePassword400",
			},
		},
	})
	if err != nil {
		// handle err
	}

	// Register with the PersonalNumber strategy.
	_, err = client.Register(context.TODO(), &gen.Input{
		Strategy: gen.Strategy_PersonalNumber,
		Data: &gen.Input_Numbers{
			Numbers: &gen.PersonalNumberInput{
				Numbers: 867497568396,
			},
		},
	})
	if err != nil {
		// handle err
	}
}
```

```json
// JSON body.
{
  "strategy": 1, // 1 is the Credentials strategy in the enum.
  "credentials": {
    "email": "email@email.com",
    "password": "password23423"
  }
}
```

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
This service serves gRPC requests. You can use GUI tools like **Insomnia** to test the endpoints.


## TODO
* Examples.
* TLS setup.
* Email sending implementation. For now a no-op stdout logger faker is used.
* Complete remaining unit tests.
* more auth strategies.
