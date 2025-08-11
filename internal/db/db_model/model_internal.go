package db_model

import (
	"database/sql"
	"fmt"
	"github.com/gogf/gf/v2/text/gstr"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"zzz_helper/internal/config"
	"zzz_helper/internal/utils/file2"
	"zzz_helper/internal/utils/slice2"
)

func (self *CommonDB[T]) initGlobalDB() (err error) {
	if *self.db != nil && (*self.db).Ping() == nil {
		return
	}
	// 创建目录
	dbAbsPath := self.dbPath
	if !filepath.IsAbs(dbAbsPath) {
		dbAbsPath = filepath.Join(config.CurrentPath, self.dbPath)
	}
	err = file2.MkdirFromFile(dbAbsPath)
	if err != nil {
		return err
	}
	// 打开数据库
	*self.db, err = sql.Open("sqlite3", dbAbsPath)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	//fmt.Printf("%#v\n", defaultDB)
	return nil
}

func initRecursion(rt reflect.Type) ([]arg, bool) {
	args := make([]arg, 0)
	primary := false
	for i := 0; i < rt.NumField(); i++ {
		argType := ""
		kind := rt.Field(i).Type.Kind()
		switch {
		case kind == reflect.String:
			length := "64"
			if tmp := rt.Field(i).Tag.Get("len"); tmp != "" {
				length = tmp
			}

			argType = fmt.Sprintf("VARCHAR(%s) DEFAULT \"\"", length)
		case kind == reflect.Bool:
			argType = "INTEGER DEFAULT -1"
		case kind >= reflect.Int && kind <= reflect.Uint64:
			argType = "INTEGER DEFAULT -1"
		case kind == reflect.Float64:
			argType = "FLOAT DEFAULT -1"
		case kind == reflect.Struct:
			if rt.Field(i).Type.String() == "time.Time" {
				argType = "TIMESTAMP DEFAULT (datetime('now', 'localtime'))"
			} else if rt.Field(i).Tag.Get(TagEmbedStruct) != "" {
				tmpArgs, tmpPrimary := initRecursion(rt.Field(i).Type)
				if !primary {
					primary = tmpPrimary
				}
				if len(tmpArgs) != 0 {
					args = append(args, tmpArgs...)
				}
				continue
			}

		}

		name := rt.Field(i).Tag.Get(fieldTag)
		if name == "" {
			continue
		}
		if strings.Contains(name, ",") {
			tmp := strings.Split(name, ",")
			name = tmp[0]
			switch tmp[1] {
			case "primary":
				argType += " PRIMARY KEY AUTOINCREMENT"
				primary = true
			case "unique":
				argType += " UNIQUE"
			}

		}

		a := arg{
			name:    name,
			argType: argType,
		}
		args = append(args, a)
	}
	return args, primary

}

func (self *CommonDB[T]) init() (err error) {
	if err = self.initGlobalDB(); err != nil {
		return
	}
	rt := reflect.TypeOf(self.crud).Elem()
	self.args = make([]arg, 0)
	self.args, self.primary = initRecursion(rt)

	// initsql
	initExpr := self.initSql()
	useInMatchInitExpr := strings.ReplaceAll(initExpr, " if not exists", "")
	useInMatchInitExpr = strings.TrimSpace(useInMatchInitExpr)
	useInMatchInitExpr = strings.TrimRight(useInMatchInitExpr, ";")
	// 查询表是否存在，并且做相似度检测，如果有改动，则重新创建表
	row := (*self.db).QueryRow(fmt.Sprintf("SELECT sql FROM sqlite_master WHERE type='table' and name='%s';", self.crud.DBName()))
	oldInitExpr := ""
	err = row.Scan(&oldInitExpr)
	if err != nil {
		oldInitExpr = useInMatchInitExpr
	} else {
		oldInitExpr = strings.TrimSpace(oldInitExpr)
	}
	var percent float64
	gstr.SimilarText(oldInitExpr, useInMatchInitExpr, &percent)
	if percent < 99.2 {
		// 先删除旧表
		(*self.db).Exec(fmt.Sprintf("DROP TABLE %s_old;", self.crud.DBName()))
		// 迁移表为_old后缀的旧表
		_, err = (*self.db).Exec(fmt.Sprintf(`ALTER TABLE '%s' RENAME TO '%s_old'; `, self.crud.DBName(), self.crud.DBName()))
		if err != nil {
			return
		}
		// 创建新表
		_, err = (*self.db).Exec(initExpr)
		if err != nil {
			return
		}
		// 将旧表数据拷贝到新表,这里获取新表字段和旧表字段，然后比较哪个短，就用哪个，然后做迁移
		tmps := regexp.MustCompile(`(?m)^\s*"(.*?)"`).FindAllStringSubmatch(oldInitExpr, -1)

		oldArgs := make([]string, 0)
		for _, tmp := range tmps {
			oldArgs = append(oldArgs, tmp[1])
		}

		newArgs := make([]string, 0)
		for _, a := range self.args {
			newArgs = append(newArgs, a.name)
		}

		newArgs = slice2.IntersectArray(oldArgs, newArgs)

		fields := ""
		for _, a := range newArgs {
			fields += fmt.Sprintf("%s,", a)
		}
		fields = strings.TrimRight(fields, ",")

		_, err = (*self.db).Exec(fmt.Sprintf(` INSERT INTO %s (%s) SELECT %s FROM %s_old;`, self.crud.DBName(), fields, fields, self.crud.DBName()))
		if err != nil {
			return
		}
		// 删除旧表
		_, err = (*self.db).Exec(fmt.Sprintf("DROP TABLE %s_old;", self.crud.DBName()))
		if err != nil {
			return
		}

	} else {
		_, err = (*self.db).Exec(initExpr)
		if err != nil {
			return
		}
	}

	if self.crud.DefaultValue() != nil {
		for _, dvs := range self.crud.DefaultValue() {
			r, err := self.Read(0, 1, dvs.CheckValue.(T))
			if err == nil && len(r) > 0 {
				continue
			}
			self.Insert(dvs.DefaultValue.(T))

		}
	}

	return
}

func getPlaceFolder(mode int) string {
	placeholder := ""
	switch mode {
	case FilterEqual:
		placeholder = "=?"
	case FilterNotEqual:
		placeholder = "!=?"
	case FilterLike:
		placeholder = "LIKE '%' || ? || '%'"
	case FilterNotLike:
		placeholder = "NOT LIKE '%' || ? || '%'"
	case FilterRegexp:
		placeholder = "REGEXP ?"
	case FilterGreater:
		placeholder = ">?"
	case FilterLess:
		placeholder = "<?"
	}
	return placeholder
}

func struct2SqlRecursion(rv reflect.Value, rt reflect.Type,
	filterModeMap map[string]int,
	sep string, filterMode bool, checkNull bool, containPrimary bool) ([]interface{}, string) {
	placeholder := "=?"

	expr := ""
	args := make([]interface{}, 0)

	for i := 0; i < rv.NumField(); i++ {
		rvf := rv.Field(i)
		// field为空并且无内嵌结构体，则跳过
		if rt.Field(i).Tag.Get(fieldTag) == "" && rt.Field(i).Tag.Get(TagEmbedStruct) == "" {
			continue
		}
		//if strings.Contains(rt.Field(i).Tag.Get(fieldTag), ",primary" ) {
		//	continue
		//}

		kind := rvf.Kind()
		switch {
		case kind == reflect.String:
			if !checkNull || !rvf.IsZero() {
				field := strings.Split(rt.Field(i).Tag.Get(fieldTag), ",")[0]
				// 判断是否有过滤模式选择，根据map使用不同过滤模式
				if filterMode {
					placeholder = getPlaceFolder(filterModeMap[field])
				}

				expr += fmt.Sprintf("%s %s %s ", field, placeholder, sep)
				value := rvf.String()
				if value == StringZeroFlag {
					value = ""
				}
				args = append(args, value)
			}
		case kind == reflect.Bool:
			if !checkNull || !rvf.IsZero() {
				field := strings.Split(rt.Field(i).Tag.Get(fieldTag), ",")[0]
				if filterMode {
					placeholder = getPlaceFolder(filterModeMap[field])
				}

				expr += fmt.Sprintf("%s %s %s ", field, placeholder, sep)
				if rvf.Bool() {
					args = append(args, 1)
				} else {
					args = append(args, -1)
				}

			}
		case kind >= reflect.Int && kind <= reflect.Uint64:
			if !checkNull || !rvf.IsZero() {
				field := rt.Field(i).Tag.Get(fieldTag)
				// 如果有primary tag，但不允许包含，那么就跳过
				if strings.Contains(field, "primary") && !containPrimary {
					continue
				}
				field = strings.Split(field, ",")[0]
				if filterMode {
					placeholder = getPlaceFolder(filterModeMap[field])
				}

				expr += fmt.Sprintf("%s %s %s ", field, placeholder, sep)
				if kind <= reflect.Int64 {
					args = append(args, rvf.Int())
				} else {
					args = append(args, rvf.Uint())
				}
			}
		case kind == reflect.Float64:
			if !checkNull || !rvf.IsZero() {
				field := rt.Field(i).Tag.Get(fieldTag)
				// 如果有primary tag，但不允许包含，那么就跳过
				if strings.Contains(field, "primary") && !containPrimary {
					continue
				}
				field = strings.Split(field, ",")[0]
				if filterMode {
					placeholder = getPlaceFolder(filterModeMap[field])
				}

				expr += fmt.Sprintf("%s %s %s ", field, placeholder, sep)
				args = append(args, rvf.Float())
			}
		case kind == reflect.Struct:
			if rt.Field(i).Type.String() == "time.Time" {
				if !checkNull || !rvf.IsZero() {
					field := strings.Split(rt.Field(i).Tag.Get(fieldTag), ",")[0]
					if filterMode {
						placeholder = getPlaceFolder(filterModeMap[field])
					}

					expr += fmt.Sprintf("%s %s %s ", field, placeholder, sep)
					args = append(args, rvf.Interface())
				}
			} else if rt.Field(i).Tag.Get(TagEmbedStruct) != "" {
				tmpArgs, tmpExpr := struct2SqlRecursion(rvf, rt.Field(i).Type, filterModeMap,
					sep, filterMode, checkNull, containPrimary)

				if tmpExpr != "" {
					expr += fmt.Sprintf("%s %s ", tmpExpr, sep)
					args = append(args, tmpArgs...)
				}

			}

		}
	}
	expr = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(expr), sep))
	return args, expr
}

// 将结构体转换成sql语句, 这里可以指定分隔符sep
/*
sep:
	" and "
	","
checkNull: 是否检查结构体值为空
containPrimary: 是否包含主键
name =? and
*/
func (self *CommonDB[T]) struct2Sql(st T, sep string, checkNull bool, containPrimary bool, filterModeMap map[string]int) (expr string, args []interface{}) {
	//switch fi := s.(type) {
	//case FilterInfo:
	//	filterMode = true
	//	s = fi.DB
	//	filterModeMap = fi.FilterModeMap
	//default:
	//
	//}
	filterMode := false
	if filterModeMap != nil {
		filterMode = true
	}

	rv := reflect.ValueOf(st)
	rt := reflect.TypeOf(st)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}

	expr = ""
	args = make([]interface{}, 0)
	args, expr = struct2SqlRecursion(rv, rt, filterModeMap, sep, filterMode, checkNull, containPrimary)

	//fmt.Println(expr)
	return
}

// where子句生成
func (self *CommonDB[T]) conditionGenerate(filterModes []map[string]int, conditions ...T) (conditionExpr string, args []interface{}) {
	conditionExpr = ""
	args = make([]interface{}, 0)
	for i, condition := range conditions {
		filterMode := make(map[string]int)
		if filterModes != nil && len(filterModes) > i {
			filterMode = filterModes[i]
		}
		subExpr, subArgs := self.struct2Sql(condition, "and", true, true, filterMode)
		if subExpr != "" {
			args = append(args, subArgs...)
			conditionExpr += fmt.Sprintf("(%s) or ", subExpr)
		}

	}

	conditionExpr = strings.TrimSpace(strings.TrimSuffix(conditionExpr, " or "))

	return
}

// 固定表达式
func (self *CommonDB[T]) initSql() string {
	fields := ""
	for _, a := range self.args {
		fields += fmt.Sprintf(`"%s" %s,
`, a.name, a.argType)
	}
	fields = strings.TrimRight(strings.TrimSpace(fields), ",")

	sqlStr := fmt.Sprintf(`CREATE TABLE if not exists "%s" (
%s
);`, self.crud.DBName(), fields)

	return sqlStr
}

// 从args里固定表达式
func (self *CommonDB[T]) insertSql(s T) (string, []interface{}) {

	expr, args := self.struct2Sql(s, ",", false, false, nil)

	expr = strings.ReplaceAll(expr, "=?", "")
	marks := strings.Repeat("?,", len(args))
	marks = strings.TrimRight(marks, ",")

	sqlStr := fmt.Sprintf(`INSERT INTO %s(%s) values(%s)`, self.crud.DBName(), expr, marks)

	return sqlStr, args
}

func (self *CommonDB[T]) appendSql(s T, keyName, valueName string) (string, []interface{}) {

	expr, args := self.struct2Sql(s, ",", false, false, nil)

	expr = strings.ReplaceAll(expr, "=?", "")
	marks := strings.Repeat("?,", len(args))
	marks = strings.TrimRight(marks, ",")

	names := strings.Split(expr, ",")
	for i, name := range names {
		if strings.TrimSpace(name) == valueName && len(args) > i {
			args = append(args, args[i])
		}
	}

	sqlStr := fmt.Sprintf(`
INSERT INTO %s(%s) 
values(%s) 
ON CONFLICT(%s) DO UPDATE SET
    %s = %s || ?
`, self.crud.DBName(), expr, marks, keyName, valueName, valueName)
	return sqlStr, args
}

type insertOnce struct {
	sqlStr string
	args   []interface{}
}

const (
	maxLineCount = 500
	maxArgsCount = 999
)

func (self *CommonDB[T]) insertMultiSql(ss []T) []insertOnce {
	expr := ""
	inserts := make([]insertOnce, 0)

	currentInsert := insertOnce{
		sqlStr: "",
		args:   make([]interface{}, 0),
	}
	currentArgsCount := 0
	currentLineCount := 0
	for _, s := range ss {
		tmpExpr, tmpArgs := self.struct2Sql(s, ",", false, false, nil)
		if expr == "" {
			expr = strings.ReplaceAll(tmpExpr, "=?", "")
		}

		currentArgsCount += len(tmpArgs)
		currentLineCount++

		tmpMarks := strings.Repeat("?,", len(tmpArgs))
		tmpMarks = strings.TrimRight(tmpMarks, ",")

		// 超过上限就重置
		if currentArgsCount > maxArgsCount || currentLineCount > maxLineCount {
			// 到上限的数据添加到全局切片中
			currentInsert.sqlStr = strings.TrimSuffix(currentInsert.sqlStr, ",\n")
			currentInsert.sqlStr = fmt.Sprintf(`INSERT INTO %s(%s) VALUES 
%s`, self.crud.DBName(), expr, currentInsert.sqlStr)
			inserts = append(inserts, currentInsert)

			// 重置
			currentArgsCount = len(tmpArgs)
			currentLineCount = 1

			currentInsert = insertOnce{
				sqlStr: fmt.Sprintf("(%s),\n", tmpMarks),
				args:   tmpArgs,
			}
		} else {
			currentInsert.sqlStr += fmt.Sprintf("(%s),\n", tmpMarks)
			currentInsert.args = append(currentInsert.args, tmpArgs...)
		}

	}
	// 最后一批数据添加
	if currentInsert.sqlStr != "" {
		currentInsert.sqlStr = strings.TrimSuffix(currentInsert.sqlStr, ",\n")
		currentInsert.sqlStr = fmt.Sprintf(`INSERT INTO %s(%s) VALUES 
%s`, self.crud.DBName(), expr, currentInsert.sqlStr)
		inserts = append(inserts, currentInsert)
	}

	return inserts
}

// 相同结构体的多个值是and关系，不同结构体之间是or关系
func (self *CommonDB[T]) readSql(orderBy string, filterModes []map[string]int, conditions ...T) (sqlStr string, args []interface{}) {
	if conditions != nil && (len(conditions) > 0) {
		conditionExpr := ""
		conditionExpr, args = self.conditionGenerate(filterModes, conditions...)
		if conditionExpr != "" {
			sqlStr = fmt.Sprintf(`select * from %s where %s %s`, self.crud.DBName(), conditionExpr, orderBy)
		} else {
			sqlStr = fmt.Sprintf(`select * from %s %s`, self.crud.DBName(), orderBy)
		}
	} else {
		sqlStr = fmt.Sprintf(`select * from %s %s`, self.crud.DBName(), orderBy)
		args = make([]interface{}, 0)
	}
	return sqlStr, args
}

func (self *CommonDB[T]) countSql(filterModes []map[string]int, conditions ...T) (sqlStr string, args []interface{}) {
	if conditions != nil && (len(conditions) > 0) {
		conditionExpr := ""
		conditionExpr, args = self.conditionGenerate(filterModes, conditions...)
		if conditionExpr != "" {
			sqlStr = fmt.Sprintf(`select count(1) from %s where %s`, self.crud.DBName(), conditionExpr)
		} else {
			sqlStr = fmt.Sprintf(`select count(1) from %s`, self.crud.DBName())
		}
	} else {
		sqlStr = fmt.Sprintf(`select count(1) from %s`, self.crud.DBName())
		args = make([]interface{}, 0)
	}
	return sqlStr, args
}

func (self *CommonDB[T]) updateSql(expr T, updateAll bool, conditions ...T) (sqlStr string, args []interface{}) {
	setExpr, subArgs1 := self.struct2Sql(expr, ",", !updateAll, false, nil)
	args = subArgs1
	if conditions != nil && (len(conditions) > 0) {
		conditionExpr, subArgs2 := self.conditionGenerate(nil, conditions...)
		args = append(args, subArgs2...)
		sqlStr = fmt.Sprintf(`update %s SET %s where %s`, self.crud.DBName(), setExpr, conditionExpr)
	} else {
		sqlStr = fmt.Sprintf(`update %s SET %s`, self.crud.DBName(), setExpr)
	}

	return sqlStr, args
}

func (self *CommonDB[T]) deleteSql(conditions ...T) (sqlStr string, args []interface{}) {
	if conditions != nil && (len(conditions) > 0) {
		conditionExpr := ""
		conditionExpr, args = self.conditionGenerate(nil, conditions...)
		if len(args) > 0 {
			sqlStr = fmt.Sprintf(`delete from %s where %s`, self.crud.DBName(), conditionExpr)
		} else {
			sqlStr = fmt.Sprintf(`delete from %s`, self.crud.DBName())
		}
	} else {
		sqlStr = fmt.Sprintf(`delete from %s`, self.crud.DBName())
		//// 这里做的比较粗糙，就先将主键默认为自增长的，如果出现全部删除情况就重置主键自增长
		//if self.primary {
		//	sqlStr = fmt.Sprintf(`%s;update sqlite_sequence SET seq = 0 where name ='%s';`, sqlStr, self.crud.DBName())
		//}
		args = make([]interface{}, 0)
	}

	return sqlStr, args
}

func (self *CommonDB[T]) equalCurrentDB(s T) bool {
	if reflect.TypeOf(s).Elem().String() == self.structName {
		return true
	}
	return false
}
