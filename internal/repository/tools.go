package repository

import (
	"database/sql"
	"errors"

	"github.com/GroVlAn/auth-user/internal/core/e"
)

func handleQueryError(err error, msg string) error {
	if errors.Is(err, sql.ErrNoRows) {
		return e.NewErrNotFound(
			err,
			msg,
		)
	}

	return e.NewErrInternal(
		err,
	)
}
