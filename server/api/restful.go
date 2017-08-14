package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-midway/midway"
)

// Link contains a HATEOAS hypermedia reference URL
type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

// FileInfo is a JSON display of a subset of os.FileInfo information
type FileInfo struct {
	Name  string    `json:"name"`
	Path  string    `json:"path,omitempty"`
	Type  string    `json:"type"`
	Mime  string    `json:"mime,omitempty"`
	Size  int64     `json:"size,omitempty"`
	MTime time.Time `json:"mtime,omitempty"`
	Links []Link    `json:"links,omitempty"`
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

func statsEndpoint(ctx context.Context, req interface{}) (stats interface{}, err error) {

	path := req.(string)

	// TODO: rewrite with FileSystem
	stat, err := os.Stat(path)

	// if file not found
	if os.IsNotExist(err) {
		err = NewStatError(http.StatusNotFound, path)
		return
	}

	// permission problem
	if err != nil {
		perr, _ := err.(*os.PathError)
		if perr.Err.Error() == os.ErrPermission.Error() {
			err = NewStatError(http.StatusForbidden, path)
		}
		return
	}

	// for files
	if stat.Mode().IsRegular() {

		// test permission
		var file *os.File
		file, err = os.OpenFile(path, os.O_RDONLY, 0444)
		if err != nil {
			perr, _ := err.(*os.PathError)
			if perr.Err.Error() == os.ErrPermission.Error() {
				err = NewStatError(http.StatusForbidden, path)
			}
			return
		}
		file.Close() // close immediately

		stats = FileStat{
			Name:  stat.Name(),
			Path:  path,
			Size:  stat.Size(),
			MTime: stat.ModTime(),
		}
		return
	}

	// for directories
	if stat.Mode().IsDir() {
		stats = DirStat{
			Name:  stat.Name(),
			Path:  path,
			MTime: stat.ModTime(),
		}
		return
	}

	return
}

func listEndpoint(ctx context.Context, req interface{}) (resp interface{}, err error) {
	path := req.(string)
	if path == "" {
		path = "."
	}

	// TODO: build the absolute file / dir path for stat and open
	stat, err := os.Stat(path)

	// if file not found
	if os.IsNotExist(err) {
		err = NewStatError(http.StatusNotFound, path)
		return
	}

	// permission problem
	if err != nil {
		perr, _ := err.(*os.PathError)
		if perr.Err.Error() == os.ErrPermission.Error() {
			err = NewStatError(http.StatusForbidden, path)
		}
		return
	}

	// for directories
	if stat.Mode().IsDir() {

		var d *os.File
		files := make([]os.FileInfo, 0, 40)
		// TODO: use FileSystem for file access
		if d, err = os.Open(path); err != nil {
			log.Printf("Error listing path %#v:%s", path, err)
			err = NewStatError(http.StatusInternalServerError, path)
			return
		}
		defer d.Close()

		files, err = d.Readdir(0)
		if err != nil {
			log.Printf("Error listing path %#v:%s", path, err)
			return
		}

		// sort according to query
		epCtx := getEndpointContext(ctx)
		s := epCtx.Sort
		if s == "" {
			s = "-mtime"
		}
		// TODO: rewrite with go-linq
		QuerySort(s, files) // TODO: add error reporting here

		listLen := len(files)
		list := make([]FileInfo, listLen)
		for i := 0; i < listLen; i++ {
			item := files[i]

			// parse item URL
			itemPath := path + "/" + item.Name()
			if path == "." {
				itemPath = item.Name()
			}

			if item.Mode().IsRegular() {
				list[i] = FileInfo{
					Name:  item.Name(),
					Type:  "file",
					Path:  itemPath,
					Size:  item.Size(),
					MTime: item.ModTime(),
					Links: []Link{
						{
							Rel:  "self",
							Href: epCtx.Scheme + "://" + epCtx.Host + "/" + itemPath,
						},
						{
							Rel:  "stat",
							Href: epCtx.Scheme + "://" + epCtx.Host + "/_goserve/api/stats/" + itemPath,
						},
					},
				}
			} else if item.IsDir() {
				list[i] = FileInfo{
					Name:  item.Name(),
					Type:  "directory",
					Path:  itemPath,
					MTime: item.ModTime(),
					Links: []Link{
						{
							Rel:  "self",
							Href: epCtx.Scheme + "://" + epCtx.Host + "/" + itemPath,
						},
						{
							Rel:  "stat",
							Href: epCtx.Scheme + "://" + epCtx.Host + "/_goserve/api/stats/" + itemPath,
						},
						{
							Rel:  "list",
							Href: epCtx.Scheme + "://" + epCtx.Host + "/_goserve/api/lists/" + itemPath,
						},
					},
				}
			} else {
				list[i] = FileInfo{
					Name: item.Name(),
					Type: "other",
					Path: itemPath,
					Links: []Link{
						{
							Rel:  "self",
							Href: epCtx.Scheme + "://" + epCtx.Host + "/" + itemPath,
						},
						{
							Rel:  "stat",
							Href: epCtx.Scheme + "://" + epCtx.Host + "/_goserve/api/stats/" + itemPath,
						},
					},
				}
			}
		}

		resp = struct {
			Items []FileInfo `json:"items"`
		}{
			Items: list,
		}
		return
	}

	err = NewStatError(http.StatusBadRequest, path)
	return
}

func handleEndpoint(endpoint func(ctx context.Context, req interface{}) (resp interface{}, err error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := context.Background()

		// prepare context
		if r != nil {
			ctx = withEndpointContext(ctx, r)
		}

		// handle path request
		resp, err := endpoint(ctx, r.URL.Path)

		// handle error
		if err != nil {
			switch serr := err.(type) {
			case *StatError:
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(serr.Code)
				jsonw := json.NewEncoder(w)
				jsonw.Encode(serr)
			default:
				statusCode := http.StatusInternalServerError
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(statusCode)
				jsonw := json.NewEncoder(w)
				jsonw.Encode(struct {
					Code    int    `json:"code"`
					Status  string `json:"status"`
					Message string `json:"message"`
				}{
					Code:    statusCode,
					Status:  "error",
					Message: err.Error(),
				})
			}
			return
		}

		// handle normal response
		w.Header().Set("Content-Type", "application/json")
		jsonw := json.NewEncoder(w)
		jsonw.Encode(resp)

		log.Printf("resp: %#v", resp)
	}
}

// ServeAPI generates a middleware to serve API for file / directory information
// query
func ServeAPI(path string, root http.FileSystem) midway.Middleware {

	path = strings.TrimRight(path, "/") // strip trailing slash
	pathWithSlash := path + "/"
	pathLen := len(pathWithSlash)

	// wrap endpoints
	handleStats := handleEndpoint(statsEndpoint)
	handleList := handleEndpoint(listEndpoint)
	handleGraphQL := GraphQLHandler()

	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// serve API endpoint
			if r.URL.Path == path {
				http.Redirect(w, r, pathWithSlash, http.StatusMovedPermanently)
				return
			}
			if r.URL.Path == path+"/graphql" {
				graphCtx := withFilesystem(withEndpointContext(r.Context(), r), root)
				handleGraphQL.ServeHTTP(w, r.WithContext(graphCtx))
				return
			}
			if strings.HasPrefix(r.URL.Path, pathWithSlash) {
				r.URL.Path = strings.TrimRight(r.URL.Path[pathLen:], "/") // strip base path

				// stats of file / directory
				if strings.HasPrefix(r.URL.Path, "stats/") {
					r.URL.Path = r.URL.Path[6:]
					handleStats(w, r)
					return
				}

				// listing files in directory
				if strings.HasPrefix(r.URL.Path, "lists/") {
					r.URL.Path = r.URL.Path[6:]
					handleList(w, r)
					return
				}
				if r.URL.Path == "lists" {
					r.URL.Path = r.URL.Path[5:]
					handleList(w, r)
					return
				}

				// if no matching endpoint
				statusCode := http.StatusNotFound
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(statusCode)
				jsonw := json.NewEncoder(w)
				jsonw.Encode(struct {
					Code    int    `json:"code"`
					Status  string `json:"status"`
					Message string `json:"message"`
				}{
					Code:    statusCode,
					Status:  "error",
					Message: "not a valid API endpoint",
				})

				return
			}
			// server file / directory info query at the URL
			if r.Header.Get("Content-Type") == "application/goserve+json" {
				// TODO: also detect the request content-type: "goserve+json/application"
				// and return file info
			}

			// defers to inner handler
			inner.ServeHTTP(w, r)
		})
	}
}
