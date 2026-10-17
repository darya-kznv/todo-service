package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	core_errors "github.com/darya-kznv/todo-service/internal/core/errors"
	"github.com/google/uuid"
)

func GetStringQueryParam(r *http.Request, key string) *string {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil
	}

	return &param
}

func GetDateQueryParam(r *http.Request, key string) (*time.Time, error) {
	param := GetStringQueryParam(r, key)
	if param == nil {
		return nil, nil
	}

	layout := "2006-01-02"

	date, err := time.Parse(layout, *param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' is not a valid date: %v: %w",
			*param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return &date, nil
}

func GetUUIDQueryParam(r *http.Request, key string) (*uuid.UUID, error) {
	param := GetStringQueryParam(r, key)
	if param == nil {
		return nil, nil
	}

	id, err := uuid.Parse(*param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' is not a valid UUID: %v: %w",
			*param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return &id, nil
}

func GetInitQueryParams(r *http.Request, key string) (*int, error) {
	param := GetStringQueryParam(r, key)
	if param == nil {
		return nil, nil
	}

	val, err := strconv.Atoi(*param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' not a valid integer: %v: %w",
			*param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return &val, nil
}
