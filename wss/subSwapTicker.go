package wss

import (
	"encoding/json"
	"fmt"
	"strings"

	root "github.com/laoliu6668/esharp_bitget_utils"
	"github.com/laoliu6668/esharp_bitget_utils/util"
	"github.com/laoliu6668/esharp_bitget_utils/util/websocketclient"
)

func SubSwapTicker(symbols []string, reciveHandle func(Ticker), logHandle func(string), errHandle func(error)) {

	args := []map[string]any{}
	for _, symbol := range symbols {
		args = append(args, map[string]any{
			"instType": "USDT-FUTURES",
			"channel":  "ticker",
			"instId":   fmt.Sprintf("%sUSDT", symbol),
		})
	}
	onConnected := func(ws *websocketclient.Wsc) {
		mp := map[string]any{
			"op":   "subscribe",
			"args": args,
		}
		buf, _ := json.Marshal(mp)
		ws.SendTextMessage(string(buf))
	}

	onRecived := func(msg string, ws *websocketclient.Wsc) {
		ms := TickerMessage{}
		err := json.Unmarshal([]byte(msg), &ms)
		if err != nil {
			go errHandle(fmt.Errorf("msg json.Unmarshal err: %s", msg))
			return
		}
		if ms.Action == "snapshot" {
			for _, m := range ms.Data {
				go reciveHandle(Ticker{
					Exchange: root.ExchangeName,
					Symbol:   strings.Replace(m.Symbol, "USDT", "", 1),
					Buy: Values{
						Price: util.ParseFloat(m.BuyPrice, 0),
						Size:  util.ParseFloat(m.BuySize, 0),
					},
					Sell: Values{
						Price: util.ParseFloat(m.SellPrice, 0),
						Size:  util.ParseFloat(m.SellSize, 0),
					},
					FundingRate: util.ParseFloat(m.FundingRate, 0),
					FundingTime: util.ParseInt(m.NextFundingTime, 0),
					UpdateAt:    root.GetTimeFloat(),
				})
			}
		} else if ms.Event == "subscribe" {
			// go logHandle("订阅成功: " + m.Symbol)
		} else {
			go logHandle("unkown msg: " + msg)
		}

	}

	SubWss(PublicGateway, onConnected, onRecived, logHandle, errHandle)

}
