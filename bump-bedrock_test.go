package main

import (
	"bytes"
	"fmt"
	"github.com/austinpray/bump-bedrock/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/austinpray/bump-bedrock/bedrock/mocks"
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

func mockTags() *httptest.Server {
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
	return ts
}

func TestGetVersion(t *testing.T) {

	// assert equality
	output := captureStdout(func() {
		GetVersion(Tags{
			Tag{
				Name: "4.2.1",
			},
		})
	})
	assert.Equal(t, "4.2.1\n", output, "they should be equal")

}

func TestGetWordPressTags(t *testing.T) {
	ts := mockTags()

	tags := GetWordPressTags(ts.URL)

	assert.NotNil(t, tags)
	assert.Equal(t, "4.2.1", tags[0].Name, "first tag should be 4.2.1")
	assert.Equal(t, 1, len(tags), "should have one array element")
}

func TestBump(t *testing.T) {
	testBedrock := new(mocks.BedrockRepo)

	testBedrock.On("UpdateWordPressVersion", "4.2.0").Return("4.2.0", nil)

	Bump(testBedrock, "4.2.0")

	testBedrock.AssertExpectations(t)

}
