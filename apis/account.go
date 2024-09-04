package apis

import (
	"encoding/json"
	"fmt"

	binance "github.com/laoliu6668/esharp_bitget_utils"
)

const gateway_binance = "https://api.bitget.com"

type TranId struct {
	TranId string `json:"transferId"`
}

// ### 现货账户向期货账户划转
// doc: https://www.bitget.fit/zh-CN/api-doc/spot/account/Wallet-Transfer
func SpotToSwapTransfer(amount float64) (id string, err error) {
	var flag = GetFlag()
	body, _, err := binance.ApiConfig.Post(Gateway, "/api/v2/spot/wallet/transfer", map[string]any{
		"coin":     "USDT",
		"amount":   amount,
		"fromType": "spot",
		"toType":   "usdt_futures",
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		fmt.Println(err)
		return
	}
	res := TranId{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", flag, err)
		fmt.Println(err)
		return
	}

	return fmt.Sprintf("%v", res.TranId), nil
}

// 期货账户向现货账户划转
// doc: https://www.bitget.fit/zh-CN/api-doc/spot/account/Wallet-Transfer
func SwapToSpotTransfer(amount float64) (id string, err error) {
	var flag = GetFlag()
	body, _, err := binance.ApiConfig.Post(Gateway, "/api/v2/spot/wallet/transfer", map[string]any{
		"coin":     "USDT",
		"amount":   amount,
		"fromType": "usdt_futures",
		"toType":   "spot",
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		fmt.Println(err)
		return
	}
	res := TranId{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", flag, err)
		fmt.Println(err)
		return
	}

	return fmt.Sprintf("%v", res.TranId), nil

}
