package pay

// 扫码支付
type MicroPay struct {
	XMLName        struct{} `xml:"xml" json:"-"`
	AppId          string   `xml:"appid"   json:"appid"`
	MchId          string   `xml:"mch_id" json:"mch_id"`
	DeviceInfo     string   `xml:"device_info" json:"device_info"`
	NonceStr       string   `xml:"nonce_str" json:"nonce_str"`
	Sign           string   `xml:"sign" json:"sign"`
	Body           string   `xml:"body" json:"body"`
	Detail         string   `xml:"detail,omitempty" json:"detail,omitempty"`
	Attach         string   `xml:"attach,omitempty" json:"attach,omitempty"`
	OutTradeNo     string   `xml:"out_trade_no" json:"out_trade_no"`
	TotalFee       int      `xml:"total_fee" json:"total_fee"`
	FeeType        string   `xml:"fee_type,omitempty" json:"fee_type,omitempty"`
	SpbillCreateIP string   `xml:"spbill_create_ip" json:"spbill_create_ip"`
	TimeStart      string   `xml:"time_start,omitempty" json:"time_start,omitempty"`
	TimeExpire     string   `xml:"time_expire,omitempty" json:"time_expire,omitempty"`
	GoodsTag       string   `xml:"goods_tag,omitempty" json:"goods_tag,omitempty"`
	AuthCode       string   `xml:"auth_code" json:"auth_code"`
}

// 提交被扫支付API.
func (clt *Client) MicroPay(req MicroPay) (resp map[string]string, err error) {
	return clt.PostXML("https://api.mch.weixin.qq.com/pay/micropay", req)
}
