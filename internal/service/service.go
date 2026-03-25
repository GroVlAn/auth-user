package service

import (
	"context"
	"fmt"
	"time"

	"github.com/GroVlAn/auth-user/internal/core"
	"github.com/GroVlAn/auth-user/internal/core/e"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Repo interface {
	Create(ctx context.Context, user core.User) error
	User(ctx context.Context, userQuery core.UserQuery) (core.User, error)
	UserInfo(ctx context.Context, userQuery core.UserQuery) (core.UserInfo, error)
	ChangePassword(ctx context.Context, userID, newPasswordHash string) error
	Exist(ctx context.Context, userQuery core.UserQuery) (bool, error)
	BanUser(ctx context.Context, userID string) error
	UnbanUser(ctx context.Context, userID string) error
	InactivateUser(ctx context.Context, userID string) error
	RestoreUser(ctx context.Context, userID string) error
	DeleteInactiveUser(ctx context.Context) error
}

type Service struct {
	repo     Repo
	hashCost int
}

func New(repo Repo, hashCost int) *Service {
	return &Service{
		repo:     repo,
		hashCost: hashCost,
	}
}

func (s *Service) Create(ctx context.Context, user core.User) error {
	if err := validateUser(user); err != nil {
		return err
	}

	exist, err := s.repo.Exist(ctx, core.UserQuery{
		Username: user.Username,
		Email:    user.Email,
	})
	if err != nil {
		return e.NewErrInternal(
			fmt.Errorf("checking if user exist: %w", err),
		)
	}
	if exist {
		return e.NewErrConflict(
			e.ErrUserAlreadyExists,
			e.ErrUserAlreadyExists.Error(),
		)
	}

	user.ID = uuid.NewString()

	passwordHash, err := passwordHash(user.Password, s.hashCost)
	if err != nil {
		return fmt.Errorf("creating password hash: %w", err)
	}

	user.PasswordHash = string(passwordHash)

	user.CreatedAt = time.Now()

	if err = s.repo.Create(ctx, user); err != nil {
		return fmt.Errorf("creating user: %w", err)
	}

	return nil
}

func (s *Service) User(ctx context.Context, userQuery core.UserQuery) (core.User, error) {
	if err := s.validateUserQuery(userQuery); err != nil {
		return core.User{}, fmt.Errorf("validating user query: %w", err)
	}

	user, err := s.repo.User(ctx, userQuery)
	if err != nil {
		return core.User{}, fmt.Errorf("getting user: %w", err)
	}

	return user, nil
}

func (s *Service) UserInfo(ctx context.Context, userQuery core.UserQuery) (core.UserInfo, error) {
	if err := s.validateUserQuery(userQuery); err != nil {
		return core.UserInfo{}, fmt.Errorf("validating user query: %w", err)
	}

	userInfo, err := s.repo.UserInfo(ctx, userQuery)
	if err != nil {
		return core.UserInfo{}, fmt.Errorf("getting user info: %w", err)
	}

	return userInfo, nil
}

func (s *Service) ChangePassword(ctx context.Context, userQueryNewPassword core.UserQueryNewPassword) error {
	if err := s.validateUserQuery(userQueryNewPassword.UserQuery); err != nil {
		return fmt.Errorf("validating user query: %w", err)
	}

	user, err := s.repo.User(ctx, userQueryNewPassword.UserQuery)
	if err != nil {
		return fmt.Errorf("getting user: %w", err)
	}

	if ok, reason := validatePassword(userQueryNewPassword.NewPassword); !ok {
		return e.NewErrValidation(reason)
	}

	if err := s.verifyOldPassword(userQueryNewPassword.OldPassword, user.PasswordHash); err != nil {
		return fmt.Errorf("verifying password: %w", err)
	}

	newPasswordHash, err := passwordHash(userQueryNewPassword.NewPassword, s.hashCost)
	if err != nil {
		return fmt.Errorf("creating password hash: %w", err)
	}

	if err = s.repo.ChangePassword(ctx, user.ID, newPasswordHash); err != nil {
		return fmt.Errorf("changing user password: %w", err)
	}

	return nil
}

func (s *Service) InactivateUser(ctx context.Context, userQuery core.UserQuery) error {
	user, err := s.User(ctx, userQuery)
	if err != nil {
		return fmt.Errorf("getting user: %w", err)
	}

	if err := s.repo.InactivateUser(ctx, user.ID); err != nil {
		return fmt.Errorf("inactivating user: %w", err)
	}

	return nil
}

func (s *Service) RestoreUser(ctx context.Context, userQuery core.UserQuery) error {
	user, err := s.User(ctx, userQuery)
	if err != nil {
		return fmt.Errorf("getting user: %w", err)
	}

	if err := s.repo.RestoreUser(ctx, user.ID); err != nil {
		return fmt.Errorf("restoring user: %w", err)
	}

	return nil
}

func (s *Service) BanUser(ctx context.Context, userQuery core.UserQuery) error {
	user, err := s.User(ctx, userQuery)
	if err != nil {
		return fmt.Errorf("getting user: %w", err)
	}

	if err := s.repo.BanUser(ctx, user.ID); err != nil {
		return fmt.Errorf("banning user: %w", err)
	}

	return nil
}

func (s *Service) UnbanUser(ctx context.Context, userQuery core.UserQuery) error {
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

func passwordHash(password string, hashCost int) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	if err != nil {
		return "", e.NewErrInternal(
			fmt.Errorf("hashing password: %w", err),
		)
	}

	return string(passwordHash), nil
}

func (s *Service) validateUserQuery(userQuery core.UserQuery) *e.ErrValidation {
	err := e.NewErrValidation("validation user query data error")

	if userQuery.ID == "" && userQuery.Username == "" && userQuery.Email == "" {
		err.AddField("id|username|email", "at least one field must be provided")
	}

	if err.IsEmpty() {
		return nil
	}

	return err

}

func (s *Service) verifyOldPassword(oldPassword, newPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(newPassword), []byte(oldPassword))
	if err == nil {
		return e.NewErrValidation("new password must be different from old password")
	}

	return nil
}
