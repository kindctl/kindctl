package ingress

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"kindctl/internal/logger"
)

func AddHostEntry(log *logger.Logger, host string) error {
	entry := fmt.Sprintf("127.0.0.1 %s", host)
	var hostsFile string

	if runtime.GOOS == "windows" {
		hostsFile = `C:\Windows\System32\drivers\etc\hosts`
	} else {
		hostsFile = "/etc/hosts"
	}

	exists, err := containsHostEntry(hostsFile, entry)
	if err != nil {
		return fmt.Errorf("failed to check hosts file: %w", err)
	}
	if exists {
		log.Info("Host entry already exists: %s", entry)
		return nil
	}

	if runtime.GOOS == "windows" {
		cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`Add-Content -Path %s -Value "%s"`, hostsFile, entry))
		if err := cmd.Run(); err != nil {
			return err
		}
	} else {
		cmd := exec.Command("sh", "-c", fmt.Sprintf(`echo "%s" | sudo tee -a %s`, entry, hostsFile))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	log.Info("✅ Added hosts entry: %s", entry)
	return nil
}

func containsHostEntry(hostsFile string, entry string) (bool, error) {
	file, err := os.Open(hostsFile)
	if err != nil {
		return false, err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			_, err := fmt.Fprintf(os.Stderr, "⚠️ Failed to close file: %v\n", cerr)
			if err != nil {
				return
			}
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == entry {
			return true, nil
		}
	}
	return false, scanner.Err()
}
