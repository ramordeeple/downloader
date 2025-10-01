package queue

type Queue struct {
	ch chan string
}

func New(size int) *Queue {
	if size <= 0 {
		size = 100
	}
	return &Queue{ch: make(chan string, size)}
}

func (q *Queue) Push(id string) {
	q.ch <- id
}

func (q *Queue) Pop() <-chan string {
	return q.ch
}

func (q *Queue) Close() {
	close(q.ch)
}
