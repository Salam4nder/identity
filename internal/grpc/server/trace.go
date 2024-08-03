package server

import (
	"errors"
	"fmt"

	"github.com/Salam4nder/user/proto/gen"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("server")

// GenerateSpanAttributes returns span attributes for generated request structs.
// Experimental solution, not the prettiest.
func GenSpanAttributes(param any) ([]attribute.KeyValue, error) {
	if param == nil {
		return nil, errors.New("param is nil")
	}

	switch t := param.(type) {
	case *gen.Input_Credentials:
		return []attribute.KeyValue{
			attribute.String("email", t.Credentials.GetEmail()),
			attribute.Int("password length", len(t.Credentials.GetEmail())),
		}, nil
	default:
		return nil, fmt.Errorf("server: span attributes, unsupported type %T", t)
	}
}
