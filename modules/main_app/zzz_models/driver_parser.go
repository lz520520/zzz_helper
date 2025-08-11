package zzz_models

type FileInfo struct {
	ID   string `json:"id"`
	Data []byte `json:"data"`
}

type DriverParserResp struct {
	Status bool     `json:"status"`
	IDs    []string `json:"ids"`
	Err    string   `json:"err"`
}
