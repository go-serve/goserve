package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"log"
)

var errNotDir = errors.New("not a directory")

// FileServer returns our custom goserve file server
func FileServer(root http.FileSystem) http.Handler {
	return &fileServer{
		root:    root,
		fileSrv: http.FileServer(root),
	}
}

// custom implementation of FileServer
type fileServer struct {
	root    http.FileSystem
	fileSrv http.Handler
}

// ServeHTTP implements http.Handler
func (fs *fileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if d, err := fs.ReadDirInfo(r.URL.Path); err == nil {
		log.Printf("%s, %#v", r.URL.Path, d)
		fmt.Fprint(w, "Hello listing")
		return
	}
	fs.fileSrv.ServeHTTP(w, r)
}

// ReadDirInfo determines if a given path is directory
// and return error if not dir (or other error in file system)
func (fs *fileServer) ReadDirInfo(path string) (f os.FileInfo, err error) {

	fh, err := fs.root.Open(path)
	if err != nil {
		return
	}

	f, err = fh.Stat()
	if err != nil {
		return
	}

	if !f.IsDir() {
		err = errNotDir
	}
	return
}
