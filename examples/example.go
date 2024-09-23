package examples

import (
	"context"
	"log/slog"
	"time"

	"github.com/Salam4nder/identity/pkg/random"
	"github.com/Salam4nder/identity/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RegisterAndAuthenticate() error {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	client := gen.NewIdentityClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, err := client.Register(ctx, &gen.RegisterRequest{
		Strategy: gen.Strategy_TypePersonalNumber,
		Data:     &gen.RegisterRequest_Empty{},
	})
	if err != nil {
		return err
	}
	n := resp.GetNumber().GetNumber()

	slog.Info("client: got number", "number", n)

	resp, err = client.Register(ctx, &gen.RegisterRequest{
		Strategy: gen.Strategy_TypeCredentials,
		Data: &gen.RegisterRequest_Credentials{
			Credentials: &gen.CredentialsInput{
				Email:    random.Email(),
				Password: "Pas5w0rd9912",
			},
		},
	})
	if err != nil {
		return err
	}

	authResp, err := client.Authenticate(ctx, &gen.AuthenticateRequest{
		Strategy: gen.Strategy_TypePersonalNumber,
		Data: &gen.AuthenticateRequest_Number{
			Number: &gen.PersonalNumber{
				Number: n,
			},
		},
	})
	if err != nil {
		return err
	}

	slog.Info(
		"client: got tokens",
		"access", authResp.GetAccessToken(),
		"refresh", authResp.GetRefreshToken(),
	)

	_, err = client.Validate(ctx, &gen.TokenRequest{
		Token: authResp.GetAccessToken(),
	})
	if err != nil {
		return err
	}

	newAccessTokenResp, err := client.Refresh(ctx, &gen.TokenRequest{
		Token: authResp.GetRefreshToken(),
	})
	if err != nil {
		return err
	}

	slog.Info("client: got new access token", "token", newAccessTokenResp.GetToken())

	return nil
}
