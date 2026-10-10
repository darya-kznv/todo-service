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
	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("validate user domain: %w", err)
	}

	userWithUUID := domain.NewUserWithUUID(user.FullName, user.PhoneNumber)

	createdUser, err := s.usersRepository.CreateUser(ctx, userWithUUID)
	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return createdUser, nil
}
