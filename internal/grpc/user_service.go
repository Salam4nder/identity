package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/Salam4nder/user/internal/proto/pb"
	"github.com/Salam4nder/user/internal/storage"
	"go.uber.org/zap"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type userService struct {
	pb.UserServer
	logger *zap.Logger
	storage.UserStorage
}

// NewUserService returns a new instance of UserService.
func NewUserService(store storage.UserStorage) *userService {
	return &userService{UserStorage: store}
}

// CreateUser creates a new user. Returns an error if the user couldn't be created
// or if the request is invalid.
func (s *userService) CreateUser(
	ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if err := validateCreateUserRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	insertOneParam := protoToInsertOneParam(req)
	insertOneParam.CreatedAt = time.Now()

	createdUserID, err := s.UserStorage.InsertOne(ctx, insertOneParam)
	if err != nil {
		s.logger.Error("failed to insert user", zap.Error(err))
		internalServerError()
	}

	return &pb.CreateUserResponse{Id: createdUserID}, nil
}

// GetUser returns a user by id. Returns an error if the user couldn't be found
// or if the request is invalid.
func (s *userService) GetUser(
	ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if req == nil {
		return nil, requestIsNilError()
	}

	if req.GetId() == "" {
		return nil, invalidIDError()
	}

	user, err := s.UserStorage.FindOneByID(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, storage.UserNotFoundErr()) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		s.logger.Error("failed to find user", zap.Error(err))

		return nil, internalServerError()
	}

	return &pb.GetUserResponse{User: userToProto(user)}, nil
}

// validateCreateUserRequest returns nil if the request is valid.
func validateCreateUserRequest(req *pb.CreateUserRequest) error {
	if req == nil {
		return errors.New("request can not be nil")
	}

	var (
		fullNameErr error
		emailErr    error
		passwordErr error
	)

	if req.GetFullName() == "" {
		fullNameErr = errors.New("full name can not be empty")
	}

	if req.GetEmail() == "" {
		emailErr = errors.New("email can not be empty")
	}

	if req.GetPassword() == "" {
		passwordErr = errors.New("password can not be empty")
	}

	return errors.Join(fullNameErr, emailErr, passwordErr)
}

func protoToInsertOneParam(req *pb.CreateUserRequest) storage.InsertOneParam {
	return storage.InsertOneParam{
		FullName: req.GetFullName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
}

func userToProto(user storage.User) *pb.UserResponse {
	return &pb.UserResponse{
		Id:        user.ID.Hex(),
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}
