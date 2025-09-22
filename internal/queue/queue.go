package queue

type Queue struct{ ch chan string }

type Config struct{ Size int }

func New(cfg Config) *Queue {
	if cfg.Size <= 0 {
		cfg.Size = 100
	}
	return &Queue{ch: make(chan string, cfg.Size)}
}

func (q *Queue) Push(id string) {
	q.ch <- id
}

func (q *Queue) Jobs() <-chan string {
	return q.ch
}
