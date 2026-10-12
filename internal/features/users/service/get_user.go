package users_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/darya-kznv/todo-service/internal/core/domain"
	core_errors "github.com/darya-kznv/todo-service/internal/core/errors"
)

func (s *UsersService) GetUser(
	ctx context.Context,
	id uuid.UUID,
) (domain.User, error) {
	user, err := s.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf(
			"get user from repository: %w",
			core_errors.ErrNotFound,
		)
	}

	return user, nil
}
