package xshell

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type resource struct {
	No            int       `json:"no"`
	Host          string    `json:"host"`
	Port          int       `json:"port"`
	Username      string    `json:"username"`
	Password      string    `json:"password"`
	Count         int64     `json:"count"`
	LastLoginTime time.Time `json:"last-login-time"`
}

func (r *resource) verify() error {
	if strings.TrimSpace(r.Username) == "" {
		return fmt.Errorf("Invalid username")
	}
	if strings.TrimSpace(r.Password) == "" {
		return fmt.Errorf("Invalid password")
	}
	if strings.TrimSpace(r.Host) == "" {
		return fmt.Errorf("Invalid host")
	}
	if r.Port >= 65535 || r.Port < 22 {
		return fmt.Errorf("Invalid port")
	}
	return nil
}

func (r *resource) String() string {
	buf, _ := json.Marshal(r)
	return string(buf)
}

func encryptPassword(password string) string {
	return ""
}

func decryptPassword(password string) string {
	return ""
}
