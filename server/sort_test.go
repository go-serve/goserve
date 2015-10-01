package server_test

import (
	"github.com/yookoala/goserve/server"

	"testing"

	"sort"
)

func TestByName(t *testing.T) {
	var s sort.Interface = server.ByName{}
	_ = s
	t.Log("server.ByName implements sort.Interface")
}
