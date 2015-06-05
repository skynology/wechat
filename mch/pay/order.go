package pay

// 查询订单
type OrderQuery struct {
	XMLName       struct{} `xml:"xml" json:"-"`
	AppId         string   `xml:"appid"   json:"appid"`
	MchId         string   `xml:"mch_id" json:"mch_id"`
	TransactionId string   `xml:"transaction_id" json:"transaction_id"`
	OutTradeNo    string   `xml:"out_trade_no" json:"out_trade_no"`
	NonceStr      string   `xml:"nonce_str" json:"nonce_str"`
	Sign          string   `xml:"sign" json:"sign"`
}

// 订单查询.
func (clt *Client) OrderQuery(req OrderQuery) (resp map[string]string, err error) {
	return clt.PostXML("https://api.mch.weixin.qq.com/pay/orderquery", req)
}
