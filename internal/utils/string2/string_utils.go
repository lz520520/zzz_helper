package string2

import (
	"encoding/hex"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	whiteChar = []rune{'\n', '\r'}
)

func CompareIgnoreCase(s1, s2 string) bool {
	return strings.ToLower(s1) == strings.ToLower(s2)
}

func ContainIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func StringSplitWithoutSpace(src string, split string) []string {
	tmp := strings.Split(src, split)
	result := make([]string, 0)
	for _, t := range tmp {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		result = append(result, t)
	}
	return result
}

// 检查是否在白名单内
func isWhiteChar(r rune) bool {
	white := false

	for _, wc := range whiteChar {
		if r == wc {
			white = true
			break
		}
	}
	return white
}

func InsertChar(str, char string) string {
	// 使用 strings.Builder 创建一个可变的字符串
	var builder strings.Builder

	// 遍历原始字符串的每一位
	for i, c := range str {
		// 将字符添加到字符串构建器中
		builder.WriteRune(c)

		// 如果不是最后一位，则插入固定字符
		if i != len(str)-1 {
			builder.WriteString(char)
		}
	}

	// 返回最终的字符串
	return builder.String()
}

func StringToPrintChar(src string) (dst string) {
	rs := []rune(src)
	result := make([]rune, 0)
	for _, r := range rs {

		if !isWhiteChar(r) && !unicode.IsPrint(r) {
			runeLen := utf8.RuneLen(r)
			if runeLen == -1 {
				continue
			}
			var buf = make([]byte, runeLen)
			utf8.EncodeRune(buf, r)
			result = append(result, []rune("\\x"+hex.EncodeToString(buf))...)
		} else {
			result = append(result, r)
		}
	}
	dst = string(result)
	return
}

// FirstUpper 字符串首字母大写
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// FirstLower 字符串首字母小写
func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}
