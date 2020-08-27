package gitlib

import (
	"fmt"
	"time"
)

type cache struct {
	path    chan string
	records chan *record
}

func newCache(data chan string) *cache {
	return &cache{
		path:    data,
		records: make(chan *record, 32),
	}
}

const (
	success string = "success"
	failure string = "failure"
	skip    string = "skip"
)

type record struct {
	Count         int       `json:"count"`
	GitLocalPath  string    `json:"git-local-path"`
	GitRemotePath string    `json:"git-remote-path"`
	Status        string    `json:"status"`
	Branch        string    `json:"branch"`
	Timstamp      time.Time `json:"timestamp"`
	Err           error     `json:"error"`
}

func (r *record) String() string {
	return fmt.Sprintf("[%s] %s %s", r.Status, r.GitRemotePath, r.GitLocalPath)
}
