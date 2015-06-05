package pay

import "net/url"

// 统一下单
type ShortURL struct {
	XMLName  struct{} `xml:"xml" json:"-"`
	AppId    string   `xml:"appid"   json:"appid"`
	MchId    string   `xml:"mch_id" json:"mch_id"`
	LongURL  string   `xml:"long_url" json:"long_url"`
	NonceStr string   `xml:"nonce_str" json:"nonce_str"`
	Sign     string   `xml:"sign" json:"sign"`
}

// 转换短链接.
func (clt *Client) ShortURL(req ShortURL) (resp map[string]string, err error) {
	return clt.PostXML("https://api.mch.weixin.qq.com/tools/shorturl", req)
}

// 生成二维码规则
func (clt *Client) GenBizPayURL(productId string, timestamp string, nonceStr string) string {
	m := make(map[string]string, 5)
	m["appid"] = clt.appId
	m["mch_id"] = clt.mchId
	m["product_id"] = productId
	m["time_stamp"] = timestamp
	m["nonce_str"] = nonceStr

	sign := clt.Sign(m)

	url := "weixin://wxpay/bizpayurl?sign=" + sign +
		"&appid=" + url.QueryEscape(clt.appId) +
		"&mch_id=" + url.QueryEscape(clt.mchId) +
		"&product_id=" + url.QueryEscape(productId) +
		"&time_stamp=" + url.QueryEscape(timestamp) +
		"&nonce_str=" + url.QueryEscape(nonceStr)

	return url
}
