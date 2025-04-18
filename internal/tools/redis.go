package tools

import (
	"os"
	"os/exec"

	"kindctl/internal/config"
	"kindctl/internal/logger"
)

// InstallRedis installs Redis and sets up ingress.
func InstallRedis(log *logger.Logger, cfg *config.Config) error {
	// Ensure Bitnami Helm repository
	cmd := exec.Command("helm", "repo", "add", "bitnami", "https://charts.bitnami.com/bitnami")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Warn("Failed to add Bitnami Helm repo, it may already exist: %v", err)
	}
	cmd = exec.Command("helm", "repo", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	log.Info("Ensured Bitnami Helm repository")

	// Apply Redis manifest
	cmd = exec.Command("helm", "install", "redis", "bitnami/redis",
		"--set", "architecture=standalone",
		"--namespace", "default", "--create-namespace")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Apply ingress
	ingressManifest := `
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: redis-ingress
  namespace: default
spec:
  rules:
  - host: ` + cfg.Redis.Ingress + `
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: redis
            port:
              number: 6379
`
	if err := os.WriteFile("redis-ingress.yaml", []byte(ingressManifest), 0644); err != nil {
		return err
	}
	cmd = exec.Command("kubectl", "apply", "-f", "redis-ingress.yaml")
	if err := cmd.Run(); err != nil {
		return err
	}
	_ = os.Remove("redis-ingress.yaml")

	log.Info("Installed Redis with ingress: %s", cfg.Redis.Ingress)
	return nil
}
