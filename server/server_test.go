package server_test

import (
	"github.com/yookoala/goserve/server"

	"log"
	"net/http"
	"net/http/httptest"

	"testing"
)

func TestFileServer(t *testing.T) {

	test_dir := "./../_example"

	th := server.FileServer(http.Dir(test_dir))

	req, err := http.NewRequest("GET", "http://example.com/", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	th.ServeHTTP(w, req)

	// test end result
	greeting := w.Body.String()
	expRes := "Hello\n"
	if string(greeting) != expRes {
		log.Fatalf("Unexpected result: %#v;\n"+
			"Expected: %#v", greeting, expRes)
	}

}
