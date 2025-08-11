package variable

import (
	"testing"
	"zzz_helper/internal/db/db_model"
)

func TestTrans(t *testing.T) {

	t.Log(TransToOptions(db_model.GetModuleDBFromName("summary_cache")))
}
