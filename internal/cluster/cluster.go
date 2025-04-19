package cluster

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"kindctl/internal/config"
	"kindctl/internal/logger"
)

// commandExists checks if a command is available in the system PATH.
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// checkDocker verifies that Docker is installed and running.
func checkDocker(log *logger.Logger) error {
	if !commandExists("docker") {
		return fmt.Errorf("Docker is not installed. Please install Docker Desktop or Rancher Desktop")
	}
	log.Info("Docker is installed.")

	// Skip Docker running check in CI environments
	if os.Getenv("CI") != "true" {
		cmd := exec.Command("docker", "ps")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("Docker is not running. Please start Docker Desktop or Rancher Desktop")
		}
	}
	log.Info("Docker is running.")
	return nil
}

// installHomebrew installs Homebrew on macOS if not present.
func installHomebrew(log *logger.Logger) error {
	if !commandExists("brew") {
		log.Info("Installing Homebrew...")
		cmd := exec.Command("/bin/bash", "-c", `$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)`)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install Homebrew: %w", err)
		}
		log.Info("✅ Homebrew installed.")
	} else {
		log.Info("Homebrew is already installed.")
	}
	return nil
}

// installWinget installs winget on Windows if not present (simplified for Go).
func installWinget(log *logger.Logger) error {
	if !commandExists("winget") {
		log.Info("Installing winget...")
		log.Info("Please install winget manually from https://github.com/microsoft/winget-cli or Microsoft Store.")
		return fmt.Errorf("winget not found; manual installation required")
	}
	log.Info("winget is already installed.")
	return nil
}

// installBinary downloads and installs a binary to /usr/local/bin or equivalent.
func installBinary(log *logger.Logger, name, url, destDir string) error {
	if commandExists(name) {
		log.Info(name, " is already installed.")
		return nil
	}

	log.Info("Installing ", name, "...")
	cmd := exec.Command("curl", "-Lo", name, url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to download %s: %w", name, err)
	}

	if err := os.Chmod(name, 0755); err != nil {
		return fmt.Errorf("failed to set permissions for %s: %w", name, err)
	}

	destPath := filepath.Join(destDir, name)
	if err := os.Rename(name, destPath); err != nil {
		data, err := os.ReadFile(name)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", name, err)
		}
		if err := os.WriteFile(destPath, data, 0755); err != nil {
			return fmt.Errorf("failed to write %s to %s: %w", name, destPath, err)
		}
		if err := os.Remove(name); err != nil {
			log.Warn("Failed to clean up temporary file: ", name)
		}
	}

	log.Info("✅ ", name, " installed.")
	return nil
}

// installKind installs Kind based on the OS.
func installKind(log *logger.Logger) error {
	osType := runtime.GOOS
	if osType == "darwin" && commandExists("brew") {
		if !commandExists("kind") {
			log.Info("Installing kind via Homebrew...")
			cmd := exec.Command("brew", "install", "kind")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to install kind via Homebrew: %w", err)
			}
			log.Info("✅ kind installed.")
		} else {
			log.Info("kind is already installed.")
		}
		return nil
	}

	url := ""
	switch osType {
	case "linux":
		url = "https://kind.sigs.k8s.io/dl/v0.23.0/kind-linux-amd64"
	case "darwin":
		url = "https://kind.sigs.k8s.io/dl/v0.23.0/kind-darwin-amd64"
	case "windows":
		url = "https://kind.sigs.k8s.io/dl/v0.23.0/kind-windows-amd64"
	default:
		return fmt.Errorf("unsupported OS: %s", osType)
	}

	destDir := "/usr/local/bin"
	if osType == "windows" {
		destDir = os.Getenv("ProgramFiles") + "\\kind"
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", destDir, err)
		}
	}
	return installBinary(log, "kind", url, destDir)
}

// installKubectl installs kubectl based on the OS.
func installKubectl(log *logger.Logger) error {
	osType := runtime.GOOS
	if osType == "darwin" && commandExists("brew") {
		if !commandExists("kubectl") {
			log.Info("Installing kubectl via Homebrew...")
			cmd := exec.Command("brew", "install", "kubectl")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to install kubectl via Homebrew: %w", err)
			}
			log.Info("✅ kubectl installed.")
		} else {
			log.Info("kubectl is already installed.")
		}
		return nil
	}

	stableVersion, err := exec.Command("curl", "-L", "-s", "https://dl.k8s.io/release/stable.txt").Output()
	if err != nil {
		return fmt.Errorf("failed to get kubectl stable version: %w", err)
	}
	version := strings.TrimSpace(string(stableVersion))
	url := ""
	switch osType {
	case "linux":
		url = fmt.Sprintf("https://dl.k8s.io/release/%s/bin/linux/amd64/kubectl", version)
	case "darwin":
		url = fmt.Sprintf("https://dl.k8s.io/release/%s/bin/darwin/amd64/kubectl", version)
	case "windows":
		url = fmt.Sprintf("https://dl.k8s.io/release/%s/bin/windows/amd64/kubectl.exe", version)
	default:
		return fmt.Errorf("unsupported OS: %s", osType)
	}

	destDir := "/usr/local/bin"
	if osType == "windows" {
		destDir = os.Getenv("ProgramFiles") + "\\kubectl"
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", destDir, err)
		}
	}
	binaryName := "kubectl"
	if osType == "windows" {
		binaryName = "kubectl.exe"
	}
	return installBinary(log, binaryName, url, destDir)
}

// installHelm installs Helm based on the OS.
func installHelm(log *logger.Logger) error {
	osType := runtime.GOOS
	if osType == "darwin" && commandExists("brew") {
		if !commandExists("helm") {
			log.Info("Installing helm via Homebrew...")
			cmd := exec.Command("brew", "install", "helm")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to install helm via Homebrew: %w", err)
			}
			log.Info("✅ helm installed.")
		} else {
			log.Info("helm is already installed.")
		}
		return nil
	}

	if osType != "windows" {
		if !commandExists("helm") {
			log.Info("Installing helm...")
			cmd := exec.Command("curl", "-fsSL", "-o", "get_helm.sh", "https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3")
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to download helm install script: %w", err)
			}
			cmd = exec.Command("chmod", "700", "get_helm.sh")
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to chmod helm script: %w", err)
			}
			cmd = exec.Command("./get_helm.sh")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to run helm install script: %w", err)
			}
			if err := os.Remove("get_helm.sh"); err != nil {
				log.Warn("Failed to clean up helm install script.")
			}
			log.Info("✅ helm installed.")
		} else {
			log.Info("helm is already installed.")
		}
		return nil
	}

	url := "https://get.helm.sh/helm-v3.15.4-windows-amd64.zip"
	if !commandExists("helm") {
		log.Info("Installing helm...")
		log.Info("Please download and extract helm from ", url, " and add to PATH.")
		return fmt.Errorf("helm installation on Windows requires manual steps")
	}
	log.Info("helm is already installed.")
	return nil
}

// generateKindConfig creates a Kind configuration file matching the provided YAML structure.
func generateKindConfig(configFile string, clusterName string, workerNodes int) (string, error) {
	configDir := filepath.Dir(configFile)
	kindConfigPath := filepath.Join(configDir, "kind-config.yaml")

	if workerNodes < 1 {
		workerNodes = 1
	}

	var nodes []string
	nodes = append(nodes, `  - role: control-plane
    extraPortMappings:
      - containerPort: 80
        hostPort: 80
      - containerPort: 443
        hostPort: 443`)
	for i := 0; i < workerNodes; i++ {
		nodes = append(nodes, "  - role: worker")
	}

	kindConfigContent := fmt.Sprintf(`kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: %s
nodes:
%s
`, clusterName, strings.Join(nodes, "\n"))

	if err := os.WriteFile(kindConfigPath, []byte(kindConfigContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write Kind config file: %w", err)
	}

	return kindConfigPath, nil
}

// waitForNodesReady waits for all nodes to be in Ready state.
func waitForNodesReady(log *logger.Logger, timeout time.Duration) error {
	log.Info("Waiting for nodes to be ready...")
	start := time.Now()
	for {
		cmd := exec.Command("kubectl", "get", "nodes", "-o", "jsonpath={.items[*].status.conditions[?(@.type=='Ready')].status}")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to check node status: %w", err)
		}

		statuses := strings.Split(strings.TrimSpace(string(output)), " ")
		allReady := true
		for _, status := range statuses {
			if status != "True" {
				allReady = false
				break
			}
		}

		if allReady {
			log.Info("All nodes are ready.")
			return nil
		}

		if time.Since(start) > timeout {
			return fmt.Errorf("timeout waiting for nodes to be ready")
		}

		time.Sleep(5 * time.Second)
	}
}

// Initialize sets up a new Kind cluster and installs dependencies.
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

	if err := checkDocker(log); err != nil {
		return err
	}

	osType := runtime.GOOS
	switch osType {
	case "darwin":
		if err := installHomebrew(log); err != nil {
			return err
		}
	case "windows":
		if err := installWinget(log); err != nil {
			log.Warn("Continuing despite winget installation issue.")
		}
	}

	if err := installKind(log); err != nil {
		return err
	}
	if err := installKubectl(log); err != nil {
		return err
	}
	if err := installHelm(log); err != nil {
		return err
	}

	log.Info("Verifying installations...")
	for _, cmd := range []string{"kind", "kubectl", "helm"} {
		if commandExists(cmd) {
			output, err := exec.Command(cmd, "version", "--client").CombinedOutput()
			if err != nil {
				log.Warn("Failed to get version for ", cmd, ": ", err)
			} else {
				log.Info(cmd, " installed: ", strings.TrimSpace(string(output)))
			}
		} else {
			return fmt.Errorf("%s installation failed", cmd)
		}
	}

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		return err
	}

	cmd := exec.Command("kind", "get", "clusters")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to get kind clusters: %w", err)
	}
	clusters := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, cluster := range clusters {
		if cluster == cfg.Cluster.Name {
			log.Info("Cluster '", cfg.Cluster.Name, "' already exists.")
			return nil
		}
	}

	kindConfigPath, err := generateKindConfig(configFile, cfg.Cluster.Name, cfg.Cluster.WorkerNodes)
	if err != nil {
		return err
	}
	defer func() {
		if err := os.Remove(kindConfigPath); err != nil {
			log.Warn("Failed to clean up Kind config file: ", kindConfigPath)
		}
	}()

	log.Info("Creating Kind cluster with 1 control-plane and ", cfg.Cluster.WorkerNodes, " worker node(s)...")
	cmd = exec.Command("kind", "create", "cluster", "--name", cfg.Cluster.Name, "--config", kindConfigPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create kind cluster: %w", err)
	}
	log.Info("✅ Created Kind cluster: ", cfg.Cluster.Name)

	// Wait for nodes to be ready
	if err := waitForNodesReady(log, 2*time.Minute); err != nil {
		return err
	}

	// Label worker node for NGINX ingress
	log.Info("Labeling worker node for NGINX ingress...")
	cmd = exec.Command("kubectl", "label", "nodes", fmt.Sprintf("%s-worker", cfg.Cluster.Name), "ingress-ready=true")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to label worker node: %w", err)
	}
	log.Info("✅ Worker node labeled with ingress-ready=true")

	fmt.Println()
	log.Info("🏗 Installing NGINX ingress controller...")
	// Use a modified manifest with toleration for not-ready nodes
	nginxManifestURL := "https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml"
	// Download and patch the manifest
	cmd = exec.Command("curl", "-s", nginxManifestURL)
	nginxManifest, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to download NGINX ingress manifest: %w", err)
	}

	// Add toleration to the deployment
	tolerationPatch := `
spec:
  template:
    spec:
      tolerations:
      - key: "node.kubernetes.io/not-ready"
        operator: "Exists"
        effect: "NoSchedule"
`
	nginxManifestStr := string(nginxManifest)
	// Insert toleration after nodeSelector
	nodeSelectorLine := "nodeSelector:\n          ingress-ready: \"true\""
	if strings.Contains(nginxManifestStr, nodeSelectorLine) {
		nginxManifestStr = strings.Replace(nginxManifestStr, nodeSelectorLine, nodeSelectorLine+"\n        "+tolerationPatch, 1)
	} else {
		log.Warn("nodeSelector not found in NGINX manifest; applying without toleration patch")
	}

	// Save patched manifest to a temporary file
	tempManifestPath := filepath.Join(filepath.Dir(configFile), "nginx-ingress-patched.yaml")
	if err := os.WriteFile(tempManifestPath, []byte(nginxManifestStr), 0644); err != nil {
		return fmt.Errorf("failed to write patched NGINX manifest: %w", err)
	}
	defer func() {
		if err := os.Remove(tempManifestPath); err != nil {
			log.Warn("Failed to clean up patched NGINX manifest: ", tempManifestPath)
		}
	}()

	// Apply the patched manifest
	cmd = exec.Command("kubectl", "apply", "-f", tempManifestPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install NGINX ingress controller: %w", err)
	}
	log.Info("✅ Installed NGINX ingress controller")

	return nil
}

func Destroy(log *logger.Logger, clusterName string) error {
	cmd := exec.Command("kind", "delete", "cluster", "--name", clusterName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete kind cluster: %w", err)
	}
	log.Info("Deleted Kind cluster: ", clusterName)
	return nil
}
