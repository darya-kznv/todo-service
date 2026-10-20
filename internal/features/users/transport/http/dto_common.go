package users_transport_http

import (
	"time"

	"github.com/darya-kznv/todo-service/internal/core/domain"
	"github.com/google/uuid"
)

type UserDTOResponse struct {
	ID          uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Version     int       `json:"version" example:"1"`
	FullName    string    `json:"full_name" example:"Ivan Ivanov"`
	PhoneNumber *string   `json:"phone_number" example:"+79998887766"`
	CreatedAt   time.Time `json:"created_at" example:"2026-01-15T10:30:00Z"`
}

func userDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:          user.ID,
		Version:     user.Version,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt,
	}
}

func usersDTOFromDomains(users []domain.User) []UserDTOResponse {
	usersDTO := make([]UserDTOResponse, len(users))

	for i, user := range users {
		usersDTO[i] = userDTOFromDomain(user)
	}

	return usersDTO
}
