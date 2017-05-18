package packit

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/c4milo/packit/fastwalk"
	"github.com/dsnet/compress/bzip2"
	"github.com/klauspost/compress/zip"
	gzip "github.com/klauspost/pgzip"
	"github.com/pkg/errors"
	"github.com/ulikunitz/xz"
)

// Zip walks the file tree rooted at root and archives it into the output stream using zip
// algorithm.
func Zip(root string, out io.ReadWriter) {
	zw := zip.NewWriter(out)
	defer func() {
		if err := zw.Close(); err != nil {
			fmt.Printf("packit/zip: failed closing archive: %v\n", err)
		}
	}()

	sep := string(os.PathSeparator)
	m := &sync.Mutex{}
	fastwalk.Walk(root, func(path string, mode os.FileMode) error {
		// This function gets invoked concurrently for each file or directory found.
		// But, the zip package does not support parallelism, so we need to make sure
		// a file header is followed by its correspondent content.
		m.Lock()
		defer m.Unlock()

		// Appending a final "/" is critical to let the decompressor know this is a directory.
		if mode.IsDir() && !strings.HasSuffix(path, sep) {
			path += sep
		}

		fi, err := os.Lstat(path)
		if err != nil {
			return errors.Wrapf(err, "packit/zip: failed getting file stats for: %s", path)
		}

		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return errors.Wrapf(err, "packit/zip: failed populating file header for: %s", path)
		}

		// fi.Name returns the base name and we need the full path.
		fh.Name = path

		fw, err := zw.CreateHeader(fh)
		if err != nil {
			return errors.Wrapf(err, "packit/zip: failed creating header for: %s", path)
		}

		// For directories, we only need to add the header in the zip file.
		if mode.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return errors.Wrapf(err, "packit/zip: failed opening file at: %s", path)
		}

		defer func() {
			if err := f.Close(); err != nil {
				fmt.Printf("packit/zip: failed closing file %q: %v\n", path, err)
			}
		}()

		_, err = io.Copy(fw, f)
		if err != nil {
			return errors.Wrapf(err, "packit/zip: failed packing data from: %s", path)
		}

		return nil
	})
}

// Tar walks the file tree rooted at root and archives it all using Tar.
func Tar(root string, out io.ReadWriter) {
	tw := tar.NewWriter(out)
	defer func() {
		if err := tw.Close(); err != nil {
			fmt.Printf("packit/tar: failed closing archive: %v\n", err)
		}
	}()

	sep := string(os.PathSeparator)
	m := &sync.Mutex{}
	fastwalk.Walk(root, func(path string, mode os.FileMode) error {
		m.Lock()
		defer m.Unlock()

		// Appending a trailing path separator is critical to let the decompressor
		// know this is a directory.
		if mode.IsDir() && !strings.HasSuffix(path, sep) {
			path += sep
		}

		fi, err := os.Lstat(path)
		if err != nil {
			return errors.Wrapf(err, "packit/tar: failed getting file stats for: %s", path)
		}

		fh, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return errors.Wrapf(err, "packit/tar: failed populating file header for: %s", path)
		}

		// fi.Name returns the base name and we need the full path.
		fh.Name = path

		tw.WriteHeader(fh)
		if err != nil {
			return errors.Wrapf(err, "packit/tar: failed creating header for: %s", path)
		}

		// For directories, we only need to add the header in the tar file.
		if mode.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return errors.Wrapf(err, "packit/tar: failed opening file at: %s", path)
		}

		defer func() {
			if err := f.Close(); err != nil {
				fmt.Printf("packit/tar: failed closing file: %s\n", path)
			}
		}()

		_, err = io.Copy(tw, f)
		if err != nil {
			return errors.Wrapf(err, "packit/tar: failed packing data from: %s", path)
		}
		return nil
	})
}

// Gzip compresses an input stream using gzip.
func Gzip(in io.Reader, out io.ReadWriter) error {
	gw := gzip.NewWriter(out)
	defer func() {
		if err := gw.Close(); err != nil {
			fmt.Printf("packit/gzip: failed closing stream: %v\n", err)
		}
	}()

	if _, err := io.Copy(gw, in); err != nil {
		return err
	}
	return nil
}

// Xz compresses an input stream using xz.
func Xz(in io.Reader, out io.ReadWriter) error {
	xw, err := xz.NewWriter(out)
	if err != nil {
		return err
	}

	defer func() {
		if err := xw.Close(); err != nil {
			fmt.Printf("packit/xz: failed closing stream: %v\n", err)
		}
	}()

	if _, err := io.Copy(xw, in); err != nil {
		return err
	}
	return nil
}

// Bzip2 compresses an input stream using bzip2.
func Bzip2(in io.Reader, out io.ReadWriter) error {
	bw, err := bzip2.NewWriter(out, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err := bw.Close(); err != nil {
			fmt.Printf("packit/bzip2: failed closing stream: %v\n", err)
		}
	}()

	if _, err := io.Copy(bw, in); err != nil {
		return err
	}
	return nil
}
