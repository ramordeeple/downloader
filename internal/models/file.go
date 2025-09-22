package models

type File struct {
	URL    string     `json:"url"`
	Name   string     `json:"name"`
	Status FileStatus `json:"status"`
	Error  string     `json:"error, omitempty"`
}
