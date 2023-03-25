package grpc

import (
	"github.com/Salam4nder/user/internal/proto/pb"
	"github.com/Salam4nder/user/internal/storage"
)

type userService struct {
	pb.UserServer
	storage.UserStorage
}

func NewUserService(store storage.UserStorage) *userService {
	return &userService{UserStorage: store}
}
