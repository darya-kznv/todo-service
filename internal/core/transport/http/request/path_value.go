package core_http_request

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"

	core_errors "github.com/darya-kznv/todo-service/internal/core/errors"
)

func GetUUIDPathValue(r *http.Request, key string) (uuid.UUID, error) {
	pathValue := r.PathValue(key)

	if pathValue == "" {
		return uuid.Nil, fmt.Errorf(
			"no key='%s' in path values: %w",
			key,
			core_errors.ErrInvalidArgument,
		)
	}

	parsedUUID, err := uuid.Parse(pathValue)
	if err != nil {
		return uuid.Nil, fmt.Errorf(
			"path value='%s' by key='%s' is not a valid UUID: %w",
			pathValue,
			key,
			core_errors.ErrInvalidArgument,
		)
	}

	return parsedUUID, nil
}
