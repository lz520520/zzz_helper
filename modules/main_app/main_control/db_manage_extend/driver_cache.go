package db_manage_extend

import (
	"fmt"
	"os"
	"path/filepath"
	"zzz_helper/internal/config"
	"zzz_helper/modules/module_common/common_model"
)

func init() {
	DBManageModuleMap["driver_cache"] = ManageMethod{
		Delete: func(conditions []common_model.DBInfo) (err error) {
			if len(conditions) == 0 {
				err = fmt.Errorf("conditions is empty")
				return
			}
			id := conditions[0].Info["id"].(string)
			err = os.Remove(filepath.Join(config.CacheDir, "driver_"+id))
			return
		},
	}
}
