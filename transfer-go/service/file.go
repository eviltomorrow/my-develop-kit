package service

import (
	"crypto/md5"
	"fmt"
	"io"
	"math"
	"os"
)

//
func createDIR(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		return os.MkdirAll(dir, 0755)
	}
	if err == nil && info.IsDir() {
		return nil
	}
	return os.MkdirAll(dir, 0755)
}

//
func createFile(path string, size int64) error {
	info, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		defer file.Close()

		if size >= 1 {
			_, err = file.Seek(size-1, os.SEEK_SET)
			if err != nil {
				return err
			}
			if size >= 1 {
				_, err = file.Write([]byte{0})
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	if err == nil && info.IsDir() {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		defer file.Close()

		if size >= 1 {
			_, err = file.Seek(size-1, os.SEEK_SET)
			if err != nil {
				return err
			}
			if size >= 1 {
				_, err = file.Write([]byte{0})
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	return nil
}

// writeFile å†™å…¥æ–‡ä»¶
func writeFile(file *os.File, offset int64, data []byte, flag bool) error {
	if file == nil {
		return fmt.Errorf("file is nil")
	}
	if offset < 0 {
		return fmt.Errorf("invalid offset")
	}
	if data == nil {
		return nil
	}

	info, err := file.Stat()
	if err != nil {
		return err
	}
	if int64(offset+int64(len(data))) > info.Size() {
		return fmt.Errorf("write data overflow")
	}
	if !flag {
		_, err = file.Seek(offset, os.SEEK_SET)
		if err != nil {
			return err
		}
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

//
func calMD5(path string) (string, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return "", fmt.Errorf("path is dir")
	}
	if info.Size() < defaultFileSizeWitchCalMD5 {
		md5h := md5.New()
		if _, err = io.Copy(md5h, file); err != nil {
			return "", err
		}
		return fmt.Sprintf("%x", md5h.Sum(nil)), nil
	}

	var blocks = uint64(math.Ceil(float64(info.Size()) / float64(defaultFileChunckWitchCalMD5)))
	md5h := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(defaultFileChunckWitchCalMD5, float64(info.Size()-int64(i*defaultFileChunckWitchCalMD5))))
		buf := make([]byte, blocksize)
		_, err = file.Read(buf)
		if err != nil {
			return "", err
		}
		_, err = io.WriteString(md5h, string(buf))
		if err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("%x", md5h.Sum(nil)), nil
}

const (
	defaultFileChunckWitchCalMD5       = 8 * 1024
	defaultFileSizeWitchCalMD5   int64 = 1024 * 1024 * 1024 * 2
)
