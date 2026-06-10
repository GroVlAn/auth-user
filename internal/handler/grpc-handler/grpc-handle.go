package grpc_handler

import (
	"context"
	"time"

	api "github.com/GroVlAn/auth-api/user"
	"github.com/GroVlAn/auth-user/internal/domain"
	"github.com/rs/zerolog"
)

type service interface {
	Create(ctx context.Context, user domain.User) error
	User(ctx context.Context, userQuery domain.UserQuery) (domain.User, error)
	UserInfo(ctx context.Context, userQuery domain.UserQuery) (domain.UserInfo, error)
	UpdatePassword(ctx context.Context, userQueryNewPassword domain.UserQueryNewPassword) error
	InactivateUser(ctx context.Context, userQuery domain.UserQuery) error
	RestoreUser(ctx context.Context, userQuery domain.UserQuery) error
	BanUser(ctx context.Context, userQuery domain.UserQuery) error
	UnbanUser(ctx context.Context, userQuery domain.UserQuery) error
}

type GRPCHandler struct {
	api.UnimplementedUserServiceServer
	l              zerolog.Logger
	s              service
	defaultTimeout time.Duration
}

func New(l zerolog.Logger, s service, defTimeout time.Duration) *GRPCHandler {
	return &GRPCHandler{
		l:              l,
		s:              s,
		defaultTimeout: defTimeout,
	}
}
