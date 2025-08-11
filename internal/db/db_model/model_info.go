package db_model

import "reflect"

type Module string

type ModuleInfo struct {
	Name   Module
	Alias  string
	DB     CRUD // 结构体
	DBFunc func(dbName string) *CommonDB[CRUD]

	GRPC bool
}

var (
	ModuleInfos = make([]interface{}, 0) // []ModuleInfo[CRUD]
)

func ModuleInfoRegister(info ModuleInfo) {
	ModuleInfos = append(ModuleInfos, info)
}

/*func UnmarshalFilterInfos(data []byte, alias string) (filterInfos []FilterInfo, err error) {
	tmpFilterInfos := make([]FilterInfo, 0)
	filterInfos = make([]FilterInfo, 0)

	modDB := GetModuleDBFromAlias(alias)

	err = json.Unmarshal(data, &tmpFilterInfos)
	if err != nil {
		return
	}

	for _, fi := range tmpFilterInfos {
		var dbData []byte
		dbData, err = json.Marshal(fi.DB)
		if err != nil {
			return
		}
		tmpDBPtr := reflect.New(reflect.TypeOf(modDB).Elem())
		tmpMap := make(map[string]interface{})
		err = json.Unmarshal(dbData, &tmpMap)
		if err != nil {
			return
		}
		for k, v := range tmpMap {
			varName := GetVarNameFromJson(tmpDBPtr.Elem().Interface(), k)
			if varName == "" {
				continue
			}
			value := tmpDBPtr.Elem().FieldByName(varName)
			if value.CanSet() {
				kind :=  value.Type().Kind()
				switch{
				case kind == reflect.String:
					if vv, ok := v.(string); ok {
						value.SetString(vv)
					}
				case kind == reflect.Bool:
					if vv, ok := v.(float64); ok {
						if vv == 1 {
							value.SetBool(true)
						} else {
							value.SetBool(false)
						}
					}
				case kind >= reflect.Int && kind <= reflect.Int64:
					if vv, ok := v.(float64); ok {
						value.SetInt(int64(vv))
					}
				case kind >= reflect.Uint && kind <= reflect.Uint64:
					if vv, ok := v.(float64); ok {
						value.SetUint(uint64(vv))
					}
				case kind == reflect.Float64:
					if vv, ok := v.(float64); ok {
						value.SetFloat(vv)
					}
				case kind == reflect.Struct:
					if value.Type().String() == "time.Time" {
						if vv, ok := v.(string); ok {
							// json序列化，对于时间的解析是按照RFC3339标准
							vTime, err := time.ParseInLocation(time.RFC3339, vv, time.Local)
							if err == nil {
								value.Set(reflect.ValueOf(vTime))
							}

						}
					}
				}
			}
		}

		fi.DB = tmpDBPtr.Elem().Interface()
		filterInfos = append(filterInfos, fi)
	}
	return
}

func UnmarshalQueryResults(data []byte, alias string) (results []interface{}, err error) {
	tmpResults := make([]interface{}, 0)
	results = make([]interface{}, 0)

	modDB := GetModuleDBFromAlias(alias)

	err = json.Unmarshal(data, &tmpResults)
	if err != nil {
		return
	}

	for _, fi := range tmpResults {
		var dbData []byte
		dbData, err = json.Marshal(fi)
		if err != nil {
			return
		}
		tmpDBPtr := reflect.New(reflect.TypeOf(modDB).Elem())
		tmpMap := make(map[string]interface{})
		err = json.Unmarshal(dbData, &tmpMap)
		if err != nil {
			return
		}
		for k, v := range tmpMap {
			varName := GetVarNameFromJson(tmpDBPtr.Elem().Interface(), k)
			if varName == "" {
				continue
			}
			value := tmpDBPtr.Elem().FieldByName(varName)
			if value.CanSet() {
				kind := value.Type().Kind()
				switch  {
				case kind == reflect.String:
					if vv, ok := v.(string); ok {
						value.SetString(vv)
					}
				case kind == reflect.Bool:
					if vv, ok := v.(float64); ok {
						if vv == 1 {
							value.SetBool(true)
						} else {
							value.SetBool(false)
						}
					}
				case kind >= reflect.Int && kind <= reflect.Int64:
					if vv, ok := v.(float64); ok {
						value.SetInt(int64(vv))
					}
				case kind >= reflect.Uint && kind <= reflect.Uint64:
					if vv, ok := v.(float64); ok {
						value.SetUint(uint64(vv))
					}
				case kind == reflect.Float64:
					if vv, ok := v.(float64); ok {
						value.SetFloat(vv)
					}
				case kind == reflect.Struct:
					if value.Type().String() == "time.Time" {
						if vv, ok := v.(string); ok {
							// json序列化，对于时间的解析是按照RFC3339标准
							vTime, err := time.ParseInLocation(time.RFC3339, vv, time.Local)
							if err == nil {
								value.Set(reflect.ValueOf(vTime))
							}

						}
					}
				}
			}
		}

		results = append(results, tmpDBPtr.Elem().Interface())
	}
	return
}
*/
// 根据数据库别名获取数据库结构体
func GetModuleDBFromAlias(alias string) CRUD {
	for _, info := range ModuleInfos {

		if info.(ModuleInfo).Alias == alias {
			return info.(ModuleInfo).DB
		}
	}
	return nil
}
func GetModuleDBFromName(name string) interface{} {
	for _, info := range ModuleInfos {
		if info.(ModuleInfo).Name == Module(name) {

			return reflect.ValueOf(info.(ModuleInfo).DB).Elem().Interface()
		}
	}
	return nil
}
func GetModuleDBFuncFromAlias(alias string) func(name string) *CommonDB[CRUD] {
	for _, info := range ModuleInfos {
		if info.(ModuleInfo).Alias == alias {
			return info.(ModuleInfo).DBFunc
		}
	}
	return nil
}
func GetModuleDBFuncFromName(name string) func(name string) *CommonDB[CRUD] {
	for _, info := range ModuleInfos {
		val := reflect.ValueOf(info)

		// 获取非指针类型的实际值
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		// 检查字段类型
		if val.Kind() == reflect.Struct {
			nameField := val.FieldByName("Name")
			dbFuncField := val.FieldByName("DBFunc")

			// 检查字段是否有效并匹配
			if nameField.IsValid() && nameField.String() == name && dbFuncField.IsValid() {
				// 断言并返回 func() *CommonDB[CRUD] 类型的 DBFunc
				if dbFunc, ok := dbFuncField.Interface().(func(name string) *CommonDB[CRUD]); ok {
					return dbFunc
				}
			}
		}
	}
	return nil
}

// 根据数据库别名获取数据库名
func GetModuleNameFromAlias[T CRUD](alias string) Module {
	for _, info := range ModuleInfos {
		if info.(ModuleInfo).Alias == alias {
			return info.(ModuleInfo).Name
		}
	}
	return ""
}

func IsGRPCModule[T CRUD](alias string) bool {
	for _, info := range ModuleInfos {
		if info.(ModuleInfo).Alias == alias {
			return info.(ModuleInfo).GRPC
		}
	}
	return false
}
