package gitlib

import "testing"

var path = "/Users/shepard/Workspace/go/project/src/github.com/eviltomorrow/gitlib-upgrade-tool"

func TestFindGitDir(t *testing.T) {
	var data = FindGitDir(path)
	for line := range data {
		t.Log(line)
	}
}
