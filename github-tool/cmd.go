package gitlib

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

type cmd struct {
	c       *exec.Cmd
	timeout time.Duration
}

func newCmd(name string, args []string, timeout time.Duration) *cmd {
	return &cmd{
		c:       exec.Command(name, args...),
		timeout: timeout,
	}
}

func (c *cmd) do() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	var stdout, stderr bytes.Buffer
	c.c.Stdout = &stdout
	c.c.Stderr = &stderr

	var ch = make(chan error)

	if err := c.c.Start(); err != nil {
		return "", err
	}

	go func() {
		ch <- c.c.Wait()
	}()

	select {
	case <-ctx.Done():
		c.c.Process.Kill()
		err := <-ch
		close(ch)
		return stdout.String(), err

	case err := <-ch:
		close(ch)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s%s", stderr.String(), stdout.String()), nil
	}
}
