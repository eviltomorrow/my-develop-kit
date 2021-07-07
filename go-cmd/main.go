package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"syscall"
	"time"
)

func main() {
	stdout, stderr, err := ExecuteCmd("pwd", 10*time.Second)
	log.Printf("stdout: %s\r\n", stdout)
	log.Printf("stderr: %s\r\n", stderr)
	log.Printf("error: %v\r\n", err)

	stdout, stderr, err = ExecuteShell("echo.sh", 10*time.Second)
	log.Printf("stdout: %s\r\n", stdout)
	log.Printf("stderr: %s\r\n", stderr)
	log.Printf("error: %v\r\n", err)
}

// ExecuteCmd 执行 command
func ExecuteCmd(c string, timeout time.Duration) (string, string, error) {
	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var ch = make(chan error)
	var cmd = exec.Command("/bin/sh", "-c", c)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return "", "", fmt.Errorf("start execute cmd failure, nest error: %v", err)
	}

	go func() {
		ch <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		cmd.Process.Signal(syscall.SIGINT)
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		<-ch
		close(ch)
		return stdout.String(), stderr.String(), fmt.Errorf("execute cmd timeout")

	case err := <-ch:
		close(ch)
		return stdout.String(), stderr.String(), err
	}
}

// ExecuteShell 执行脚本
func ExecuteShell(shell string, timeout time.Duration) (string, string, error) {
	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var ch = make(chan error)

	var cmd = exec.Command("/bin/sh", shell)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return "", "", fmt.Errorf("start execute shell failure, nest error: %v", err)
	}

	go func() {
		ch <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		<-ch
		close(ch)
		return stdout.String(), stderr.String(), fmt.Errorf("execute shell timeout")

	case err := <-ch:
		close(ch)
		return stdout.String(), stderr.String(), err
	}
}
