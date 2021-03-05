package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"math"
	"os"
)

// CalMD5 md5
func CalMD5(path string) (string, error) {
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
		return "", fmt.Errorf("Path is dir")
	}
	if info.Size() < 1024*1024*1024*2 {
		md5h := md5.New()
		if _, err = io.Copy(md5h, file); err != nil {
			return "", err
		}
		return fmt.Sprintf("%x", md5h.Sum(nil)), nil
	}

	var blocks = uint64(math.Ceil(float64(info.Size()) / float64(8*1024)))
	md5h := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(8*1024, float64(info.Size()-int64(i*8*1024))))
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
