package models

import "time"

type Task struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	Status    TaskStatus `json:"status"`
	Files     []File     `json:"files"`
}
