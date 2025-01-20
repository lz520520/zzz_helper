package lang

import "testing"

func TestLang(t *testing.T) {
	t.Log(LanguageCode("111", LangUTF8, LangGBK))
}
