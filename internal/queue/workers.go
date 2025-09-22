package queue

import (
	"context"
)

type Runner func(ctx context.Context, id string)

func Start(ctx context.Context, n int, q *Queue, run Runner) {
	if n <= 0 {
		n = 2
	}
	for i := 0; i < n; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case id := <-q.ch:
					run(ctx, id)
				}
			}
		}()
	}
}
