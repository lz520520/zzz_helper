package models

type Config struct {
	AliyunAuth     AliyunAuth `yaml:"aliyun_auth"`
	MiyousheCookie string     `yaml:"miyoushe_cookie"`
	MihoyoCookie   string     `yaml:"mihoyo_cookie"`
}
type AliyunAuth struct {
	AK string `yaml:"access_key"`
	SK string `yaml:"secret_key"`
}
