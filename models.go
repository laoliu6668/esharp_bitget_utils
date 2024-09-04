package bitget

import "encoding/json"

type ApiConfigModel struct {
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key"`
	PassPhrase string `json:"pass_phrase"`
}

type SpotBalanceTicker struct {
	Symbol string      `json:"symbol"`
	Trade  json.Number `json:"trade"`
	Frozen json.Number `json:"frozen"`
}
