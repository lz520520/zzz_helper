package db_manage_extend

import (
	"zzz_helper/modules/module_common/common_model"
)

var (
	DBManageModuleMap = make(map[string]ManageMethod)
)

type ManageMethod struct {
	List   func() (resp common_model.DBManageResp, err error)
	Create func(info interface{}) (newInfo interface{}, err error)
	Update func(info interface{}) (newInfo interface{}, err error)
	Delete func(conditions []common_model.DBInfo) (err error)
	Extend map[string]func(req common_model.DBManageReq) (resp common_model.DBManageResp, err error)
}
