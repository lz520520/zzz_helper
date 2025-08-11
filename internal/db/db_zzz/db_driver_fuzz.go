package db_zzz

import "zzz_helper/internal/db/db_model"

func GetDriverFuzzDB() *db_model.CommonDB[*DriverFuzzDB] {
	c, _ := db_model.NewCacheCRUD(new(DriverFuzzDB))
	return c
}
func refGetDriverFuzzDB(name string) *db_model.CommonDB[db_model.CRUD] {
	c, _ := db_model.NewRefCacheCRUD(new(DriverFuzzDB))
	return c
}

type DriverFuzzDB struct {
	ID         int    `json:"id,primary" title:"ID" col_width:"81"`
	FuzzParam  string `json:"fuzz_param" title:"Fuzz参数"`
	FuzzResult string `json:"fuzz_result" title:"Fuzz结果"`

	Disk1 string `json:"disk1"`
	Disk2 string `json:"disk2"`
	Disk3 string `json:"disk3"`
	Disk4 string `json:"disk4"`
	Disk5 string `json:"disk5"`
	Disk6 string `json:"disk6"`

	Timestamp string `json:"timestamp" title:"时间戳"`
}

func (self *DriverFuzzDB) DBName() string {
	return "driver_fuzz"
}

func (self *DriverFuzzDB) DefaultValue() []db_model.DefaultValueStruct {
	return []db_model.DefaultValueStruct{}

}

func init() {
	db_model.ModuleInfoRegister(db_model.ModuleInfo{
		Name:   "driver_fuzz",
		Alias:  "驱动盘计算",
		DB:     &DriverFuzzDB{},
		DBFunc: refGetDriverFuzzDB,
	})
}
