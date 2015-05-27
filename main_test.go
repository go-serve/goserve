package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestFileServer(t *testing.T) {

	test_dir := os.Getenv("TEST_DIR")
	if test_dir == "" {
		log.Fatal("Please provide the TEST_DIR environment parameter for testing")
	}

	ts := httptest.NewServer(fileServer(test_dir))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}

	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	// test end result
	if string(greeting) != "Hello\n" {
		log.Fatalf("Unexpected result: \"%s\"", greeting)
	}

}
