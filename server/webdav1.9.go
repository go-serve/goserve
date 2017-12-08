// +build go1.9

package server

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/go-midway/midway"

	"golang.org/x/net/webdav"
)

// readonlyFile implements a readonly wrapper for http.File
type readonlyFile struct {
	http.File
}

func (f readonlyFile) Write(p []byte) (n int, err error) {
	return 0, os.ErrPermission
}

// webdavFS wraps an http.FileSystem and implement webdav.FileSystem for it
type webdavFS struct {
	fs http.FileSystem
}

func (fs *webdavFS) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return &os.PathError{
		Op:   "mkdir",
		Path: name,
		Err:  os.ErrPermission,
	}
}

func (fs *webdavFS) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	file, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return &readonlyFile{file}, nil
}

func (fs *webdavFS) RemoveAll(ctx context.Context, name string) error {
	return &os.PathError{
		Op:   "rm",
		Path: name,
		Err:  os.ErrPermission,
	}
}

func (fs *webdavFS) Rename(ctx context.Context, oldName, newName string) error {
	return &os.PathError{
		Op:   "rename",
		Path: oldName,
		Err:  os.ErrPermission,
	}
}

func (fs *webdavFS) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	file, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return file.Stat()
}

// WebdavFS returns a readonly webdav.FileSystem with reference to
// an http.FileSystem
func WebdavFS(fs http.FileSystem) webdav.FileSystem {
	return &webdavFS{fs}
}

// ServeWebdav serves a readonly webdav file system
// created from a given http.FileSystem
func ServeWebdav(root http.FileSystem) midway.Middleware {
	webdavHandler := &webdav.Handler{
		Prefix:     "",
		FileSystem: WebdavFS(root),
		LockSystem: webdav.NewMemLS(),
		// TODO: add Logger implementation
	}
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch strings.ToUpper(r.Method) {
			case "OPTIONS":
				fallthrough
			case "PROPFIND":
				fallthrough
			case "PROPPATCH":
				fallthrough
			case "MKCOL":
				fallthrough
			case "POST":
				fallthrough
			case "DELETE":
				fallthrough
			case "COPY":
				fallthrough
			case "MOVE":
				fallthrough
			case "LOCK":
				fallthrough
			case "UNLOCK":
				// TODO: properly logging
				webdavHandler.ServeHTTP(w, r)
			default:
				inner.ServeHTTP(w, r)
			}
		})
	}
}
