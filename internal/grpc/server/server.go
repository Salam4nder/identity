package server

import (
	"database/sql"

	"github.com/Salam4nder/user/internal/auth"
	"github.com/Salam4nder/user/internal/token"
	"github.com/Salam4nder/user/proto/gen"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc/health"
)

// Identity contains all necessary dependencies to serve gRPC requests.
type Identity struct {
	gen.IdentityServer

	db         *sql.DB
	health     *health.Server
	natsConn   *nats.Conn
	strategy   auth.Strategy
	tokenMaker token.Maker
}

// NewUserServer returns a new UserService.
func NewUserServer(
	db *sql.DB,
	health *health.Server,
	natsConn *nats.Conn,
	strategy auth.Strategy,
	tokenMaker token.Maker,
) (*Identity, error) {
	return &Identity{
		strategy:   strategy,
		tokenMaker: tokenMaker,
		health:     health,
		natsConn:   natsConn,
		db:         db,
	}, nil
}
