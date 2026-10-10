package domain

import (
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"

	core_errors "github.com/darya-kznv/todo-service/internal/core/errors"
)

type User struct {
	ID          uuid.UUID
	Version     int
	FullName    string
	PhoneNumber *string
	CreatedAt   time.Time
}

func NewUser(
	id uuid.UUID,
	version int,
	fullName string,
	phoneNumber *string,
	createdAt time.Time,
) User {
	return User{
		ID:          id,
		Version:     version,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
		CreatedAt:   createdAt,
	}
}

func NewUserWithUUID(fullName string, phoneNumber *string) User {
	return User{
		ID:          uuid.New(),
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func NewUserUninitialized(
	fullName string,
	phoneNumber *string,
) User {
	return NewUser(
		UninitializedID,
		UninitializedVersion,
		fullName,
		phoneNumber,
		UninitializedCreatedAt,
	)
}

func (u *User) Validate() error {
	fullNameLength := len([]rune(u.FullName))

	if fullNameLength < 3 || fullNameLength > 100 {
		return fmt.Errorf(
			"invalid `FullName` len: %d, %w",
			fullNameLength,
			core_errors.ErrInvalidArgument,
		)
	}

	if u.PhoneNumber != nil {
		PhoneNumberLength := len([]rune(*u.PhoneNumber))

		if PhoneNumberLength < 10 || PhoneNumberLength > 15 {
			return fmt.Errorf(
				"invalid `PhoneNumber` len: %d, %w",
				PhoneNumberLength,
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
