package lang

import "github.com/gogf/gf/v2/encoding/gcharset"

const (
	LangUTF8   = "UTF-8"
	LangGBK    = "GBK"
	LangGB2312 = "GB2312"
	LangEUCJP  = "EUC-JP"
	LangEUCKR  = "EUC-KR"
)

func LanguageCode(src, srcCharset, dstCharset string) (dst string) {
	dstStr, err := gcharset.Convert(dstCharset, srcCharset, src)
	if err != nil {
		dst = src
		return
	}

	dst = dstStr
	return
}
