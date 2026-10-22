package web_fs_repository

import (
	"errors"
	"fmt"
	"os"

	"github.com/darya-kznv/todo-service/internal/core/domain"
	core_errors "github.com/darya-kznv/todo-service/internal/core/errors"
)

func (r *WebRepository) ReadFile(path string) (domain.File, error) {
	buffer, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return domain.File{}, fmt.Errorf(
				"read web file %q: %w",
				path,
				core_errors.ErrNotFound,
			)
		}

		return domain.File{}, fmt.Errorf("read web file %q: %w", path, err)
	}

	return domain.NewFile(buffer), nil
}
