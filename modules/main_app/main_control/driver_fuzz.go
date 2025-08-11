package main_control

import (
	"encoding/base64"
	"fmt"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"time"
	"zzz_helper/internal/config"
	"zzz_helper/internal/db/db_zzz"
	"zzz_helper/internal/mylog"
	"zzz_helper/internal/utils/file2"
	"zzz_helper/modules/main_app/zzz_models"
	"zzz_helper/modules/zzz/calc"
	"zzz_helper/modules/zzz/data"
	"zzz_helper/modules/zzz/models"
)

func (this *Control) GetProxyBuff(req zzz_models.TestProxyBuffReq) (resp zzz_models.TestProxyBuffResp) {
	if req.Proxy1.Name != "" {
		proxyInfo, err := data.AgentInfos.GetInfo(req.Proxy1.Name)
		if err != nil {
			resp.Err = err.Error()
			return
		}
		resp.Attribute.Add(proxyInfo.Buff)
		for _, star := range proxyInfo.Stars {
			if star.Level <= req.Proxy1.Star {
				resp.Attribute.Add(star.Buff)
			}
		}
	}

	if req.Proxy1.Engine != "" {
		engineInfo, err := data.EngineInfos.GetInfo(req.Proxy1.Engine, req.Proxy1.EngineStar)
		if err != nil {
			resp.Err = err.Error()
			return
		}
		resp.Attribute.Add(engineInfo.Buff)
	}

	if req.Proxy1.DriverSet != "" {
		driverSet, err := models.DriversInfos.GetInfo(req.Proxy1.DriverSet)
		if err != nil {
			resp.Err = err.Error()
			return
		}
		resp.Attribute.Add(driverSet.Piece4.Buff)
	}

	if req.Proxy2.Name != "" {
		proxyInfo, err := data.AgentInfos.GetInfo(req.Proxy2.Name)
		if err != nil {
			resp.Err = err.Error()
			return
		}
		resp.Attribute.Add(proxyInfo.Buff)
		for _, star := range proxyInfo.Stars {
			if star.Level <= req.Proxy2.Star {
				resp.Attribute.Add(star.Buff)
			}
		}
	}

	if req.Proxy2.Engine != "" {
		engineInfo, err := data.EngineInfos.GetInfo(req.Proxy2.Engine, req.Proxy2.EngineStar)
		if err != nil {
			resp.Err = err.Error()
			return
		}
		resp.Attribute.Add(engineInfo.Buff)
	}

	if req.Proxy2.DriverSet != "" {
		driverSet, err := models.DriversInfos.GetInfo(req.Proxy2.DriverSet)
		if err != nil {
			resp.Err = err.Error()
			return
		}
		resp.Attribute.Add(driverSet.Piece4.Buff)
	}

	resp.Status = true
	return
}

func (this *Control) DriverFuzz(eventID string, params models.DamageFuzzParam) (resp zzz_models.DriverFuzzResp) {
	emitWriter := NewEmitWriter(func(msg string) {
		mylog.FileLogger.Debug().Msgf("%s\n", msg)
		this.eventEmitCallBack(fmt.Sprintf("driver_fuzz_%s", eventID), msg)
	})

	info, err := data.AgentInfos.GetInfo(params.Name)
	if err != nil {
		resp.Err = err.Error()
		return
	}
	params.Attribute = info.SelfAttribute

	if params.Name == "仪玄" {
		params.AgentFeatures = models.AgentFeatures{
			LifeDestroy: true,
			Attribute2Sheer: func(attr models.AgentAttribute) float64 {
				return attr.Attack*0.3 + attr.HP*0.1
			},
		}
	}

	result, err := calc.DamageFuzz(params, emitWriter)
	if err != nil {
		resp.Err = err.Error()
		return
	}

	b, _ := yaml.Marshal(result.OutGameAttr)
	resp.OutGame = string(b)

	b, _ = yaml.Marshal(result.InGameAttr)
	resp.InGame = string(b)

	fuzzDB := db_zzz.DriverFuzzDB{
		Timestamp: time.Now().Format(time.DateTime),
	}
	b, _ = yaml.Marshal(params)
	fuzzDB.FuzzParam = string(b)

	b, _ = yaml.Marshal(result)
	fuzzDB.FuzzResult = string(b)

	for _, disk := range result.Set.Disks {
		filename := filepath.Join(config.CacheDir, "driver_"+disk.Hash())
		b, err := file2.ReadFileBytes(filename)
		if err != nil {
			resp.Err = err.Error()
			return
		}
		b64 := base64.StdEncoding.EncodeToString(b)
		switch disk.Position {
		case 1:
			fuzzDB.Disk1 = disk.Hash()
			resp.Disk1 = b64
		case 2:
			fuzzDB.Disk2 = disk.Hash()
			resp.Disk2 = b64
		case 3:
			fuzzDB.Disk3 = disk.Hash()
			resp.Disk3 = b64
		case 4:
			fuzzDB.Disk4 = disk.Hash()
			resp.Disk4 = b64
		case 5:
			fuzzDB.Disk5 = disk.Hash()
			resp.Disk5 = b64
		case 6:
			fuzzDB.Disk6 = disk.Hash()
			resp.Disk6 = b64
		}
	}
	resp.Status = true

	db_zzz.GetDriverFuzzDB().Insert(&fuzzDB)

	return
}
