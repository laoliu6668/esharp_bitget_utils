package wss

import (
	"encoding/json"
	"fmt"
	"time"

	root "github.com/laoliu6668/esharp_bitget_utils"
	"github.com/laoliu6668/esharp_bitget_utils/util"
	"github.com/laoliu6668/esharp_bitget_utils/util/websocketclient"
)

const (
	PublicGateway  = "wss://ws.bitget.com/v2/ws/public"
	PrivateGateway = "wss://ws.bitget.com/v2/ws/private"
)

func SubWss(url string, OnConnected func(*websocketclient.Wsc), onTextMessageReceived func(string, *websocketclient.Wsc), logHandle func(string), errHandle func(error)) {
	var flag = util.GetFuncName(1)
	proxyUrl := ""
	if root.UseProxy {
		go logHandle(fmt.Sprintf("proxyUrl: %v\n", root.ProxyUrl))
		proxyUrl = fmt.Sprintf("http://%s", root.ProxyUrl)

	}
	ws := websocketclient.New(url, proxyUrl)
	ws.OnConnectError(func(err error) {
		go errHandle(err)
	})
	ws.OnDisconnected(func(err error) {
		go errHandle(err)
	})
	ws.OnConnected(func() {
		go logHandle(fmt.Sprintf("%s connected", flag))
		OnConnected(ws)
	})
	ticker := time.NewTicker(time.Second * 25)
	go func() {
		for range ticker.C {
			ws.SendTextMessage("ping")
		}
	}()
	ws.OnTextMessageReceived(func(message string) {
		if message == "pong" {
			return
		}
		onTextMessageReceived(message, ws)
	})

	ws.OnClose(func(code int, text string) {
		go errHandle(fmt.Errorf("close: %v, %v", code, text))
	})

	ws.Connect()

}

func SendAuth(ws *websocketclient.Wsc) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	mp := map[string]any{
		"op": "login",
		"args": []map[string]string{
			{
				"apiKey":     root.ApiConfig.AccessKey,
				"passphrase": root.ApiConfig.PassPhrase,
				"timestamp":  timestamp,
				"sign":       root.Signature(nil, "GET", "/user/verify", timestamp, root.ApiConfig.SecretKey),
			},
		},
	}
	buf, _ := json.Marshal(mp)
	ws.SendTextMessage(string(buf))
}
