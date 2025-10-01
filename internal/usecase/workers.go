package usecase

import (
	"context"
	"sync"
)

// Start запускает воркеров, которые слушают очередь задач.
func (s *TaskService) Start() {
	for i := 0; i < s.cfg.Workers; i++ {
		go s.worker()
	}
}

// worker — отдельный воркер, который берёт задачи из очереди и обрабатывает их.
func (s *TaskService) worker() {
	for id := range s.q.Pop() {
		s.runTask(context.Background(), id)
	}
}

// Stop завершает работу очереди и ждёт завершения воркеров.
func (s *TaskService) Stop(wg *sync.WaitGroup) {
	defer wg.Done()
	s.q.Close()
}
