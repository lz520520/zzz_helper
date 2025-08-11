package reflect2

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func GetBoolParam(v interface{}) (value bool, err error) {
	switch reflect.ValueOf(v).Kind() {
	case reflect.Bool:
		value = reflect.ValueOf(v).Bool()
	case reflect.String:
		str := v.(string)
		if str == "true" {
			value = true
		} else {
			value = false
		}
	default:
		err = fmt.Errorf("value is not bool")
	}
	return
}
func GetIntParam(v interface{}, bitSize int) (value int64, err error) {
	switch reflect.ValueOf(v).Kind() {
	case reflect.Float32, reflect.Float64:
		value = int64(reflect.ValueOf(v).Float())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value = int64(reflect.ValueOf(v).Int())
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value = int64(reflect.ValueOf(v).Uint())
	case reflect.String:
		var tmpValue int64
		tmpValue, err = strconv.ParseInt(reflect.ValueOf(v).String(), 10, bitSize)
		if err != nil {
			return
		}
		value = tmpValue
	case reflect.Slice:
		var tmpValue int64
		tmpValue, err = strconv.ParseInt(string(reflect.ValueOf(v).Bytes()), 10, bitSize)
		if err != nil {
			return
		}
		value = tmpValue
	default:
		err = fmt.Errorf("value is not int")
	}
	return
}

func GetFloatParam(v interface{}, bitSize int) (value float64, err error) {
	switch reflect.ValueOf(v).Kind() {
	case reflect.Float32, reflect.Float64:
		value = reflect.ValueOf(v).Float()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value = float64(reflect.ValueOf(v).Int())
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value = float64(reflect.ValueOf(v).Uint())
	case reflect.String:
		var tmpValue float64
		tmpValue, err = strconv.ParseFloat(reflect.ValueOf(v).String(), bitSize)
		if err != nil {
			return
		}
		value = tmpValue
	case reflect.Slice:
		var tmpValue float64
		tmpValue, err = strconv.ParseFloat(reflect.ValueOf(v).String(), bitSize)
		if err != nil {
			return
		}
		value = tmpValue
	default:
		err = fmt.Errorf("value is not float")
	}
	return
}
func GetUintParam(v interface{}, bitSize int) (value uint64, err error) {
	switch reflect.ValueOf(v).Kind() {
	case reflect.Float32, reflect.Float64:
		value = uint64(reflect.ValueOf(v).Float())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value = uint64(reflect.ValueOf(v).Int())
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value = uint64(reflect.ValueOf(v).Uint())
	case reflect.String:
		var tmpValue uint64
		tmpValue, err = strconv.ParseUint(reflect.ValueOf(v).String(), 10, bitSize)
		if err != nil {
			return
		}
		value = tmpValue
	case reflect.Slice:
		var tmpValue uint64
		tmpValue, err = strconv.ParseUint(string(reflect.ValueOf(v).Bytes()), 10, bitSize)
		if err != nil {
			return
		}
		value = tmpValue
	default:
		err = fmt.Errorf("value is not uint")
	}
	return
}

func Map2Struct(m map[string]interface{}, s interface{}) (interface{}, error) {
	rfInfo, err := NewReflectInfo(s, "json")
	if err != nil {
		return nil, err
	}
	err = rfInfo.Map2Struct(m)
	if err != nil {
		return nil, err
	}
	return rfInfo.StructInterface(), nil
}

func Struct2Map(s interface{}) (map[string]interface{}, error) {
	rfInfo, err := NewReflectInfo(s, "json")
	if err != nil {
		return nil, err
	}
	return rfInfo.Struct2Map(), nil
}

func NewReflectInfo(origin interface{}, tag string) (*ReflectInfo, error) {
	// 获取输入参数的类型
	value := reflect.TypeOf(origin)

	if value.Kind() != reflect.Ptr {
		if value.Kind() == reflect.Struct {
			return &ReflectInfo{
				tag:    tag,
				origin: origin,
			}, nil
		}
		return nil, fmt.Errorf("input is neither a pointer nor a struct")
	}

	// 如果输入是指针但不是结构体指针，则返回错误
	if value.Kind() != reflect.Ptr && value.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("input is not a pointer to struct")
	}

	info := &ReflectInfo{
		tag:    tag,
		origin: origin,
		isPtr:  true,
	}
	return info, nil
}

type ReflectInfo struct {
	origin interface{}
	tag    string
	isPtr  bool
}

func (this *ReflectInfo) Interface() interface{} {
	return this.origin
}

func (this *ReflectInfo) StructInterface() interface{} {
	origin := this.origin
	if this.isPtr {
		origin = reflect.ValueOf(origin).Elem().Interface()
	}
	return origin
}
func (this *ReflectInfo) Struct2Map() map[string]interface{} {
	rt := reflect.TypeOf(this.origin)
	rv := reflect.ValueOf(this.origin)
	if this.isPtr {
		rt = rt.Elem()
		rv = rv.Elem()
	}
	result := make(map[string]interface{})

	for i := 0; i < rt.NumField(); i++ {
		key := strings.Split(rt.Field(i).Tag.Get(this.tag), ",")[0]
		result[key] = rv.Field(i).Interface()

	}

	return result
}

func (this *ReflectInfo) ReNew() error {
	if this.isPtr {
		this.origin = reflect.New(reflect.TypeOf(this.origin).Elem()).Interface()
	} else {
		this.origin = reflect.New(reflect.TypeOf(this.origin)).Interface()
		this.isPtr = true
	}
	return nil
}
func (this *ReflectInfo) Map2Struct(m map[string]interface{}) error {
	if !this.isPtr {
		this.ReNew()
	}
	for k, v := range m {
		err := this.SetValueByPath(k, v)
		if err != nil {
			return err
		}
	}
	return nil

}

func (this *ReflectInfo) GetFieldByNameOrTag(value reflect.Value, name string) reflect.Value {
	//field := value.FieldByName(name)
	//if field.IsValid() {
	//    return field
	//}
	// 如果字段名不是有效的，尝试使用 yaml 标签
	for i := 0; i < value.Type().NumField(); i++ {
		fieldType := value.Type().Field(i)
		if strings.Split(fieldType.Tag.Get(this.tag), ",")[0] == name {
			return value.Field(i)
		}
	}

	return reflect.Value{}
}
func (this *ReflectInfo) GetValueByPath(path string) (interface{}, error) {
	value := reflect.ValueOf(this.origin)
	if this.isPtr {
		value = value.Elem()
	}
	if path != "" {
		parts := strings.Split(path, ".")
		for _, part := range parts {
			value = this.GetFieldByNameOrTag(value, part)
			if !value.IsValid() {
				return "", fmt.Errorf("field %s not found", part)
			}
		}
	}
	return value.Interface(), nil
}

func (this *ReflectInfo) SetValueByPath(path string, newValue interface{}) error {
	if !this.isPtr {
		return fmt.Errorf("input is not struct pointer")
	}
	value := reflect.ValueOf(this.origin)
	value = value.Elem()

	parts := strings.Split(path, ".")
	for i, part := range parts {
		if i == len(parts)-1 {
			field := this.GetFieldByNameOrTag(value, part)
			if !field.IsValid() {
				return fmt.Errorf("field %s not found", part)
			}
			if !field.CanSet() {
				return fmt.Errorf("field %s cannot be set", part)
			}
			kind := field.Type().Kind()
			switch {
			case kind >= reflect.Int && kind <= reflect.Int64:
				intValue, err := GetIntParam(newValue, 64)
				if err != nil {
					return err
				}
				field.SetInt(intValue)
			case kind >= reflect.Uint && kind <= reflect.Uint64:
				uintValue, err := GetUintParam(newValue, 64)
				if err != nil {
					return err
				}
				field.SetUint(uintValue)
			case kind >= reflect.Float32 && kind <= reflect.Float64:
				floatValue, err := GetFloatParam(newValue, 64)
				if err != nil {
					return err
				}
				field.SetFloat(floatValue)
			case kind == reflect.String:
				field.SetString(fmt.Sprintf("%v", newValue))
			case kind == reflect.Bool:
				boolValue, err := GetBoolParam(newValue)
				if err != nil {
					return err
				}
				field.SetBool(boolValue)
			case kind == reflect.Slice && field.Type().Elem().Kind() == reflect.Uint8:
				if vv, ok := newValue.([]byte); ok {
					field.SetBytes(vv)
				}

				//boolValue, err := GetBoolParam(newValue)
				//if err != nil {
				//    return err
				//}
				//field.SetBool(boolValue)
			case kind == reflect.Slice && reflect.TypeOf(newValue).Elem() != nil && field.Type().Elem().Kind() == reflect.TypeOf(newValue).Elem().Kind():
				field.Set(reflect.ValueOf(newValue))

			default:
				return fmt.Errorf("not support modified field")
			}
			return nil
		}
		value = this.GetFieldByNameOrTag(value, part)
		if !value.IsValid() {
			return fmt.Errorf("field %s not found", part)
		}
	}

	return nil
}

func MustToStructPtr(val interface{}) interface{} {
	// 获取类型和值
	v := reflect.ValueOf(val)
	t := reflect.TypeOf(val)

	// 检查是否为结构体
	if t.Kind() != reflect.Struct {
		panic(fmt.Errorf("expected struct, got %s", t.Kind()))
	}

	// 创建指向该结构体的指针
	ptr := reflect.New(t)
	ptr.Elem().Set(v)

	// 返回指针作为 interface{}
	return ptr.Interface()
}

func ToStructPtr(val interface{}) (interface{}, error) {
	// 获取类型和值
	v := reflect.ValueOf(val)
	t := reflect.TypeOf(val)

	// 检查是否为结构体
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %s", t.Kind())
	}

	// 创建指向该结构体的指针
	ptr := reflect.New(t)
	ptr.Elem().Set(v)

	// 返回指针作为 interface{}
	return ptr.Interface(), nil
}
