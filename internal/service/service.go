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

type repo interface {
	Create(ctx context.Context, user core.User) error
	User(ctx context.Context, userQuery core.UserQuery) (core.User, error)
	BanUser(ctx context.Context, userID string) error
	Exist(ctx context.Context, userQuery core.UserQuery) (bool, error)
	UnbanUser(ctx context.Context, userID string) error
	InactivateUser(ctx context.Context, userID string) error
	RestoreUser(ctx context.Context, userID string) error
	DeleteInactiveUser(ctx context.Context) error
}

type Service struct {
	repo     repo
	hashCost int
}

func New(repo repo, hashCost int) *Service {
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
