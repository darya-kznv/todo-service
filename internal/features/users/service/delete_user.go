package users_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	core_errors "github.com/darya-kznv/todo-service/internal/core/errors"
)

func (s *UsersService) DeleteUser(
	ctx context.Context,
	id uuid.UUID,
) error {
	err := s.usersRepository.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf(
			"delete user from repository: %w",
			core_errors.ErrNotFound,
		)
	}

	return nil
}
