package server

import (
	"os"
)

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
