package main_control

import (
	"fmt"
	"reflect"
	"zzz_helper/internal/db/db_model"
	"zzz_helper/internal/utils/reflect2"
	"zzz_helper/modules/main_app/main_control/db_manage_extend"
	"zzz_helper/modules/module_common/common_model"
	"zzz_helper/modules/module_common/variable"
)

func (this *Control) GetClientConfigOptionsWithName(name string) (resp common_model.CommonConfigOptionsResp) {
	db := db_model.GetModuleDBFromName(name)
	if db == nil {
		resp.Err = fmt.Sprintf("[DB] %s not found", name)
		return
	}
	var err error
	resp.Options, err = variable.TransToOptions(db)
	if err != nil {
		resp.Err = err.Error()
		return
	}
	resp.Status = true
	return
}

func (this *Control) ClientDBManage(req common_model.DBManageReq) (resp common_model.DBManageResp) {
	db := db_model.GetModuleDBFromName(req.Module)
	if db == nil {
		resp.Err = "module is not found"
		return
	}
	dbFunc := db_model.GetModuleDBFuncFromName(req.Module)
	if dbFunc == nil {
		resp.Err = "module function is not found"
		return
	}
	resp.Infos = make([]common_model.DBInfo, 0)
	if req.Conditions == nil {
		req.Conditions = make([]common_model.DBInfo, 0)
	}
	var err error
	var conditions = make([]db_model.CRUD, 0)
	var dbInfo db_model.CRUD
	defer func() {
		if err != nil {
			resp.Err = err.Error()
		}
	}()
	switch req.Operation {
	case "list":
		for _, c := range req.Conditions {
			var condition db_model.CRUD
			var tmpCondition interface{}
			tmpCondition, err = reflect2.Map2Struct(c.Info, db)
			if err != nil {
				return
			}
			condition = reflect2.MustToStructPtr(tmpCondition).(db_model.CRUD)
			conditions = append(conditions, condition)
		}
		var results []db_model.CRUD
		results, err = dbFunc("").Read(-1, -1, conditions...)
		if err != nil {
			return
		}
		for _, result := range results {
			var rtInfo *reflect2.ReflectInfo
			rtInfo, err = reflect2.NewReflectInfo(result, "json")
			if err != nil {
				return
			}
			resp.Infos = append(resp.Infos, common_model.DBInfo{Info: rtInfo.Struct2Map()})
		}
		resp.Status = true
	case "create":
		if v, ok := reflect.New(reflect.TypeOf(db)).Interface().(variable.IFormHelper); ok {
			err = v.FormAutoSetting(req.Info)
			if err != nil {
				return
			}
		}
		var tmp interface{}
		tmp, err = reflect2.Map2Struct(req.Info, db)
		if err != nil {
			return
		}
		dbInfo = reflect2.MustToStructPtr(tmp).(db_model.CRUD)
		err = dbFunc("").Insert(dbInfo)
		if err != nil {
			return
		}
		resp.Status = true
	case "delete":
		for _, c := range req.Conditions {
			var condition db_model.CRUD
			var tmpCondition interface{}
			tmpCondition, err = reflect2.Map2Struct(c.Info, db)
			if err != nil {
				return
			}
			condition = reflect2.MustToStructPtr(tmpCondition).(db_model.CRUD)
			conditions = append(conditions, condition)
		}
		err = dbFunc("").Delete(conditions...)
		if err != nil {
			return
		}
		if m, ok := db_manage_extend.DBManageModuleMap[req.Module]; ok && m.Delete != nil {
			err = m.Delete(req.Conditions)
			if err != nil {
				return
			}
		}

		resp.Status = true
	case "update":
		if v, ok := reflect.New(reflect.TypeOf(db)).Interface().(variable.IFormHelper); ok {
			err = v.FormUpdateValue(req.Info)
			if err != nil {
				return
			}
		}
		var tmp interface{}
		tmp, err = reflect2.Map2Struct(req.Info, db)
		if err != nil {
			return
		}
		dbInfo = reflect2.MustToStructPtr(tmp).(db_model.CRUD)
		for _, c := range req.Conditions {
			var condition db_model.CRUD
			var tmpCondition interface{}
			tmpCondition, err = reflect2.Map2Struct(c.Info, db)
			if err != nil {
				return
			}
			condition = reflect2.MustToStructPtr(tmpCondition).(db_model.CRUD)
			conditions = append(conditions, condition)
		}
		err = dbFunc("").Update(dbInfo, true, conditions...)
		if err != nil {
			return
		}
		resp.Status = true
	case "query":
		for _, c := range req.Conditions {
			var condition db_model.CRUD
			var tmpCondition interface{}
			tmpCondition, err = reflect2.Map2Struct(c.Info, db)
			if err != nil {
				return
			}
			condition = reflect2.MustToStructPtr(tmpCondition).(db_model.CRUD)
			conditions = append(conditions, condition)
		}

		var result interface{}
		result, err = dbFunc("").ReadOne(conditions...)
		if err != nil {
			return
		}
		var info map[string]interface{}
		info, err = reflect2.Struct2Map(result)
		if err != nil {
			return
		}
		resp.Infos = append(resp.Infos, common_model.DBInfo{Info: info})
		resp.Status = true
	default:
		if m, ok := db_manage_extend.DBManageModuleMap[req.Module]; ok && m.Extend != nil {
			if fn, okk := m.Extend[req.Operation]; okk {
				resp, err = fn(req)
			}
		}
	}

	return
}
