package pay

type Transfer struct {
	XMLName        struct{} `xml:"xml" json:"-"`
	AppId          string   `xml:"mch_appid"   json:"mch_appid"`
	MchId          string   `xml:"mchid" json:"mchid"`
	DeviceInfo     string   `xml:"device_info" json:"device_info"`
	NonceStr       string   `xml:"nonce_str" json:"nonce_str"`
	Sign           string   `xml:"sign" json:"sign"`
	ParnetTradeNo  string   `xml:"partner_trade_no" json:"partner_trade_no"`
	OpenId         string   `xml:"openid" json:"openid"`
	CheckName      string   `xml:"check_name" json:"check_name"`
	ReUserName     string   `xml:"re_user_name" json:"re_user_name"`
	Amount         int      `xml:"amount" json:"amount"`
	Description    string   `xml:"desc" json:"desc"`
	SpbillCreateIP string   `xml:"spbill_create_ip" json:"spbill_create_ip"`
}

// 企业付款
//  NOTE: 请求需要双向证书.
func (clt *Client) Transfer(req Transfer) (resp map[string]string, err error) {
	return clt.PostXMLWithoutSign("https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers", req)
}

type TransferQuery struct {
	XMLName       struct{} `xml:"xml" json:"-"`
	AppId         string   `xml:"mch_appid"   json:"mch_appid"`
	MchId         string   `xml:"mchid" json:"mchid"`
	NonceStr      string   `xml:"nonce_str" json:"nonce_str"`
	ParnetTradeNo string   `xml:"partner_trade_no" json:"partner_trade_no"`
	Sign          string   `xml:"sign" json:"sign"`
}

// 查询企业付款
//  NOTE: 请求需要双向证书.
func (clt *Client) FindTransfer(req TransferQuery) (resp map[string]string, err error) {
	return clt.PostXMLWithoutSign("https://api.mch.weixin.qq.com/mmpaymkttransfers/gettransferinfo", req)
}
