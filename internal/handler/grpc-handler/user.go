package grpc_handler

import (
	"context"

	api "github.com/GroVlAn/auth-api/user"
	"github.com/GroVlAn/auth-base/ew/grpcx"
	"github.com/GroVlAn/auth-user/internal/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *GRPCHandler) Register(ctx context.Context, req *api.User) (*api.Success, error) {
	user := domain.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Fullname: req.Fullname,
	}

	ctx, cancel := context.WithTimeout(ctx, h.defaultTimeout)
	defer cancel()

	if err := h.s.Create(ctx, user); err != nil {
		return nil, grpcx.HandleError(h.l, err)
	}

	return &api.Success{
		Success: true,
	}, nil
}

func (h *GRPCHandler) GetUser(ctx context.Context, req *api.UserQuery) (*api.User, error) {
	userQuery := domain.UserQuery{
		ID:       &req.ID,
		Username: &req.Username,
		Email:    &req.Email,
	}

	ctx, cancel := context.WithTimeout(ctx, h.defaultTimeout)
	defer cancel()

	user, err := h.s.User(ctx, userQuery)
	if err != nil {
		return nil, grpcx.HandleError(h.l, err)
	}

	return &api.User{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Fullname:     user.Fullname,
		IsSuperuser:  user.IsSuperuser,
		IsActive:     user.IsActive,
		IsBanned:     user.IsBanned,
		CreatedAt:    timestamppb.New(user.CreatedAt),
	}, nil
}

func (h *GRPCHandler) GetUserInfo(ctx context.Context, req *api.UserQuery) (*api.UserInfo, error) {
	userQuery := domain.UserQuery{
		ID:       &req.ID,
		Username: &req.Username,
		Email:    &req.Email,
	}

	ctx, cancel := context.WithTimeout(ctx, h.defaultTimeout)
	defer cancel()

	userInfo, err := h.s.UserInfo(ctx, userQuery)
	if err != nil {
		return nil, grpcx.HandleError(h.l, err)
	}

	return &api.UserInfo{
		Username: userInfo.Username,
		Email:    userInfo.Email,
		Fullname: userInfo.Fullname,
	}, nil
}

func (h *GRPCHandler) ChangePassword(ctx context.Context, req *api.UserQueryNewPassword) (*api.Success, error) {
	userQueryNewPassword := domain.UserQueryNewPassword{
		UserQuery: domain.UserQuery{
			ID:       &req.UserQuery.ID,
			Username: &req.UserQuery.Username,
			Email:    &req.UserQuery.Email,
		},
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}

	ctx, cancel := context.WithTimeout(ctx, h.defaultTimeout)
	defer cancel()

	if err := h.s.UpdatePassword(ctx, userQueryNewPassword); err != nil {
		return nil, grpcx.HandleError(h.l, err)
	}

	return &api.Success{
		Success: true,
	}, nil
}

func (h *GRPCHandler) InactivateUser(ctx context.Context, req *api.UserQuery) (*api.Success, error) {
	return h.changeUserStatus(ctx, req, h.s.InactivateUser)
}

func (h *GRPCHandler) RestoreUser(ctx context.Context, req *api.UserQuery) (*api.Success, error) {
	return h.changeUserStatus(ctx, req, h.s.RestoreUser)
}

func (h *GRPCHandler) BanUser(ctx context.Context, req *api.UserQuery) (*api.Success, error) {
	return h.changeUserStatus(ctx, req, h.s.BanUser)
}

func (h *GRPCHandler) UnbanUser(ctx context.Context, req *api.UserQuery) (*api.Success, error) {
	return h.changeUserStatus(ctx, req, h.s.UnbanUser)
}

func (h *GRPCHandler) changeUserStatus(
	ctx context.Context,
	uQr *api.UserQuery,
	fn func(context.Context, domain.UserQuery) error,
) (*api.Success, error) {
	userQuery := domain.UserQuery{
		ID:       &uQr.ID,
		Username: &uQr.Username,
		Email:    &uQr.Email,
	}

	ctx, cancel := context.WithTimeout(ctx, h.defaultTimeout)
	defer cancel()

	if err := fn(ctx, userQuery); err != nil {
		return nil, grpcx.HandleError(h.l, err)
	}

	return &api.Success{
		Success: true,
	}, nil
}
