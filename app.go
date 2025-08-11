package main

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gcharset"
	"github.com/shirou/gopsutil/v4/process"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"os/exec"
	runtime2 "runtime"
	"runtime/debug"
	"time"
	"zzz_helper/internal/mylog"
	"zzz_helper/internal/utils/file2"
	"zzz_helper/internal/utils/mem2"
	"zzz_helper/modules/module_common/common_model"
)

type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
// NewApp 创建一个新的 App 应用程序
func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	// 在这里执行初始化设置
	a.ctx = ctx
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		//runtime.EventsEmit(a.ctx, "app_status123", "test")
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				runtime.EventsEmit(a.ctx, "app_status", a.GetAppStatus())
			case <-ctx.Done():
				return
			}
		}
	}()
}

// domReady is called after the front-end dom has been loaded
// domReady 在前端Dom加载完毕后调用
func (a *App) domReady(ctx context.Context) {
	// Add your action here
	// 在这里添加你的操作
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue,
// false will continue shutdown as normal.
// beforeClose在单击窗口关闭按钮或调用runtime.Quit即将退出应用程序时被调用.
// 返回 true 将导致应用程序继续，false 将继续正常关闭。
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
// 在应用程序终止时被调用
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
	// 在此处做一些资源释放的操作
}

func (a *App) WrapOpenFileDialog(DefaultFilename string) (resp common_model.CommonResp) {
	filepath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		DefaultDirectory:           "",
		DefaultFilename:            DefaultFilename,
		Title:                      "打开",
		Filters:                    nil,
		ShowHiddenFiles:            true,
		CanCreateDirectories:       true,
		ResolvesAliases:            true,
		TreatPackagesAsDirectories: true,
	})
	if err != nil {
		resp.Err = err.Error()
		return
	}
	resp.Msg = filepath
	resp.Status = true
	return
}

func (a *App) WrapOpenDirectoryDialog(DefaultFilename string) (resp common_model.CommonResp) {
	filepath, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		DefaultDirectory:           "",
		DefaultFilename:            "",
		Title:                      "选择文件夹",
		Filters:                    nil,
		ShowHiddenFiles:            true,
		CanCreateDirectories:       true,
		ResolvesAliases:            true,
		TreatPackagesAsDirectories: true,
	})
	if err != nil {
		resp.Err = err.Error()
		return
	}
	resp.Msg = filepath
	resp.Status = true
	return
}

func (a *App) WrapSaveFileDialog(DefaultFilename string) (resp common_model.CommonResp) {
	filepath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultDirectory:           "",
		DefaultFilename:            DefaultFilename,
		Title:                      "保存",
		Filters:                    nil,
		ShowHiddenFiles:            true,
		CanCreateDirectories:       true,
		TreatPackagesAsDirectories: true,
	})
	if err != nil {
		resp.Err = err.Error()
		return
	}
	resp.Msg = filepath
	resp.Status = true
	return
}
func (a *App) LocalFileWrite(filename string, data []byte, position int64) (resp common_model.CommonResp) {
	var err error
	if position == 0 {
		err = file2.WriteFile(filename, data)
	} else {
		err = file2.AppendFile(filename, data)
	}
	if err != nil {
		resp.Err = err.Error()
		return
	}
	resp.Status = true
	return

}
func (a *App) WriteFile(filename string, data []byte) (resp common_model.CommonResp) {
	err := file2.WriteFile(filename, data)
	if err != nil {
		resp.Err = err.Error()
		return
	}
	resp.Status = true
	return
}

func (a *App) ReadFile(filename string) (resp common_model.CommonBytesResp) {
	b, err := file2.ReadFileBytes(filename)
	if err != nil {
		resp.Err = err.Error()
		return
	}
	resp.Bytes = b
	resp.Status = true
	return
}

func (this *App) LanguageEncode(req common_model.LanguageEncodeReq) (resp common_model.LanguageEncodeResp) {
	dstStr, err := gcharset.Convert(req.DstCharset, req.SrcCharset, req.Data)
	if err == nil {
		resp.Data = dstStr
	} else {
		resp.Data = req.Data
	}
	resp.Status = true
	return
}

func (a *App) GC() (resp common_model.CommonResp) {
	runtime2.GC()
	debug.FreeOSMemory()
	runtime.EventsOffAll(a.ctx)
	//runtime.WindowReload(a.ctx)
	resp.Status = true
	return
}

func (a *App) Restart() (resp common_model.CommonResp) {
	executable, err := os.Executable()
	if err != nil {
		resp.Err = err.Error()
		return
	}
	cmd := exec.Command(executable, os.Args[1:]...)
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	//cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	// 启动新进程，替代当前进程
	if err = cmd.Start(); err != nil {
		resp.Err = err.Error()
		return
	}
	time.Sleep(3 * time.Second)
	resp.Status = true
	os.Exit(0)
	return
}

type AppStatus struct {
	RunMode     string `json:"run_mode"`
	MemoryUsage string `json:"memory_usage"`
	CPUPercent  string `json:"cpu_percent"`
}

func (a *App) GetAppStatus() (status AppStatus) {

	newProcess, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		mylog.CommonLogger.Err(err).Send()
		return
	}
	cpuPercent, _ := newProcess.CPUPercent()
	status.CPUPercent = fmt.Sprintf("%.2f%%", cpuPercent)

	status.MemoryUsage, err = mem2.GetProcMemFormat()
	if err != nil {
		mylog.CommonLogger.Err(err).Send()
		return
	}
	return
}
