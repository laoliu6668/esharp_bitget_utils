package apis

import (
	"encoding/json"
	"fmt"

	root "github.com/laoliu6668/esharp_bitget_utils"
)

type SpotFilterSymbol struct {
	Status            string `json:"status"`                   // online
	QuoteCoin         string `json:"quoteCoin"`                // USDT
	BaseCoin          string `json:"baseCoin"`                 // BTC
	MinTradeAmount    string `json:"minTradeAmount"`           // 最小交易数量
	MaxTradeAmount    string `json:"maxTradeAmount"`           // 最大交易数量
	MinTradeUSDT      string `json:"minTradeUSDT"`             // 最小USDT交易额
	QuantityPrecision string `json:"quantityPrecision"`        // 数量小数位
	PricePrecision    string `json:"prpricePrecisionicePlace"` // 价格小数位
	QuotePrecision    string `json:"quotePrecision"`           // 右币精度(金额精度)
}

// 获取现货全局过滤器
// https://bybit-exchange.github.io/docs/zh-TW/v5/market/instrument
func GetSpotFiliters() (data []SpotFilterSymbol, err error) {

	body, _, err := root.ApiConfig.Request(
		"GET",
		Gateway,
		"/api/v2/spot/public/symbols",
		nil,
		0,
		false,
	)
	if err != nil {
		err = fmt.Errorf("%s err: %v", GetFlag(), err)
		fmt.Println(err)
		return
	}
	// util.WriteTestJsonFile(GetFlag(), body)
	res := []SpotFilterSymbol{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", GetFlag(), err)
		return
	}
	data = []SpotFilterSymbol{}
	for _, v := range res {
		if v.Status == "online" && v.QuoteCoin == "USDT" {
			data = append(data, v)
		}
	}
	return data, nil
}

type SwapFilterSymbol struct {
	SymbolStatus string `json:"symbolStatus"`
	QuoteCoin    string `json:"quoteCoin"`    // USDT
	BaseCoin     string `json:"baseCoin"`     // BTC
	SymbolType   string `json:"symbolType"`   // perpetual 永续  delivery交割
	MinTradeUSDT string `json:"minTradeUSDT"` // 最小USDT交易额
	VolumePlace  string `json:"volumePlace"`  // 数量小数位
	PricePlace   string `json:"pricePlace"`   // 价格小数位
}

// 获取期货全局过滤器
// https://www.bitget.fit/zh-CN/api-doc/contract/market/Get-All-Symbols-Contracts
func GetSwapFiliters() (data []SwapFilterSymbol, err error) {
	body, _, err := root.ApiConfig.Request(
		"GET",
		Gateway,
		"/api/v2/mix/market/contracts",
		map[string]any{"productType": "USDT-FUTURES"},
		0,
		false,
	)
	if err != nil {
		err = fmt.Errorf("%s err: %v", GetFlag(), err)
		fmt.Println(err)
		return
	}
	// util.WriteTestJsonFile(GetFlag(), body)
	res := []SwapFilterSymbol{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", GetFlag(), err)
		return
	}
	data = []SwapFilterSymbol{}
	for _, v := range res {
		if v.SymbolStatus == "normal" && v.QuoteCoin == "USDT" && v.SymbolType == "perpetual" {
			data = append(data, v)
		}
	}
	return data, nil
}
