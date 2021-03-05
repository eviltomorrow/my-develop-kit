package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

var (
	username = "root"
	password = ""
	host     = "192.168.180.244"
	port     = 22
	timeout  = 5 * time.Second
	path     = "/home/shepard/.ssh/id_rsa"
)

func main() {
	password = ""
	var authMethods = make([]ssh.AuthMethod, 0, 4)
	if password == "" {
		pem, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf("Read private key failure, nest error: %v\r\n", err)
		}

		signer, err := ssh.ParsePrivateKey(pem)
		if err != nil {
			log.Fatalf("Parse private key failure, nest error: %v", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	} else {
		authMethods = append(authMethods, ssh.KeyboardInteractive(setKeyboard(password)))
		authMethods = append(authMethods, ssh.Password(password))
	}

	config := ssh.ClientConfig{
		User: username,
		Auth: authMethods,
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
	if err != nil {
		log.Fatalf("set stderr failure, nest error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	var ch = make(chan error)
	var stdout, stderr bytes.Buffer

	wg.Add(1)
	go func() {
		io.Copy(&stdout, stdoutpipe)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		io.Copy(&stderr, stderrpipe)
		wg.Done()
	}()

	if err := session.Start("cd /tmp; ls -l; cat hsperfdata_root"); err != nil {
		log.Fatalf("start failure, nest error: %v", err)
	}

	wg.Add(1)
	go func() {
		ch <- session.Wait()
		wg.Done()
	}()

	fmt.Println("====================================  start execute  ====================================")
	select {
	case <-ctx.Done():
		session.Close()
		<-ch
		wg.Wait()
		close(ch)
		log.Printf("部分结果：\r\n%v", stdout.String())
		log.Printf("部分错误：\r\n%v", stderr.String())

	case err := <-ch:
		session.Close()
		wg.Wait()
		close(ch)
		if err != nil {
			log.Fatalf("Execute failure, nest error: %v, stderr: %v, stdout: %v", err, stderr.String(), stdout.String())
		}
		fmt.Println(stdout.String())
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
