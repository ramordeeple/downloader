package models

type Task struct {
	ID        string     `json:"id"`
	CreatedAt string     `json:"created_at"`
	Status    TaskStatus `json:"status"`
	Files     []File     `json:"files"`
}
