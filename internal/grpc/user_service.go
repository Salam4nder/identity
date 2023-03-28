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
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type userService struct {
	pb.UserServer
	logger *zap.Logger
	storage.UserStorage
}

// NewUserService returns a new instance of UserService.
func NewUserService(store storage.UserStorage, log *zap.Logger) *userService {
	return &userService{UserStorage: store, logger: log}
}

// CreateUser creates a new user. Returns an error if the user couldn't be created
// or if the request is invalid.
func (s *userService) CreateUser(
	ctx context.Context, req *pb.CreateUserRequest) (*pb.UserID, error) {
	if err := validateCreateUserRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	insertOneParam := protoToInsertOneParam(req)
	insertOneParam.CreatedAt = time.Now()

	createdUserID, err := s.UserStorage.InsertOne(ctx, insertOneParam)
	if err != nil {
		s.logger.Error("failed to insert user", zap.Error(err))

		return nil, internalServerError()
	}

	return &pb.UserID{Id: createdUserID}, nil
}

// GetUser returns a user by id. Returns an error if the user couldn't be found
// or if the request is invalid.
func (s *userService) GetUser(
	ctx context.Context, req *pb.UserID) (*pb.UserResponse, error) {
	if req == nil {
		return nil, requestIsNilError()
	}

	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "ID can not be empty")
	}

	user, err := s.UserStorage.FindOneByID(ctx, req.GetId())
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, storage.ErrInvalidID):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			s.logger.Error("failed to find user", zap.Error(err))
			return nil, internalServerError()
		}
	}

	return userToProto(user), nil
}

// GetByEmail returns a user by email. Returns an error if the user couldn't be not
// found or if the request is invalid.
func (s *userService) GetByEmail(
	ctx context.Context, req *pb.GetByEmailRequest) (*pb.UserResponse, error) {
	if req == nil {
		return nil, requestIsNilError()
	}

	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email can not be empty")
	}

	user, err := s.UserStorage.FindOneByEmail(ctx, req.GetEmail())
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			s.logger.Error("failed to find user", zap.Error(err))
			return nil, internalServerError()
		}
	}

	return userToProto(user), nil
}

// GetByFilter returns a list of users by filter. Returns an error if the request is invalid
// or no users were found.
func (s *userService) GetByFilter(
	ctx context.Context, req *pb.GetByFilterRequest) (*pb.GetByFilterResponse, error) {
	if err := validateGetByFilterRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	filter := storage.Filter{
		FullName:  req.GetFullName(),
		Email:     req.GetEmail(),
		CreatedAt: req.GetCreatedAt().AsTime(),
	}

	fetchedUsers, err := s.UserStorage.FindByFilter(ctx, filter)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			s.logger.Error("failed to find user", zap.Error(err))
			return nil, internalServerError()
		}
	}

	var users []*pb.UserResponse

	for _, user := range fetchedUsers {
		users = append(users, userToProto(user))
	}

	if len(users) < 1 {
		return nil, status.Error(codes.NotFound, "no users found")
	}

	return &pb.GetByFilterResponse{Users: users}, nil
}

// UpdateUser updates a user by id. Returns an error if the user couldn't be updated
// or if the request is invalid.
func (s *userService) UpdateUser(
	ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	if err := validateUpdateUserRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	updateOneParam := protoToUpdateParam(req)

	updatedUser, err := s.UserStorage.UpdateOne(ctx, updateOneParam)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, storage.ErrInvalidID):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			s.logger.Error("failed to find user", zap.Error(err))
			return nil, internalServerError()
		}
	}

	return userToProto(updatedUser), nil
}

// DeleteUser deletes a user by id. Returns an error if the user couldn't be deleted,
// if the user doesn't exist or if the request is invalid.
func (s *userService) DeleteUser(
	ctx context.Context, req *pb.UserID) (*emptypb.Empty, error) {
	if req == nil {
		return nil, requestIsNilError()
	}

	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "ID can not be empty")
	}

	err := s.UserStorage.DeleteOne(ctx, req.GetId())
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, storage.ErrInvalidID):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			s.logger.Error("failed to find user", zap.Error(err))
			return nil, internalServerError()
		}
	}

	return &emptypb.Empty{}, nil
}

// validateUpdateUserRequest returns nil if the request is valid.
func validateUpdateUserRequest(req *pb.UpdateUserRequest) error {
	if req == nil {
		return errors.New("request can not be nil")
	}

	if req.Id == "" {
		return errors.New("invalid id")
	}

	if req.GetFullName() == "" && req.GetEmail() == "" {
		return errors.New("at least one field must be provided for an update")
	}

	return nil
}

// validateGetByFilterRequest returns nil if the request is valid.
func validateGetByFilterRequest(req *pb.GetByFilterRequest) error {
	if req == nil {
		return errors.New("request can not be nil")
	}

	if req.GetFullName() == "" &&
		req.GetEmail() == "" &&
		req.CreatedAt.AsTime().IsZero() {
		return errors.New("at least one field must be provided for a filter")
	}

	return nil
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

func protoToUpdateParam(req *pb.UpdateUserRequest) storage.UpdateParam {
	return storage.UpdateParam{
		ID:       req.GetId(),
		FullName: req.GetFullName(),
		Email:    req.GetEmail(),
	}
}

func protoToInsertOneParam(req *pb.CreateUserRequest) storage.InsertParam {
	return storage.InsertParam{
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
