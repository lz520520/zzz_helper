package db_model

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func GetReflectValueFromField(dbPtr reflect.Value, field string) reflect.Value {
	rv := dbPtr
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	count := rv.NumField()
	for i := 0; i < count; i++ {
		tmp := rv.Field(i)
		tmpField := strings.Split(rv.Type().Field(i).Tag.Get(fieldTag), ",")[0]
		if tmpField != "" && tmpField == field {
			return tmp
		} else if rv.Type().Field(i).Tag.Get(TagEmbedStruct) != "" {
			tmp = GetReflectValueFromField(tmp, field)
			if tmp.IsValid() {
				return tmp
			}
		}
	}
	return reflect.Value{}
}

//func GenerateDBPtrFromAlias(alias string) reflect.Value {
//	return reflect.New(reflect.TypeOf(GetModuleDBFromAlias(alias)))
//}

func SetDBValueForAlias(dbPtr reflect.Value, alias string, text string) error {
	field := GetFieldFromAlias(dbPtr.Elem().Interface(), alias)

	return SetDBValueForField(dbPtr, field, text)
}

func GetValueMapsForAlias(modDBs []interface{}) []map[string]string {
	datas := make([]map[string]string, 0)
	for _, modDB := range modDBs {
		tmp := GetValueMapForAlias(modDB)
		if len(tmp) != 0 {
			datas = append(datas, tmp)
		}
	}
	return datas
}

func GetValueMapForAlias(modDB interface{}) map[string]string {
	rv := reflect.ValueOf(modDB)
	count := rv.NumField()
	valueMap := make(map[string]string)
	for i := 0; i < count; i++ {
		tmp := rv.Field(i)
		tmpAlias := strings.Split(rv.Type().Field(i).Tag.Get("alias"), ",")[0]
		if tmpAlias != "" {
			valueMap[tmpAlias] = fmt.Sprint(tmp.Interface())
		} else if rv.Type().Field(i).Tag.Get(TagEmbedStruct) != "" {
			tmpMap := GetValueMapForAlias(tmp.Interface())
			for k, v := range tmpMap {
				valueMap[k] = v
			}
		}
	}
	return valueMap
}

func SetDBValueForField(dbPtr reflect.Value, field string, text string) error {
	value := GetReflectValueFromField(dbPtr, field)
	if !value.IsValid() {
		return fmt.Errorf("field to value is invalid")
	}
	if value.CanSet() {
		kind := value.Type().Kind()
		switch {
		case kind == reflect.String:
			value.SetString(text)
		case kind == reflect.Bool:
			var filterInt = 0
			var err error
			if text != "" {
				filterInt, err = strconv.Atoi(text)
				if err != nil {
					return err
				}
			}
			if filterInt == 1 {
				value.SetBool(true)
			} else {
				value.SetBool(false)
			}
		case kind >= reflect.Int && kind <= reflect.Uint64:
			var filterInt = 0
			var err error
			if text != "" {
				filterInt, err = strconv.Atoi(text)
				if err != nil {
					return err
				}
			}
			if kind <= reflect.Int64 {
				value.SetInt(int64(filterInt))
			} else {
				value.SetUint(uint64(filterInt))
			}
		case kind == reflect.Float64:
			var filterFloat float64
			var err error
			if text != "" {
				filterFloat, err = strconv.ParseFloat(text, 64)
				if err != nil {
					return err
				}
			}

			value.SetFloat(filterFloat)
		case kind == reflect.Struct:
			if value.Type().String() == "time.Time" {
				t, err := time.ParseInLocation("2006-01-02 15:04:05", text, time.Local)
				if err != nil {
					return err
				}
				value.Set(reflect.ValueOf(t))
			}
		}
	}
	return nil
}

func rowScan(rows *sql.Rows, dest interface{}) error {
	// 获取dest的反射值对象，确保它是一个指向结构体的指针
	v := reflect.ValueOf(dest).Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("dest must be a pointer to a struct")
	}

	// 为每个结构体字段构造一个指向该字段地址的reflect.Value
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	pointers := make([]interface{}, len(columns))
	fitMap := make(map[int]reflect.Value)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		kind := field.Kind()
		switch {
		case kind > reflect.Int && kind <= reflect.Uint64:
			fitValue := reflect.New(reflect.TypeOf(0))
			fitMap[i] = fitValue
			pointers[i] = fitValue.Interface()
		case kind == reflect.Bool:
			fitValue := reflect.New(reflect.TypeOf(0))
			fitMap[i] = fitValue
			pointers[i] = fitValue.Interface()
		default:
			pointers[i] = v.Field(i).Addr().Interface()
		}
	}

	// 使用构造的指针切片调用原生的Scan方法
	err = rows.Scan(pointers...)
	if err != nil {
		return err
	}
	for i, fitValue := range fitMap {
		fitValue = fitValue.Elem()
		origin := v.Field(i)
		switch fitValue.Kind() {
		case reflect.Int:
			fitInt := fitValue.Int()
			originKind := origin.Kind()
			switch {
			case originKind == reflect.Bool:
				if fitInt == 1 {
					origin.SetBool(true)
				} else {
					origin.SetBool(false)
				}
			case originKind >= reflect.Int && originKind <= reflect.Int64:
				origin.SetInt(fitInt)
			case originKind >= reflect.Uint && originKind <= reflect.Uint64:
				origin.SetUint(uint64(fitInt))
			}
		}
	}
	return nil
}
