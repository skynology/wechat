package pay

// 退款申请
type Refund struct {
	XMLName       struct{} `xml:"xml" json:"-"`
	AppId         string   `xml:"appid"   json:"appid"`
	MchId         string   `xml:"mch_id" json:"mch_id"`
	DeviceInfo    string   `xml:"device_info" json:"device_info"`
	NonceStr      string   `xml:"nonce_str" json:"nonce_str"`
	Sign          string   `xml:"sign" json:"sign"`
	TransactionId string   `xml:"transaction_id" json:"transaction_id"`
	OutTradeNo    string   `xml:"out_trade_no" json:"out_trade_no"`
	OutRefundNo   string   `xml:"out_refund_no" json:"out_refund_no"`
	FeeType       string   `xml:"fee_type,omitempty" json:"fee_type,omitempty"`
	TotalFee      int      `xml:"total_fee" json:"total_fee"`
	RefundFee     int      `xml:"refund_fee" json:"refund_fee"`
	OpUserId      string   `xml:"op_user_id" json:"op_user_id"`
}

// 申请退款.
//  NOTE: 请求需要双向证书.
func (clt *Client) Refund(req Refund) (resp map[string]string, err error) {
	return clt.PostXML("https://api.mch.weixin.qq.com/secapi/pay/refund", req)
}

type RefundQuery struct {
	XMLName       struct{} `xml:"xml" json:"-"`
	AppId         string   `xml:"appid"   json:"appid"`
	MchId         string   `xml:"mch_id" json:"mch_id"`
	DeviceInfo    string   `xml:"device_info" json:"device_info"`
	NonceStr      string   `xml:"nonce_str" json:"nonce_str"`
	Sign          string   `xml:"sign" json:"sign"`
	TransactionId string   `xml:"transaction_id" json:"transaction_id"`
	OutTradeNo    string   `xml:"out_trade_no" json:"out_trade_no"`
	OutRefundNo   string   `xml:"out_refund_no" json:"out_refund_no"`
	RefundId      string   `xml:"refund_id,omitempty" json:"refund_id,omitempty"`
}

// 退款查询.
func (clt *Client) RefundQuery(req RefundQuery) (resp map[string]string, err error) {
	return clt.PostXML("https://api.mch.weixin.qq.com/pay/refundquery", req)
}
