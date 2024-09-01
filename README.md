# Identity: A lightweight identity service.

## Quickstart

*Identity* operates on given authentication strategies.

A `strategy` must implement the following interface:

```go
// WIP.
type Strategy interface {
	// Register an entry with the configured strategy.
	// Outputs from this method are stored in the
	// returned context.
	Register(context.Context) (context.Context, error)
	// Authenticate the user with the configured strategy.
	Authenticate(context.Context) error
}
```


## Usage

```go
// gRPC client stubs.
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		// handle err
	}
	client := gen.NewIdentityClient(conn)

	resp, err := client.Register(context.TODO(), &gen.RegisterRequest{
		Strategy: gen.Strategy_TypePersonalNumber,
		Data:     &gen.RegisterRequest_Empty{},
	})
	if err != nil {
		// handle err
	}
	slog.Info("client got resp number", "number", resp.GetNumber().Numbers)

	resp, err = client.Register(context.TODO(), &gen.RegisterRequest{
		Strategy: gen.Strategy_TypeCredentials,
		Data: &gen.RegisterRequest_Credentials{
			Credentials: &gen.CredentialsInput{
				Email:    random.Email(),
				Password: "pasSwe22i3rj",
			},
		},
	})
	if err != nil {
		// handle err
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
