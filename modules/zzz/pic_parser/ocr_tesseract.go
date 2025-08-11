package pic_parser

import (
	"fmt"
	"os/exec"
	"syscall"
	"zzz_helper/internal/config"
)

func ParseWithTesseract(src string) (res string, err error) {
	cmd := exec.Command(config.ConsoleConfigInst.TesseractPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CmdLine:    fmt.Sprintf(`tesseract.exe %s stdout  -l chi_sim --psm 6 -c preserve_interword_spaces=1`, src),
		HideWindow: true,
	}
	ret, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(ret), nil
}
