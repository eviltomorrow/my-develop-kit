package server

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config config
type Config struct {
	System System `json:"system" toml:"system"`
	Log    Log    `json:"log" toml:"log"`
}

// System system
type System struct {
	CertsDir string `json:"certs-dir" toml:"certs-dir"`
	Port     int    `json:"port" toml:"port"`
}

// Log log
type Log struct {
	DisableTimestamp bool   `json:"disable-timestamp" toml:"disable-timestamp"`
	Level            string `json:"level" toml:"level"`
	Format           string `json:"format" toml:"format"`
	FileName         string `json:"filename" toml:"filename"`
	MaxSize          int    `json:"maxsize" toml:"maxsize"`
}

// Load 加载配置文件
func (cg *Config) Load(f func(*Config)) error {
	dir, err := os.Getwd()
	if err == nil {
		cg.System.CertsDir = filepath.Join(dir, "certs")
	}
	f(cg)
	return nil
}

func (cg *Config) String() string {
	buf, err := json.Marshal(cg)
	if err != nil {
		return fmt.Sprintf("Marshal config to json failure, nest error: %v", err)
	}
	return string(buf)
}

// DefaultGlobalConfig 默认配置
var DefaultGlobalConfig = &Config{
	Log: Log{
		DisableTimestamp: false,
		Level:            "info",
		Format:           "text",
		FileName:         "/var/log/z-console-server/data.log",
		MaxSize:          200,
	},
	System: System{
		Port:     9090,
		CertsDir: "/var/run/z-console-server/certs",
	},
}
