package httpapi

import "test-task/internal/domain"

type TaskUsecase interface {
	NewTask(urls []string) (string, error)
	GetTask(id string) (*domain.Task, error)
}
