package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the kindctl configuration.
type Config struct {
	Logging struct {
		Level string `yaml:"level"`
	} `yaml:"logging"`
	Cluster struct {
		Name string `yaml:"name"`
	} `yaml:"cluster"`
	Postgres struct {
		Enabled  bool   `yaml:"enabled"`
		Ingress  string `yaml:"ingress"`
		Version  string `yaml:"version"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"postgres"`
	Redis struct {
		Enabled bool   `yaml:"enabled"`
		Ingress string `yaml:"ingress"`
	} `yaml:"redis"`
	PgAdmin struct {
		Enabled  bool   `yaml:"enabled"`
		Ingress  string `yaml:"ingress"`
		Email    string `yaml:"email"`
		Password string `yaml:"password"`
	} `yaml:"pgadmin"`
	Adminer struct {
		Enabled bool   `yaml:"enabled"`
		Ingress string `yaml:"ingress"`
	} `yaml:"adminer"`
	RabbitMQ struct {
		Enabled  bool   `yaml:"enabled"`
		Ingress  string `yaml:"ingress"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"rabbitmq"`
	Mailpit struct {
		Enabled  bool   `yaml:"enabled"`
		Ingress  string `yaml:"ingress"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"mailpit"`
	Dashboard struct {
		Enabled bool   `yaml:"enabled"`
		Ingress string `yaml:"ingress"`
	} `yaml:"dashboard"`
}

// LoadConfig reads and parses the YAML configuration file.
func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// DefaultConfig returns a default configuration for initialization.
func DefaultConfig() *Config {
	return &Config{
		Logging: struct {
			Level string `yaml:"level"`
		}{
			Level: "info",
		},
		Cluster: struct {
			Name string `yaml:"name"`
		}{
			Name: "kind-cluster",
		},
		Dashboard: struct {
			Enabled bool   `yaml:"enabled"`
			Ingress string `yaml:"ingress"`
		}{
			Enabled: true,
			Ingress: "dashboard.local",
		},
	}
}

// SaveConfig writes the configuration to a file.
func SaveConfig(filePath string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}
