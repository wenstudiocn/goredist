package utils

type JsonReq struct {
	Cmd    string      `json:"cmd"`
	Params interface{} `json:"params"`
}

type JsonResp struct {
	Cmd    string      `json:"cmd"`
	Ec     int32       `json:"ec"`
	Result interface{} `json:"result"`
}