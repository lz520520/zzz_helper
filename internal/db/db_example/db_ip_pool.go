package db_example

import "zzz_helper/internal/db/db_model"

func GetIpPoolDB() *db_model.CommonDB[*IpPoolDB] {
	c, _ := db_model.NewCacheCRUD(new(IpPoolDB))
	return c
}
func refGetIpPoolDB(name string) *db_model.CommonDB[db_model.CRUD] {
	c, _ := db_model.NewRefCacheCRUD(new(IpPoolDB))
	return c
}

type IpPoolDB struct {
	ID           int    `json:"id,primary" title:"序号"  col_width:"81"`
	IP           string `json:"ip" title:"IP"  default_value:"" form_required:"true" form_component_type:"Edit" form_config_type:"基础配置"`
	Protocol     string `json:"protocol" title:"协议"  default_value:"tcp" default_options:"tcp,udp,icmp" form_required:"true" form_component_type:"Edit" form_config_type:"基础配置"`
	Describe     string `json:"describe" title:"描述"  default_value:""  form_required:"false" form_component_type:"TextArea" form_config_type:"基础配置"`
	OnlineCount  int    `json:"online_count" title:"上线使用次数"   col_width:"130"`
	OnlineClient string `json:"online_client" title:"上线客户端" `

	CompileCount  int    `json:"compile_count" title:"编译使用次数"   col_width:"130"`
	CompileClient string `json:"compile_client" title:"编译客户端" `
}

func (self *IpPoolDB) DBName() string {
	return "ip_pool"
}

func (self *IpPoolDB) DefaultValue() []db_model.DefaultValueStruct {
	return []db_model.DefaultValueStruct{
		{
			CheckValue: &IpPoolDB{
				IP:            "1.1.1.1",
				Protocol:      "",
				Describe:      "",
				OnlineCount:   0,
				OnlineClient:  "",
				CompileCount:  0,
				CompileClient: "",
			},
			DefaultValue: &IpPoolDB{
				IP:            "1.1.1.1",
				Protocol:      "TCP",
				Describe:      "",
				OnlineCount:   0,
				OnlineClient:  "",
				CompileCount:  0,
				CompileClient: "",
			},
		},
	}

}

func init() {
	db_model.ModuleInfoRegister(db_model.ModuleInfo{
		Name:   "ip_pool",
		Alias:  "IP池",
		DB:     &IpPoolDB{},
		DBFunc: refGetIpPoolDB,
	})
}
