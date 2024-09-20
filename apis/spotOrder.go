package apis

import (
	"encoding/json"
	"fmt"

	root "github.com/laoliu6668/esharp_bitget_utils"
)

type ApiResponseListData struct {
	Status  string           `json:"status"`
	Message string           `json:"err_msg"`
	Data    []map[string]any `json:"data"`
}

func (a *ApiResponseListData) Success() bool {
	return a.Status == "ok"
}

// ### 现货下单
// doc: https://www.bitget.fit/zh-CN/api-doc/spot/trade/Place-Order

type OrderRes struct {
	OrderID string `json:"orderId"`
}

func SpotBuyMarket(symb string, amount float64) (orderId string, err error) {
	// 市价买入
	var flag = GetFlag()
	body, _, err := root.ApiConfig.Post(Gateway, "/api/v2/spot/trade/place-order", map[string]any{
		"symbol":    symb + "USDT",
		"side":      "buy",
		"orderType": "market",
		"size":      amount, // 金额(USDT)
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		fmt.Println(err)
		return
	}
	// fmt.Printf("string(body): %v\n", string(body))
	res := OrderRes{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", flag, err)
		fmt.Println(err)
		return
	}

	return fmt.Sprintf("%v", res.OrderID), nil
}

func SpotSellMarket(symb string, volume float64) (data string, err error) {
	// 市价卖出
	var flag = GetFlag()
	body, _, err := root.ApiConfig.Post(Gateway, "/api/v2/spot/trade/place-order", map[string]any{
		"symbol":    symb + "USDT",
		"side":      "sell",
		"orderType": "market",
		"size":      volume, // 数量
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		fmt.Println(err)
		return
	}
	// fmt.Printf("string(body): %v\n", string(body))
	res := OrderRes{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", flag, err)
		fmt.Println(err)
		return
	}

	return fmt.Sprintf("%v", res.OrderID), nil
}
