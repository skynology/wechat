package mp

// 添加客服会话.
//  account:    完整客服账号，格式为：账号前缀@公众号微信号，账号前缀最多10个字符，必须是英文或者数字字符。
//  openid:     客户openid
//  text:       附加信息，文本会展示在客服人员的多客服客户端
func (clt *Client) CreateKfSession(account, openid, text string) (err error) {
	request := struct {
		Account string `json:"kf_account"`
		OpenId  string `json:"openid"`
		Text    string `json:"text,omitempty"`
	}{
		Account: account,
		OpenId:  openid,
		Text:    text,
	}

	var result Error

	incompleteURL := "https://api.weixin.qq.com/customservice/kfsession/create?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result
		return
	}
	return
}

// 关闭客服会话.
//  account:    完整客服账号，格式为：账号前缀@公众号微信号，账号前缀最多10个字符，必须是英文或者数字字符。
//  openid:     客户openid
//  text:       附加信息，文本会展示在客服人员的多客服客户端
func (clt *Client) CloseKfSession(account, openid, text string) (err error) {
	request := struct {
		Account string `json:"kf_account"`
		OpenId  string `json:"openid"`
		Text    string `json:"text,omitempty"`
	}{
		Account: account,
		OpenId:  openid,
		Text:    text,
	}

	var result Error

	incompleteURL := "https://api.weixin.qq.com/customservice/kfsession/close?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result
		return
	}
	return
}

// 客户的会话状态
type KFSessionInfo struct {
	OpenId     string `json:"openid,omitempty"`     //客户openid
	Account    string `json:"kf_account,omitempty"` // 完整客服账号，格式为：账号前缀@公众号微信号
	CreateTime int64  `json:"createtime"`           // 会话接入的时间
}

// 获取客户的会话状态
//  openid:     客户openid
func (clt *Client) GetKfSession(openid string) (kfSessionInfo KFSessionInfo, err error) {
	var result struct {
		Error
		KFSessionInfo
	}

	incompleteURL := "https://api.weixin.qq.com/customservice/kfsession/getsession?openid=" + openid
	incompleteURL += "&access_token="
	if err = clt.GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	kfSessionInfo = result.KFSessionInfo
	return
}

// 获取客服的会话列表
//  openid:     客户openid
func (clt *Client) KfSessionList(kfaccount string) (kfSessionList []KFSessionInfo, err error) {
	var result struct {
		Error
		SessionList []KFSessionInfo `json:"sessionlist"`
	}

	incompleteURL := "https://api.weixin.qq.com/customservice/kfsession/getsessionlist?kf_account=" + kfaccount
	incompleteURL += "&access_token="
	if err = clt.GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	kfSessionList = result.SessionList
	return
}

type KFSessionWait struct {
	Count        int             `json:"count"`        //未接入会话数量
	WaitCaseList []KFSessionInfo `json:"waitcaselist"` //未接入会话列表，最多返回100条数据
}

// 获取未接入会话列表
func (clt *Client) KfSessionWaitList() (kfSessionWaitList KFSessionWait, err error) {
	var result struct {
		Error
		Count       int             `json:"count"`        //未接入会话数量
		SessionList []KFSessionInfo `json:"waitcaselist"` // 未接入会话列表，最多返回100条数据
	}

	incompleteURL := "https://api.weixin.qq.com/customservice/kfsession/getwaitcase?access_token="
	if err = clt.GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	kfSessionWaitList.Count = result.Count
	kfSessionWaitList.WaitCaseList = result.SessionList
	return
}
