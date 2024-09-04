package bitget

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/laoliu6668/esharp_bitget_utils/util"
)

func (c *ApiConfigModel) Get(gateway, path string, data map[string]any) (body []byte, resp *http.Response, err error) {
	return c.Request("GET", gateway, path, data, 0, true)
}

func (c *ApiConfigModel) Post(gateway, path string, data map[string]any) (body []byte, resp *http.Response, err error) {
	return c.Request("POST", gateway, path, data, 0, true)
}

func (c *ApiConfigModel) Delete(gateway, path string, data map[string]any) (body []byte, resp *http.Response, err error) {
	return c.Request("DELETE", gateway, path, data, 0, true)
}

const proto = "https://"

// 获取TRONSCAN API数据
func (c *ApiConfigModel) Request(method, gateway, path string, data map[string]any, timeout time.Duration, sign bool) (body []byte, resp *http.Response, err error) {

	if timeout == 0 {
		timeout = time.Second * 10
	}

	// 创建http client
	client := &http.Client{
		Timeout: timeout,
	}
	if UseProxy {
		uri, _ := url.Parse(fmt.Sprintf("http://%s", ProxyUrl))
		fmt.Printf("uri: %v\n", uri)
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(uri),
		}
	}
	if data == nil {
		data = make(map[string]any)
	}

	// 构造query
	url := proto + gateway + path

	// 声明 body
	var reqBody io.Reader
	if method == "POST" || method == "PUT" || method == "DELETE" {
		if len(data) != 0 {
			buf, _ := json.Marshal(data)
			// 添加body
			reqBody = bytes.NewReader(buf)
		}

	} else if method == "GET" {
		url = GetQueryUrl(proto, gateway, path, data)
	} else {
		err = errors.New("不支持的http方法")
		return
	}
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return
	}
	if sign {
		req.Header.Add("ACCESS-KEY", ApiConfig.AccessKey)
		timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
		req.Header.Add("ACCESS-TIMESTAMP", timestamp)
		req.Header.Add("ACCESS-SIGN", Signature(data, method, path, timestamp, ApiConfig.SecretKey))
		req.Header.Add("locale", "zh-CN")
		req.Header.Add("ACCESS-PASSPHRASE", ApiConfig.PassPhrase)

	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		msg := string(body)
		mp := map[string]any{}
		err = json.Unmarshal(body, &mp)
		if err == nil {
			if _, ok := mp["msg"]; ok {
				msg = fmt.Sprintf("%v", mp["msg"])
			}
		}
		return nil, nil, fmt.Errorf("http %v %v", resp.StatusCode, string(msg))
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	type Resb struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Data    any    `json:"data"`
	}
	resb := Resb{}
	err = json.Unmarshal(body, &resb)
	if err != nil {
		return
	}
	if resb.Code != "00000" {
		err = fmt.Errorf("code %s %s", resb.Code, resb.Message)
		return
	}
	body, err = json.Marshal(resb.Data)
	if err != nil {
		return
	}
	return
}

func GetQueryUrl(proto, gateway, path string, queryMap map[string]any) string {
	return fmt.Sprintf("%s%s%s?%s", proto, gateway, path, util.HttpBuildQuery(queryMap))

}

func Signature(args map[string]any, method, path, timestamp, secretKey string) string {
	p1 := timestamp
	p2 := strings.ToUpper(method)
	p3 := path
	p4 := ""
	if method == "GET" {
		if len(args) > 0 {
			p4 += "?" + util.HttpBuildQuery(args)
		}
	} else {
		buf, _ := json.Marshal(args)
		p4 = string(buf)
	}
	str := p1 + p2 + p3 + p4
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(str))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func UTCTimeNow() string {
	return time.Now().In(time.UTC).Format("2006-01-02T15:04:05")
}
