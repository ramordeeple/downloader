package downloader

import (
	"context"
	"time"
)

func sleepOrCancel(ctx context.Context, attempt int) bool {
	delay := time.Duration(attempt*attempt) * time.Millisecond * 300

	select {
	case <-time.After(delay):
		return true

	case <-ctx.Done():
		return false
	}
}
