package server

import (
	"errors"

	"github.com/Salam4nder/user/proto/gen"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("rpc")

// GenerateSpanAttributes returns span attributes for generated request structs.
// Experimental solution, not the prettiest.
func GenSpanAttributes(param any) ([]attribute.KeyValue, error) {
	if param == nil {
		return nil, errors.New("param is nil")
	}

	switch p := param.(type) {
	case *gen.UpdateUserRequest:
		return []attribute.KeyValue{
			attribute.String("id", p.GetId()),
			attribute.String("full_name", p.GetFullName()),
			attribute.String("email", p.GetEmail()),
		}, nil
	case *gen.CreateUserRequest:
		return []attribute.KeyValue{
			attribute.String("full_name", p.FullName),
			attribute.String("email", p.Email),
			attribute.Int("password_length", len(p.Password)),
		}, nil
	case *gen.UserID:
		return []attribute.KeyValue{
			attribute.String("id", p.GetId()),
		}, nil
	case *gen.LoginUserRequest:
		return []attribute.KeyValue{
			attribute.String("email", p.GetEmail()),
			attribute.Int("password_length", len(p.GetPassword())),
		}, nil
	case *gen.ReadByEmailRequest:
		return []attribute.KeyValue{
			attribute.String("email", p.GetEmail()),
		}, nil

	default:
		return nil, errors.New("unsupported type")
	}
}
