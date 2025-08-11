package common_model

type DynamicFormReq struct {
	ChangeKey string `json:"change_key"`
	TrickKey  string `json:"trick_key"`
}

type DynamicFormResp struct {
	Status bool     `json:"status"`
	Err    string   `json:"err"`
	Values []string `json:"values"`
}
