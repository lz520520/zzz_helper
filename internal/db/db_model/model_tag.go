package db_model

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	TagEmbedStruct = "sub_section"

	PrimaryFlag  = "(P) "
	UniqueFlag   = "(U) "
	FakeFlag     = "* "
	RequiredFlag = FakeFlag
)

type FieldAlias struct {
	Field string // 字段名
	Alias string // 别名
}

func GetFieldWithAlias(s interface{}) []FieldAlias {
	rt := reflect.TypeOf(s)
	fas := make([]FieldAlias, 0)

	for i := 0; i < rt.NumField(); i++ {
		if alias := rt.Field(i).Tag.Get("alias"); alias != "" {
			fas = append(fas, FieldAlias{
				Field: rt.Field(i).Tag.Get(fieldTag),
				Alias: alias,
			})
		}
	}
	return fas
}

func getVarNameFromFieldRecur(rt reflect.Type, field string) string {
	varName := ""
	for i := 0; i < rt.NumField(); i++ {
		tmp := strings.Split(rt.Field(i).Tag.Get(fieldTag), ",")[0]
		if tmp != "" && tmp == field {
			varName = rt.Field(i).Name
			break
		} else if rt.Field(i).Tag.Get(TagEmbedStruct) != "" {
			varName = getVarNameFromFieldRecur(rt.Field(i).Type, field)
			if varName != "" {
				break
			}
		}
	}
	return varName
}

// 根据数据库字段名获取结构体字段名
func GetVarNameFromField(s interface{}, field string) string {
	rt := reflect.TypeOf(s)
	varName := getVarNameFromFieldRecur(rt, field)
	return varName
}

func getVarNameFromJsonRecur(rt reflect.Type, json string) string {
	varName := ""
	for i := 0; i < rt.NumField(); i++ {
		tmp := strings.Split(rt.Field(i).Tag.Get("json"), ",")[0]

		if tmp != "" && tmp == json {
			varName = rt.Field(i).Name
			break
		} else if rt.Field(i).Tag.Get(TagEmbedStruct) != "" {
			varName = getVarNameFromJsonRecur(rt.Field(i).Type, json)
			if varName != "" {
				break
			}
		}
	}
	return varName
}

func GetVarNameFromJson(s interface{}, json string) string {
	rt := reflect.TypeOf(s)
	varName := getVarNameFromJsonRecur(rt, json)
	return varName
}

func getFieldFromAliasRecur(rt reflect.Type, alias string) string {
	field := ""
	for i := 0; i < rt.NumField(); i++ {
		if tmp := rt.Field(i).Tag.Get("alias"); tmp != "" && tmp == alias {
			field = strings.Split(rt.Field(i).Tag.Get(fieldTag), ",")[0]
			break
		} else if tmp = rt.Field(i).Tag.Get("data_alias"); tmp != "" && tmp == alias {
			field = strings.Split(rt.Field(i).Tag.Get(fieldTag), ",")[0]
			break
		} else if rt.Field(i).Tag.Get(TagEmbedStruct) != "" {
			field = getFieldFromAliasRecur(rt.Field(i).Type, alias)
			if field != "" {
				break
			}
		}
	}
	return field
}

// 根据字段别名获取字段名
func GetFieldFromAlias(s interface{}, alias string) string {
	alias = removeAliasFlag(alias)
	rt := reflect.TypeOf(s)
	field := getFieldFromAliasRecur(rt, alias)
	return field
}

func getStructAliasRecur(rt reflect.Type) []string {
	aliases := make([]string, 0)

	for i := 0; i < rt.NumField(); i++ {
		if alias := rt.Field(i).Tag.Get("alias"); alias != "" {
			if rt.Field(i).Tag.Get("search_only") != "true" {
				aliases = append(aliases, addAliasFlag(rt.Field(i).Tag.Get(fieldTag), alias))
			}
		} else if rt.Field(i).Tag.Get(TagEmbedStruct) != "" {
			aliases = append(aliases, getStructAliasRecur(rt.Field(i).Type)...)
		}

	}
	return aliases
}

func GetStructAlias(s interface{}) []string {
	rt := reflect.TypeOf(s)
	aliases := getStructAliasRecur(rt)
	return aliases
}

func getStructAliasWithSearchRecur(rt reflect.Type) []string {
	aliases := make([]string, 0)

	for i := 0; i < rt.NumField(); i++ {
		if alias := rt.Field(i).Tag.Get("alias"); alias != "" {
			aliases = append(aliases, addAliasFlag(rt.Field(i).Tag.Get(fieldTag), alias))
		} else if rt.Field(i).Tag.Get(TagEmbedStruct) != "" {
			aliases = append(aliases, getStructAliasWithSearchRecur(rt.Field(i).Type)...)
		}
	}

	return aliases
}
func GetStructAliasWithSearch(s interface{}) []string {
	rt := reflect.TypeOf(s)
	aliases := getStructAliasWithSearchRecur(rt)
	return aliases
}

func getStructAliasWithDataRecur(rt reflect.Type) []string {
	aliases := make([]string, 0)

	for i := 0; i < rt.NumField(); i++ {
		if alias := rt.Field(i).Tag.Get("alias"); alias != "" {
			aliases = append(aliases, addAliasFlag(rt.Field(i).Tag.Get(fieldTag), alias))
		} else if alias = rt.Field(i).Tag.Get("data_alias"); alias != "" {
			aliases = append(aliases, addAliasFlag(rt.Field(i).Tag.Get(fieldTag), alias))
		} else if rt.Field(i).Tag.Get(TagEmbedStruct) != "" {
			aliases = append(aliases, getStructAliasWithDataRecur(rt.Field(i).Type)...)

		}
	}

	return aliases
}
func GetStructAliasWithData(s interface{}) []string {
	rt := reflect.TypeOf(s)
	aliases := getStructAliasWithDataRecur(rt)
	return aliases
}

func hasPrimaryRecur(rt reflect.Type) bool {
	for i := 0; i < rt.NumField(); i++ {
		if strings.Contains(rt.Field(i).Tag.Get(fieldTag), "primary") {
			return true
		}
		if rt.Field(i).Tag.Get(TagEmbedStruct) != "" {
			if hasPrimaryRecur(rt.Field(i).Type) {
				return true
			}
		}

	}
	return false
}
func HasPrimary(s interface{}) bool {
	rt := reflect.TypeOf(s)
	return hasPrimaryRecur(rt)
}

func isPrimaryOrUniqueRecur(rt reflect.Type, alias string) bool {
	for i := 0; i < rt.NumField(); i++ {
		if tmp := rt.Field(i).Tag.Get("alias"); tmp != "" && tmp == alias {
			field := rt.Field(i).Tag.Get(fieldTag)
			if strings.Contains(field, ",primary") ||
				strings.Contains(field, ",unique") {
				return true
			}
		} else if tmp = rt.Field(i).Tag.Get("data_alias"); tmp != "" && tmp == alias {
			field := rt.Field(i).Tag.Get(fieldTag)
			if strings.Contains(field, ",primary") ||
				strings.Contains(field, ",unique") {
				return true
			}
		} else if rt.Field(i).Tag.Get(TagEmbedStruct) != "" {
			if isPrimaryOrUniqueRecur(rt.Field(i).Type, alias) {
				return true
			}
		}
	}
	return false
}

// 判断该别名是否为唯一值
func IsPrimaryOrUnique(s interface{}, alias string) bool {
	alias = removeAliasFlag(alias)
	return isPrimaryOrUniqueRecur(reflect.TypeOf(s), alias)
}

// 别名添加flag前缀，主要是表的特殊属性，比如primary和unique
func addAliasFlag(field, alias string) (flagAlias string) {
	tmp := strings.SplitN(field, ",", 2)
	if len(tmp) == 2 {
		switch tmp[1] {
		case "primary":
			flagAlias = PrimaryFlag + alias
		case "unique":
			flagAlias = UniqueFlag + alias
		default:
			flagAlias = alias

		}
	} else {
		flagAlias = alias
	}
	return
}

// 移除别名前缀flag
func removeAliasFlag(flagAlias string) (alias string) {
	if strings.HasPrefix(flagAlias, FakeFlag) {
		flagAlias = strings.TrimPrefix(flagAlias, FakeFlag)
	}
	if strings.HasPrefix(flagAlias, PrimaryFlag) {
		flagAlias = strings.TrimPrefix(flagAlias, PrimaryFlag)
	}
	if strings.HasPrefix(flagAlias, UniqueFlag) {
		flagAlias = strings.TrimPrefix(flagAlias, UniqueFlag)
	}
	return flagAlias
}

func getUniqueValueRecur(rt reflect.Type, rv reflect.Value) (string, string) {
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i).Tag.Get(fieldTag)
		if strings.Contains(field, ",unique") {
			return strings.Split(field, ",")[0], fmt.Sprint(rv.Field(i).Interface())
		} else if rt.Field(i).Tag.Get(TagEmbedStruct) != "" {
			if subField, subValue := getUniqueValueRecur(rt.Field(i).Type, rv.Field(i)); subField != "" {
				return subField, subValue
			}
		}
	}
	return "", ""
}

// 获取结构体实例中 unique字段的field和value
func GetUniqueValue(s interface{}) (string, string) {
	return getUniqueValueRecur(reflect.TypeOf(s), reflect.ValueOf(s))
}
