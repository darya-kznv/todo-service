package domain

import (
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"

	core_errors "github.com/darya-kznv/todo-service/internal/core/errors"
)

var phoneRegex = regexp.MustCompile(`^\+[0-9]+$`)

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

func NewUserWithUUID(fullName string, phoneNumber *string) User {
	return User{
		ID:          uuid.New(),
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
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

		if !phoneRegex.MatchString(*u.PhoneNumber) {
			return fmt.Errorf(
				"invalid `PhoneNumber` format: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}

type UserPatch struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
}

func NewUserPatch(
	fullName Nullable[string],
	phoneNumber Nullable[string],
) UserPatch {
	return UserPatch{
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func (u *UserPatch) Validate() error {
	if u.FullName.Set && u.FullName.Value == nil {
		return fmt.Errorf(
			"'FullName' can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
	}

	tmp := *u

	if patch.FullName.Set {
		tmp.FullName = *patch.FullName.Value
	}

	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = patch.PhoneNumber.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}

	*u = tmp

	return nil
}
