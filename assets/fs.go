package assets

import (
	"net/http"

	"github.com/go-serve/bindatafs"
	"golang.org/x/tools/godoc/vfs/httpfs"
)

// FileSystem returns a http.FileSystem of the assets
func FileSystem() (fs http.FileSystem) {
	return httpfs.New(bindatafs.New("assets://", Asset, AssetDir, AssetInfo))
}
