package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Target  TargetConfig  `yaml:"target"`
	Logging LoggingConfig `yaml:"logging"`
	Proxy   ProxyConfig   `yaml:"proxy"`
}

type ServerConfig struct {
	Host      string `yaml:"host"`
	HTTPPort  int    `yaml:"http_port"`
	HTTPSPort int    `yaml:"https_port"`
	CertFile  string `yaml:"cert_file"`
	KeyFile   string `yaml:"key_file"`
}

type TargetConfig struct {
	Host     string    `yaml:"host"`
	Port     int       `yaml:"port"`
	Username string    `yaml:"username"`
	Password string    `yaml:"password"`
	SSH      SSHConfig `yaml:"ssh"`
}

type SSHConfig struct {
	Enabled bool `yaml:"enabled"`
	Port    int  `yaml:"port"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type ProxyConfig struct {
	Timeout        int      `yaml:"timeout"`
	MaxConnections int      `yaml:"max_connections"`
	PathRewrite    bool     `yaml:"path_rewrite"`
	AllowedPaths   []string `yaml:"allowed_paths"`
	BlockedPaths   []string `yaml:"blocked_paths"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	cfg.setDefaults()

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.Server.HTTPPort < 0 || c.Server.HTTPPort > 65535 {
		return fmt.Errorf("invalid http_port: %d", c.Server.HTTPPort)
	}
	if c.Server.HTTPSPort < 0 || c.Server.HTTPSPort > 65535 {
		return fmt.Errorf("invalid https_port: %d", c.Server.HTTPSPort)
	}
	if c.Target.Host == "" {
		return fmt.Errorf("target.host is required")
	}
	if c.Target.Port <= 0 || c.Target.Port > 65535 {
		return fmt.Errorf("invalid target.port: %d", c.Target.Port)
	}
	if c.Target.SSH.Port <= 0 || c.Target.SSH.Port > 65535 {
		return fmt.Errorf("invalid target.ssh.port: %d", c.Target.SSH.Port)
	}
	if c.Proxy.Timeout < 0 {
		return fmt.Errorf("proxy.timeout cannot be negative")
	}
	return nil
}

func (c *Config) setDefaults() {
	if c.Server.Host == "" {
		c.Server.Host = "0.0.0.0"
	}
	if c.Server.HTTPPort == 0 {
		c.Server.HTTPPort = 8080
	}
	if c.Server.HTTPSPort == 0 {
		c.Server.HTTPSPort = 8443
	}
	if c.Target.SSH.Port == 0 {
		c.Target.SSH.Port = 22
	}
	if c.Logging.Level == "" {
		c.Logging.Level = "info"
	}
	if c.Logging.Format == "" {
		c.Logging.Format = "text"
	}
	if c.Proxy.Timeout <= 0 {
		c.Proxy.Timeout = 30
	}
	if c.Proxy.MaxConnections <= 0 {
		c.Proxy.MaxConnections = 100
	}
}
