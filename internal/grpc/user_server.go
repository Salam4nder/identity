package grpc

import (
	"github.com/Salam4nder/user/internal/config"
	"github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/internal/proto/gen"
	"github.com/Salam4nder/user/pkg/token"
	"github.com/rs/zerolog"
)

type userServer struct {
	gen.UserServer

	storage    *db.SQL
	tokenMaker token.Maker
	logger     *zerolog.Logger
	config     config.UserService
}

// NewUserService returns a new instance of UserService.
func NewUserService(
	store *db.SQL,
	log *zerolog.Logger,
	cfg config.UserService) (*userServer, error) {
	tokenMaker, err := token.NewPasetoMaker(cfg.SymmetricKey)
	if err != nil {
		return nil, err
	}

	return &userServer{
		storage:    store,
		tokenMaker: tokenMaker,
		logger:     log,
		config:     cfg,
	}, nil
}
