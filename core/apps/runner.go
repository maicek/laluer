package apps

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/maicek/laluer/core/history"
)

func (a *Application) Run() error {
	fmt.Printf("Running %s\n", a.Name)

	history.Service.Push(history.ENTRY_TYPE_APP, a.Path)

	if a.Terminal {
		// not implemented
	} else {
		// run app in fork and close parent
		cmd := exec.Command("sh", "-c", a.Exec)
		cmd.Env = os.Environ()
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setsid: true,
		}

		err := cmd.Start()
		if err != nil {
			return err
		}
	}

	os.Exit(0)

	return nil
}
