package usecase

import (
	"context"
	"sync"
)

func (s *TaskService) Start(workers int) {
	if workers <= 0 {
		workers = 1
	}
	for i := 0; i < workers; i++ {
		go s.worker()
	}
}

func (s *TaskService) worker() {
	for id := range s.q.Pop() {
		s.runTask(context.Background(), id)
	}
}

func (s *TaskService) Stop(wg *sync.WaitGroup) {
	defer wg.Done()
	s.q.Close()
}
