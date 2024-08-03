package server

import (
	"database/sql"

	"github.com/Salam4nder/user/internal/auth"
	"github.com/Salam4nder/user/internal/token"
	"github.com/Salam4nder/user/proto/gen"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc/health"
)

// UserServer contains all necessary dependencies to serve user requests.
type UserServer struct {
	gen.UserServer

	strategy   auth.Strategy
	tokenMaker token.Maker
	health     *health.Server
	natsConn   *nats.Conn
	// TODO(kg): This is kind of ass.
	// It is here for now so we can monitor it's health with the [MoniterHealth] method.
	db *sql.DB
}

// NewUserServer returns a new UserService.
func NewUserServer(
	health *health.Server,
	natsConn *nats.Conn,
	strategy auth.Strategy,
	tokenMaker token.Maker,
	db *sql.DB,
) (*UserServer, error) {
	return &UserServer{
		strategy:   strategy,
		tokenMaker: tokenMaker,
		health:     health,
		natsConn:   natsConn,
		db:         db,
	}, nil
}
