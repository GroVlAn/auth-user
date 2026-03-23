package repository

import (
	"context"
	"fmt"

	"github.com/GroVlAn/auth-user/internal/core"
	"github.com/GroVlAn/auth-user/internal/core/e"
	"github.com/jmoiron/sqlx"
)

const (
	userTable = "auth_user"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, user core.User) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (id, username, email, password_hash, fullname, is_superuser, 
		created_at) VALUES (:id, :username, :email, :password_hash, :fullname, :is_superuser,
		:created_at)`,
		userTable,
	)

	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return e.NewErrInternal(
			fmt.Errorf("creating new user: %w", err),
		)
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, userQuery core.UserQuery) (core.User, error) {
	query := fmt.Sprintf(
		`SELECT id, username, email, password_hash, fullname, created_at FROM %s
		WHERE id = $1 OR username = $2 OR email = $3`,
		userTable,
	)

	var user core.User
	err := r.db.GetContext(ctx, &user, query, userQuery.ID, userQuery.Username, userQuery.Email)
	if err != nil {
		return core.User{}, handleQueryError(fmt.Errorf(
			"getting user: %w", err),
			"user not found",
		)
	}

	return user, nil
}

func (r *Repository) Exist(ctx context.Context, userQuery core.UserQuery) (bool, error) {
	query := fmt.Sprintf(
		`SELECT EXISTS(SELECT 1 FROM %s WHERE id=$1 OR username=$2 OR email=$3)`,
		userTable,
	)

	var exist bool
	err := r.db.GetContext(ctx, &exist, query, userQuery.ID, userQuery.Username, userQuery.Email)
	if err != nil {
		return false, e.NewErrInternal(err)
	}

	return exist, nil
}

func (r *Repository) BanUser(ctx context.Context, userID string) error {
	query := fmt.Sprintf(
		`UPDATE %s SET is_banned=true WHERE id=$1`,
		userTable,
	)

	if _, err := r.db.ExecContext(ctx, query, userID); err != nil {
		return e.NewErrInternal(err)
	}

	return nil
}

func (r *Repository) UnbanUser(ctx context.Context, userID string) error {
	query := fmt.Sprintf(
		`UPDATE %s SET is_banned=false WHERE id=$1`,
		userTable,
	)

	if _, err := r.db.ExecContext(ctx, query, userID); err != nil {
		return e.NewErrInternal(err)
	}

	return nil
}

func (r *Repository) InactivateUser(ctx context.Context, userID string) error {
	query := fmt.Sprintf(
		`UPDATE %s SET is_active=false WHERE id=$1`,
		userTable,
	)

	if _, err := r.db.ExecContext(ctx, query, userID); err != nil {
		return e.NewErrInternal(err)
	}

	return nil
}

func (r *Repository) RestoreUser(ctx context.Context, userID string) error {
	query := fmt.Sprintf(
		`UPDATE %s SET is_active=true WHERE id=$1`,
		userTable,
	)

	if _, err := r.db.ExecContext(ctx, query, userID); err != nil {
		return e.NewErrInternal(err)
	}

	return nil
}

func (r *Repository) DeleteInactiveUser(ctx context.Context) error {
	query := fmt.Sprintf(
		`DELETE FROM %s WHERE is_active=false`,
		userTable,
	)

	if _, err := r.db.ExecContext(ctx, query); err != nil {
		return e.NewErrInternal(err)
	}

	return nil
}
