package users_service

import (
	"context"

	"github.com/google/uuid"

	core_clock "github.com/darya-kznv/todo-service/internal/core/clock"
	"github.com/darya-kznv/todo-service/internal/core/domain"
)

type UsersService struct {
	usersRepository UsersRepository
	clock           core_clock.Clock
}

type UsersRepository interface {
	CreateUser(
		ctx context.Context,
		domain domain.User,
	) (domain.User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)

	GetUser(
		ctx context.Context,
		id uuid.UUID,
	) (domain.User, error)

	DeleteUser(
		ctx context.Context,
		id uuid.UUID,
	) error

	PatchUser(
		ctx context.Context,
		id uuid.UUID,
		user domain.User,
	) (domain.User, error)
}

func NewUsersService(
	usersRepository UsersRepository,
	clock core_clock.Clock,
) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
		clock:           clock,
	}
}
