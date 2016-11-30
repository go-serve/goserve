package assets_test

import (
	"net/http"

	"github.com/go-serve/goserve/assets"

	"testing"
)

func TestFile(t *testing.T) {
	var f http.File = &assets.File{}
	_ = f
	t.Logf("*assets.File implements http.File interface")
}
