package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/yookoala/goserve/server"
)

var port *uint64
var dir string

func init() {

	var envPort uint64 = 8080 // default port
	var err error

	// parse env for PORT parameter
	if envPortStr := os.Getenv("PORT"); envPortStr != "" {
		envPort, err = strconv.ParseUint(envPortStr, 10, 16)
		if err != nil {
			log.Fatalf(
				"Cannot parse \"%s\" as PORT. Must be unsigned integer", envPortStr)
		}
	}

	// flags, if any provided, may override the default
	port = flag.Uint64("port", envPort, "Determine the port to serve")
	flag.Parse()

	// read directory from remaining argument
	// or use current directory
	if flag.NArg() == 1 {
		log.Printf("run here!!!")
		dir = flag.Arg(0)
	} else if flag.NArg() > 1 {
		log.Fatalf("Too many argument. goserve can only serve one directory")
	} else {
		// get current dir
		dir, err = os.Getwd()
		if err != nil {
			log.Fatalf("Failed to parse current path: %s", err.Error())
		}
	}

}

// test if a path is a valid directory
func validDir(path string) (err error) {
	var f *os.File
	var fi os.FileInfo
	if f, err = os.Open(dir); err != nil {
		return
	}
	if fi, err = f.Stat(); err != nil {
		return
	}
	if !fi.IsDir() {
		err = fmt.Errorf("Path \"%s\" is not a valid directory", path)
		return
	}
	return
}

func fileServer(root string) http.Handler {
	return server.FileServer(http.Dir(root))
}

func main() {

	// check if provided dir is a valid dir
	if err := validDir(dir); err != nil {
		log.Fatal(err)
	}

	// port string for server
	portStr := fmt.Sprintf(":%d", *port)

	// some logs before starting
	log.Printf("Listening to port %d", *port)
	log.Printf("Serving path: %s", dir)
	log.Fatal(http.ListenAndServe(portStr, fileServer(dir)))
}
