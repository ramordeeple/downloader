package storefs

import (
	"encoding/json"
	"test-task/internal/domain"
)

func encode(t *domain.Task) ([]byte, error) {
	return json.MarshalIndent(t, "", "  ")
}

func decode(data []byte) (*domain.Task, error) {
	var t domain.Task
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}
	return &t, nil
}
