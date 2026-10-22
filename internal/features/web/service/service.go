package web_service

import "github.com/darya-kznv/todo-service/internal/core/domain"

type WebRepository interface {
	ReadFile(path string) (domain.File, error)
}

type WebService struct {
	webRepository WebRepository
}

func NewWebService(
	webRepository WebRepository,
) *WebService {
	return &WebService{
		webRepository: webRepository,
	}
}
