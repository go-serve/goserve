package api_test

import (
	"os"
	"time"
)

func testList() []os.FileInfo {
	return []os.FileInfo{
		dummyFileInfo{
			name:    "B",
			modTime: time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
			isDir:   false,
		},
		dummyFileInfo{
			name:    "A",
			modTime: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			isDir:   false,
		},
		dummyFileInfo{
			name:    "C",
			modTime: time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
			isDir:   false,
		},
		dummyFileInfo{
			name:    "D",
			modTime: time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
			isDir:   false,
		},
		dummyFileInfo{
			name:    "F",
			modTime: time.Date(2005, 1, 1, 0, 0, 0, 0, time.UTC),
			isDir:   true,
		},
		dummyFileInfo{
			name:    "G",
			modTime: time.Date(2006, 1, 1, 0, 0, 0, 0, time.UTC),
			isDir:   true,
		},
		dummyFileInfo{
			name:    "E",
			modTime: time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
			isDir:   true,
		},
	}
}

type dummyFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
	sys     interface{}
}

// base name of the file
func (fi dummyFileInfo) Name() string {
	return fi.name
}

// length in bytes for regular files; system-dependent for others
func (fi dummyFileInfo) Size() int64 {
	return fi.size
}

// file mode bits
func (fi dummyFileInfo) Mode() os.FileMode {
	return fi.mode
}

// modification time
func (fi dummyFileInfo) ModTime() time.Time {
	return fi.modTime
}

// abbreviation for Mode().IsDir()
func (fi dummyFileInfo) IsDir() bool {
	return fi.isDir
}

// underlying data source (can return nil)
func (fi dummyFileInfo) Sys() interface{} {
	return fi.sys
}
