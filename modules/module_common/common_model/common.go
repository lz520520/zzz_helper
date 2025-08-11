package common_model

type CommonReq struct {
	Msg string `json:"msg"`
}
type CommonResp struct {
	Status bool   `json:"status"`
	Msg    string `json:"msg"`
	Err    string `json:"err"`
}

type CommonBytesResp struct {
	Status bool   `json:"status"`
	Bytes  []byte `json:"bytes"`
	Err    string `json:"err"`
}

type LanguageEncodeReq struct {
	Data       string `json:"data"`
	SrcCharset string `json:"src_charset"`
	DstCharset string `json:"dst_charset"`
}

type LanguageEncodeResp struct {
	Status bool   `json:"status"`
	Msg    string `json:"msg"`
	Err    string `json:"err"`
	Data   string `json:"data"`
}
