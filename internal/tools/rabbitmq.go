package tools

import (
	"os"
	"os/exec"

	"kindctl/internal/config"
	"kindctl/internal/logger"
)

// InstallRabbitMQ installs RabbitMQ and sets up ingress.
func InstallRabbitMQ(log *logger.Logger, cfg *config.Config) error {
	// Apply RabbitMQ manifest
	cmd := exec.Command("helm", "install", "rabbitmq", "bitnami/rabbitmq",
		"--set", "auth.username="+cfg.RabbitMQ.Username,
		"--set", "auth.password="+cfg.RabbitMQ.Password,
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
  name: rabbitmq-ingress
  namespace: default
spec:
  rules:
  - host: ` + cfg.RabbitMQ.Ingress + `
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: rabbitmq
            port:
              number: 15672
`
	if err := os.WriteFile("rabbitmq-ingress.yaml", []byte(ingressManifest), 0644); err != nil {
		return err
	}
	cmd = exec.Command("kubectl", "apply", "-f", "rabbitmq-ingress.yaml")
	if err := cmd.Run(); err != nil {
		return err
	}
	_ = os.Remove("rabbitmq-ingress.yaml")

	log.Info("Installed RabbitMQ with ingress: %s", cfg.RabbitMQ.Ingress)
	return nil
}
