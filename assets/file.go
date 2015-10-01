package assets

import (
	"bytes"
	"os"
	"path"
)

// File implments http.File
type File struct {
	*bytes.Reader
	name string
	t    int
}

// Close is a dummy method to implment io.Closer
func (f *File) Close() error {
	return nil
}

// Readdir reads the contents of the directory associated with
// file and returns a slice of up to n FileInfo values, as would
// be returned by Lstat, in directory order. Subsequent calls on
// the same file will yield further FileInfos.
func (f *File) Readdir(count int) (lfi []os.FileInfo, err error) {
	names, err := AssetDir(f.name)
	lfi = make([]os.FileInfo, 0)
	var fi os.FileInfo

	i := 1
	for _, name := range names {
		fi, err = AssetInfo(path.Join(f.name, name))
		lfi = append(lfi, fi)
		if i == count {
			break
		}
		i++
	}

	return
}

// Stat returns the FileInfo structure describing file. If
// there is an error, it will be of type *PathError.
func (f *File) Stat() (fi os.FileInfo, err error) {
	if f.t == TypeDir {
		fi = &assetDirInfo{}
		return
	}
	return AssetInfo(f.name)
}
