package config

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"zzz_helper/internal/utils/file2"
)

var (
	DataPath    = "db/data.db"
	CachePath   = "db/cache.db"
	CurrentPath = ""
	CacheDir    = "cache"
	TmpDir      = "tmp"
)

func init() {
	CurrentPath = filepath.Dir(file2.GetAbsPath(os.Args[0]))
	if strings.Contains(CurrentPath, "/Contents/MacOS") {
		CurrentPath = filepath.Join(CurrentPath, "../../../")
	}

	CacheDir = path.Join(CurrentPath, CacheDir)
	DataPath = path.Join(CurrentPath, DataPath)
	TmpDir = path.Join(CurrentPath, TmpDir)

}
