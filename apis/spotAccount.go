package apis

import (
	"encoding/json"
	"fmt"

	binance "github.com/laoliu6668/esharp_bitget_utils"
)

// # MODEL 获取用户账户
type AccountData struct {
	Coin      string `json:"coin"`
	Available string `json:"available"`
}

// 获取现货账户信息
// doc:  https://www.bitget.fit/zh-CN/api-doc/spot/account/Get-Account-Assets
func GetSpotAccount() (data []AccountData, err error) {
	var flag = GetFlag()
	body, _, err := binance.ApiConfig.Get(Gateway, "/api/v2/spot/account/assets", nil)
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		fmt.Println(err)
		return
	}
	// util.WriteTestJsonFile(flag, body)
	res := []AccountData{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", flag, err)
		fmt.Println(err)
		return
	}
	return res, nil
}
