package cluster

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"kindctl/internal/config"
	"kindctl/internal/logger"
)

func Initialize(log *logger.Logger, configFile string) error {
	if _, err := os.Stat(configFile); err == nil {
		log.Info("kindctl.yaml file already exists.")
	} else if os.IsNotExist(err) {
		cfg := config.DefaultConfig()
		if err := config.SaveConfig(configFile, cfg); err != nil {
			return err
		}
		log.Info("✅ Created default configuration file: ", configFile)
	} else {
		return err
	}

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		return err
	}

	cmd := exec.Command("kind", "get", "clusters")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	clusters := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, cluster := range clusters {
		if cluster == cfg.Cluster.Name {
			log.Info("Cluster '", cfg.Cluster.Name, "' already exists.")
			return nil
		}
	}

	cmd = exec.Command("kind", "create", "cluster", "--name", cfg.Cluster.Name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	log.Info("✅ Created Kind cluster: ", cfg.Cluster.Name)

	fmt.Println()
	log.Info("🏗 Installing NGINX ingress controller...")
	cmd = exec.Command("kubectl", "apply", "-f", "https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	log.Info("✅ Installed NGINX ingress controller")

	return nil
}

func Destroy(log *logger.Logger, clusterName string) error {
	cmd := exec.Command("kind", "delete", "cluster", "--name", clusterName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	log.Info("Deleted Kind cluster: ", clusterName)
	return nil
}
