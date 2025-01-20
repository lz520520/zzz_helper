package sub

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"zzz_helper/internal/config"
	"zzz_helper/internal/data"
	"zzz_helper/internal/img"
	"zzz_helper/internal/models"
	"zzz_helper/internal/mylog"
	"zzz_helper/internal/pic_parser"
	"zzz_helper/internal/utils/file2"
	"zzz_helper/internal/utils/string2"
)

var (
	filename string
	output   string
)

func GetSubAttribute(disk *models.DriverDiskStat, attrName string, add int, value string) {
	switch true {
	case strings.Contains(attrName, "暴击率"):
		disk.Sub.CriticalRate = data.BaseDriverDiskSubStat.CriticalRate * float64(add+1)
	case strings.Contains(attrName, "暴击伤害"):
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
	}

}
func ParseDriveInfo(src string) (*models.DriverDiskStat, error) {
	disk := &models.DriverDiskStat{}
	src = strings.TrimSpace(src)
	offset := 0
	positionStart := strings.Index(src, "[")
	if positionStart == -1 {
		return nil, fmt.Errorf("not found position")
	}
	disk.Name = src[:positionStart]

	positionEnd := strings.Index(src[positionStart:], "]")
	if positionEnd == -1 {
		return nil, fmt.Errorf("not found position")
	}
	position, err := strconv.Atoi(src[positionStart+1 : positionStart+positionEnd])
	if err != nil {
		return nil, fmt.Errorf("not found position")
	}
	offset += positionStart + positionEnd
	disk.Position = position

	mainAttrStart := strings.Index(src[offset:], "主属性")
	if mainAttrStart == -1 {
		return nil, fmt.Errorf("not found main attribute")
	}
	offset += mainAttrStart
	tmp := string2.StringSplitWithoutSpace(src[offset:], " ")
	mainAttr := tmp[1]

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
	case strings.Contains(mainAttr, "暴击率"):
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
	case strings.Contains(mainAttr, "以太属性伤害加成"):
		disk.Main.EtherDamageBonus = data.BaseDriverDiskMainStat.EtherDamageBonus
	case strings.Contains(mainAttr, "物理伤害加成"):
		disk.Main.PhysicalDamageBonus = data.BaseDriverDiskMainStat.PhysicalDamageBonus
	case strings.Contains(mainAttr, "异常掌控"):
		disk.Main.AnomalyMastery = data.BaseDriverDiskMainStat.AnomalyMastery
	default:
		return nil, fmt.Errorf("not found main attribute")
	}
	if tmp[3] != "副属性" {
		return nil, fmt.Errorf("not found sub attribute")
	}
	tmp = tmp[4:]
	newTmp := make([]string, 0)
	pattern := regexp.MustCompile(`^\p{Han}\+$`)
	for _, s := range tmp {
		if pattern.MatchString(s) {
			continue
		}
		newTmp = append(newTmp, s)
	}
	tmp = newTmp
	count := 0
	digitalPattern := regexp.MustCompile(`\+\d+`)
	for i := 0; i < len(tmp); i++ {
		name := tmp[i]
		add := 0
		value := tmp[i+1]

		if digitalPattern.MatchString(name) {
			names := strings.SplitN(name, "+", 2)
			name = strings.TrimSpace(names[0])
			add, err = strconv.Atoi(names[1])
			if err != nil {
				return nil, fmt.Errorf("%s parser error", names[1])
			}
			i += 1
		} else if strings.Contains(value, "+") {
			index := strings.Index(value, "+")
			add, err = strconv.Atoi(value[index:])
			if err != nil {
				return nil, fmt.Errorf("%s parser error", value[index:])
			}
			value = tmp[i+2]
			i += 2
		} else {
			i += 1
		}
		count += add

		GetSubAttribute(disk, name, add, value)
	}
	if count < 4 || count > 5 {
		return nil, fmt.Errorf("parser sub attribute count %v", count)
	}

	return disk, nil

}
func isPicture(path string) bool {
	if strings.HasSuffix(path, ".png") ||
		strings.HasSuffix(path, ".jpg") ||
		strings.HasSuffix(path, ".jpeg") {
		return true
	}
	return false
}
func parserDriver(client *pic_parser.Client, path string) error {
	if !isPicture(path) {
		return fmt.Errorf("not found picture")
	}
	png, err := file2.ReadFileBytes(path)
	if err != nil {
		return err
	}
	var disk *models.DriverDiskStat
	for i := 0; i < 2; i++ {
		var result string
		var png2 []byte
		if i == 0 {
			png2, err = img.BinarizeImageWithBytes(png, 127)
			if err != nil {
				return err
			}
		} else if i == 1 {
			png2 = png
		}
		result, err = client.Parse(png2)
		if err != nil {
			return err
		}
		mylog.CommonLogger.Info().Msgf("%s: %v", path, result)

		disk, err = ParseDriveInfo(result)
		if err != nil {
			mylog.CommonLogger.Err(err).Send()
			continue
		}
		break
	}
	if err != nil {
		return err
	}

	m := map[string]interface{}{
		uuid.New().String(): disk,
	}
	r, _ := yaml.Marshal(m)
	file2.AppendFile(output, r)
	return nil
}

var ocrCmd = &cobra.Command{
	Use:   "ocr",
	Short: "use ocr generate driver info",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := pic_parser.NewClient(config.GlobalConfig.AliyunAuth.AK, config.GlobalConfig.AliyunAuth.SK)
		if err != nil {
			return err
		}
		if isPicture(filename) {
			err = parserDriver(client, filename)
			if err != nil {
				mylog.CommonLogger.Error().Msgf("%s: %s", filename, err.Error())
			}
			return nil
		} else {
			filepath.Walk(filename, func(path string, info os.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}
				err = parserDriver(client, path)
				if err != nil {
					mylog.CommonLogger.Error().Msgf("%s: %s", path, err.Error())
				}
				return nil
			})
		}

		return nil
	},
}

func init() {
	ocrCmd.Flags().StringVarP(&filename, "filename", "f", "", "filename")
	ocrCmd.Flags().StringVarP(&output, "output", "o", "out/driver.yml", "output path")

	ocrCmd.MarkFlagRequired("filename")
	RootCmd.AddCommand(ocrCmd)
}
