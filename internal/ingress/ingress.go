package ingress

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"kindctl/internal/logger"
)

// AddHostEntry adds an entry to /etc/hosts (or equivalent on Windows).
func AddHostEntry(log *logger.Logger, host string) error {
	entry := fmt.Sprintf("127.0.0.1 %s", host)
	if runtime.GOOS == "windows" {
		hostsFile := `C:\Windows\System32\drivers\etc\hosts`
		cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`Add-Content -Path %s -Value "%s"`, hostsFile, entry))
		if err := cmd.Run(); err != nil {
			return err
		}
	} else {
		hostsFile := "/etc/hosts"
		cmd := exec.Command("sh", "-c", fmt.Sprintf(`echo "%s" | sudo tee -a %s`, entry, hostsFile))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	log.Info("Added /etc/hosts entry: %s", entry)
	return nil
}
