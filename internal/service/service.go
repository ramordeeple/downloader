package service

import (
	"sync"
	"test-task/internal/downloader"
	"test-task/internal/models"
	"test-task/internal/queue"
	"test-task/internal/store"
)

type Config struct {
	QueueSize int
	Workers   int
}

type Service struct {
	mutex sync.Mutex
	tasks map[string]*models.Task
	st    store.Store
	dl    downloader.Downloader
	q     *queue.Queue
}
