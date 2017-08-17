package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Link contains a HATEOAS hypermedia reference URL
type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

// FileInfo is a JSON display of a subset of os.FileInfo information
type FileInfo struct {
	Name     string    `json:"name"`
	Path     string    `json:"path,omitempty"`
	Type     string    `json:"type"`
	Mime     string    `json:"mime,omitempty"`
	HasIndex bool      `json:"hasIndex,omitempty"`
	Size     int64     `json:"size,omitempty"`
	MTime    time.Time `json:"mtime,omitempty"`
	Links    []Link    `json:"links,omitempty"`
}

// FileStat stores and display a file's information as JSON
type FileStat struct {
	Name  string
	Path  string
	Size  int64
	MTime time.Time
}

// MarshalJSON implements encoding/json.Marshaler
func (file FileStat) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type  string    `json:"type"`
		Name  string    `json:"name"`
		Path  string    `json:"path"`
		Size  int64     `json:"size"`
		MTime time.Time `json:"mtime"`
	}{
		Type:  "file",
		Name:  file.Name,
		Path:  file.Path,
		Size:  file.Size,
		MTime: file.MTime,
	})
}

// DirStat stores and display a directory's information as JSON
type DirStat struct {
	Name  string
	Path  string
	MTime time.Time
}

// MarshalJSON implements encoding/json.Marshaler
func (file DirStat) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type  string    `json:"type"`
		Name  string    `json:"name"`
		Path  string    `json:"path"`
		MTime time.Time `json:"mtime"`
	}{
		Type:  "directory",
		Name:  file.Name,
		Path:  file.Path,
		MTime: file.MTime,
	})
}

// StatError represents an error in JSON format
type StatError struct {
	Code int
	Path string
}

// Message return message for a given error
func (err StatError) Message() string {
	msg := http.StatusText(err.Code)
	if msg == "" {
		return "unknown error"
	}
	return msg
}

// Error implements error interface
func (err StatError) Error() string {
	return fmt.Sprintf("error %d: %s", err.Code, err.Message())
}

// NewStatError returns a new StatError
func NewStatError(code int, path string) *StatError {
	return &StatError{
		Code: code,
		Path: path,
	}
}

// MarshalJSON implements encoding/json.Marshaler
func (err StatError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Status  string `json:"status"`
		Code    int    `json:"code"`
		Path    string `json:"path"`
		Message string `json:"message"`
	}{
		Status:  "error",
		Code:    err.Code,
		Path:    err.Path,
		Message: err.Message(),
	})
}
