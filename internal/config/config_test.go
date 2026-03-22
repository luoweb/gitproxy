package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	content := `
server:
  host: "localhost"
  http_port: 8080
  https_port: 8443
target:
  host: "example.com"
  port: 443
  ssh:
    enabled: true
    port: 22
logging:
  level: "debug"
  format: "text"
proxy:
  timeout: 30
  max_connections: 50
`
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Server.Host != "localhost" {
		t.Errorf("Expected host 'localhost', got '%s'", cfg.Server.Host)
	}
	if cfg.Server.HTTPPort != 8080 {
		t.Errorf("Expected http_port 8080, got %d", cfg.Server.HTTPPort)
	}
	if cfg.Target.Host != "example.com" {
		t.Errorf("Expected target host 'example.com', got '%s'", cfg.Target.Host)
	}
	if cfg.Logging.Level != "debug" {
		t.Errorf("Expected log level 'debug', got '%s'", cfg.Logging.Level)
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				Server: ServerConfig{HTTPPort: 8080},
				Target: TargetConfig{Host: "example.com", Port: 443},
			},
			wantErr: false,
		},
		{
			name: "invalid http port",
			config: Config{
				Server: ServerConfig{HTTPPort: 70000},
				Target: TargetConfig{Host: "example.com", Port: 443},
			},
			wantErr: true,
		},
		{
			name: "missing target host",
			config: Config{
				Server: ServerConfig{HTTPPort: 8080},
				Target: TargetConfig{Host: "", Port: 443},
			},
			wantErr: true,
		},
		{
			name: "invalid target port",
			config: Config{
				Server: ServerConfig{HTTPPort: 8080},
				Target: TargetConfig{Host: "example.com", Port: 0},
			},
			wantErr: true,
		},
		{
			name: "negative timeout",
			config: Config{
				Server: ServerConfig{HTTPPort: 8080},
				Target: TargetConfig{Host: "example.com", Port: 443},
				Proxy:  ProxyConfig{Timeout: -1},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetDefaults(t *testing.T) {
	cfg := &Config{}
	cfg.setDefaults()

	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Expected default host '0.0.0.0', got '%s'", cfg.Server.Host)
	}
	if cfg.Server.HTTPPort != 8080 {
		t.Errorf("Expected default http_port 8080, got %d", cfg.Server.HTTPPort)
	}
	if cfg.Logging.Level != "info" {
		t.Errorf("Expected default log level 'info', got '%s'", cfg.Logging.Level)
	}
	if cfg.Proxy.Timeout != 30 {
		t.Errorf("Expected default timeout 30, got %d", cfg.Proxy.Timeout)
	}
	if cfg.Proxy.MaxConnections != 100 {
		t.Errorf("Expected default max_connections 100, got %d", cfg.Proxy.MaxConnections)
	}
}

func TestLoadConfigFileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}
