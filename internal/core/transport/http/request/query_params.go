package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/darya-kznv/todo-service/internal/core/errors"
)

func GetInitQueryParams(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)

	if param == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' not a valid integer: %v: %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	if val < 0 {
		return nil, fmt.Errorf(
			"param must be non-negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return &val, nil

}
