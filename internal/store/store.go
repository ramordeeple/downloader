package store

import "test-task/internal/models"

type Store interface {
	SaveTask(*models.Task) error
	LoadTasks() (map[string]*models.Task, error)
}
