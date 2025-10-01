package usecase

//
//import (
//	"context"
//	"test-task/internal/domain"
//	"time"
//)
//
//type TaskRepo interface {
//	SaveTask(*domain.Task) error
//	LoadTask(id string) (*domain.Task, error)
//	LoadAll() (map[string]*domain.Task, error)
//}
//
//type FileFetcher interface {
//	Fetch(ctx context.Context, url, suggestedName, outDir string) (name string, n int64, err error)
//}
//
//type Queue interface {
//	Push(id string)
//	Pop() <-chan string
//	Close()
//}
//
//type Logger interface {
//	Infof(string, ...any)
//	Errorf(string, ...any)
//}
//
//type Clock interface{ Now() time.Time }
