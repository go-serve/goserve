package server_test

import (
	"os"
	"strings"

	"github.com/go-serve/goserve/server"

	"testing"

	"sort"
)

func listNames(l []os.FileInfo) string {
	results := make([]string, len(l))
	for i, fileInfo := range l {
		results[i] = fileInfo.Name()
	}
	return strings.Join(results, ", ")
}

func TestByName(t *testing.T) {
	l := testList()
	var s sort.Interface = server.ByName(l)
	t.Log("server.ByName implements sort.Interface")

	// test sorting
	sort.Stable(s)

	// expectation
	expNames := []string{"A", "B", "C", "D", "E", "F", "G"}
	for i, expName := range expNames {
		if l[i].Name() != expName {
			t.Errorf("sorted result: l[%d].Name expected %#v, got %#v", i, expName, l[i].Name())
		}
	}

}

func TestByModTime(t *testing.T) {
	l := testList()
	var s sort.Interface = server.ByModTime(l)
	t.Log("server.ByModTime implements sort.Interface")

	// test sorting
	sort.Stable(s)

	// expectation
	expNames := []string{"A", "B", "C", "D", "E", "F", "G"}
	for i, expName := range expNames {
		if l[i].Name() != expName {
			t.Errorf("sorted result: l[%d].Name expected %#v, got %#v", i, expName, l[i].Name())
		}
	}

}

func TestByType(t *testing.T) {
	l := testList()

	var s sort.Interface = server.ByType(l)
	t.Log("server.ByType implements sort.Interface")

	// should sort all directory before files
	// test sorting (file first)
	sort.Stable(s)

	// expectation
	expNames := []string{"F", "G", "E", "B", "A", "C", "D"}
	for i, expName := range expNames {
		if l[i].Name() != expName {
			t.Errorf("sorted result: l[%d].Name expected %#v, got %#v", i, expName, l[i].Name())
		}
	}

}

func TestByNameAndType(t *testing.T) {
	l := testList()

	var s1 sort.Interface = server.ByName(l)
	var s2 sort.Interface = server.ByType(l)

	// test sorting
	sort.Stable(s1)
	sort.Stable(s2)

	// expectation
	expNames := []string{"E", "F", "G", "A", "B", "C", "D"}
	for i, expName := range expNames {
		if l[i].Name() != expName {
			t.Errorf("sorted result: l[%d].Name expected %#v, got %#v", i, expName, l[i].Name())
		}
	}
}
