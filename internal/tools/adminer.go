package tools

import (
	"os"
	"os/exec"

	"kindctl/internal/config"
	"kindctl/internal/logger"
)

// InstallAdminer installs Adminer and sets up ingress.
func InstallAdminer(log *logger.Logger, cfg *config.Config) error {
	// Apply Adminer manifest
	manifest := `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: adminer
  namespace: default
spec:
  selector:
    matchLabels:
      app: adminer
  template:
    metadata:
      labels:
        app: adminer
    spec:
      containers:
      - name: adminer
        image: adminer:4.8.1
        ports:
        - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: adminer
  namespace: default
spec:
  selector:
    app: adminer
  ports:
  - port: 80
    targetPort: 8080
`
	if err := os.WriteFile("adminer.yaml", []byte(manifest), 0644); err != nil {
		return err
	}
	cmd := exec.Command("kubectl", "apply", "-f", "adminer.yaml")
	if err := cmd.Run(); err != nil {
		return err
	}
	_ = os.Remove("adminer.yaml")

	// Apply ingress
	ingressManifest := `
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: adminer-ingress
  namespace: default
spec:
  rules:
  - host: ` + cfg.Adminer.Ingress + `
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: adminer
            port:
              number: 80
`
	if err := os.WriteFile("adminer-ingress.yaml", []byte(ingressManifest), 0644); err != nil {
		return err
	}
	cmd = exec.Command("kubectl", "apply", "-f", "adminer-ingress.yaml")
	if err := cmd.Run(); err != nil {
		return err
	}
	_ = os.Remove("adminer-ingress.yaml")

	log.Info("Installed Adminer with ingress: %s", cfg.Adminer.Ingress)
	return nil
}
