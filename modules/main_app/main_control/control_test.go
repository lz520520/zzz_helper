package main_control

import (
	"testing"
	"zzz_helper/internal/db/db_zzz"
)

func TestDriverFuzz(t *testing.T) {
	db := db_zzz.GetDriverFuzzDB()
	db.Insert(&db_zzz.DriverFuzzDB{})

}
