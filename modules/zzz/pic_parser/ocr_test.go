package pic_parser

import (
	"github.com/otiai10/gosseract/v2"
	"testing"
)

func TestOCR(t *testing.T) {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage(`E:\code\go\zzz_helper\test\2025年8月5日215933\Snipaste_2025-08-05_21-59-23_bin.png`)
	text, _ := client.Text()
	t.Logf(text)
}
