package mem2

import (
	"fmt"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
	"os"
	runtime2 "runtime"
)

func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func GetProcMemFormat() (string, error) {
	m, err := GetProcMem()
	if err != nil {
		return "", err
	}
	return FormatBytes(m), nil
}

func GetSysMemFormat() (string, error) {
	s, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}
	return FormatBytes(s.Total), nil
}

func GetProcMem() (uint64, error) {
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return 0, err
	}
	memInfo, err := p.MemoryInfo()
	if err != nil {
		return 0, err
	}
	var procMem uint64
	if runtime2.GOOS == "windows" {
		procMem = memInfo.VMS
	} else {
		procMem = memInfo.RSS
	}
	return procMem, nil
}
