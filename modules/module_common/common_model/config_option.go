package common_model

type FormComponentType int

const (
	Edit FormComponentType = iota + 1
	Combobox
	TextArea
	CheckBox
	CheckEdit
	CodeEditor

	Process
	TreeSelect
	File
	Password
)

type FormConfigType string

type CommonConfigOption struct {
	Title               string      `json:"title"` // 标题，也可以说是别名
	Key                 string      `json:"key"`   // 唯一标识符
	WebKey              string      `json:"web_key"`
	DefaultValue        interface{} `json:"default_value"` // 默认值
	DefaultValueDynamic bool        `json:"default_value_dynamic,omitempty"`

	/*
	   ShellFormValue
	    * 自动设置选项，如果为true表示不需要通过配置界面设置，也就是不需要显示， 默认false
	    * 如果未设置该项，则formType和configType是必须设置的
	    * */
	Tips           string   `json:"tips,omitempty"` // 提示语
	DefaultOptions []string `json:"default_options,omitempty"`

	AutoSetting bool               `json:"auto_setting,omitempty"`
	AutoSetFunc func() interface{} `json:"-" copier:"-"` // 在autoSetting为true时，用于自动生成参数使用

	FormRequired      bool              `json:"form_required,omitempty"` // 是否为表单必选项,默认为false
	FormHidden        bool              `json:"form_hidden"`
	FormComponentType FormComponentType `json:"form_component_type,omitempty"` // 表单组件类型
	FormConfigType    FormConfigType    `json:"form_config_type,omitempty"`    // 表单配置类型

	//    ShellTableValue
	ShowInTable bool   `json:"show_in_table"`        // 是否在table组件中显示
	ColWidth    int    `json:"col_width,omitempty"`  // 列宽
	ColFixed    string `json:"col_fixed,omitempty"`  // 列固定
	ColHidden   bool   `json:"col_hidden,omitempty"` // 列是否隐藏 TODO: 目前是通过设置colWidth=0来实现，后续在调整
	Edit        bool   `json:"edit,omitempty"`
	Sort        string `json:"sort,omitempty"`

	CustomFilterDropdown bool `json:"custom_filter_dropdown,omitempty"` // 是否设置过滤下拉菜单

	RowGroup bool `json:"row_group,omitempty"` // 行分组
}

type CommonConfigOptionsResp struct {
	Status  bool                 `json:"status"`
	Msg     string               `json:"msg"`
	Options []CommonConfigOption `json:"options"`
	Err     string               `json:"err"`
}
