package service

import (
	"context"
	"fmt"
	"time"

	"github.com/GroVlAn/auth-base/ew"
	"github.com/GroVlAn/auth-user/internal/domain"
	"github.com/GroVlAn/auth-user/internal/domain/e"
	"github.com/google/uuid"
)

const (
	invalidPassword = "invalid password"
)

type repo interface {
	Create(ctx context.Context, user domain.User) error
	User(ctx context.Context, userQuery domain.UserQuery) (domain.User, error)
	UserInfo(ctx context.Context, userQuery domain.UserQuery) (domain.UserInfo, error)
	UpdatePassword(ctx context.Context, userID, newPasswordHash string) error
	Exist(ctx context.Context, userQuery domain.UserQuery) (bool, error)
	BanUser(ctx context.Context, userID string) error
	UnbanUser(ctx context.Context, userID string) error
	InactivateUser(ctx context.Context, userID string) error
	RestoreUser(ctx context.Context, userID string) error
	DeleteInactiveUser(ctx context.Context) error
}

type hasher interface {
	Hash(password string) (string, error)
	Compare(encodedHash, password string) error
}

type accessGRPCClient interface {
	BindUserRole(ctx context.Context, userID string) error
}

type Service struct {
	repo         repo
	accessClient accessGRPCClient
	hasher       hasher
}

func New(repo repo, hasher hasher, accessClient accessGRPCClient) *Service {
	return &Service{
		repo:         repo,
		accessClient: accessClient,
		hasher:       hasher,
	}
}

func (s *Service) Create(ctx context.Context, user domain.User) error {
	if err := validateUser(user); err != nil {
		return err
	}

	exist, err := s.repo.Exist(ctx, domain.UserQuery{
		Username: &user.Username,
		Email:    &user.Email,
	})
	if err != nil {
		return ew.New(
			ew.ErrorTypeInternal,
			fmt.Errorf("checking if user exist: %w", err),
		)
	}
	if exist {
		return ew.New(
			ew.ErrorTypeConflict,
			e.ErrUserAlreadyExists,
		).Msg(e.ErrUserAlreadyExists.Error())
	}

	user.ID = uuid.NewString()

	passwordHash, err := s.hasher.Hash(user.Password)
	if err != nil {
		return fmt.Errorf("creating password hash: %w", err)
	}

	user.PasswordHash = string(passwordHash)

	user.CreatedAt = time.Now()

	if err = s.repo.Create(ctx, user); err != nil {
		return fmt.Errorf("creating user: %w", err)
	}

	if err := s.accessClient.BindUserRole(ctx, user.ID); err != nil {
		return fmt.Errorf("binding user role: %w", err)
	}

	return nil
}

func (s *Service) User(ctx context.Context, userQuery domain.UserQuery) (domain.User, error) {
	if err := userQuery.Validation(); err != nil {
		return domain.User{}, fmt.Errorf("validating user query: %w", err)
	}

	user, err := s.repo.User(ctx, userQuery)
	if err != nil {
		return domain.User{}, fmt.Errorf("getting user: %w", err)
	}

	return user, nil
}

func (s *Service) UserInfo(ctx context.Context, userQuery domain.UserQuery) (domain.UserInfo, error) {
	if err := userQuery.Validation(); err != nil {
		return domain.UserInfo{}, fmt.Errorf("validating user query: %w", err)
	}

	userInfo, err := s.repo.UserInfo(ctx, userQuery)
	if err != nil {
		return domain.UserInfo{}, fmt.Errorf("getting user info: %w", err)
	}

	return userInfo, nil
}

func (s *Service) UpdatePassword(ctx context.Context, userQueryNewPassword domain.UserQueryNewPassword) error {
	if err := userQueryNewPassword.UserQuery.Validation(); err != nil {
		return fmt.Errorf("validating user query: %w", err)
	}

	user, err := s.repo.User(ctx, userQueryNewPassword.UserQuery)
	if err != nil {
		return fmt.Errorf("getting user: %w", err)
	}

	if err := s.verifyOldPassword(user.PasswordHash, userQueryNewPassword.OldPassword); err != nil {
		return fmt.Errorf("verifying password: %w", err)
	}

	if ok, reason := validatePassword(userQueryNewPassword.NewPassword); !ok {
		return ew.NewErrValidation(reason)
	}

	if err := s.verifyNewPassword(user.PasswordHash, userQueryNewPassword.NewPassword); err != nil {
		return fmt.Errorf("verifying password: %w", err)
	}

	newPasswordHash, err := s.hasher.Hash(userQueryNewPassword.NewPassword)
	if err != nil {
		return fmt.Errorf("creating password hash: %w", err)
	}

	if err = s.repo.UpdatePassword(ctx, user.ID, newPasswordHash); err != nil {
		return fmt.Errorf("changing user password: %w", err)
	}

	return nil
}

func (s *Service) InactivateUser(ctx context.Context, userQuery domain.UserQuery) error {
	user, err := s.User(ctx, userQuery)
	if err != nil {
		return fmt.Errorf("getting user: %w", err)
	}

	if err := s.repo.InactivateUser(ctx, user.ID); err != nil {
		return fmt.Errorf("inactivating user: %w", err)
	}

	return nil
}

func (s *Service) RestoreUser(ctx context.Context, userQuery domain.UserQuery) error {
	user, err := s.User(ctx, userQuery)
	if err != nil {
		return fmt.Errorf("getting user: %w", err)
	}

	if err := s.repo.RestoreUser(ctx, user.ID); err != nil {
		return fmt.Errorf("restoring user: %w", err)
	}

	return nil
}

func (s *Service) BanUser(ctx context.Context, userQuery domain.UserQuery) error {
	user, err := s.User(ctx, userQuery)
	if err != nil {
		return fmt.Errorf("getting user: %w", err)
	}

	if err := s.repo.BanUser(ctx, user.ID); err != nil {
		return fmt.Errorf("banning user: %w", err)
	}

	return nil
}

func (s *Service) UnbanUser(ctx context.Context, userQuery domain.UserQuery) error {
	user, err := s.User(ctx, userQuery)
	if err != nil {
		return fmt.Errorf("getting user: %w", err)
	}

	if err := s.repo.UnbanUser(ctx, user.ID); err != nil {
		return fmt.Errorf("unbanning user: %w", err)
	}

	return nil
}

func (s *Service) DeleteInactiveUser(ctx context.Context) error {
	if err := s.repo.DeleteInactiveUser(ctx); err != nil {
		return fmt.Errorf("deleting inactive users: %w", err)
	}

	return nil
}

func (s *Service) verifyNewPassword(oldHash, newPassword string) error {
	err := s.hasher.Compare(oldHash, newPassword)
	if err == nil {
		return ew.NewErrValidation("new password must be different from old password")
	}

	return nil
}

func (s *Service) verifyOldPassword(passwordHash, oldPassword string) error {
	err := s.hasher.Compare(passwordHash, oldPassword)
	if err != nil {
		return ew.New(
			ew.ErrorTypeUnauthorized,
			fmt.Errorf("comparing hash and password: %w", err),
		).Msg(invalidPassword)
	}

	return nil
}
