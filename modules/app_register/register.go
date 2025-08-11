package app_register

type IAppControl interface {
	Startup(ctx interface{})
}

var (
	AppControlList = make([]IAppControl, 0)
)

func RegisterAppControl(control IAppControl) {
	AppControlList = append(AppControlList, control)
}
