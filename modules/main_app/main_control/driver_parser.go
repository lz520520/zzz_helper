package main_control

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"zzz_helper/internal/config"
	"zzz_helper/internal/db/db_zzz"
	"zzz_helper/internal/utils/file2"
	"zzz_helper/internal/utils/string2"
	"zzz_helper/modules/main_app/zzz_models"
	"zzz_helper/modules/module_common/common_model"
	"zzz_helper/modules/zzz/data"
	"zzz_helper/modules/zzz/img"
	"zzz_helper/modules/zzz/models"
	pic_parser2 "zzz_helper/modules/zzz/pic_parser"
)

func NewEmitWriter(callback func(msg string)) *EmitWriter {
	return &EmitWriter{callback: callback}
}

type EmitWriter struct {
	callback func(msg string)
}

func (this *EmitWriter) WriteString(s string) (n int, err error) {
	this.callback(fmt.Sprintf("[%s] %s", time.Now().Format(time.DateTime), s))
	return len(s), nil
}

type FileInfo struct {
	ID   string `json:"id"`
	Data []byte `json:"data"`
}

func (this *Control) DriverParser(eventID string, files []zzz_models.FileInfo, ocr string) (resp zzz_models.DriverParserResp) {
	emitWriter := NewEmitWriter(func(msg string) {
		this.eventEmitCallBack(fmt.Sprintf("driver_parser_%s", eventID), msg)
	})
	cacheDB := db_zzz.GetDriverCacheDB()
	resp.IDs = make([]string, 0)

	if ocr != "" && len(files) > 0 {
		disk, err := ParseDriveInfo(ocr)
		if err != nil {
			emitWriter.WriteString(fmt.Sprintf("[-] %s", err.Error()))
			return
		}
		dataBytes, _ := yaml.Marshal(disk)
		id := disk.Hash()
		if cacheDB.IsExist(&db_zzz.DriverCacheDB{ID: id}) {
			emitWriter.WriteString(fmt.Sprintf("[!] disk already exist"))
		} else {
			err = cacheDB.Insert(&db_zzz.DriverCacheDB{
				ID:        id,
				Name:      disk.Name,
				Position:  disk.Position,
				Data:      string(dataBytes),
				Timestamp: time.Now().Format(time.DateTime),
			})
			if err != nil {
				emitWriter.WriteString(fmt.Sprintf("[-] %s", err.Error()))
				return
			}
			file2.WriteFile(filepath.Join(config.CacheDir, "driver_"+id), files[0].Data)
		}
		emitWriter.WriteString(fmt.Sprintf("[+] 解析结果: \n%s", string(dataBytes)))

		resp.IDs = append(resp.IDs, files[0].ID)
		resp.Status = true
		return
	}

	for _, file := range files {
		tmpPath := filepath.Join(config.TmpDir, file.ID)
		file2.WriteFile(tmpPath, file.Data)
		disk, err := parserDriverWithTesseract(tmpPath, emitWriter)
		if err != nil {
			emitWriter.WriteString(fmt.Sprintf("[-] %s", err.Error()))

			emitWriter.WriteString("[*] 尝试二值化")
			tmpData, err := img.BinarizeImageWithBytes(file.Data, 127)
			if err != nil {
				os.Remove(tmpPath)
				emitWriter.WriteString(fmt.Sprintf("[-] 二值化：%s", err.Error()))
				continue
			}
			file2.WriteFile(tmpPath, tmpData)
			disk, err = parserDriverWithTesseract(tmpPath, emitWriter)
			if err != nil {
				os.Remove(tmpPath)
				emitWriter.WriteString(fmt.Sprintf("[-] 二值化：%s", err.Error()))
				continue
			}
			os.Remove(tmpPath)
		}
		dataBytes, _ := yaml.Marshal(disk)
		id := disk.Hash()

		if cacheDB.IsExist(&db_zzz.DriverCacheDB{ID: id}) {
			emitWriter.WriteString(fmt.Sprintf("[!] disk already exist"))
		} else {
			err = cacheDB.Insert(&db_zzz.DriverCacheDB{
				ID:       id,
				Name:     disk.Name,
				Position: disk.Position,
				Data:     string(dataBytes),
			})
			if err != nil {
				emitWriter.WriteString(fmt.Sprintf("[-] %s", err.Error()))
				continue
			}
			file2.WriteFile(filepath.Join(config.CacheDir, "driver_"+id), file.Data)
		}
		emitWriter.WriteString(fmt.Sprintf("[+] 解析结果: \n%s", string(dataBytes)))
		resp.IDs = append(resp.IDs, file.ID)
	}
	resp.Status = true
	return
}

func (this *Control) ReadDriverCache(id string) (resp common_model.CommonBytesResp) {
	filename := filepath.Join(config.CacheDir, "driver_"+id)
	b, err := file2.ReadFileBytes(filename)
	if err != nil {
		resp.Err = err.Error()
		return
	}
	resp.Bytes = b
	resp.Status = true
	return
}

func nameFix(name string) string {
	if regexp.MustCompile(`云\p{Han}如我`).MatchString(name) {
		name = "云岿如我"
	} else if regexp.MustCompile(`河\p{Han}电音`).MatchString(name) {
		name = "河豚电音"
	} else if regexp.MustCompile(`\p{Han}木鸟电音`).MatchString(name) {
		name = "啄木鸟电音"
	}
	return strings.TrimSpace(name)
}
func ParseDriveInfo(src string) (*models.DriverDiskStat, error) {
	lines := string2.StringSplitWithoutSpace(src, "\n")
	if len(lines) < 8 {
		return nil, fmt.Errorf("lines is error")
	}
	disk := &models.DriverDiskStat{}

	// 驱动盘名字
	positionStart := strings.Index(lines[0], "[")
	if positionStart == -1 {
		return nil, fmt.Errorf("name not found position")
	}
	disk.Name = nameFix(lines[0][:positionStart])
	positionEnd := strings.Index(lines[0][positionStart:], "]")
	if positionEnd == -1 {
		positionEnd = strings.Index(lines[0][positionStart:], " ")
		if positionEnd == -1 {
			return nil, fmt.Errorf("not found position")
		}
	}
	position, err := strconv.Atoi(strings.TrimSpace(src[positionStart+1 : positionStart+positionEnd]))
	if err != nil {
		return nil, err
	}
	if position < 1 || position > 6 {
		return nil, fmt.Errorf("invalid position")
	}
	disk.Position = position

	// 主属性
	mainAttr := lines[3]
	switch true {
	case strings.Contains(mainAttr, "攻击力"):
		if disk.Position == 2 {
			disk.Main.Attack = data.BaseDriverDiskMainStat.Attack
		} else {
			disk.Main.AttackBonus = data.BaseDriverDiskMainStat.AttackBonus
		}
	case strings.Contains(mainAttr, "防御力"):
		if disk.Position == 3 {
			disk.Main.Defense = data.BaseDriverDiskMainStat.Defense
		} else {
			disk.Main.DefenseBonus = data.BaseDriverDiskMainStat.DefenseBonus
		}
	case strings.Contains(mainAttr, "生命值"):
		if disk.Position == 1 {
			disk.Main.HP = data.BaseDriverDiskMainStat.HP
		} else {
			disk.Main.HPBonus = data.BaseDriverDiskMainStat.HPBonus
		}
	case strings.Contains(mainAttr, "暴击率") || strings.Contains(mainAttr, "击率"):
		disk.Main.CriticalRate = data.BaseDriverDiskMainStat.CriticalRate
	case strings.Contains(mainAttr, "暴击伤害"):
		disk.Main.CriticalDamage = data.BaseDriverDiskMainStat.CriticalDamage
	case strings.Contains(mainAttr, "穿透率"):
		disk.Main.PenetrationRadio = data.BaseDriverDiskMainStat.PenetrationRadio
	case strings.Contains(mainAttr, "电属性伤害加成"):
		disk.Main.ElectricDamageBonus = data.BaseDriverDiskMainStat.ElectricDamageBonus
	case strings.Contains(mainAttr, "火属性伤害加成"):
		disk.Main.FireDamageBonus = data.BaseDriverDiskMainStat.FireDamageBonus
	case strings.Contains(mainAttr, "冰属性伤害加成"):
		disk.Main.IceDamageBonus = data.BaseDriverDiskMainStat.IceDamageBonus
	case strings.Contains(mainAttr, "以太伤害加成"):
		disk.Main.EtherDamageBonus = data.BaseDriverDiskMainStat.EtherDamageBonus
	case strings.Contains(mainAttr, "物理伤害加成"):
		disk.Main.PhysicalDamageBonus = data.BaseDriverDiskMainStat.PhysicalDamageBonus
	case strings.Contains(mainAttr, "异常掌控"):
		disk.Main.AnomalyMastery = data.BaseDriverDiskMainStat.AnomalyMastery
	default:
		return nil, fmt.Errorf("not found main attribute")
	}
	pattern := regexp.MustCompile(`^(\p{Han}+)\s*\+[\s|\.]*(\d)\s+([\d|\.]+\%*)`)
	pattern2 := regexp.MustCompile(`^(\p{Han}+)\s+([\d|\.]+\%*)`)
	count := 0
	for _, line := range lines[5:] {
		name := ""
		add := 0
		value := ""
		if strings.Contains(line, "+") {
			groups := pattern.FindStringSubmatch(line)
			if groups == nil {
				return nil, fmt.Errorf("sub attribute parse error: %s", line)
			}
			add, err = strconv.Atoi(groups[2])
			if err != nil {
				return nil, fmt.Errorf("%s parser error: %s", groups[2], line)
			}
			name = groups[1]
			value = groups[3]

		} else {
			groups := pattern2.FindStringSubmatch(line)
			if groups == nil {
				return nil, fmt.Errorf("sub attribute parse error: %s", line)
			}
			name = groups[1]
			value = groups[2]
		}
		err = GetSubAttribute(disk, name, add, value)
		if err != nil {
			return nil, err
		}
		count += add
	}

	if count < 4 || count > 5 {
		return nil, fmt.Errorf("parser sub attribute count %v", count)
	}

	return disk, nil

}

func GetSubAttribute(disk *models.DriverDiskStat, attrName string, add int, value string) error {
	switch true {
	case strings.Contains(attrName, "暴击率") || strings.Contains(attrName, "击率"):
		disk.Sub.CriticalRate = data.BaseDriverDiskSubStat.CriticalRate * float64(add+1)
	case strings.Contains(attrName, "暴击伤害") || strings.Contains(attrName, "暴击伤"):
		disk.Sub.CriticalDamage = data.BaseDriverDiskSubStat.CriticalDamage * float64(add+1)
	case strings.Contains(attrName, "攻击力"):
		if strings.HasSuffix(value, "%") {
			disk.Sub.AttackBonus = data.BaseDriverDiskSubStat.AttackBonus * float64(add+1)
		} else {
			disk.Sub.Attack = data.BaseDriverDiskSubStat.Attack * float64(add+1)
		}
	case strings.Contains(attrName, "穿透值"):
		disk.Sub.Penetration = data.BaseDriverDiskSubStat.Penetration * float64(add+1)
	case strings.Contains(attrName, "异常精通"):
		disk.Sub.AnomalyProficiency = data.BaseDriverDiskSubStat.AnomalyProficiency * float64(add+1)
	case strings.Contains(attrName, "防御力"):
		if strings.HasSuffix(value, "%") {
			disk.Sub.DefenseBonus = data.BaseDriverDiskSubStat.DefenseBonus * float64(add+1)
		} else {
			disk.Sub.Defense = data.BaseDriverDiskSubStat.Defense * float64(add+1)
		}
	case strings.Contains(attrName, "生命值"):
		if strings.HasSuffix(value, "%") {
			disk.Sub.HPBonus = data.BaseDriverDiskSubStat.HPBonus * float64(add+1)
		} else {
			disk.Sub.HP = data.BaseDriverDiskSubStat.HP * float64(add+1)
		}
	default:
		return fmt.Errorf("sub attribute parse error: %s", attrName)
	}
	return nil

}

func parserDriverWithTesseract(path string, writer io.StringWriter) (*models.DriverDiskStat, error) {
	var disk *models.DriverDiskStat
	var result string
	result, err := pic_parser2.ParseWithTesseract(path)
	if err != nil {
		return nil, err
	}
	writer.WriteString(fmt.Sprintf("[*] OCR识别 %s: \n%v", path, result))

	disk, err = ParseDriveInfo(result)
	if err != nil {
		return nil, err
	}
	return disk, nil
}
