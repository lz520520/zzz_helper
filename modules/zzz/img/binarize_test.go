package img

import (
	"regexp"
	"strings"
	"testing"
	"zzz_helper/internal/utils/file2"
)

func TestBinarizeImageWithBytes(t *testing.T) {
	b, _ := file2.ReadFileBytes(`C:\Users\Administrator\Downloads\下载.png`)
	b, _ = BinarizeImageWithBytes(b, 127)
	file2.WriteFile(`C:\Users\Administrator\Downloads\下载-bin.png`, b)

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

func TestBinarizeImageWithString(t *testing.T) {
	t.Log(nameFix("啦木鸟电音"))
}
