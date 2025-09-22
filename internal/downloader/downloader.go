package downloader

import "context"

type Downloader interface {
	Fetch(ctx context.Context, rawURL, suggestedName string) (savedName string, size int64, err error)
}
