package vclexec

import (
	"fmt"
	"os/exec"
)

func VclExec(filename string) {
	exec.Command("/bin/bash", "-c", fmt.Sprintf(`open "%s"`, filename)).Start()
}

func OpenDirectory(dir string) {
	exec.Command("/bin/bash", "-c", fmt.Sprintf(`open "%s"`, dir)).Start()
}

func ShellExec(cmd string) {
	exec.Command("/bin/bash", "-c", cmd).Start()
}
