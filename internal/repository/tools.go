package repository

import (
	"database/sql"
	"errors"

	"github.com/GroVlAn/auth-base/ew"
)

func handleQueryError(err error, msg string) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ew.New(
			ew.ErrorTypeNotFound,
			err,
		).Msg(msg)
	}

	return ew.New(
		ew.ErrorTypeInternal,
		err,
	)
}
