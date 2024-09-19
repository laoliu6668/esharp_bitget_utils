package wss

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	root "github.com/laoliu6668/esharp_bitget_utils"
	"github.com/laoliu6668/esharp_bitget_utils/util"
	"github.com/laoliu6668/esharp_bitget_utils/util/websocketclient"
)

type SpotOrderMessageData struct {
	InstId        string `json:"instId"`        // 产品id 例如：ETHUSDT
	OrderId       string `json:"orderId"`       // 订单id
	Price         string `json:"price"`         // 	委托价格
	Size          string `json:"size"`          // 委托数量 side=buy 时，该值为计价币数量  side=sell 时，该值为基础币数量
	Notional      string `json:"notional"`      // 买入金额，市价买入时返回
	OrdType       string `json:"ordType"`       // 订单类型，market：市价单 limit：限价单
	Side          string `json:"side"`          // 订单方向 buy|sell
	AccBaseVolume string `json:"accBaseVolume"` // 累计已成交数量
	PriceAvg      string `json:"priceAvg"`      // 累计成交均价
	Status        string `json:"status"`        // filled
	CTime         string `json:"cTime"`
	FillTime      string `json:"fillTime"`
}
type SpotOrderMessage struct {
	Event  string                 `json:"event"`
	Arg    map[string]any         `json:"arg"`
	Action string                 `json:"action"`
	Data   []SpotOrderMessageData `json:"data"`
}

// 订阅现货账户变化
// 同步:ReciveBalanceMsg 异步:reciveOrderHandle、logHandle、errHandle
func SubSpotOrder(symbols []string, reciveOrderHandle func(ReciveSpotOrderMsg), logHandle func(string), errHandle func(error)) {

	onConnected := func(ws *websocketclient.Wsc) {
		SendAuth(ws)
	}

	onRecived := func(msg string, ws *websocketclient.Wsc) {
		ms := SpotOrderMessage{}
		err := json.Unmarshal([]byte(msg), &ms)
		if err != nil {
			go errHandle(fmt.Errorf("msg json.Unmarshal err: %s", msg))
			return
		}
		if ms.Action == "snapshot" {
			if ms.Arg["channel"] == "order" {
				// 订单频道
				for _, m := range ms.Data {
					if m.OrdType == "market" && m.Status == "filled" {
						var (
							orderVolume float64
							orderValue  float64
							tradePrice  = util.ParseFloat(m.PriceAvg, 0)
							tradeVolume = util.ParseFloat(m.AccBaseVolume, 0)
						)
						if m.Side == "buy" {
							orderValue = util.ParseFloat(m.Size, 0)
						} else {
							orderVolume = util.ParseFloat(m.Size, 0)
						}
						ctime, _ := strconv.Atoi(m.CTime)
						ftime, _ := strconv.Atoi(m.FillTime)
						go reciveOrderHandle(ReciveSpotOrderMsg{
							Exchange:    root.ExchangeName,
							Symbol:      strings.Replace(m.InstId, "USDT", "", 1),
							OrderId:     m.OrderId,
							OrderType:   fmt.Sprintf("%s-market", m.Side),
							OrderPrice:  util.ParseFloat(m.Price, 0),
							OrderValue:  orderValue,
							OrderVolume: orderVolume,
							TradePrice:  tradePrice,
							TradeValue:  tradePrice * tradeVolume,
							TradeVolume: tradeVolume,
							Status:      2,
							CreateAt:    int64(ctime),
							FilledAt:    int64(ftime),
						})
					}
				}

			}
		} else if ms.Event == "subscribe" {
			buf, _ := json.Marshal(ms.Arg)
			fmt.Printf("订阅成功: %s\n", buf)
			// go logHandle(fmt.Sprintf("订阅成功: %s", buf))
		} else if ms.Event == "login" {
			go logHandle("ws登录成功: " + time.Now().Format("2006-01-02 15:04:05"))
			subSpotOrder(ws, symbols)
		} else {
			go logHandle("unkown msg: " + msg)
		}

	}

	SubWss(PrivateGateway, onConnected, onRecived, logHandle, errHandle)

}

func subSpotOrder(ws *websocketclient.Wsc, symbols []string) {
	args := []map[string]string{}
	for _, s := range symbols {
		args = append(args, map[string]string{
			"instType": "SPOT",
			"channel":  "orders",
			"instId":   fmt.Sprintf("%sUSDT", s),
		})
	}
	mp := map[string]any{
		"op":   "subscribe",
		"args": args,
	}
	buf, _ := json.Marshal(mp)
	ws.SendTextMessage(string(buf))
}
