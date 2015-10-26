// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package corp

import (
	"fmt"
	"net/url"

	"github.com/skynology/go-crypto"
)

// 构造获取code的URL.
//  corpId:      企业的CorpID
//  redirectURL: 授权后重定向的回调链接地址, 员工点击后，页面将跳转至
//               redirect_uri/?code=CODE&state=STATE，企业可根据code参数获得员工的userid。
//  scope:       应用授权作用域，此时固定为：snsapi_base
//  state:       重定向后会带上state参数，企业可以填写a-zA-Z0-9的参数值，长度不可超过128个字节
func AuthCodeURL(corpId, redirectURL, scope string) (authUrl string, state string) {
	if scope == "" {
		scope = "snsapi_base"
	}
	state = crypto.GetRandomKey()
	authUrl = "https://open.weixin.qq.com/connect/oauth2/authorize" +
		"?appid=" + url.QueryEscape(corpId) +
		"&redirect_uri=" + url.QueryEscape(redirectURL) +
		"&response_type=code&scope=" + url.QueryEscape(scope) +
		"&state=" + url.QueryEscape(state) +
		"#wechat_redirect"
	return
}

type AuthUserInfo struct {
	UserId   string `json:"UserId"`   // 员工UserID
	DeviceId string `json:"DeviceId"` // 手机设备号(由微信在安装时随机生成)
}

// 根据code获取成员信息.
//  agentId: 跳转链接时所在的企业应用ID
//  code:    通过员工授权获取到的code，每次员工授权带上的code将不一样，
//           code只能使用一次，5分钟未被使用自动过期
func (clt *Client) GetUserIdByCode(code string) (info *AuthUserInfo, err error) {
	var result struct {
		Error
		AuthUserInfo
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?code=" + url.QueryEscape(code) +
		"&access_token="
	fmt.Println("url:", incompleteURL)
	if err = clt.GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	info = &result.AuthUserInfo
	return
}
