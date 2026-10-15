package users_service

import (
	"context"
	"fmt"

	"github.com/darya-kznv/todo-service/internal/core/domain"
)

func (s *UsersService) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	initializedUser, err := user.Initialized(s.clock.Now())
	if err != nil {
		return domain.User{}, fmt.Errorf("initialize user domain: %w", err)
	}

	createdUser, err := s.usersRepository.CreateUser(ctx, initializedUser)
	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return createdUser, nil
}
