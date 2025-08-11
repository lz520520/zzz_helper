package db_zzz

import "zzz_helper/internal/db/db_model"

func GetDriverCacheDB() *db_model.CommonDB[*DriverCacheDB] {
	c, _ := db_model.NewCacheCRUD(new(DriverCacheDB))
	return c
}
func refGetDriverCacheDB(name string) *db_model.CommonDB[db_model.CRUD] {
	c, _ := db_model.NewRefCacheCRUD(new(DriverCacheDB))
	return c
}

type DriverCacheDB struct {
	ID       string `json:"id" title:"ID"  col_width:"81"`
	Name     string `json:"name" title:"名称"`
	Position int    `json:"position" title:"位置" col_width:"102"`

	Data      string `json:"data" title:"数据" col_hidden:"true"`
	Timestamp string `json:"timestamp" title:"时间戳"  col_hidden:"true"`
}

func (self *DriverCacheDB) DBName() string {
	return "driver_cache"
}

func (self *DriverCacheDB) DefaultValue() []db_model.DefaultValueStruct {
	return []db_model.DefaultValueStruct{}

}

func init() {
	db_model.ModuleInfoRegister(db_model.ModuleInfo{
		Name:   "driver_cache",
		Alias:  "驱动盘",
		DB:     &DriverCacheDB{},
		DBFunc: refGetDriverCacheDB,
	})
}
