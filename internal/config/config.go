package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Logging defines the logging configuration.
type Logging struct {
	Level string `yaml:"level"`
}

// Cluster defines the Kind cluster configuration.
type Cluster struct {
	Name        string `yaml:"name"`
	WorkerNodes int    `yaml:"workerNodes"`
}

// Postgres defines the PostgreSQL configuration.
type Postgres struct {
	Enabled  bool   `yaml:"enabled"`
	Ingress  string `yaml:"ingress"`
	Version  string `yaml:"version"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// Redis defines the Redis configuration.
type Redis struct {
	Enabled bool   `yaml:"enabled"`
	Ingress string `yaml:"ingress"`
}

// PgAdmin defines the PgAdmin configuration.
type PgAdmin struct {
	Enabled  bool   `yaml:"enabled"`
	Ingress  string `yaml:"ingress"`
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

// Adminer defines the Adminer configuration.
type Adminer struct {
	Enabled bool   `yaml:"enabled"`
	Ingress string `yaml:"ingress"`
}

// RabbitMQ defines the RabbitMQ configuration.
type RabbitMQ struct {
	Enabled  bool   `yaml:"enabled"`
	Ingress  string `yaml:"ingress"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Mailpit defines the Mailpit configuration.
type Mailpit struct {
	Enabled  bool   `yaml:"enabled"`
	Ingress  string `yaml:"ingress"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Dashboard defines the dashboard configuration.
type Dashboard struct {
	Enabled bool   `yaml:"enabled"`
	Ingress string `yaml:"ingress"`
}

// Config represents the kindctl configuration.
type Config struct {
	Logging   Logging   `yaml:"logging"`
	Cluster   Cluster   `yaml:"cluster"`
	Postgres  Postgres  `yaml:"postgres"`
	Redis     Redis     `yaml:"redis"`
	PgAdmin   PgAdmin   `yaml:"pgadmin"`
	Adminer   Adminer   `yaml:"adminer"`
	RabbitMQ  RabbitMQ  `yaml:"rabbitmq"`
	Mailpit   Mailpit   `yaml:"mailpit"`
	Dashboard Dashboard `yaml:"dashboard"`
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
		Logging: Logging{
			Level: "info",
		},
		Cluster: Cluster{
			Name:        "kind-cluster",
			WorkerNodes: 1, // Default to 1 worker node to match provided YAML
		},
		Postgres: Postgres{
			Enabled:  false,
			Ingress:  "postgres.local",
			Version:  "15",
			Username: "postgres",
			Password: "password",
			Database: "app",
		},
		Redis: Redis{
			Enabled: false,
			Ingress: "redis.local",
		},
		PgAdmin: PgAdmin{
			Enabled:  false,
			Ingress:  "pgadmin.local",
			Email:    "admin@pgadmin.local",
			Password: "admin",
		},
		Adminer: Adminer{
			Enabled: false,
			Ingress: "adminer.local",
		},
		RabbitMQ: RabbitMQ{
			Enabled:  false,
			Ingress:  "rabbitmq.local",
			Username: "guest",
			Password: "guest",
		},
		Mailpit: Mailpit{
			Enabled:  false,
			Ingress:  "mailpit.local",
			Username: "",
			Password: "",
		},
		Dashboard: Dashboard{
			Enabled: true,
			Ingress: "dashboard.local",
		},
	}
}

// SaveConfig writes a filtered configuration to a file, including only logging, cluster, and enabled sections.
func SaveConfig(filePath string, cfg *Config) error {
	// Create a filtered configuration map
	filteredConfig := make(map[string]interface{})

	// Always include logging and cluster
	filteredConfig["logging"] = cfg.Logging
	filteredConfig["cluster"] = cfg.Cluster

	// Include other sections only if Enabled is true
	if cfg.Postgres.Enabled {
		filteredConfig["postgres"] = cfg.Postgres
	}
	if cfg.Redis.Enabled {
		filteredConfig["redis"] = cfg.Redis
	}
	if cfg.PgAdmin.Enabled {
		filteredConfig["pgadmin"] = cfg.PgAdmin
	}
	if cfg.Adminer.Enabled {
		filteredConfig["adminer"] = cfg.Adminer
	}
	if cfg.RabbitMQ.Enabled {
		filteredConfig["rabbitmq"] = cfg.RabbitMQ
	}
	if cfg.Mailpit.Enabled {
		filteredConfig["mailpit"] = cfg.Mailpit
	}
	if cfg.Dashboard.Enabled {
		filteredConfig["dashboard"] = cfg.Dashboard
	}

	// Serialize the filtered configuration to YAML
	data, err := yaml.Marshal(filteredConfig)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}
