package variable

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"zzz_helper/modules/module_common/common_model"
)

type IFormHelper interface {
	FormDefaultValue(key string) interface{}
	FormAutoSetting(values map[string]interface{}) (err error)
	FormUpdateValue(values map[string]interface{}) (err error)
}

func TransToOptions(s interface{}) (options []common_model.CommonConfigOption, err error) {
	options = make([]common_model.CommonConfigOption, 0)
	rf := reflect.TypeOf(s)
	if rf.Kind() != reflect.Struct {
		err = fmt.Errorf("input is not struct")
		return
	}
	ptr := reflect.New(rf).Interface()
	count := rf.NumField()
	//rv := reflect.ValueOf(s)
	orders := make(map[int]string)
	for i := 0; i < count; i++ {
		field := rf.Field(i)
		fieldName := strings.Split(field.Tag.Get("json"), ",")[0]
		title := field.Tag.Get("title")
		if fieldName != "" && title != "" {

			showInTable := true
			if v := field.Tag.Get("show_in_table"); v == "false" {
				showInTable = false
			}
			colWidth := 150
			if v := field.Tag.Get("col_width"); v != "" {
				colWidth, _ = strconv.Atoi(v)
			}
			colHidden := false
			if v := field.Tag.Get("col_hidden"); v == "true" {
				colHidden = true
			}
			edit := false
			if v := field.Tag.Get("edit"); v == "true" {
				edit = true
			}
			sort := field.Tag.Get("sort")

			customFilterDropdown := true
			if v := field.Tag.Get("custom_filter_dropdown"); v == "false" {
				customFilterDropdown = false
			}
			var defaultValue interface{}
			if v := field.Tag.Get("default_value"); v != "" {
				if v == "true" {
					defaultValue = true
				} else if v == "false" {
					defaultValue = false
				} else if strings.Contains(v, ",") {
					defaultValue = strings.Split(v, ",")
				} else {
					defaultValue = v
				}
			} else {
				defaultValue = reflect.Zero(field.Type).Interface()
			}
			if helper, ok := ptr.(IFormHelper); ok {
				if v := helper.FormDefaultValue(fieldName); v != nil {
					defaultValue = v
				}
			}

			default_value_dynamic := false
			if v := field.Tag.Get("default_value_dynamic"); v == "true" {
				default_value_dynamic = true
			}
			defaultOptions := make([]string, 0)
			if v := field.Tag.Get("default_options"); v != "" {

				defaultOptions = strings.Split(v, ",")
			}

			autoSetting := false
			if v := field.Tag.Get("auto_setting"); v == "true" {
				autoSetting = true
			}

			formRequired := false
			var formComponentType common_model.FormComponentType = 0
			formConfigType := ""
			if v := field.Tag.Get("form_required"); v == "true" {
				formRequired = true
			}

			formHidden := false
			if v := field.Tag.Get("form_hidden"); v == "true" {
				formHidden = true
			}
			rowGroup := false
			if v := field.Tag.Get("row_group"); v == "true" {
				rowGroup = true
			}

			switch field.Tag.Get("form_component_type") {
			case "Edit":
				formComponentType = common_model.Edit
			case "Combobox":
				formComponentType = common_model.Combobox
			case "TextArea":
				formComponentType = common_model.TextArea
			case "CheckBox":
				formComponentType = common_model.CheckBox
			case "CheckEdit":
				formComponentType = common_model.CheckEdit
			case "CodeEditor":
				formComponentType = common_model.CodeEditor
			case "Process":
				formComponentType = common_model.Process
			case "TreeSelect":
				formComponentType = common_model.TreeSelect
			case "File":
				formComponentType = common_model.File
			case "Password":
				formComponentType = common_model.Password

			}
			formConfigType = field.Tag.Get("form_config_type")

			if v := field.Tag.Get("order"); v != "" {
				order, _ := strconv.Atoi(v)
				orders[order] = fieldName
			}
			option := common_model.CommonConfigOption{
				Title:  title,
				Key:    fieldName,
				WebKey: field.Name,
				Tips:   field.Tag.Get("tips"),

				DefaultValue:        defaultValue,
				DefaultValueDynamic: default_value_dynamic,
				DefaultOptions:      defaultOptions,
				AutoSetting:         autoSetting,

				ShowInTable:          showInTable,
				Edit:                 edit,
				Sort:                 sort,
				ColWidth:             colWidth,
				ColFixed:             field.Tag.Get("col_fixed"),
				ColHidden:            colHidden,
				CustomFilterDropdown: customFilterDropdown,

				FormRequired:      formRequired,
				FormHidden:        formHidden,
				FormComponentType: formComponentType,
				FormConfigType:    common_model.FormConfigType(formConfigType),
				RowGroup:          rowGroup,
			}
			options = append(options, option)
		}
	}
	if len(orders) > 0 {
		newOptions := make([]common_model.CommonConfigOption, 0)
		for i := 0; i < len(orders); i++ {
			for _, option := range options {
				if option.Key == orders[i] {
					newOptions = append(newOptions, option)
					break
				}
			}
		}
		options = newOptions
	}
	return

}
