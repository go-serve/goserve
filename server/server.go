package server

import (
	"github.com/go-midway/midway"
	"github.com/go-serve/goserve/assets"

	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"
)

const assetsPath = "/_goserve/assets"

var errNotDir = errors.New("not a directory")
var errIsDir = errors.New("is a directory")
var tplIndex *template.Template

var stylesheets []string
var scripts []string

func init() {

	fs := assets.FileSystem()
	fh, err := fs.Open("/html/index.html")
	if err != nil {
		log.Print("Failed to load template")
		panic(err)
	}

	b, err := ioutil.ReadAll(fh)
	if err != nil {
		log.Print("Failed to read template file")
		panic(err)
	}

	tplIndex = template.New("index.html")

	// add utility functions to templates
	tplIndex = tplIndex.Funcs(template.FuncMap{
		"joinPath": func(parts ...string) string {
			return path.Join(parts...)
		},
	})

	// parse template
	tplIndex, err = tplIndex.Parse(string(b))
	if err != nil {
		log.Print("Failed to parse index.html into template")
		panic(err)
	}

	// common stylesheets to use
	stylesheets = []string{
		assetsPath + "/css/app.css",
	}

	// common scripts to use
	scripts = []string{
		assetsPath + "/js/app.js",
	}

	// NODE_ENV check
	if os.Getenv("NODE_ENV") == "development" {
		stylesheets = []string{}
		scripts = []string{
			"http://localhost:8081" + assetsPath + "/js/app.js",
		}
	}

}

// ServeAssets generates a middleware that serves file assets
func ServeAssets(path string, root http.FileSystem) midway.Middleware {

	path = strings.TrimRight(path, "/")
	pathWithSlash := path + "/"
	pathLen := len(pathWithSlash)
	fileAssets := http.FileServer(root)

	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// serve assets
			if r.URL.Path == path {
				http.Redirect(w, r, pathWithSlash, http.StatusMovedPermanently)
				return
			}
			if strings.HasPrefix(r.URL.Path, pathWithSlash) {
				r.URL.Path = r.URL.Path[pathLen:] // strip base path
				fileAssets.ServeHTTP(w, r)
				return
			}

			// defers to inner handler
			inner.ServeHTTP(w, r)
		})
	}
}

// FileServer returns our custom goserve file server
func FileServer(root http.FileSystem) http.Handler {
	fserver := &fileServer{
		root:    root,
		fileSrv: http.FileServer(root),
	}
	middlewares := midway.Chain(
		ServeAPI("/_goserve/api", root),
		ServeAssets("/_goserve/assets", assets.FileSystem()),
		ServeVideo(root),
	)
	return middlewares(fserver)
}

// custom implementation of FileServer
type fileServer struct {
	root    http.FileSystem
	fileSrv http.Handler
}

// ServeHTTP implements http.Handler
func (fs *fileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Printf("access %#v", r.URL.Path)

	// serve directory indexes
	if d, err := fs.ReadDirInfo(r.URL.Path); err == nil {

		if _, err := fs.ReadIndex(r.URL.Path); err != nil {

			files, err := d.Readdir(0)
			if err != nil {
				log.Printf("Error listing path %#v:%s", r.URL.Path, err)
				return
			}

			// sort according to query
			s := r.URL.Query().Get("sort")
			if s == "" {
				// default sort order: by mtime, desc
				s = "-mtime"
			}
			QuerySort(s, files) // TODO: add error reporting here

			// list the files
			listFiles(w, r.URL.Path, files)
			return

		}

	} else if err != errNotDir {
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

func listFiles(w http.ResponseWriter, base string, files []os.FileInfo) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	tplIndex.Execute(w, map[string]interface{}{
		"Stylesheets": stylesheets,
		"Scripts":     scripts,
		"Files":       files,
		"Base":        base,
	})
}
