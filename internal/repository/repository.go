package repository

import "github.com/jmoiron/sqlx"

const (
	userTable = "auth_user"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}
