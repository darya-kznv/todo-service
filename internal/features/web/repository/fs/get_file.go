package web_fs_repository

import (
	"fmt"
	"os"

	"github.com/darya-kznv/todo-service/internal/core/domain"
	core_errors "github.com/darya-kznv/todo-service/internal/core/errors"
)

func (r *WebRepository) ReadFile(path string) (domain.File, error) {
	buffer, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
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
