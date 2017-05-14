package packit

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/hooklift/assert"
)

func TestZip(t *testing.T) {
	buf := new(bytes.Buffer)
	Zip("fixtures/myfiles", buf)

	f, err := os.Create("boom.zip")
	assert.Ok(t, err)
	defer func() {
		err := f.Close()
		assert.Ok(t, err)
	}()

	_, err = io.Copy(f, buf)
	assert.Ok(t, err)
}

func TestTar(t *testing.T) {
	buf := new(bytes.Buffer)
	Tar("fixtures/myfiles", buf)

	f, err := os.Create("boom.tar")
	assert.Ok(t, err)
	defer func() {
		err := f.Close()
		assert.Ok(t, err)
	}()

	_, err = io.Copy(f, buf)
	assert.Ok(t, err)
}

func TestGzip(t *testing.T) {
	tar := new(bytes.Buffer)
	targz := new(bytes.Buffer)
	Tar("fixtures/myfiles", tar)
	Gzip(tar, targz)
}
