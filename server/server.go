package server

import (
	"github.com/yookoala/goserve/assets"

	"errors"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
)

var errNotDir = errors.New("not a directory")
var errIsDir = errors.New("is a directory")

// FileServer returns our custom goserve file server
func FileServer(root http.FileSystem) http.Handler {
	return &fileServer{
		root:    root,
		fileSrv: http.FileServer(root),
		assets:  http.FileServer(assets.FileSystem()),
	}
}

// custom implementation of FileServer
type fileServer struct {
	root    http.FileSystem
	fileSrv http.Handler
	assets  http.Handler
}

// ServeHTTP implements http.Handler
func (fs *fileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Printf("access %#v", r.URL.Path)

	// serve assets
	if r.URL.Path == "/_goserve" {
		http.Redirect(w, r, "/_goserve/", http.StatusMovedPermanently)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/_goserve/") {
		r.URL.Path = r.URL.Path[9:]
		fs.assets.ServeHTTP(w, r)
		return
	}

	// serve directory indexes
	if d, err := fs.ReadDirInfo(r.URL.Path); err == nil {

		if _, err := fs.ReadIndex(r.URL.Path); err != nil {

			files, err := d.Readdir(0)
			if err != nil {
				log.Printf("Error listing path %#v:%s", r.URL.Path, err)
				return
			}

			w.Header().Add("Content-Type", "text/html")
			fmt.Fprint(w, "<h1>Index</h1>")
			for _, file := range files {
				fmt.Fprintf(w, "<ul>")
				fmt.Fprintf(w, "<li><a href=\"%s\">%s</a></li>",
					path.Join(r.URL.Path, file.Name()), file.Name())
				fmt.Fprintf(w, "</ul>")
			}
			return

		}

	} else {
		log.Printf("Error reading path %#v: %s", r.URL.Path, err)
	}

	// fallback to normal file server
	fs.fileSrv.ServeHTTP(w, r)
}

// ReadDirInfo determines if a given path is directory
// and return error if not dir (or other error in file system)
func (fs *fileServer) ReadDirInfo(path string) (f http.File, err error) {

	f, err = fs.root.Open(path)
	if err != nil {
		return
	}

	fi, err := f.Stat()
	if err != nil {
		return
	}

	if !fi.IsDir() {
		err = errNotDir
	}

	return
}

func (fs *fileServer) ReadIndex(path string) (f http.File, err error) {

	const indexPage = "/index.html"

	f, err = fs.root.Open(path + indexPage)
	if err != nil {
		return
	}

	fi, err := f.Stat()
	if err != nil {
		return
	}

	if fi.IsDir() {
		err = errIsDir
	}

	return
}
