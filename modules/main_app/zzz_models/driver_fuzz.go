package zzz_models

import "zzz_helper/modules/zzz/models"

type DriverFuzzResp struct {
	Err    string `json:"err"`
	Status bool   `json:"status"`

	OutGame string `json:"out_game"`
	InGame  string `json:"in_game"`

	Disk1 string `json:"disk1"`
	Disk2 string `json:"disk2"`
	Disk3 string `json:"disk3"`
	Disk4 string `json:"disk4"`
	Disk5 string `json:"disk5"`
	Disk6 string `json:"disk6"`
}

type TestProxyBuffReq struct {
	Proxy1 TestProxyInfo `json:"proxy1"`
	Proxy2 TestProxyInfo `json:"proxy2"`
}
type TestProxyInfo struct {
	Name       string `json:"name"`
	Star       int    `json:"star"`
	Engine     string `json:"engine"`
	EngineStar int    `json:"engine_star"`
	DriverSet  string `json:"driver_set"`
}

type TestProxyBuffResp struct {
	Err       string                `json:"err"`
	Status    bool                  `json:"status"`
	Attribute models.AgentAttribute `json:"attribute"`
}
