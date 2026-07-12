package domain

import (
	"time"

	"github.com/GroVlAn/auth-base/ew"
)

type User struct {
	ID           string    `json:"-" db:"id" valid:"require"`
	Username     string    `json:"username" db:"username" valid:"require"`
	Email        string    `json:"email" db:"email" valid:"require"`
	Password     string    `json:"password,omitempty" db:"-" valid:"require"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Fullname     string    `json:"fullname" db:"fullname" valid:"require"`
	IsSuperuser  bool      `json:"is_superuser" db:"is_superuser"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	IsBanned     bool      `json:"is_banned" db:"is_banned"`
	CreatedAt    time.Time `json:"create_at" db:"created_at"`
}

type UserInfo struct {
	Username string `json:"username" db:"username" valid:"require"`
	Email    string `json:"email" db:"email" valid:"require"`
	Fullname string `json:"fullname" db:"fullname" valid:"require"`
}

type UserQuery struct {
	ID       *string
	Username *string
	Email    *string
}

type UserQueryNewPassword struct {
	UserQuery
	OldPassword string `json:"old_password" valid:"require"`
	NewPassword string `json:"new_password" valid:"require"`
}

func (uq UserQuery) Validation() error {
	err := ew.NewErrValidation("validation user query data error")

	var count int

	if uq.ID != nil {
		count++
	}
	if uq.Username != nil {
		count++
	}
	if uq.Email != nil {
		count++
	}

	if count == 0 {
		err.AddField("id|username|email", "at least one field must be provided")
	}

	if err.IsEmpty() {
		return nil
	}

	return err
}
