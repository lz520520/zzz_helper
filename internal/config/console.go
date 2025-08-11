package config

import (
	"gopkg.in/yaml.v3"
	"path/filepath"
	"zzz_helper/internal/utils/file2"
)

var (
	defaultConsoleConfig = ConsoleConfig{
		Ws:       false,
		LogLevel: "info",
		APIMode:  false,
		TLSMode:  true,

		TesseractPath: "Tesseract-OCR/tesseract.exe",
	}
	ConsoleConfigInst = &defaultConsoleConfig
)

type ConsoleConfig struct {
	Ws       bool   `yaml:"ws,omitempty"`
	LogLevel string `yaml:"log_level,omitempty"`
	TLSMode  bool   `yaml:"tls_mode,omitempty"`
	APIMode  bool   `yaml:"api_mode,omitempty"`

	TesseractPath  string     `yaml:"tesseract_path"`
	AliyunAuth     AliyunAuth `yaml:"aliyun_auth"`
	MiyousheCookie string     `yaml:"miyoushe_cookie"`
	MihoyoCookie   string     `yaml:"mihoyo_cookie"`
}

type AliyunAuth struct {
	AK string `yaml:"access_key"`
	SK string `yaml:"secret_key"`
}

func GetConsoleConfig() *ConsoleConfig {
	conPath := filepath.Join(CurrentPath, "conf/console.yml")
	b, err := file2.ReadFileBytes(conPath)
	if err != nil {
		b, err = yaml.Marshal(&defaultConsoleConfig)
		file2.WriteFile(conPath, b)
		return &defaultConsoleConfig
	}

	conf := &ConsoleConfig{}
	err = yaml.Unmarshal(b, conf)
	if err != nil {
		b, err = yaml.Marshal(&defaultConsoleConfig)
		file2.WriteFile(conPath, b)
		return &defaultConsoleConfig
	}
	return conf
}

func init() {
	ConsoleConfigInst = GetConsoleConfig()
}
