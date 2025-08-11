package main

import (
	"context"
	"crypto/tls"
	"embed"
	_ "embed"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	runtime2 "runtime"
	"strings"
	"zzz_helper/internal/config"
	vclexec "zzz_helper/internal/system/exec"
	_ "zzz_helper/modules"
	"zzz_helper/modules/app_register"
	"zzz_helper/res"
)

//go:embed build/appicon.png
var icon []byte

//go:embed frontend/dist
var assets embed.FS

func main() {
	os.Setenv("loglevel", config.ConsoleConfigInst.LogLevel)
	if config.ConsoleConfigInst.Ws {
		os.Setenv("ws", "true")
		if conns, err := net.Connections("tcp"); err == nil {
			for _, conn := range conns {
				if conn.Laddr.Port != 34115 {
					continue
				}
				proc, err := process.NewProcess(conn.Pid)
				if err != nil {
					continue
				}
				proc.Kill()

			}
		}
	}

	// Create an instance of the app structure
	// 创建一个App结构体实例
	app := NewApp()
	var AppMenu *menu.Menu
	if runtime2.GOOS == "windows" {
		AppMenu = menu.NewMenu()
		helpMenu := AppMenu.AddSubmenu("帮助")
		// keys.Key("f5")
		helpMenu.AddText("刷新", keys.Key("f5"), func(data *menu.CallbackData) {
			runtime.WindowReload(app.ctx)
		})
		// keys.Key("f5")
		helpMenu.AddText("打开目录", nil, func(data *menu.CallbackData) {
			vclexec.OpenDirectory(config.CurrentPath)
		})
	}

	//appControl := mental_control.NewMentalControl()
	//accControl := acc_control.NewAcceleratorControl()
	//railgunControl := railgun_control.NewRailgunControl()
	binds := []interface{}{app}
	for _, control := range app_register.AppControlList {
		binds = append(binds, control)
	}
	// Create application with options
	// 使用选项创建应用
	err := wails.Run(&options.App{
		Title:     strings.TrimSpace(res.TITLE) + "-" + res.VERSION,
		Width:     1420,
		Height:    900,
		MinWidth:  900,
		MinHeight: 600,
		//MaxWidth:          1200,
		//MaxHeight:         800,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 255, G: 255, B: 255, A: 0},
		Menu:              AppMenu,
		Logger:            nil,
		LogLevel:          logger.DEBUG,
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			for _, control := range app_register.AppControlList {
				control.Startup(ctx)
			}
		},
		OnDomReady:       app.domReady,
		OnBeforeClose:    app.beforeClose,
		OnShutdown:       app.shutdown,
		WindowStartState: options.Normal,
		WebSocket: options.WebSocket{
			WsOnly: config.ConsoleConfigInst.Ws,
			Server: &http.Server{
				Addr: "0.0.0.0:34115",
				TLSConfig: func() *tls.Config {
					return nil
				}(),
				//TLSConfig: mtls.MustNewServerTLSConfigFromBytes(nil,nil,nil),

			}},
		AssetServer: &assetserver.Options{
			Assets:     assets,
			Handler:    nil,
			Middleware: nil,
		},
		Bind: binds,

		// Windows platform specific options
		// Windows平台特定选项
		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               false,
			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: false,
			WebviewUserDataPath:               "",
			WebviewBrowserPath:                "",
			Theme:                             windows.SystemDefault,
		},
		// Mac platform specific options
		// Mac平台特定选项
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: false,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       false,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "Wails Template Vue",
				Message: "A Wails template based on Vue and Vue-Router",
				Icon:    icon,
			},
		},
		Linux: &linux.Options{
			Icon: icon,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
