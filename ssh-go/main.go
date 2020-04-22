package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

var (
	username = "root"
	password = "CentOS_7"
	host     = "192.168.180.67"
	port     = 22
	timeout  = 5 * time.Second
)

func main() {
	config := ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.KeyboardInteractive(setKeyboard(password)),
			ssh.Password(password),
		},
		Config: ssh.Config{
			Ciphers: []string{
				"aes128-ctr",
				"aes192-ctr",
				"aes256-ctr",
				"aes128-gcm@openssh.com",
				"arcfour256",
				"arcfour128",
				"aes128-cbc",
			},
		},
		Timeout: timeout,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	connection, err := ssh.Dial("tcp", net.JoinHostPort(host, fmt.Sprintf("%d", port)), &config)
	if err != nil {
		log.Fatalf("dail failure, nest error: %v", err)
	}

	session, err := connection.NewSession()
	if err != nil {
		log.Fatalf("create session failure, nest error: %v", err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("vt220", 500, 900, modes); err != nil {
		log.Fatalf("request pty failure, nest error: %v", err)
	}

	stdoutpipe, err := session.StdoutPipe()
	if err != nil {
		log.Fatalf("set stdout failure, nest error: %v", err)
	}

	stderrpipe, err := session.StderrPipe()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var complete = make(chan struct{})
	var ch = make(chan error)
	var stdout, stderr bytes.Buffer

	go func() {
		io.Copy(&stdout, stdoutpipe)
		complete <- struct{}{}
	}()
	go io.Copy(&stderr, stderrpipe)

	if err := session.Start("ps -ef | grep mysql | grep -v 'grep mysql'"); err != nil {
		log.Fatalf("start failure, nest error: %v", err)
	}

	go func() {
		ch <- session.Wait()
	}()

	fmt.Println("====================================  start execute  ====================================")
	select {
	case <-ctx.Done():
		session.Close()
		err = <-ch
		<-complete
		close(ch)
		close(complete)
		log.Printf("部分结果：%v\r\n", stdout.String())
		log.Fatalf("Execute timeout, nest error: %v", err)

	case err := <-ch:
		<-complete
		close(ch)
		close(complete)
		if err != nil {
			log.Fatalf("Execute failure, nest error: %v", err)
		}
		if stderr.String() != "" {
			log.Fatalf("Execute error: %v", stderr.String())
		}
		fmt.Println(stdout.String())
		session.Close()
	}
	fmt.Println("====================================  end execute  ====================================")
}

func setKeyboard(password string) func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
	return func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
		answers = make([]string, len(questions))
		for n := range questions {
			answers[n] = password
		}
		return answers, nil
	}
}
