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
	userWithUUID := domain.NewUserWithUUID(user.FullName, user.PhoneNumber)
	if err := userWithUUID.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("validate user domain: %w", err)
	}

	createdUser, err := s.usersRepository.CreateUser(ctx, userWithUUID)
	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return createdUser, nil
}
