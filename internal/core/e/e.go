package e

import (
	"errors"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrPasswordMismatch  = errors.New("password mismatch")
)
