package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/Salam4nder/user/pkg/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

func (x *UserServer) authorizeUser(ctx context.Context) (*token.Payload, error) {
	metadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("grpc: missing metadata")
	}

	values := metadata.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("grpc: missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("grpc: invalid authorization header format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationBearer {
		return nil, fmt.Errorf("grpc: unsupported authorization type: %s", authType)
	}

	accessToken := fields[1]
	payload, err := x.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("grpc: invalid access token: %w", err)
	}

	return payload, nil
}
