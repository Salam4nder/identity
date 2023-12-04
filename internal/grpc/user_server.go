package grpc

import (
	"github.com/Salam4nder/user/internal/config"
	"github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/internal/grpc/gen"
	"github.com/Salam4nder/user/internal/task"
	"github.com/Salam4nder/user/pkg/token"
)

// UserServer contains all necessary dependencies to serve user requests.
type UserServer struct {
	gen.UserServer

	storage     db.Storage
	taskCreator task.Creator
	tokenMaker  token.Maker
	config      config.UserService
}

// NewUserServer returns a new UserService.
func NewUserServer(
	store db.Storage,
	task task.Creator,
	cfg config.UserService,
) (*UserServer, error) {
	tokenMaker, err := token.NewPasetoMaker(cfg.SymmetricKey)
	if err != nil {
		return nil, err
	}

	return &UserServer{
		storage:     store,
		taskCreator: task,
		tokenMaker:  tokenMaker,
		config:      cfg,
	}, nil
}
