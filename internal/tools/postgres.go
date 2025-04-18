package tools

import (
	"os"
	"os/exec"

	"kindctl/internal/config"
	"kindctl/internal/logger"
)

// InstallPostgres installs PostgreSQL and sets up ingress.
func InstallPostgres(log *logger.Logger, cfg *config.Config) error {
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

	// Apply PostgreSQL manifest (using Bitnami Helm chart)
	cmd = exec.Command("helm", "install", "postgres", "bitnami/postgresql",
		"--set", "global.postgresql.auth.username="+cfg.Postgres.Username,
		"--set", "global.postgresql.auth.password="+cfg.Postgres.Password,
		"--set", "global.postgresql.auth.database="+cfg.Postgres.Database,
		"--set", "image.tag="+cfg.Postgres.Version,
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
  name: postgres-ingress
  namespace: default
spec:
  rules:
  - host: ` + cfg.Postgres.Ingress + `
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: postgres-postgresql
            port:
              number: 5432
`
	if err := os.WriteFile("postgres-ingress.yaml", []byte(ingressManifest), 0644); err != nil {
		return err
	}
	cmd = exec.Command("kubectl", "apply", "-f", "postgres-ingress.yaml")
	if err := cmd.Run(); err != nil {
		return err
	}
	_ = os.Remove("postgres-ingress.yaml")

	log.Info("Installed PostgreSQL with ingress: %s", cfg.Postgres.Ingress)
	return nil
}
