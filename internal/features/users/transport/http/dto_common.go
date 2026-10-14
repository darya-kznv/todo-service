package users_transport_http

import (
	"time"

	"github.com/darya-kznv/todo-service/internal/core/domain"
	"github.com/google/uuid"
)

type UserDTOResponse struct {
	ID          uuid.UUID `json:"id"`
	Version     int       `json:"version"`
	FullName    string    `json:"full_name"`
	PhoneNumber *string   `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
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
