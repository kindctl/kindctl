package tools

import (
	"os"
	"os/exec"

	"kindctl/internal/config"
	"kindctl/internal/logger"
)

// InstallDashboard installs the Kubernetes Dashboard.
func InstallDashboard(log *logger.Logger, cfg *config.Config) error {
	// Apply Kubernetes Dashboard manifests (simplified example)
	cmd := exec.Command("kubectl", "apply", "-f", "https://raw.githubusercontent.com/kubernetes/dashboard/v2.7.0/aio/deploy/recommended.yaml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	log.Info("Installed Kubernetes Dashboard")
	return nil
}
