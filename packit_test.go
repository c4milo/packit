package packit

import (
	"bytes"
	"testing"

	"github.com/c4milo/unzipit"
	"github.com/hooklift/assert"
)

func TestZip(t *testing.T) {
	buf := new(bytes.Buffer)
	Zip("fixtures/myfiles", buf)

	_, err := unpackit.Unpack(buf, "")
	assert.Ok(t, err)
}

func TestTar(t *testing.T) {
	buf := new(bytes.Buffer)
	Tar("fixtures/myfiles", buf)

	_, err := unpackit.Unpack(buf, "")
	assert.Ok(t, err)
}

func TestGzip(t *testing.T) {
	tar := new(bytes.Buffer)
	tarGz := new(bytes.Buffer)

	Tar("fixtures/myfiles", tar)
	err := Gzip(tar, tarGz)
	assert.Ok(t, err)

	_, err = unpackit.Unpack(tarGz, "")
	assert.Ok(t, err)
}

func TestBzip2(t *testing.T) {
	tar := new(bytes.Buffer)
	tarBzip2 := new(bytes.Buffer)
	Tar("fixtures/myfiles", tar)
	err := Bzip2(tar, tarBzip2)
	assert.Ok(t, err)

	_, err = unpackit.Unpack(tarBzip2, "")
	assert.Ok(t, err)
}

func TestXz(t *testing.T) {
	tar := new(bytes.Buffer)
	tarXz := new(bytes.Buffer)
	Tar("fixtures/myfiles", tar)
	err := Xz(tar, tarXz)
	assert.Ok(t, err)

	_, err = unpackit.Unpack(tarXz, "")
	assert.Ok(t, err)
}
