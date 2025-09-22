package models

type FileStatus string
type TaskStatus string

const (
	Pending   FileStatus = "pending"
	Running   FileStatus = "running"
	Completed FileStatus = "completed"
	Failed    FileStatus = "failed"

	TaskPending   TaskStatus = "pending"
	TaskRunning   TaskStatus = "running"
	TaskCompleted TaskStatus = "completed"
	TaskFailed    TaskStatus = "failed"
)
