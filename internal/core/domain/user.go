package domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/KyoshiBlame/TodoKy/internal/core/errors"
)

type User struct {
	ID      int
	Version int

	FullName    string
	PhoneNumber *string //ссылка т.к номер не обязателен
}

func NewUser(
	id int,
	version int,
	fullname string,
	phone_number *string,
) User {
	return User{
		ID:          id,
		Version:     version,
		FullName:    fullname,
		PhoneNumber: phone_number,
	}
}

func NewUserUninitialized(
	fullname string,
	phone_number *string,
) User {
	return NewUser(
		UninitializedID,
		UninitializedVersion,
		fullname,
		phone_number,
	)
}

func (u *User) Validate() error {
	fullNameLen := len([]rune(u.FullName))

	if fullNameLen < 3 || fullNameLen > 100 {
		return fmt.Errorf(
			"invalid `FullName` length: %d: %w",
			fullNameLen,
			core_errors.ErrInvalidArgument,
		)
	}

	if u.PhoneNumber != nil {
		phoneNumberLen := len([]rune(*u.PhoneNumber))

		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf(
				"invalid `PhoneNumber` length: %d: %w",
				phoneNumberLen,
				core_errors.ErrInvalidArgument,
			)
		}

		re := regexp.MustCompile(`^\+[0-9]+$`)
		if !re.MatchString(*u.PhoneNumber) {
			return fmt.Errorf(
				"invalid `PhoneNumber` format: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}
