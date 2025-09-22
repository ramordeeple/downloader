package store

import "test-task/internal/models"

type Store interface {
	SaveTask(*models.Task) error
	LoadTask() (map[string]*models.Task, error)
}
