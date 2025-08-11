package db_model

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"
	"reflect"
	"zzz_helper/internal/config"
	"zzz_helper/internal/mylog"
	"zzz_helper/internal/utils/sync2"
)

// https://blog.csdn.net/qianbo042311/article/details/122519180
var (
	// 全局数据库连接变量
	defaultDB = new(*sql.DB)

	cacheDB = new(*sql.DB)

	crudMap    = sync2.SyncMap{} // 用来存储已初始化的crud
	refCrudMap = sync2.SyncMap{}

	diffErr = fmt.Errorf("different databases")
)

const (
	OperationInsert = iota
	OperationDelete
	OperationUpdate
	OperationQuery
)

const (
	StringZeroFlag = "string-zero-flag"
)

/*
备注tag
field: sql查询字段，两种格式, [field/field,primary] 逗号之后可以扩展字段属性，比如primary或unique
len: sql字段的长度

json: 用于json序列化和反序列化

alias: 用于stringgrid显示字段
data_alias: 这个独立显示的字段，用于memo打印
search_only: 基于alias的扩展tag，只在历史查询时候显示
sub_section: 内嵌结构体，需要递归

overwrite: 在更新字段时，如果有此tag，表示可以覆盖。
	无tag 则只有空的时候才会覆盖
	有tag 则会比较长度，优先选择更长的
*/

const (
	FilterEqual = iota
	FilterNotEqual
	FilterLike
	FilterNotLike
	FilterGreater
	FilterLess
	FilterRegexp
)

type FilterInfo[T CRUD] struct {
	DB            T              `json:"defaultDB"`
	FilterModeMap map[string]int `json:"filter_mode_map"` // 每个字段可能匹配不同过滤模式
}

type DefaultValueStruct struct {
	CheckValue   CRUD
	DefaultValue CRUD
}

type CRUD interface {
	DBName() string
	DefaultValue() []DefaultValueStruct
	//ReadCallBack(rows *sql.Rows) (interface{}, error)
}
type DBOperation interface {
}

type arg struct {
	name    string
	argType string
}

type CommonDB[T CRUD] struct {
	db        **sql.DB
	dbPath    string
	cacheName string

	crud T // 一定要是结构体指针

	structName string // curd的结构体名

	args    []arg
	primary bool
}

func NewCRUD[T CRUD](crud T) (*CommonDB[T], error) {
	d, _, err := newCRUD(crud, defaultDB, config.DataPath, "", false)
	return d, err

}

func NewCacheCRUD[T CRUD](crud T) (*CommonDB[T], error) {
	d, _, err := newCRUD(crud, cacheDB, config.CachePath, "", false)
	return d, err
}

func NewSingleDBCRUD[T CRUD](crud T, name string) (*CommonDB[T], bool, error) {
	return newCRUD(crud, new(*sql.DB), filepath.Join(config.CacheDir, name+".db"), name, false)
}

func NewRefCRUD(crud CRUD) (*CommonDB[CRUD], error) {
	d, _, err := newCRUD(crud, defaultDB, config.DataPath, "", true)
	return d, err
}

func NewRefCacheCRUD(crud CRUD) (*CommonDB[CRUD], error) {
	d, _, err := newCRUD(crud, cacheDB, config.CachePath, "", true)
	return d, err
}

func NewRefSingleCRUD(crud CRUD, name string) (*CommonDB[CRUD], bool, error) {
	return newCRUD(crud, new(*sql.DB), filepath.Join(config.CacheDir, name+".db"), name, true)
}

func newCRUD[T CRUD](crud T, sqlDB **sql.DB, dbPath string, prefix string, ref bool) (*CommonDB[T], bool, error) {
	t := reflect.TypeOf(crud)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	cacheName := prefix + t.Name()
	if ref {
		if v, ok := refCrudMap.Load(t.Name()); ok {
			return v.(*CommonDB[T]), false, nil
		}
	} else {
		if v, ok := crudMap.Load(cacheName); ok {
			return v.(*CommonDB[T]), false, nil
		}
	}

	commonDB := &CommonDB[T]{
		crud:       crud,
		db:         sqlDB,
		dbPath:     dbPath,
		cacheName:  cacheName,
		structName: t.String(),
	}
	err := commonDB.init()
	if err != nil {
		mylog.FileLogger.Error().Msg(err.Error())
		return nil, false, err
	}
	if ref {
		refCrudMap.Store(cacheName, commonDB)
	} else {
		crudMap.Store(cacheName, commonDB)
	}
	return commonDB, true, nil
}

func (self *CommonDB[T]) GetCRUD() T {
	return self.crud
}

func (self *CommonDB[T]) Insert(s T) (err error) {
	if !self.equalCurrentDB(s) {
		return diffErr
	}
	var sqlStr string
	var args []interface{}
	sqlStr, args = self.insertSql(s)
	var stmt *sql.Stmt
	stmt, err = (*self.db).Prepare(sqlStr)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	return
}

func (self *CommonDB[T]) Append(s T, keyName, valueName string) (err error) {
	if !self.equalCurrentDB(s) {
		return diffErr
	}
	var sqlStr string
	var args []interface{}
	sqlStr, args = self.appendSql(s, keyName, valueName)
	var stmt *sql.Stmt
	stmt, err = (*self.db).Prepare(sqlStr)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	return
}

func (self *CommonDB[T]) InsertMulti(ss []T) (err error) {
	var tx *sql.Tx
	tx, err = (*self.db).Begin()
	if err != nil {
		return
	}
	defer tx.Commit()

	sqlStr, _ := self.insertSql(ss[0])
	var stmt *sql.Stmt
	stmt, err = tx.Prepare(sqlStr)
	if err != nil {
		return
	}
	defer stmt.Close()

	for _, s := range ss {
		_, args := self.insertSql(s)
		_, err = stmt.Exec(args...)
		if err != nil {
			return
		}
	}
	return
}

func (self *CommonDB[T]) Read(offset int, count int, conditions ...T) (results []T, err error) {
	return self.ReadWithOrder(offset, count, "", nil, conditions...)
}
func (self *CommonDB[T]) ReadWithFilter(offset int, count int, filterModes []map[string]int, conditions ...T) (results []T, err error) {
	return self.ReadWithOrder(offset, count, "", filterModes, conditions...)
}

func (self *CommonDB[T]) ReadWithOrder(offset int, count int, orderBy string, filterModes []map[string]int, conditions ...T) (results []T, err error) {
	results = make([]T, 0)
	var rows *sql.Rows
	var stmt *sql.Stmt

	query, args := self.readSql(orderBy, filterModes, conditions...)
	if count == 0 {
		count = -1
	}
	query += fmt.Sprintf(" LIMIT %v OFFSET %v", count, offset)

	stmt, err = (*self.db).Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err = stmt.Query(args...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		//var i interface{}
		v := reflect.New(reflect.TypeOf(self.crud).Elem())
		err = rowScan(rows, v.Interface())
		//i, err = self.crud.ReadCallBack(rows)
		if err != nil {
			mylog.FileLogger.Error().Msg(err.Error())
			return
		}

		results = append(results, v.Interface().(T))
		if count != -1 && len(results) >= count {
			break
		}
	}
	return
}

func (self *CommonDB[T]) IsExist(conditions ...T) bool {
	results, err := self.ReadWithOrder(-1, 1, "", nil, conditions...)
	if err != nil {
		return false
	}
	if len(results) == 0 {
		return false
	}
	return true
}
func (self *CommonDB[T]) Exec(query string) error {
	_, err := (*self.db).Exec(query)
	return err
}

// 判断表中unique字段是否存在
func (self *CommonDB[T]) IsUniqueExist(s T) bool {
	field, value := GetUniqueValue(s)
	if field == "" {
		return false
	}
	sPtr := reflect.New(reflect.TypeOf(s).Elem())
	err := SetDBValueForField(sPtr, field, value)
	if err != nil {
		mylog.FileLogger.Error().Msg(err.Error())
		return false
	}
	return self.IsExist(sPtr.Interface().(T))
}

func (self *CommonDB[T]) MustReadOne(conditions ...T) (result T) {
	results, err := self.ReadWithOrder(-1, 1, "", nil, conditions...)
	if err != nil {
		panic(err)
	}
	if len(results) == 0 {
		panic(errors.New("no rows"))
	}
	return results[0]
}

func (self *CommonDB[T]) ReadOne(conditions ...T) (result T, err error) {
	var results []T
	results, err = self.ReadWithOrder(-1, 1, "", nil, conditions...)
	if err != nil {
		return
	}
	if len(results) == 0 {
		err = errors.New("no rows")
		return
	}
	result = results[0]
	return
}

func (self *CommonDB[T]) UpdateOrInsert(expr T, updateAll bool, conditions ...T) (err error) {
	_, err = self.ReadOne(conditions...)
	if err == nil {
		err = self.Update(expr, updateAll, conditions...)
	} else {
		err = self.Insert(expr)
	}
	return
}

func (self *CommonDB[T]) Update(expr T, updateAll bool, conditions ...T) (err error) {
	query, args := self.updateSql(expr, updateAll, conditions...)
	var stmt *sql.Stmt
	stmt, err = (*self.db).Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(args...)
	return
}
func (self *CommonDB[T]) Count(filterModes []map[string]int, conditions ...T) (count int, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	query, args := self.countSql(filterModes, conditions...)

	stmt, err = (*self.db).Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err = stmt.Query(args...)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&count)
	}
	return
}

func (self *CommonDB[T]) Delete(conditions ...T) (err error) {
	query, args := self.deleteSql(conditions...)
	var stmt *sql.Stmt
	stmt, err = (*self.db).Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)

	if self.primary && len(conditions) == 0 {
		_, err = (*self.db).Exec(fmt.Sprintf(`update sqlite_sequence SET seq = 0 where name ='%s';`, self.crud.DBName()))
	}
	return
}
func (self *CommonDB[T]) CreateIndex(name string) (err error) {
	return self.Exec(fmt.Sprintf("CREATE INDEX idx_%s ON %s(%s);", name, self.crud.DBName(), name))
}
func (self *CommonDB[T]) DropIndex(name string) (err error) {
	return self.Exec(fmt.Sprintf("DROP INDEX IF EXISTS idx_%s;", name))
}

func (self *CommonDB[T]) Close() (err error) {
	crudMap.Delete(self.cacheName)
	refCrudMap.Delete(self.cacheName)
	if *self.db != nil {
		err = (*self.db).Close()
	}

	return
}
