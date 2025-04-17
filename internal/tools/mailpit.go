package tools

import (
	"os"
	"os/exec"

	"kindctl/internal/config"
	"kindctl/internal/logger"
)

// InstallMailpit installs Mailpit and sets up ingress.
func InstallMailpit(log *logger.Logger, cfg *config.Config) error {
	// Apply Mailpit manifest
	manifest := `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mailpit
  namespace: default
spec:
  selector:
    matchLabels:
      app: mailpit
  template:
    metadata:
      labels:
        app: mailpit
    spec:
      containers:
      - name: mailpit
        image: axllent/mailpit:latest
        ports:
        - containerPort: 8025
        env:
        - name: MP_SMTP_AUTH_ACCEPT_ANY
          value: "1"
        - name: MP_SMTP_AUTH_ALLOW_INSECURE
          value: "1"
---
apiVersion: v1
kind: Service
metadata:
  name: mailpit
  namespace: default
spec:
  selector:
    app: mailpit
  ports:
  - port: 80
    targetPort: 8025
`
	if err := os.WriteFile("mailpit.yaml", []byte(manifest), 0644); err != nil {
		return err
	}
	cmd := exec.Command("kubectl", "apply", "-f", "mailpit.yaml")
	if err := cmd.Run(); err != nil {
		return err
	}
	_ = os.Remove("mailpit.yaml")

	// Apply ingress
	ingressManifest := `
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: mailpit-ingress
  namespace: default
spec:
  rules:
  - host: ` + cfg.Mailpit.Ingress + `
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: mailpit
            port:
              number: 80
`
	if err := os.WriteFile("mailpit-ingress.yaml", []byte(ingressManifest), 0644); err != nil {
		return err
	}
	cmd = exec.Command("kubectl", "apply", "-f", "mailpit-ingress.yaml")
	if err := cmd.Run(); err != nil {
		return err
	}
	_ = os.Remove("mailpit-ingress.yaml")

	log.Info("Installed Mailpit with ingress: %s", cfg.Mailpit.Ingress)
	return nil
}
