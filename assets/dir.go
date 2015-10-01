package assets

import (
	"os"
	"time"
)

type assetDirInfo struct {
	name string
	size int64
}

// Name gives base name of the file
func (fi *assetDirInfo) Name() string {
	return fi.name
}

// Size gives length in bytes for regular files;
// system-dependent for others
func (fi *assetDirInfo) Size() int64 {
	return fi.size
}

// Mode gives file mode bits
func (fi *assetDirInfo) Mode() os.FileMode {
	return os.ModeDir
}

// ModTime gives modification time
func (fi *assetDirInfo) ModTime() (t time.Time) {
	return t
}

// IsDir is abbreviation for Mode().IsDir()
func (fi *assetDirInfo) IsDir() bool {
	return fi.Mode().IsDir()
}

// Sys gives underlying data source (can return nil)
func (fi *assetDirInfo) Sys() interface{} {
	return nil
}
