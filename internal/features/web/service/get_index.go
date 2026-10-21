package web_service

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/darya-kznv/todo-service/internal/core/domain"
)

func (s *WebService) GetIndex() (domain.File, error) {
	projectRoot := os.Getenv("PROJECT_ROOT")
	if projectRoot == "" {
		return domain.File{}, fmt.Errorf("PROJECT_ROOT is not set")
	}

	filePath := filepath.Join(projectRoot, "public", "index.html")
	file, err := s.webRepository.ReadFile(filePath)
	if err != nil {
		return domain.File{}, fmt.Errorf("get index page from repository: %w", err)
	}

	return file, nil
}
