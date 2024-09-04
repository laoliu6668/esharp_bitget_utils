package wss

import (
	"encoding/json"
	"fmt"
	"time"

	root "github.com/laoliu6668/esharp_bitget_utils"
	"github.com/laoliu6668/esharp_bitget_utils/util"
	"github.com/laoliu6668/esharp_bitget_utils/util/websocketclient"
)

type SpotAccMessageData struct {
	Coin      string `json:"coin"`
	Available string `json:"available"`
	UTime     string `json:"uTime"`
}
type SpotAccMessage struct {
	Event  string               `json:"event"`
	Arg    map[string]any       `json:"arg"`
	Action string               `json:"action"`
	Data   []SpotAccMessageData `json:"data"`
}

// 订阅现货账户变化
// 同步:ReciveBalanceMsg 异步:reciveOrderHandle、logHandle、errHandle
func SubSpotAccount(reciveAccHandle func(ReciveBalanceMsg), logHandle func(string), errHandle func(error)) {

	onConnected := func(ws *websocketclient.Wsc) {
		SendAuth(ws)
	}

	onRecived := func(msg string, ws *websocketclient.Wsc) {
		ms := SpotAccMessage{}
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
			subSpotAccount(ws)
		} else if ms.Action == "snapshot" || ms.Action == "update" {
			if ms.Arg["channel"] == "account" {
				// 账户频道
				for _, m := range ms.Data {
					reciveAccHandle(ReciveBalanceMsg{
						Exchange: root.ExchangeName,
						Symbol:   m.Coin,
						Free:     util.ParseFloat(m.Available, 0),
					})
				}
			}
		} else {
			go logHandle("unkown msg: " + msg)
		}

	}

	SubWss(PrivateGateway, onConnected, onRecived, logHandle, errHandle)

}

func subSpotAccount(ws *websocketclient.Wsc) {
	mp := map[string]any{
		"op": "subscribe",
		"args": []map[string]string{
			{
				"instType": "SPOT",
				"channel":  "account",
				"coin":     "default",
			},
		},
	}
	buf, _ := json.Marshal(mp)
	ws.SendTextMessage(string(buf))
}
