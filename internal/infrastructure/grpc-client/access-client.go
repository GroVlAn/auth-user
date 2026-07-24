package grpc_client

import (
	"context"
	"fmt"

	api "github.com/GroVlAn/auth-api/access"
	"github.com/GroVlAn/auth-base/ew"
	"google.golang.org/grpc"
)

type UserGRPCClient struct {
	client api.AccessServiceClient
}

func New(conn *grpc.ClientConn) *UserGRPCClient {
	return &UserGRPCClient{
		client: api.NewAccessServiceClient(conn),
	}
}

func (uc *UserGRPCClient) BindUserRole(ctx context.Context, userID string) error {
	success, err := uc.client.BindUserRole(ctx, &api.UserID{User_ID: userID})
	if err != nil {
		return ew.New(
			ew.ErrorTypeInternal,
			fmt.Errorf("binding user role: %w", err),
		)
	}
	if !success.Success {
		return ew.New(
			ew.ErrorTypeInternal,
			fmt.Errorf("failed binding user role"),
		)
	}

	return nil
}
