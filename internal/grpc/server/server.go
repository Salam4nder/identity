package server

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Salam4nder/identity/internal/auth"
	"github.com/Salam4nder/identity/internal/auth/strategy/credentials"
	"github.com/Salam4nder/identity/internal/auth/strategy/personalnumber"
	"github.com/Salam4nder/identity/internal/token"
	"github.com/Salam4nder/identity/proto/gen"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc/health"
)

// Identity contains all necessary dependencies to serve gRPC requests.
type Identity struct {
	gen.IdentityServer

	db         *sql.DB
	health     *health.Server
	natsConn   *nats.Conn
	tokenMaker token.Maker

	strategies map[gen.Strategy]auth.Strategy
}

// NewIdentity returns a new [Identity] gRPC server.
func NewIdentity(
	db *sql.DB,
	health *health.Server,
	natsConn *nats.Conn,
	tokenMaker token.Maker,
) *Identity {
	return &Identity{
		tokenMaker: tokenMaker,
		health:     health,
		natsConn:   natsConn,
		db:         db,
	}
}

// MountStrategies will parse the configured strategy string representations and mount them on the server.
// Aborts and returns an error if any of them fails to parse.
func (x *Identity) MountStrategies(s ...string) error {
	m := make(map[gen.Strategy]auth.Strategy)

	for _, v := range s {
		strategy, err := auth.StrategyFromString(v)
		if err != nil {
			return err
		}
		slog.Info(fmt.Sprintf("mounted strategy %s", v))
		switch strategy {
		case gen.Strategy_TypeCredentials:
			m[strategy] = credentials.New(x.db, x.natsConn)
		case gen.Strategy_TypePersonalNumber:
			m[strategy] = personalnumber.New(x.db)
		default:
			return errors.New("server: unsupported strategy")
		}
	}

	x.strategies = m

	return nil
}
