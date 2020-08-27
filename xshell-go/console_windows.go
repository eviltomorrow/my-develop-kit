package xshell

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

var (
	termType = "vt220"
)

// console 控制台
type console struct {
	signal  chan struct{}
	exitext string
	session *sessionSSH

	stdout io.Reader
	stderr io.Reader
	stdin  io.WriteCloser
}

// SetTerminalType set term type
func SetTerminalType(ttype string) {
	termType = ttype
}

// Buildconsole build
func buildConsole(session *sessionSSH) *console {
	return &console{
		signal:  make(chan struct{}),
		session: session,
	}
}

func (c *console) getSession() *ssh.Session {
	return c.session.session
}

func (c *console) updateTerminalSize() {
	// go func() {
	// 	var signalch = make(chan os.Signal, 1)
	// 	signal.Notify(signalch, syscall.SIGWINCH)

	// 	var fd = int(os.Stdin.Fd())
	// 	oriWidth, oriHeight, err := terminal.GetSize(fd)
	// 	if err != nil {
	// 		log.Printf("Warning: Get terminal size failure, nest error: %v", err)
	// 	}

	// 	for {
	// 		select {
	// 		case sig := <-signalch:
	// 			if sig == nil {
	// 				return
	// 			}

	// 			curWidth, curHeight, err := terminal.GetSize(fd)

	// 			if curWidth == oriWidth && curHeight == oriHeight {
	// 				continue
	// 			}

	// 			c.getSession().WindowChange(curHeight, curWidth)
	// 			if err != nil {
	// 				log.Printf("Warning: Unable to send windows change request: %v", err)
	// 				continue
	// 			}

	// 			oriHeight, oriWidth = curHeight, curWidth
	// 		}
	// 	}
	// }()
}

// InteractiveSession interactive
func (c *console) interactiveSession() error {
	defer func() {
		if c.exitext == "" {
			fmt.Fprintln(os.Stdout, "The connection was closed on the remote side on", time.Now().Format(time.RFC822))
		} else {
			fmt.Fprintln(os.Stdout, c.exitext)
		}

		fmt.Fprintln(os.Stdout, "	------ Press enter to continue ------  ")
	}()

	// var fd = int(os.Stdin.Fd())
	// state, err := terminal.MakeRaw(fd)
	// if err != nil {
	// 	return err
	// }

	// defer terminal.Restore(fd, state)

	// curWidth, curHeight, err := terminal.GetSize(fd)
	// if err != nil {
	// 	return err
	// }

	var curHeight = 500
	var curWidth = 900
	err := c.getSession().RequestPty(termType, curHeight, curWidth, ssh.TerminalModes{})
	if err != nil {
		return err
	}

	c.updateTerminalSize()

	c.stdin, err = c.getSession().StdinPipe()
	if err != nil {
		return err
	}

	c.stdout, err = c.getSession().StdoutPipe()
	if err != nil {
		return err
	}

	c.stderr, err = c.getSession().StderrPipe()
	if err != nil {
		return err
	}

	go func() {
		_, err := io.Copy(os.Stdout, c.stdout)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Printf("Warning: stdout copy failure, nest error: %v", err)
		}
	}()

	go func() {
		_, err := io.Copy(os.Stderr, c.stderr)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Printf("Warning: stderr copy failure, nest error: %v", err)
		}
	}()

	go func() {
		buf := make([]byte, 128)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				log.Printf("Warning: Read from stdin failure, nest error: %v", err)
				return
			}
			if n > 0 {
				_, err = c.stdin.Write(buf[:n])
				if err == io.EOF {
					c.signal <- struct{}{}
					return
				}
				if err != nil {
					c.exitext = err.Error()
					return
				}
			}
		}
	}()

	err = c.getSession().Shell()
	if err != nil {
		return err
	}

	err = c.getSession().Wait()
	if err != nil {
		return err
	}
	go func() {

	}()

	return nil
}

func (c *console) complete() {
	<-c.signal
}

// Close close
func (c *console) close() error {
	return c.stdin.Close()
}
