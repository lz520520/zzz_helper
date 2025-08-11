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
	"zzz_helper/internal/mylog"
	"zzz_helper/internal/utils/file2"
	"zzz_helper/internal/utils/string2"
	"zzz_helper/modules/zzz/data"
	"zzz_helper/modules/zzz/img"
	"zzz_helper/modules/zzz/models"
	pic_parser2 "zzz_helper/modules/zzz/pic_parser"
)

var (
	filename string
	output   string
)

func GetSubAttribute(disk *models.DriverDiskStat, attrName string, add int, value string) {
	switch true {
	case strings.Contains(attrName, "暴击率") || strings.Contains(attrName, "击率"):
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

func nameFix(name string) string {
	if regexp.MustCompile(`云.*如我`).MatchString(name) {
		name = "云岿如我"
	} else if regexp.MustCompile(`河\w电音`).MatchString(name) {
		name = "河豚电音"
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
	position, err := strconv.Atoi(src[positionStart+1 : positionStart+positionEnd])
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
	pattern := regexp.MustCompile(`^(\p{Han}+)\s*\+[\s|\.]*(\d)\s+([\d|\.]+\%*)$`)
	pattern2 := regexp.MustCompile(`^(\p{Han}+)\s+([\d|\.]+\%*)$`)
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
		GetSubAttribute(disk, name, add, value)
		count += add
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
func parserDriver(client *pic_parser2.Client, path string) error {
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

func parserDriverWithTesseract(path string) error {
	if !isPicture(path) {
		return fmt.Errorf("not found picture")
	}
	var disk *models.DriverDiskStat
	for i := 0; i < 2; i++ {
		var result string
		result, err := pic_parser2.ParseWithTesseract(path)
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
		//client, err := pic_parser.NewClient(config.GlobalConfig.AliyunAuth.AK, config.GlobalConfig.AliyunAuth.SK)
		//if err != nil {
		//    return err
		//}
		if isPicture(filename) {
			err := parserDriverWithTesseract(filename)
			if err != nil {
				mylog.CommonLogger.Error().Msgf("%s: %s", filename, err.Error())
			}
			return nil
		} else {
			filepath.Walk(filename, func(path string, info os.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}
				err = parserDriverWithTesseract(path)
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
