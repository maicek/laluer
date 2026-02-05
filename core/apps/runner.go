package apps

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
)

// stripFieldCodes removes .desktop Exec field codes like %u, %U, %f, %F, %i, %c, %k
var fieldCodeRegex = regexp.MustCompile(`%[uUfFdDnNickvm]`)

func cleanExec(execStr string) string {
	cleaned := fieldCodeRegex.ReplaceAllString(execStr, "")
	cleaned = strings.Join(strings.Fields(cleaned), " ")
	return strings.TrimSpace(cleaned)
}

func (a *Application) Run() error {
	fmt.Printf("Running %s\n", a.Name)

	execCmd := cleanExec(a.Exec)
	if execCmd == "" {
		return fmt.Errorf("empty exec command for %s", a.Name)
	}

	if a.Terminal {
		// try common terminal emulators
		terminals := []string{"kitty", "alacritty", "foot", "wezterm", "xterm"}
		for _, term := range terminals {
			if _, err := exec.LookPath(term); err == nil {
				execCmd = fmt.Sprintf("%s -e %s", term, execCmd)
				break
			}
		}
	}

	cmd := exec.Command("sh", "-c", execCmd)
	cmd.Env = os.Environ()
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	err := cmd.Start()
	if err != nil {
		return err
	}

	os.Exit(0)

	return nil
}
