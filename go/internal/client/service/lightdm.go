package service

import (
	"fmt"
	"os/exec"
	"strings"
)

func RestartLightDM() error {
	user, err := currentDesktopUser()
	if err != nil {
		return nil
	}
	if user != "" {
		return fmt.Errorf("cannot restart LightDM, user %s is logged in", user)
	}

	return exec.Command("systemctl", "restart", "lightdm").Run()

}

// returns the logged-in desktop user on seat0 ("" if none)
func currentDesktopUser() (string, error) {
	out, err := exec.Command("loginctl", "list-sessions", "--no-legend").Output()
	if err != nil {
		return "", fmt.Errorf("failed to run loginctl: %w", err)
	}
	for line := range strings.SplitSeq(string(out), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		cols := strings.Fields(line)
		if len(cols) < 4 {
			continue
		}
		user := cols[2]
		seat := cols[3]
		if seat == "seat0" && user != "gdm" && user != "lightdm" {
			return user, nil
		}
	}
	return "", nil
}
