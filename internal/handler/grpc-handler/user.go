package grpc_handler

import (
	"context"

	api "github.com/GroVlAn/auth-api/user"
	"github.com/GroVlAn/auth-user/internal/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *GRPCHandler) Register(ctx context.Context, u *api.User) (*api.Success, error) {
	user := domain.User{
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
		Fullname: u.Fullname,
	}

	ctx, cancel := context.WithTimeout(ctx, h.DefaultTimeout)
	defer cancel()

	if err := h.s.Create(ctx, user); err != nil {
		return nil, h.handleError(err)
	}

	return &api.Success{
		Success: true,
	}, nil
}

func (h *GRPCHandler) GetUser(ctx context.Context, uQr *api.UserQuery) (*api.User, error) {
	userQuery := domain.UserQuery{
		ID:       uQr.ID,
		Username: uQr.Username,
		Email:    uQr.Email,
	}

	ctx, cancel := context.WithTimeout(ctx, h.DefaultTimeout)
	defer cancel()

	user, err := h.s.User(ctx, userQuery)
	if err != nil {
		return nil, h.handleError(err)
	}

	return &api.User{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Password:    user.Password,
		Fullname:    user.Fullname,
		IsSuperuser: user.IsSuperuser,
		IsActive:    user.IsActive,
		IsBanned:    user.IsBanned,
		CreatedAt:   timestamppb.New(user.CreatedAt),
	}, nil
}

func (h *GRPCHandler) GetUserInfo(ctx context.Context, uQr *api.UserQuery) (*api.UserInfo, error) {
	userQuery := domain.UserQuery{
		ID:       uQr.ID,
		Username: uQr.Username,
		Email:    uQr.Email,
	}

	ctx, cancel := context.WithTimeout(ctx, h.DefaultTimeout)
	defer cancel()

	userInfo, err := h.s.UserInfo(ctx, userQuery)
	if err != nil {
		return nil, h.handleError(err)
	}

	return &api.UserInfo{
		Username: userInfo.Username,
		Email:    userInfo.Email,
		Fullname: userInfo.Fullname,
	}, nil
}

func (h *GRPCHandler) ChangePassword(ctx context.Context, uQrNP *api.UserQueryNewPassword) (*api.Success, error) {
	userQueryNewPassword := domain.UserQueryNewPassword{
		UserQuery: domain.UserQuery{
			ID:       uQrNP.UserQuery.ID,
			Username: uQrNP.UserQuery.Username,
			Email:    uQrNP.UserQuery.Email,
		},
		OldPassword: uQrNP.OldPassword,
		NewPassword: uQrNP.NewPassword,
	}

	ctx, cancel := context.WithTimeout(ctx, h.DefaultTimeout)
	defer cancel()

	if err := h.s.UpdatePassword(ctx, userQueryNewPassword); err != nil {
		return nil, h.handleError(err)
	}

	return &api.Success{
		Success: true,
	}, nil
}

func (h *GRPCHandler) InactivateUser(ctx context.Context, uQr *api.UserQuery) (*api.Success, error) {
	return h.changeUserStatus(ctx, uQr, h.s.InactivateUser)
}

func (h *GRPCHandler) RestoreUser(ctx context.Context, uQr *api.UserQuery) (*api.Success, error) {
	return h.changeUserStatus(ctx, uQr, h.s.RestoreUser)
}

func (h *GRPCHandler) BanUser(ctx context.Context, uQr *api.UserQuery) (*api.Success, error) {
	return h.changeUserStatus(ctx, uQr, h.s.BanUser)
}

func (h *GRPCHandler) UnbanUser(ctx context.Context, uQr *api.UserQuery) (*api.Success, error) {
	return h.changeUserStatus(ctx, uQr, h.s.UnbanUser)
}

func (h *GRPCHandler) changeUserStatus(
	ctx context.Context,
	uQr *api.UserQuery,
	fn func(context.Context, domain.UserQuery) error,
) (*api.Success, error) {
	userQuery := domain.UserQuery{
		ID:       uQr.ID,
		Username: uQr.Username,
		Email:    uQr.Email,
	}

	ctx, cancel := context.WithTimeout(ctx, h.DefaultTimeout)
	defer cancel()

	if err := fn(ctx, userQuery); err != nil {
		return nil, h.handleError(err)
	}

	return &api.Success{
		Success: true,
	}, nil
}
