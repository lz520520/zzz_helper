package config

import (
	"gopkg.in/yaml.v3"
	"zzz_helper/internal/utils/file2"
	"zzz_helper/modules/zzz/models"
)

var (
	ConfigPath  = "conf/config.yml"
	CurrentPath = ""

	GlobalConfig models.Config
)

func init() {
	//CurrentPath = filepath.Dir(file2.GetAbsPath(os.Args[0]))
	//ConfigPath = path.Join(CurrentPath, ConfigPath)

	content, err := file2.ReadFileBytes(ConfigPath)
	if err != nil {
		c, _ := yaml.Marshal(GlobalConfig)
		err = file2.WriteFile(ConfigPath, c)
		if err != nil {
			panic(err)
		}
	} else {
		err = yaml.Unmarshal(content, &GlobalConfig)
		if err != nil {
			panic(err)
		}
	}

}
