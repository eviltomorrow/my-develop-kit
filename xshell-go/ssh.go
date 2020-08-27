package xshell

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

// SessionSSH session ssh
type sessionSSH struct {
	connection *ssh.Client
	session    *ssh.Session
}

// NewSessionSSH new session ssh
func newSessionSSH(host string, port int, username, password string, timeout time.Duration) (*sessionSSH, error) {
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
		return nil, err
	}

	session, err := connection.NewSession()
	if err != nil {
		return nil, err
	}

	var client = &sessionSSH{
		connection: connection,
		session:    session,
	}
	return client, nil
}

// Close close
func (s *sessionSSH) Close() error {
	if s == nil {
		return nil
	}

	if s.session != nil {
		s.session.Close()
	}

	if s.connection != nil {
		s.connection.Close()
	}

	return nil
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
