package pay

// 红包发放API.
//  NOTE: 请求需要双向证书
func (clt *Client) SendRedPack(req map[string]string) (resp map[string]string, err error) {
	return clt.PostXML("https://api.mch.weixin.qq.com/mmpaymkttransfers/sendredpack", req)
}
