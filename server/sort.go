package server

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// QuerySort sort the files by the provided query string
func QuerySort(sortq string, files []os.FileInfo) (err error) {

	if sortq == "" {
		return
	}

	sorts := strings.Split(sortq, ",")
	n := len(sorts) - 1
	var si sort.Interface

	for i := n; i >= 0; i-- {
		si, err = SortBy(sorts[i], files)
		if err != nil {
			return
		}
		sort.Sort(si)
	}
	return
}

// SortBy sorts the files
func SortBy(by string, files []os.FileInfo) (s sort.Interface, err error) {

	asc := true // default order

	if by[0] == '-' {
		by = by[1:]
		asc = false
	}

	switch by {
	case "name":
		s = ByName(files)
	case "mtime":
		s = ByModTime(files)
	case "type":
		s = ByType(files)
	default:
		err = fmt.Errorf("unsupported sorting %#v", by)
		return
	}

	if asc == false {
		s = sort.Reverse(s)
	}
	return
}

// ByName implements sort.Interface
type ByName []os.FileInfo

// Len is the number of elements in the collection.
func (fi ByName) Len() int {
	return len(fi)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (fi ByName) Less(i, j int) bool {
	return fi[i].Name() < fi[j].Name()
}

// Swap swaps the elements with indexes i and j.
func (fi ByName) Swap(i, j int) {
	tmp := fi[i]
	fi[i] = fi[j]
	fi[j] = tmp
}

// ByModTime
type ByModTime []os.FileInfo

// Len is the number of elements in the collection.
func (fi ByModTime) Len() int {
	return len(fi)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (fi ByModTime) Less(i, j int) bool {
	return fi[i].ModTime().Before(fi[j].ModTime())
}

// Swap swaps the elements with indexes i and j.
func (fi ByModTime) Swap(i, j int) {
	tmp := fi[i]
	fi[i] = fi[j]
	fi[j] = tmp
}

// ByType
type ByType []os.FileInfo

// Len is the number of elements in the collection.
func (fi ByType) Len() int {
	return len(fi)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (fi ByType) Less(i, j int) bool {
	if fi[i].IsDir() != fi[j].IsDir() {
		return fi[i].IsDir()
	}
	return false
}

// Swap swaps the elements with indexes i and j.
func (fi ByType) Swap(i, j int) {
	tmp := fi[i]
	fi[i] = fi[j]
	fi[j] = tmp
}
