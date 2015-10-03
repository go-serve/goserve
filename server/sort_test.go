package server_test

import (
	"github.com/yookoala/goserve/server"

	"testing"

	"sort"
)

func TestByName(t *testing.T) {
	l := testList()
	var s sort.Interface = server.ByName(l)
	t.Log("server.ByName implements sort.Interface")

	// test sorting
	sort.Sort(s)

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
	sort.Sort(s)

	// expectation
	expNames := []string{"A", "B", "C", "D", "E", "F", "G"}
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

	t.Log("server.ByType implements sort.Interface")

	// test sorting
	sort.Sort(s1)
	sort.Sort(s2)

	// expectation
	expNames := []string{"E", "F", "G", "A", "B", "C", "D"}
	for i, expName := range expNames {
		if l[i].Name() != expName {
			t.Errorf("sorted result: l[%d].Name expected %#v, got %#v", i, expName, l[i].Name())
		}
	}

}
