package gitlib

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// FindGitDir find dir
func FindGitDir(root string) chan string {
	var data = make(chan string, 10)
	go func() {
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() && strings.HasSuffix(path, ".git") {
				var dir = filepath.Dir(path)
				data <- dir
			}
			return nil
		})
		if err != nil {
			log.Printf("[Warning] Walk dir failure, nest error: %v", err)
		}
		close(data)
	}()
	return data
}
