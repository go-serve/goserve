package assets

import (
	"bytes"
	"net/http"
	"os"
	"strings"
)

const (
	TypeFile = iota
	TypeDir
)

// FileSystem returns a http.FileSystem of the assets
func FileSystem() (fs http.FileSystem) {
	return &fileSystem{}
}

type fileSystem struct {
}

// Open the given string or return error
func (fs *fileSystem) Open(name string) (f http.File, err error) {

	if name[0] == '/' {
		name = name[1:]
	}

	// test if is an unempty dir
	names, _ := AssetDir(name)
	if len(names) != 0 {
		f = &File{
			bytes.NewReader([]byte{}),
			name,
			TypeDir,
		}
		return
	}

	// test if is a file
	buf, err := Asset(name)
	if err != nil {
		if strings.HasSuffix(err.Error(), "not found") {
			err = os.ErrNotExist
		}
		return
	}
	f = &File{
		bytes.NewReader(buf),
		name,
		TypeFile,
	}
	return
}
