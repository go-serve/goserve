package assets_test

import (
	"net/http"
	"os"

	"github.com/go-serve/goserve/assets"

	"testing"
)

func TestFileSystem(t *testing.T) {
	var fs http.FileSystem = assets.FileSystem()
	_ = fs
	t.Log("assets.FileSystem() returns a proper http.FileSystem")
}

func TestOpenSuccess(t *testing.T) {
	fs := assets.FileSystem()
	_, err := fs.Open("/js/jquery/jquery-1.11.3.min.js")
	if err != nil {
		t.Errorf("Failed opening file \"%s\"", err)
	}
}

func TestOpenFail(t *testing.T) {
	fs := assets.FileSystem()
	_, err := fs.Open("no-such-file")
	if err == nil {
		t.Error("Expected to see error but didn't get it")
	}
	if !os.IsNotExist(err) {
		t.Errorf("Error did not satisify os.IsNotExists(): %#v", err)
	}
}

func TestList1(t *testing.T) {
	fs := assets.FileSystem()
	dir := "/js/jquery"
	d, err := fs.Open(dir)
	if err != nil {
		t.Errorf("Failed openning %#v: %s", dir, err)
	}

	fl, err := d.Readdir(1)
	if err != nil {
		t.Errorf("Failed listing %#v: %s", dir, err)
	}

	if len(fl) != 1 {
		t.Errorf("len(fl) is %d, expected %d", len(fl), 1)
	}
}

func TestListAll(t *testing.T) {
	fs := assets.FileSystem()
	dir := "/js/jquery"
	d, err := fs.Open(dir)
	if err != nil {
		t.Errorf("Failed openning %#v: %s", dir, err)
	}

	fl, err := d.Readdir(0)
	if err != nil {
		t.Errorf("Failed listing %#v: %s", dir, err)
	}

	expLen := 2
	if len(fl) != expLen {
		t.Errorf("len(fl) is %d, expected %d", len(fl), expLen)
	}
}
