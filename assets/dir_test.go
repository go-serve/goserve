package assets

import (
	"os"

	"testing"
)

func Test_assetDirInfo(t *testing.T) {
	var i os.FileInfo = &assetDirInfo{}
	_ = i
	t.Log("*assetDirInfo{} implements os.FileInfo interface")
}
