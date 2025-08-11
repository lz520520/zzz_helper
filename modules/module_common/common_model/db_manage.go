package common_model

type DBManageReq struct {
	Module     string                 `json:"module"`          // db模块，如db_portscan_cache
	Operation  string                 `json:"operation"`       // 操作，包括list/update/get/del
	Conditions []DBInfo               `json:"conditions"`      // 条件，用再del/get/update操作下
	Info       map[string]interface{} `json:"info" copier:"-"` // db信息
}

type DBManageResp struct {
	Status bool     `json:"status"`
	Msg    []byte   `json:"msg"`
	Err    string   `json:"err"`
	Infos  []DBInfo `json:"infos" copier:"-"`
}

type DBInfo struct {
	Info map[string]interface{} `json:"info"`
}
