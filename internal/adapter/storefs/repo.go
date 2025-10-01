package storefs

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"test-task/internal/domain"
	"test-task/internal/usecase"
)

// FileRepo реализует usecase.TaskRepo поверх файловой системы.
type FileRepo struct {
	root string
}

// New возвращает новый репозиторий.
func New(root string) *FileRepo {
	return &FileRepo{root: root}
}

// SaveTask сохраняет задачу в файл <root>/<taskID>/task.json
func (r *FileRepo) SaveTask(t *domain.Task) error {
	if t == nil || t.ID == "" {
		return errors.New("invalid task: nil or no ID")
	}
	dir := filepath.Join(r.root, t.ID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	file := filepath.Join(dir, "task.json")

	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0o644)
}

// LoadTask загружает одну задачу по ID.
func (r *FileRepo) LoadTask(id string) (*domain.Task, error) {
	file := filepath.Join(r.root, id, "task.json")
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var t domain.Task
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// LoadAll загружает все задачи из root.
func (r *FileRepo) LoadAll() (map[string]*domain.Task, error) {
	res := make(map[string]*domain.Task)

	entries, err := os.ReadDir(r.root)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return res, nil
		}
		return nil, err
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		file := filepath.Join(r.root, e.Name(), "task.json")
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		var t domain.Task
		if json.Unmarshal(data, &t) == nil && t.ID != "" {
			res[t.ID] = &t
		}
	}
	return res, nil
}

var _ usecase.TaskRepo = (*FileRepo)(nil) // compile-time check
