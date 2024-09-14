package wss

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	root "github.com/laoliu6668/esharp_bitget_utils"
	"github.com/laoliu6668/esharp_bitget_utils/util"
	"github.com/laoliu6668/esharp_bitget_utils/util/websocketclient"
)

type SwapAccMessageData struct {
	MarginCoin string `json:"marginCoin"`
	Available  string `json:"available"`
	Frozen     string `json:"frozen"`
}
type SwapPositionMessageData struct {
	InstId             string `json:"instId"`
	MarginCoin         string `json:"marginCoin"`         // ETHUSDT
	MarginSize         string `json:"marginSize"`         // 保证金
	Total              string `json:"total"`              // 持仓数量
	Available          string `json:"available"`          // 可平仓数量
	Leverage           int    `json:"leverage"`           // 杠杆倍数
	HoldSide           string `json:"holdSide"`           // 持仓方向
	UnrealizedPL       string `json:"unrealizedPL"`       // 未实现盈亏
	IsolatedMarginRate string `json:"isolatedMarginRate"` // 逐仓时，实际保证金率
}
type SwapOrderMessageData struct {
	InstId        string `json:"instId"`        // 产品id 例如：ETHUSDT
	OrderId       string `json:"orderId"`       // 订单id
	Price         string `json:"price"`         // 委托价格
	Size          string `json:"size"`          // 委托数量 side=buy 时，该值为计价币数量  side=sell 时，该值为基础币数量
	Notional      string `json:"notional"`      // 买入金额，市价买入时返回
	OrdType       string `json:"ordType"`       // 订单类型，market：市价单 limit：限价单
	PosSide       string `json:"posSide"`       // long|short
	Side          string `json:"side"`          // 订单方向 buy|sell
	AccBaseVolume string `json:"accBaseVolume"` // 累计已成交数量
	PriceAvg      string `json:"priceAvg"`      // 累计成交均价
	Status        string `json:"status"`        // filled
	CTime         string `json:"cTime"`
	FillTime      string `json:"fillTime"`
}
type SwapAccMessage struct {
	Event  string         `json:"event"`
	Arg    map[string]any `json:"arg"`
	Action string         `json:"action"`
}

var (
	SwapPositionCache     = map[string]SwapPositionMessageData{}
	SwapPositionCacheLock = sync.Mutex{}
)

// 订阅期货账户变化
// 余额+持仓+订单
func SubSwapAccount(reciveAccHandle func(ReciveBalanceMsg), recivePositionHandle func(RecivePositionMsg), reciveOrderHandle func(ReciveSwapOrderMsg), logHandle func(string), errHandle func(error)) {
	onConnected := func(ws *websocketclient.Wsc) {
		SendAuth(ws)
	}
	onRecived := func(msg string, ws *websocketclient.Wsc) {
		ms := SwapAccMessage{}
		err := json.Unmarshal([]byte(msg), &ms)
		if err != nil {
			go errHandle(fmt.Errorf("msg json.Unmarshal err: %s", msg))
			return
		}
		if ms.Event == "subscribe" {
			buf, _ := json.Marshal(ms.Arg)
			go logHandle(fmt.Sprintf("订阅成功: %s", buf))
		} else if ms.Event == "login" {
			go logHandle("ws登录成功: " + time.Now().Format("2006-01-02 15:04:05"))
			subSwapAccount(ws)
			subSwapOrder(ws)
			subSwapPosition(ws)
		} else if ms.Action == "snapshot" || ms.Action == "update" {
			if ms.Arg["channel"] == "account" {
				type Data struct {
					SwapAccMessage
					Data []SwapAccMessageData `json:"data"`
				}
				ms := Data{}
				err := json.Unmarshal([]byte(msg), &ms)
				if err != nil {
					go errHandle(fmt.Errorf("msg json.Unmarshal err: %s", msg))
					return
				}
				// 账户频道
				for _, m := range ms.Data {
					reciveAccHandle(ReciveBalanceMsg{
						Exchange: root.ExchangeName,
						Symbol:   m.MarginCoin,
						Free:     util.ParseFloat(m.Available, 0),
					})
				}
			} else if ms.Arg["channel"] == "positions" {
				// 持仓频道
				type Data struct {
					SwapAccMessage
					Data []SwapPositionMessageData `json:"data"`
				}
				ms := Data{}
				err := json.Unmarshal([]byte(msg), &ms)
				if err != nil {
					go errHandle(fmt.Errorf("msg json.Unmarshal err:%v %s", err, msg))
					return
				}

				SwapPositionCacheLock.Lock()
				for k, v := range SwapPositionCache {
					for i, m := range ms.Data {
						coin := strings.Replace(m.InstId, "USDT", "", 1)
						if k == coin {
							// 匹配
							// 判断是否变化
							if v.MarginSize == m.MarginSize && v.UnrealizedPL == m.UnrealizedPL && v.IsolatedMarginRate == m.IsolatedMarginRate && v.Leverage == m.Leverage && v.Total == m.Total {
								// 无变化 不推送
								break
							}
							// 有变化 推送
							margin := util.ParseFloat(m.MarginSize, 0) + util.ParseFloat(m.UnrealizedPL, 0)
							marginRate := util.FixedFloat(util.ParseFloat(m.IsolatedMarginRate, 0)*float64(m.Leverage)*100, 2)
							volume := util.ParseFloat(m.Total, 0)
							position := RecivePositionMsg{
								Exchange:    root.ExchangeName,
								Symbol:      coin,
								Margin:      margin,
								MarginRatio: marginRate,
							}
							if m.HoldSide == "long" {
								position.BuyVolume = volume
							} else {
								position.SellVolume = volume
							}
							recivePositionHandle(position)
							SwapPositionCache[k] = m
							break
						}
						if i == len(ms.Data)-1 {
							// 不匹配
							recivePositionHandle(RecivePositionMsg{
								Exchange:    root.ExchangeName,
								Symbol:      coin,
								BuyVolume:   0,
								SellVolume:  0,
								Margin:      0,
								MarginRatio: 0,
							})
							delete(SwapPositionCache, coin)

						}
					}
				}
				SwapPositionCacheLock.Unlock()

				// 账户频道
				for _, m := range ms.Data {
					position := RecivePositionMsg{
						Exchange:    root.ExchangeName,
						Symbol:      strings.Replace(m.InstId, "USDT", "", 1),
						Margin:      util.ParseFloat(m.MarginSize, 0) + util.ParseFloat(m.UnrealizedPL, 0),
						MarginRatio: util.FixedFloat(util.ParseFloat(m.IsolatedMarginRate, 0)*float64(m.Leverage)*100, 2),
					}
					if m.HoldSide == "long" {
						position.BuyVolume = util.ParseFloat(m.Total, 0)
					} else {
						position.SellVolume = util.ParseFloat(m.Total, 0)
					}
					recivePositionHandle(position)
				}
			} else if ms.Arg["channel"] == "orders" {
				// 订单频道
				// 持仓频道
				type Data struct {
					SwapAccMessage
					Data []SwapOrderMessageData `json:"data"`
				}
				ms := Data{}
				err := json.Unmarshal([]byte(msg), &ms)
				if err != nil {
					go errHandle(fmt.Errorf("msg json.Unmarshal err: %s", msg))
					return
				}
				for _, m := range ms.Data {
					var (
						orderVolume = util.ParseFloat(m.Size, 0)
						tradePrice  = util.ParseFloat(m.PriceAvg, 0)
						tradeVolume = util.ParseFloat(m.AccBaseVolume, 0)
						orderType   string
					)
					if m.PosSide == "short" {
						if m.Side == "buy" {
							orderType = "sell-open"
						} else {
							orderType = "buy-close"
						}
					} else if m.PosSide == "long" {
						if m.Side == "buy" {
							orderType = "buy-open"
						} else {
							orderType = "sell-close"
						}
					}
					ctime, _ := strconv.Atoi(m.CTime)
					ftime, _ := strconv.Atoi(m.FillTime)
					go reciveOrderHandle(ReciveSwapOrderMsg{
						Exchange:    root.ExchangeName,
						Symbol:      strings.Replace(m.InstId, "USDT", "", 1),
						OrderId:     m.OrderId,
						OrderType:   orderType,
						OrderPrice:  util.ParseFloat(m.Price, 0),
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
		} else {
			go logHandle("unkown msg: " + msg)
		}

	}

	SubWss(PrivateGateway, onConnected, onRecived, logHandle, errHandle)

}

func subSwapAccount(ws *websocketclient.Wsc) {
	mp := map[string]any{
		"op": "subscribe",
		"args": []map[string]string{
			{
				"instType": "USDT-FUTURES",
				"channel":  "account",
				"coin":     "default",
			},
		},
	}
	buf, _ := json.Marshal(mp)
	ws.SendTextMessage(string(buf))
}
func subSwapOrder(ws *websocketclient.Wsc) {
	mp := map[string]any{
		"op": "subscribe",
		"args": []map[string]string{
			{
				"instType": "USDT-FUTURES",
				"channel":  "orders",
				"instId":   "default",
			},
		},
	}
	buf, _ := json.Marshal(mp)
	ws.SendTextMessage(string(buf))
}
func subSwapPosition(ws *websocketclient.Wsc) {
	mp := map[string]any{
		"op": "subscribe",
		"args": []map[string]string{
			{
				"instType": "USDT-FUTURES",
				"channel":  "positions",
				"instId":   "default",
			},
		},
	}
	buf, _ := json.Marshal(mp)
	ws.SendTextMessage(string(buf))
}
