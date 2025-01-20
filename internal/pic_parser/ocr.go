package pic_parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ocr_api20210707 "github.com/alibabacloud-go/ocr-api-20210707/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"strings"
)

func NewClient(ak string, sk string) (*Client, error) {
	config := &openapi.Config{
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
		AccessKeyId: tea.String(ak),
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
		AccessKeySecret: tea.String(sk),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/ocr-api
	config.Endpoint = tea.String("ocr-api.cn-hangzhou.aliyuncs.com")
	ocrClient, err := ocr_api20210707.NewClient(config)
	if err != nil {
		return nil, err
	}
	if ocrClient.Credential == nil {
		return nil, fmt.Errorf("credential is nil")
	}

	return &Client{ocrClient: ocrClient}, nil
}

type Client struct {
	ocrClient *ocr_api20210707.Client
}

func (this *Client) Parse(src []byte) (res string, err error) {
	body := bytes.NewBuffer(src)
	recognizeAdvancedRequest := &ocr_api20210707.RecognizeAdvancedRequest{
		Body: body,
	}
	resp, err := this.ocrClient.RecognizeAdvanced(recognizeAdvancedRequest)
	if err != nil {
		return "", err
	}
	if resp.Body == nil {
		return "", fmt.Errorf("body is nil")
	}

	var data interface{}
	d := json.NewDecoder(strings.NewReader(tea.StringValue(resp.Body.Data)))
	err = d.Decode(&data)
	if err != nil {
		return "", err
	}
	if m, ok := data.(map[string]interface{}); ok {
		if content, okk := m["content"]; okk {
			return content.(string), nil
		}
	}

	return "", fmt.Errorf("not found content")
}
