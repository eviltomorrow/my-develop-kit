package gitlib

import (
	"testing"
	"time"
)

func TestCmdExecute(t *testing.T) {
	var c = newCmd("/bin/sh", []string{"-c", "ls -l"}, 10*time.Second)
	result, err := c.do()
	t.Log(err)
	t.Log(result)

	c = newCmd("/bin/sh", []string{"-c", "ping www.baidu.com"}, 5*time.Second)
	result, err = c.do()
	t.Log(err)
	t.Log(result)

	c = newCmd("/bin/sh", []string{"-c", "git pull origin master"}, 10*time.Second)
	result, err = c.do()
	t.Log(err)
	t.Log(result)

}
