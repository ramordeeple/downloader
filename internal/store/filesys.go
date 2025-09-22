package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"test-task/internal/models"
)

type FileSys struct{ Root string }

func NewFileSys(root string) *FileSys { return &FileSys{Root: root} }

func (fsys *FileSys) SaveTask(t *models.Task) error {
	if t == nil || t.ID == "" {
		return errors.New("Task is either nil or has no ID")
	}

	dir := filepath.Join(fsys.Root, t.ID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	file := filepath.Join(dir, "task.json")
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", " ")
	enc.SetEscapeHTML(false)
	return enc.Encode(t)
}

func (fsys *FileSys) LoadTask() (map[string]*models.Task, error) {
	res := map[string]*models.Task{}
	entries, err := os.ReadDir(fsys.Root)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return res, nil
		}
		return nil, err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		b, err := os.ReadFile(filepath.Join(fsys.Root, entry.Name(), "task.json"))
		if err != nil {
			continue
		}
		var t models.Task
		if json.Unmarshal(b, &t) == nil && t.ID != "" {
			res[t.ID] = &t
		}
	}
	return res, nil
}
