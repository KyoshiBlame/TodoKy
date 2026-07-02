package domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/KyoshiBlame/TodoKy/internal/core/errors"
)



type User struct {
	ID int
	Version int

	FullName string
	PhoneNumber *string//ссылка т.к номер не обязателен
}

func NewUser(
	id int,
	version int,
	fullname string,
	phone_number *string,
) User {
	return User{
		ID: id,
		Version: version,
		FullName: fullname,
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
	fullNameLens := len([]rune(u.FullName))

	if fullNameLens < 3 || fullNameLens > 100 {
		return fmt.Errorf(
			"invalid `Fullnames` lenght: %d: %w",
			fullNameLens,
			core_errors.ErrInvalidArgument,
		)
	}

	if u.PhoneNumber != nil {
		phomeNumberLen := len([]rune(*u.PhoneNumber))
		if phomeNumberLen < 10 || phomeNumberLen > 15 {
			return fmt.Errorf(
				"invalid `PhoneNumbers` lenght %d:%w",
				phomeNumberLen,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	re := regexp.MustCompile(`^\+[0-9]+$`)

	if !re.MatchString(*u.PhoneNumber) {
		return fmt.Errorf(
			"invalid `PhoneNumbers` format: %w",
				core_errors.ErrInvalidArgument,
		)
	}

	return nil
}