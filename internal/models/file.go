package models

type File struct {
	URL    string     `json:"url"`
	Name   string     `json:"name"`
	Status FileStatus `json:"status"`
	Error  string     `json:"error,omitempty"`

	SizeBytes int64  `json:"size_bytes,omitempty"`
	BytesDone int64  `json:"bytes_done,omitempty"`
	ETag      string `json:"etag,omitempty"`
	LastMod   string `json:"last_mod,omitempty"`
}
