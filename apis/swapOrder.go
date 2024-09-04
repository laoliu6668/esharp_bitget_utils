package apis

import (
	"encoding/json"
	"fmt"
	"strings"

	root "github.com/laoliu6668/esharp_bitget_utils"
)

var SwapOrderUrl = "/api/v2/mix/order/place-order"

// 期货卖出开空
// doc: https://www.bitget.fit/zh-CN/api-doc/contract/trade/Place-Order
func SwapSellOpen(symb string, volume float64) (orderId string, err error) {
	body, _, err := root.ApiConfig.Post(Gateway, SwapOrderUrl, map[string]any{
		"symbol":      fmt.Sprintf("%susdt", strings.ToLower(symb)),
		"productType": "USDT-FUTURES",
		"marginCoin":  "USDT",
		"marginMode":  "isolated",
		"orderType":   "market",
		"side":        "sell",
		"tradeSide":   "open",
		"size":        volume,
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", GetFlag(), err)
		fmt.Println(err)
		return
	}
	res := OrderRes{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", GetFlag(), err)
		fmt.Println(err)
		return
	}

	return fmt.Sprintf("%v", res.OrderID), nil
}

// 期货买入平空
// doc: https://www.bitget.fit/zh-CN/api-doc/contract/trade/Place-Order
func SwapBuyClose(symb string, volume float64) (orderId string, err error) {
	body, _, err := root.ApiConfig.Post(Gateway, SwapOrderUrl, map[string]any{
		"symbol":      fmt.Sprintf("%susdt", strings.ToLower(symb)),
		"productType": "USDT-FUTURES",
		"marginCoin":  "USDT",
		"marginMode":  "isolated",
		"orderType":   "market",
		"side":        "sell",
		"tradeSide":   "close",
		"size":        volume,
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", GetFlag(), err)
		fmt.Println(err)
		return
	}
	res := OrderRes{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", GetFlag(), err)
		fmt.Println(err)
		return
	}
	return fmt.Sprintf("%v", res.OrderID), nil
}

// 期货买入开多
// doc: https://www.bitget.fit/zh-CN/api-doc/contract/trade/Place-Order
func SwapBuyOpen(symb string, volume float64) (orderId string, err error) {
	body, _, err := root.ApiConfig.Post(Gateway, SwapOrderUrl, map[string]any{
		"symbol":      fmt.Sprintf("%susdt", strings.ToLower(symb)),
		"productType": "USDT-FUTURES",
		"marginCoin":  "USDT",
		"marginMode":  "isolated",
		"orderType":   "market",
		"side":        "buy",
		"tradeSide":   "open",
		"size":        volume,
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", GetFlag(), err)
		fmt.Println(err)
		return
	}
	// fmt.Printf("string(body): %v\n", string(body))
	res := OrderRes{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", GetFlag(), err)
		fmt.Println(err)
		return
	}
	return fmt.Sprintf("%v", res.OrderID), nil
}

// 期货卖出平多
// doc: https://www.bitget.fit/zh-CN/api-doc/contract/trade/Place-Order
func SwapSellClose(symb string, volume float64) (orderId string, err error) {
	body, _, err := root.ApiConfig.Post(Gateway, SwapOrderUrl, map[string]any{
		"symbol":      fmt.Sprintf("%susdt", strings.ToLower(symb)),
		"productType": "USDT-FUTURES",
		"marginCoin":  "USDT",
		"marginMode":  "isolated",
		"orderType":   "market",
		"side":        "buy",
		"tradeSide":   "close",
		"size":        volume,
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", GetFlag(), err)
		fmt.Println(err)
		return
	}
	res := OrderRes{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", GetFlag(), err)
		fmt.Println(err)
		return
	}
	return fmt.Sprintf("%v", res.OrderID), nil
}

const SetMarginUrl = "/api/v2/mix/account/set-margin"

// 增加空头逐仓保证金
// doc: https://www.bitget.fit/zh-CN/api-doc/contract/account/Change-Margin
func SwapIncShortPositionMargin(symb string, amount float64) (err error) {
	_, _, err = root.ApiConfig.Post(Gateway, SetMarginUrl, map[string]any{
		"productType": "USDT-FUTURES",
		"symbol":      fmt.Sprintf("%sUSDT", strings.ToLower(symb)),
		"amount":      amount,
		"marginCoin":  "USDT",
		"holdSide":    "short",
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", GetFlag(), err)
		fmt.Println(err)
		return
	}
	return nil
}

// 减少空头逐仓保证金
// doc: https://www.bitget.fit/zh-CN/api-doc/contract/account/Change-Margin
func SwapDecShortPositionMargin(symb string, amount float64) (err error) {
	_, _, err = root.ApiConfig.Post(Gateway, SetMarginUrl, map[string]any{
		"productType": "USDT-FUTURES",
		"symbol":      fmt.Sprintf("%sUSDT", strings.ToLower(symb)),
		"amount":      -amount,
		"marginCoin":  "USDT",
		"holdSide":    "short",
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", GetFlag(), err)
		fmt.Println(err)
		return
	}
	return nil
}

// 增加多头逐仓保证金
// doc: https://www.bitget.fit/zh-CN/api-doc/contract/account/Change-Margin
func SwapIncLongPositionMargin(symb string, amount float64) (err error) {
	_, _, err = root.ApiConfig.Post(Gateway, SetMarginUrl, map[string]any{
		"productType": "USDT-FUTURES",
		"symbol":      fmt.Sprintf("%sUSDT", strings.ToLower(symb)),
		"amount":      amount,
		"marginCoin":  "USDT",
		"holdSide":    "long",
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", GetFlag(), err)
		fmt.Println(err)
		return
	}
	return nil
}

// 减少多头逐仓保证金
// doc: https://www.bitget.fit/zh-CN/api-doc/contract/account/Change-Margin
func SwapDecLongPositionMargin(symb string, amount float64) (err error) {
	_, _, err = root.ApiConfig.Post(Gateway, SetMarginUrl, map[string]any{
		"productType": "USDT-FUTURES",
		"symbol":      fmt.Sprintf("%sUSDT", strings.ToLower(symb)),
		"amount":      -amount,
		"marginCoin":  "USDT",
		"holdSide":    "long",
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", GetFlag(), err)
		fmt.Println(err)
		return
	}
	return nil
}
