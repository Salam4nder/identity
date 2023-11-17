package grpc

import (
	"github.com/Salam4nder/user/internal/config"
	"github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/internal/grpc/gen"
	"github.com/Salam4nder/user/pkg/token"
)

// UserServer is a gRPC server for user service.
type UserServer struct {
	gen.UserServer

	storage    db.Storage
	tokenMaker token.Maker
	config     config.UserService
}

// NewUserServer returns a new instance of UserService.
func NewUserServer(
	store db.Storage,
	cfg config.UserService,
) (*UserServer, error) {
	tokenMaker, err := token.NewPasetoMaker(cfg.SymmetricKey)
	if err != nil {
		return nil, err
	}

	return &UserServer{
		storage:    store,
		tokenMaker: tokenMaker,
		config:     cfg,
	}, nil
}
