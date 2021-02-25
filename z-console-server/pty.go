package server

import (
	"io"
	"net"
	"os/exec"

	"github.com/creack/pty"
)

func buildPTY(conn net.Conn) error {
	// Create arbitrary command.
	c := exec.Command("bash")

	// Start the command with a pty.
	ptmx, err := pty.Start(c)
	if err != nil {
		return err
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.
	if err = pty.Setsize(ptmx, &pty.Winsize{Rows: 9999, Cols: 1280}); err != nil {
		return err
	}

	// Copy stdin to the pty and the pty to stdout.
	go func() { _, _ = io.Copy(ptmx, conn) }()
	_, err = io.Copy(conn, ptmx)

	return err
}
