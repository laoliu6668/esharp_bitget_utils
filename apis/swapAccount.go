package apis

import (
	"encoding/json"
	"fmt"
	"strings"

	root "github.com/laoliu6668/esharp_bitget_utils"
)

type Res struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// 变换逐全仓模式 (TRADE)
// param margin_type: isolated(逐仓), crossed(全仓)
// doc: https://www.bitget.fit/zh-CN/api-doc/contract/account/Change-Margin-Mode
func ChangeSwapMarginType(symbol string, ISOLATED bool) (err error) {
	marginType := "crossed"
	if ISOLATED {
		marginType = "isolated"
	}
	body, _, err := root.ApiConfig.Post(Gateway, "/api/v2/mix/account/set-margin-mode", map[string]any{
		"marginMode":  marginType,
		"productType": "USDT-FUTURES",
		"marginCoin":  "USDT",
		"symbol":      fmt.Sprintf("%sUSDT", strings.ToLower(symbol)),
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", GetFlag(), err)
		return
	}
	fmt.Printf("body: %s\n", body)

	return nil
}

// 更改持仓模式(TRADE)
// 限速规则: 5次/1s (uid)
// param dual_side_position: "true": 双向持仓模式；"false": 单向持仓模式
// doc: https://www.bitget.fit/zh-CN/api-doc/contract/account/Change-Hold-Mode
func ChangeSwapPositionSideDual(dual_side_position bool) (err error) {
	var flag = GetFlag()
	posMode := "one_way_mode"
	if dual_side_position {
		posMode = "hedge_mode"
	}
	_, _, err = root.ApiConfig.Post(Gateway, "/api/v2/mix/account/set-position-mode", map[string]any{
		"productType": "USDT-FUTURES",
		"posMode":     posMode,
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		return
	}
	// fmt.Printf("body: %s\n", body)
	return nil
}

// 调整开仓杠杆(TRADE)
// 限速规则: 5次/1s (uid)
// param leverage: 目标杠杆倍数：1 到 125 整数
// doc: https://www.bitget.fit/zh-CN/api-doc/contract/account/Change-Leverage
func ChangeSwapLeverage(symbol string, leverage int, holdSide string) (err error) {
	_, _, err = root.ApiConfig.Post(Gateway, "/api/v2/mix/account/set-leverage", map[string]any{
		"symbol":      fmt.Sprintf("%susdt", strings.ToLower(symbol)),
		"productType": "USDT-FUTURES",
		"marginCoin":  "USDT",
		"leverage":    leverage,
		"holdSide":    holdSide,
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", GetFlag(), err)
		return
	}
	return nil
}

type SwapBalance struct {
	MarginCoin string `json:"marginCoin"` // 保证金币种
	Locked     string `json:"locked"`     // 锁定数量(保证金币种)
	Available  string `json:"available"`  // 账户可用数量
}

// 账户信息 持仓
// doc: https://www.bitget.fit/zh-CN/api-doc/contract/account/Get-Account-List
func GetSwapBalance() (data []SwapBalance, err error) {
	body, _, err := root.ApiConfig.Get(Gateway, "/api/v2/mix/account/accounts", map[string]any{
		"productType": "USDT-FUTURES",
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", GetFlag(), err)
		return
	}
	// fmt.Printf("body: %s\n", body)
	err = json.Unmarshal(body, &data)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", GetFlag(), err)
		fmt.Println(err)
		return
	}
	return
}

type SwapFunding struct {
	Symbol      string `json:"symbol"`          // 交易对 "BTCUSDT"
	FundingRate string `json:"lastFundingRate"` // 最近更新的资金费率
	// "nextFundingTime": 1597392000000,   // 下次资金费时间
	// "interestRate": "0.00010000",       // 标的资产基础利率
	Time int64 `json:"time"` // 更新时间 1597370495002
}

// 期货资金费率
// doc: https://www.bitget.fit/zh-CN/api-doc/contract/market/Get-Current-Funding-Rate
func GetSwapFunding() (data []SwapFunding, err error) {
	body, _, err := root.ApiConfig.Request("GET", Gateway, "/api/v2/mix/market/current-fund-rate", nil, 0, false)
	if err != nil {
		err = fmt.Errorf("%s err: %v", GetFlag(), err)
		return
	}
	fmt.Printf("body: %s\n", body)
	err = json.Unmarshal(body, &data)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", GetFlag(), err)
		fmt.Println(err)
		return
	}
	return
}

// // 账户信息 持仓风险
// // doc: https://binance-docs.github.io/apidocs/futures/cn/#v2-user_data-2
// func GetPositionRisk() (data []SwapPosition, err error) {
// 	const flag = "binance GetPositionRisk"
// 	body, _, err := root.ApiConfig.Get(Gateway, "/fapi/v3/positionRisk", nil)
// 	if err != nil {
// 		err = fmt.Errorf("%s err: %v", flag, err)
// 		return
// 	}
// 	err = json.Unmarshal(body, &data)
// 	if err != nil {
// 		err = fmt.Errorf("%s jsonDecodeErr: %v", flag, err)
// 		fmt.Println(err)
// 		return
// 	}
// 	return
// }
