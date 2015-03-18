package mp

func (clt *Client) SendCustomMessage(msg interface{}) (err error) {
	var result Error

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token="
	if err = clt.PostJSON(incompleteURL, msg, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result
		return
	}
	return
}
