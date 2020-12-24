package compress

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// DeCompressZip decompress zip
func DeCompressZip(path string, dir string) error {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer reader.Close()

	fi, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	} else {
		if err != nil {
			return err
		}
		if fi != nil && !fi.IsDir() {
			return fmt.Errorf("Not a dir, path: %v", dir)
		}
	}

	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}

		var fd = filepath.Join(dir, file.Name)
		if err := os.MkdirAll(filepath.Dir(fd), 0755); err != nil {
			rc.Close()
			return err
		}

		fo, err := os.OpenFile(fd, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
		if err != nil {
			rc.Close()
			return err
		}

		if _, err := io.Copy(fo, rc); err != nil {
			fo.Close()
			rc.Close()
			return err
		}

		fo.Close()
		rc.Close()
	}
	return nil
}

// DeCompressTargz decompress tar.gz
func DeCompressTargz(path string, dir string) error {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	} else {
		if err != nil {
			return err
		}
		if fi != nil && !fi.IsDir() {
			return fmt.Errorf("Not a dir, path: %v", dir)
		}
	}

	gzipReader, err := gzip.NewReader(f)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(filepath.Join(dir, header.Name), 0755); err != nil {
				return err
			}

		case tar.TypeReg:
			var fd = filepath.Join(dir, header.Name)
			if err := os.MkdirAll(filepath.Dir(fd), 0755); err != nil {
				return err
			}
			fo, err := os.Create(fd)
			if err != nil {
				return err
			}
			if _, err := io.Copy(fo, tarReader); err != nil {
				return err
			}
			fo.Close()

		default:
			return fmt.Errorf("Unknown file type, header: %s, type: %v", header.Name, header.Typeflag)
		}

	}
	return nil
}
