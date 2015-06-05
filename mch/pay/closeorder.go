package pay

// 关闭订单
type CloseOrder struct {
	XMLName    struct{} `xml:"xml" json:"-"`
	AppId      string   `xml:"appid"   json:"appid"`
	MchId      string   `xml:"mch_id" json:"mch_id"`
	OutTradeNo string   `xml:"out_trade_no" json:"out_trade_no"`
	NonceStr   string   `xml:"nonce_str" json:"nonce_str"`
	Sign       string   `xml:"sign" json:"sign"`
}

// 关闭订单.
func (clt *Client) CloseOrder(req CloseOrder) (resp map[string]string, err error) {
	return clt.PostXML("https://api.mch.weixin.qq.com/pay/closeorder", req)
}
