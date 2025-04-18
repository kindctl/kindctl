package cluster

import (
	"os"
	"os/exec"

	"kindctl/internal/config"
	"kindctl/internal/logger"
)

// Initialize creates a new Kind cluster and writes a default config file.
func Initialize(log *logger.Logger, configFile string) error {
	// Create default config
	cfg := config.DefaultConfig()
	if err := config.SaveConfig(configFile, cfg); err != nil {
		return err
	}
	log.Info("Created default configuration file: %s", configFile)

	// Create Kind cluster
	cmd := exec.Command("kind", "create", "cluster", "--name", cfg.Cluster.Name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	log.Info("Created Kind cluster: %s", cfg.Cluster.Name)

	// Install NGINX ingress controller
	cmd = exec.Command("kubectl", "apply", "-f", "https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	log.Info("Installed NGINX ingress controller")

	return nil
}

// Destroy deletes the Kind cluster.
func Destroy(log *logger.Logger, clusterName string) error {
	cmd := exec.Command("kind", "delete", "cluster", "--name", clusterName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	log.Info("Deleted Kind cluster: %s", clusterName)
	return nil
}
