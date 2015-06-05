package pay

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/skynology/wechat/util"
)

type Client struct {
	appId      string
	mchId      string
	apiKey     string
	httpClient *http.Client
}

func (cli *Client) SetHttpClient(c *http.Client) {
	cli.httpClient = c
}

// 创建一个新的 Client.
//  如果 httpClient == nil 则默认用 http.DefaultClient.
func NewClient(appId, mchId, apiKey string) *Client {
	return &Client{
		appId:      appId,
		mchId:      mchId,
		apiKey:     apiKey,
		httpClient: http.DefaultClient,
	}
}

// NewTLSHttpClient 创建支持双向证书认证的 http.Client
func NewTLSHttpClient(certFile, keyFile string) (httpClient *http.Client, err error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	httpClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig:     tlsConfig,
		},
		Timeout: 60 * time.Second,
	}
	return
}

// 微信支付通用请求方法.
//  注意: err == nil 表示协议状态都为 SUCCESS.
func (clt *Client) PostXML(url string, request interface{}) (resp map[string]string, err error) {
	b, err := xml.Marshal(request)
	if err != nil {
		return
	}

	httpResp, err := clt.httpClient.Post(url, "text/xml; charset=utf-8", bytes.NewReader(b))
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http.Status: %s", httpResp.Status)
		return
	}

	if resp, err = util.ParseXMLToMap(httpResp.Body); err != nil {
		return
	}

	// 判断协议状态
	ReturnCode, ok := resp["return_code"]
	if !ok {
		err = errors.New("no return_code parameter")
		return
	}
	if ReturnCode != ReturnCodeSuccess {
		err = &Error{
			ReturnCode: ReturnCode,
			ReturnMsg:  resp["return_msg"],
		}
		return
	}

	// 认证签名
	signature1, ok := resp["sign"]
	if !ok {
		err = errors.New("no sign parameter")
		return
	}
	signature2 := clt.Sign(resp)
	if signature1 != signature2 {
		err = fmt.Errorf("check signature failed, \r\ninput: %q, \r\nlocal: %q", signature1, signature2)
		return
	}
	return
}

// 测速上报.
func (clt *Client) Report(req map[string]string) (resp map[string]string, err error) {
	return clt.PostXML("https://api.mch.weixin.qq.com/payitil/report", req)
}

// 下载对账单.
func (clt *Client) DownloadBill(req map[string]string) (data []byte, err error) {
	bodyBuf := textBufferPool.Get().(*bytes.Buffer)
	bodyBuf.Reset()
	defer textBufferPool.Put(bodyBuf)

	if err = util.FormatMapToXML(bodyBuf, req); err != nil {
		return
	}

	url := "https://api.mch.weixin.qq.com/pay/downloadbill"
	httpResp, err := clt.httpClient.Post(url, "text/xml; charset=utf-8", bodyBuf)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http.Status: %s", httpResp.Status)
		return
	}

	respBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return
	}

	var result Error
	if err = xml.Unmarshal(respBody, &result); err == nil {
		err = &result
		return
	}

	data = respBody
	err = nil
	return
}

// 撤销支付API.
//  NOTE: 请求需要双向证书.
func (clt *Client) Reverse(req map[string]string) (resp map[string]string, err error) {
	return clt.PostXML("https://api.mch.weixin.qq.com/secapi/pay/reverse", req)
}
