package config

import (
	"encoding/json"

	"github.com/BurntSushi/toml"
	"github.com/eviltomorrow/tools/plog"
)

// Config config
type Config struct {
	Server *Server `toml:"server" json:"server"`
	System *System `toml:"system" json:"system"`
}

// Server server
type Server struct {
	Port int `toml:"port" json:"port"`
}

// System system
type System struct {
	LogFileDir      string `toml:"log-file-dir" json:"log-file-dir"`
	LogLevel        string `toml:"log-level" json:"log-level"`
	LogFormat       string `toml:"log-format" json:"log-format"`
	PProfListenPort int    `toml:"pprof-listen-port" json:"pprof-listen-port"`
}

var defaultConf = &Config{
	Server: &Server{
		Port: 8080,
	},
}

var globalConf = defaultConf

// GetGlobalConfig returns the global configuration
func GetGlobalConfig() *Config {
	return globalConf
}

// Load loads config options from a toml file
func (c *Config) Load(path string) error {
	_, err := toml.DecodeFile(path, defaultConf)
	return err
}

// ToLogConfig get log config
func (c *Config) ToLogConfig() *plog.LogConfig {
	return &plog.LogConfig{
		Level:      c.System.LogLevel,
		Format:     c.System.LogFormat,
		LogFileDir: c.System.LogFileDir,
	}
}

// String to json
func (c *Config) String() string {
	buf, _ := json.Marshal(c)
	return string(buf)
}
