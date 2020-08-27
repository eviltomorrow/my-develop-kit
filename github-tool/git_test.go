package gitlib

import (
	"testing"
	"time"
)

func TestGitPull(t *testing.T) {
	result, err := GitPull(path, 30*time.Second)
	t.Log(err)
	t.Log(result)
}

func TestGitRemoteWithV(t *testing.T) {
	result, err := GitRemoteWithV(path, 30*time.Second)
	t.Log(err)
	t.Log(result)
}

func TestGitListBranch(t *testing.T) {
	result, err := GitListBranch(path, 10*time.Second)
	t.Log(err)
	t.Log(result)
}

func TestGitCurrentBranch(t *testing.T) {
	result, err := GitCurrentBranch(path, 10*time.Second)
	t.Log(err)
	t.Log(result)
}
