package users_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/darya-kznv/todo-service/internal/core/domain"
)

func (s *UsersService) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("validate user domain: %w", err)
	}

	readyUser := domain.NewUser(
		uuid.New(),
		1,
		user.FullName,
		user.PhoneNumber,
	)

	user, err := s.usersRepository.CreateUser(ctx, readyUser)
	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}
