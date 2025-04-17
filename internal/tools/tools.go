package tools

import (
	"kindctl/internal/config"
	"kindctl/internal/ingress"
	"kindctl/internal/logger"
)

// UpdateCluster installs or updates tools in the Kind cluster based on the config.
func UpdateCluster(log *logger.Logger, cfg *config.Config) error {
	if cfg.Dashboard.Enabled {
		if err := InstallDashboard(log, cfg); err != nil {
			return err
		}
		if err := ingress.AddHostEntry(log, cfg.Dashboard.Ingress); err != nil {
			log.Warn("Failed to add /etc/hosts entry for %s: %v", cfg.Dashboard.Ingress, err)
		}
	}
	if cfg.Postgres.Enabled {
		if err := InstallPostgres(log, cfg); err != nil {
			return err
		}
		if err := ingress.AddHostEntry(log, cfg.Postgres.Ingress); err != nil {
			log.Warn("Failed to add /etc/hosts entry for %s: %v", cfg.Postgres.Ingress, err)
		}
	}
	if cfg.Redis.Enabled {
		if err := InstallRedis(log, cfg); err != nil {
			return err
		}
		if err := ingress.AddHostEntry(log, cfg.Redis.Ingress); err != nil {
			log.Warn("Failed to add /etc/hosts entry for %s: %v", cfg.Redis.Ingress, err)
		}
	}
	if cfg.PgAdmin.Enabled {
		if err := InstallPgAdmin(log, cfg); err != nil {
			return err
		}
		if err := ingress.AddHostEntry(log, cfg.PgAdmin.Ingress); err != nil {
			log.Warn("Failed to add /etc/hosts entry for %s: %v", cfg.PgAdmin.Ingress, err)
		}
	}
	if cfg.Adminer.Enabled {
		if err := InstallAdminer(log, cfg); err != nil {
			return err
		}
		if err := ingress.AddHostEntry(log, cfg.Adminer.Ingress); err != nil {
			log.Warn("Failed to add /etc/hosts entry for %s: %v", cfg.Adminer.Ingress, err)
		}
	}
	if cfg.RabbitMQ.Enabled {
		if err := InstallRabbitMQ(log, cfg); err != nil {
			return err
		}
		if err := ingress.AddHostEntry(log, cfg.RabbitMQ.Ingress); err != nil {
			log.Warn("Failed to add /etc/hosts entry for %s: %v", cfg.RabbitMQ.Ingress, err)
		}
	}
	if cfg.Mailpit.Enabled {
		if err := InstallMailpit(log, cfg); err != nil {
			return err
		}
		if err := ingress.AddHostEntry(log, cfg.Mailpit.Ingress); err != nil {
			log.Warn("Failed to add /etc/hosts entry for %s: %v", cfg.Mailpit.Ingress, err)
		}
	}
	return nil
}
