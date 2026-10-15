package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/darya-kznv/todo-service/internal/core/domain"
)

func (r *UsersRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
    INSERT INTO todoservice.users (id, full_name, phone_number, created_at)
    VALUES ($1, $2, $3, $4)
    RETURNING id, version, full_name, phone_number, created_at;`

	row := r.pool.QueryRow(ctx, query, user.ID, user.FullName, user.PhoneNumber, user.CreatedAt)

	var userModel UserModel

	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
		&userModel.CreatedAt,
	)

	if err != nil {
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
		userModel.CreatedAt,
	)

	return userDomain, nil
}
