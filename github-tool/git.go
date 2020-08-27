package gitlib

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"
)

// GitFetch fetch
func GitFetch(path string, timeout time.Duration) (string, error) {
	var cmd = fmt.Sprintf("cd %s; git fetch", path)
	var c = newCmd("/bin/sh", []string{"-c", cmd}, timeout)
	return c.do()
}

// GitPull pull cmd
func GitPull(path string, timeout time.Duration) (string, error) {
	var cmd = fmt.Sprintf("cd %s; git pull origin master", path)
	var c = newCmd("/bin/sh", []string{"-c", cmd}, timeout)
	return c.do()
}

// GitRemoteWithV git remote -v
func GitRemoteWithV(path string, timeout time.Duration) (string, error) {
	var cmd = fmt.Sprintf("cd %s; git remote -v", path)
	var c = newCmd("/bin/sh", []string{"-c", cmd}, timeout)
	result, err := c.do()
	if err != nil {
		return "", err
	}
	var reader = bufio.NewReader(bytes.NewBufferString(result))
	for {
		line, _, err := reader.ReadLine()
		if err == nil {
			return string(line), nil
		}
		if err == io.EOF {
			return string(line), nil
		}
		return "", err
	}
}

// GitCurrentBranch get current branch
func GitCurrentBranch(path string, timeout time.Duration) (string, error) {
	branches, err := GitListBranch(path, timeout)
	if err != nil {
		return "", err
	}
	for _, branch := range branches {
		if strings.HasPrefix(strings.TrimSpace(branch), "*") {
			return strings.TrimSpace(strings.Replace(branch, "*", "", -1)), nil
		}
	}
	return "", fmt.Errorf("Not found branch information")
}

// GitListBranch list all branch
func GitListBranch(path string, timeout time.Duration) ([]string, error) {
	var cmd = fmt.Sprintf("cd %s; git branch -a", path)
	var c = newCmd("/bin/sh", []string{"-c", cmd}, timeout)
	result, err := c.do()
	if err != nil {
		return nil, err
	}

	var branches = make([]string, 0, 8)
	var reader = bufio.NewReader(bytes.NewBufferString(result))
	for {
		line, _, err := reader.ReadLine()
		if err == nil {
			branches = append(branches, strings.TrimSpace(string(line)))
			continue
		}
		if err == io.EOF {
			branches = append(branches, strings.TrimSpace(string(line)))
			break
		}
		return nil, err
	}
	return branches, nil
}
