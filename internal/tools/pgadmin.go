package tools

import (
	"os"
	"os/exec"

	"kindctl/internal/config"
	"kindctl/internal/logger"
)

// InstallPgAdmin installs pgAdmin and sets up ingress.
func InstallPgAdmin(log *logger.Logger, cfg *config.Config) error {
	// Ensure Runix Helm repository
	cmd := exec.Command("helm", "repo", "add", "runix", "https://helm.runix.net")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Warn("Failed to add Runix Helm repo, it may already exist: %v", err)
	}
	cmd = exec.Command("helm", "repo", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	log.Info("Ensured Runix Helm repository")

	// Apply pgAdmin manifest
	cmd = exec.Command("helm", "install", "pgadmin", "runix/pgadmin4",
		"--set", "env.email="+cfg.PgAdmin.Email,
		"--set", "env.password="+cfg.PgAdmin.Password,
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
  name: pgadmin-ingress
  namespace: default
spec:
  rules:
  - host: ` + cfg.PgAdmin.Ingress + `
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: pgadmin4
            port:
              number: 80
`
	if err := os.WriteFile("pgadmin-ingress.yaml", []byte(ingressManifest), 0644); err != nil {
		return err
	}
	cmd = exec.Command("kubectl", "apply", "-f", "pgadmin-ingress.yaml")
	if err := cmd.Run(); err != nil {
		return err
	}
	_ = os.Remove("pgadmin-ingress.yaml")

	log.Info("Installed pgAdmin with ingress: %s", cfg.PgAdmin.Ingress)
	return nil
}
