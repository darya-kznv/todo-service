package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/darya-kznv/todo-service/internal/core/domain"
)

const (
	defaultLimit  = 10
	defaultOffset = 0
)

func (r *UsersRepository) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	limitVal := defaultLimit
	if limit != nil {
		limitVal = *limit
	}

	offsetVal := defaultOffset
	if offset != nil {
		offsetVal = *offset
	}

	query := `
	SELECT id, version, full_name, phone_number, created_at
	FROM todoservice.users
	ORDER BY created_at DESC 
	LIMIT $1 OFFSET $2;`

	rows, err := r.pool.Query(
		ctx,
		query,
		limitVal,
		offsetVal,
	)

	if err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}

	defer rows.Close()

	var userModels []UserModel

	for rows.Next() {
		var userModel UserModel

		err := rows.Scan(
			&userModel.ID,
			&userModel.Version,
			&userModel.FullName,
			&userModel.PhoneNumber,
			&userModel.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("scan users: %w", err)
		}

		userModels = append(userModels, userModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	userDomains := userDomainsFromModels(userModels)

	return userDomains, nil
}
