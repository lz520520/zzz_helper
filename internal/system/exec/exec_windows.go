package vclexec

import (
	"fmt"
	"os/exec"
	"syscall"
)

//	func VclExec(filename string) {
//		win.ShellExecute(vcl.Application.MainFormHandle(), "", filename, "", "", win.SW_SHOWNORMAL)
//
// }
func OpenDirectory(dir string) {
	c := exec.Command("cmd.exe")
	c.SysProcAttr = &syscall.SysProcAttr{
		CmdLine: fmt.Sprintf("cmd.exe /c cd /d \"%s\" && start .", dir),
	}
	c.Start()
}

func ShellExec(cmd string) {
	c := exec.Command("cmd.exe", "/c", cmd)
	c.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	c.Start()
}
