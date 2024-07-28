package grpc

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"log/slog"
	"time"

	"github.com/Salam4nder/user/internal/db"
	"github.com/Salam4nder/user/internal/email"
	"github.com/Salam4nder/user/internal/event"
	"github.com/Salam4nder/user/internal/grpc/gen"
	"github.com/Salam4nder/user/internal/observability/metrics"
	"github.com/Salam4nder/user/pkg/password"
	"github.com/Salam4nder/user/pkg/validation"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (x *UserServer) CreateUser(ctx context.Context, req *gen.CreateUserRequest) (*emptypb.Empty, error) {
	ctx, span := tracer.Start(ctx, "rpc.CreateUser")
	defer span.End()

	if req == nil {
		return nil, requestIsNilError(span)
	}

	attrs, err := GenSpanAttributes(req)
	if err == nil {
		span.SetAttributes(attrs...)
	} else {
		slog.WarnContext(ctx, "grpc.CreateUser: GenSpanAttributes", "err", err.Error())
	}

	if err = validateCreateUserRequest(req); err != nil {
		return nil, invalidArgumentError(err, span, err.Error())
	}

	pw, err := password.FromString(req.GetPassword())
	if err != nil {
		return nil, invalidArgumentError(err, span, err.Error())
	}
	params := db.CreateUserParams{
		ID:        uuid.New(),
		FullName:  req.GetFullName(),
		Email:     req.GetEmail(),
		Password:  pw,
		CreatedAt: time.Now(),
	}
	if err = x.storage.CreateUser(ctx, params); err != nil {
		if errors.Is(err, db.ErrDuplicateEmail) {
			return nil, alreadyExistsError(err, span, "user with the provided email already exists")
		}
		return nil, internalServerError(err, span)
	}

	m := email.Email{
		To:      req.GetEmail(),
		From:    email.TestFrom,
		Subject: email.TestSubject,
		Body:    email.TestBody,
	}
	var buf bytes.Buffer
	if err = gob.NewEncoder(&buf).Encode(m); err != nil {
		slog.WarnContext(ctx, "grpc.CreateUser: gob.NewEncoder", "err", err.Error())
	}
	if err = x.natsConn.Publish(event.UserRegistered, buf.Bytes()); err != nil {
		slog.WarnContext(ctx, "grpc.CreateUser: nats.Publish", "err", err.Error())
	}
	metrics.UsersActive.Inc()
	metrics.UsersRegistered.Inc()

	return &emptypb.Empty{}, nil
}

// validateCreateUserRequest returns nil if the request is valid.
func validateCreateUserRequest(req *gen.CreateUserRequest) error {
	var (
		fullNameErr error
		emailErr    error
	)

	if err := validation.FullName(req.GetFullName()); err != nil {
		fullNameErr = err
	}

	if err := validation.Email(req.GetEmail()); err != nil {
		emailErr = err
	}

	return errors.Join(fullNameErr, emailErr)
}
