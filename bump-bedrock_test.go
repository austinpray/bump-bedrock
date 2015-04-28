package main

import (
	"bytes"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestGetVersion(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[
{
    "name": "4.2.1",
    "zipball_url": "https://api.github.com/repos/johnpbloch/wordpress/zipball/4.2.1",
    "tarball_url": "https://api.github.com/repos/johnpbloch/wordpress/tarball/4.2.1",
    "commit": {
      "sha": "c1cefa55c50dadb75b5e9f0e4844e420c794ab48",
      "url": "https://api.github.com/repos/johnpbloch/wordpress/commits/c1cefa55c50dadb75b5e9f0e4844e420c794ab48"
    }
	}
		]`)
	}))

	APITagsUrl = ts.URL

	// assert equality
	c := cli.NewContext(nil, nil, nil)
	output := captureStdout(func() {
		GetVersion(c)
	})
	assert.Equal(t, "4.2.1\n", output, "they should be equal")

}
