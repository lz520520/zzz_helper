package main_control

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"zzz_helper/internal/utils/sync2"
	"zzz_helper/modules/app_register"
	"zzz_helper/modules/module_common/common_model"

	_ "zzz_helper/internal/db/db_example"
	_ "zzz_helper/internal/db/db_zzz"
)

type Control struct {
	ctx        context.Context
	taskCaches *sync2.OrderMap[string, common_model.TaskStatus]
}

// NewApp 创建一个新的 AppControl 应用程序
func NewControl() *Control {
	return &Control{}
}

func (this *Control) Startup(ctx interface{}) {
	// Perform your setup here
	// 在这里执行初始化设置
	this.ctx = ctx.(context.Context)
	this.taskCaches = sync2.NewOrderMap[string, common_model.TaskStatus]()

}

func (a *Control) eventEmitCallBack(eventName string, optionalData ...interface{}) {
	runtime.EventsEmit(a.ctx, eventName, optionalData...)
}

func init() {
	app_register.RegisterAppControl(NewControl())
}
