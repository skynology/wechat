package mp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/skcloud/crypto"
)

const (
	Language_zh_CN = "zh_CN" // 简体中文
	Language_zh_TW = "zh_TW" // 繁体中文
	Language_en    = "en"    // 英文
)

const (
	SexUnknown = 0 // 未知
	SexMale    = 1 // 男性
	SexFemale  = 2 // 女性
)

type OauthUserInfo struct {
	OpenId   string `json:"openid"`   // 用户的唯一标识
	Nickname string `json:"nickname"` // 用户昵称
	Sex      int    `json:"sex"`      // 用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	City     string `json:"city"`     // 普通用户个人资料填写的城市
	Province string `json:"province"` // 用户个人资料填写的省份
	Country  string `json:"country"`  // 国家，如中国为CN

	// 用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），
	// 用户没有头像时该项为空
	HeadImageURL string `json:"headimgurl,omitempty"`

	// 用户特权信息，json 数组，如微信沃卡用户为（chinaunicom）
	Privilege []string `json:"privilege"`

	// 用户统一标识。针对一个微信开放平台帐号下的应用，同一用户的unionid是唯一的。
	UnionId string `json:"unionid"`
}

// 获取用户图像的大小, 如果用户没有图像则返回 ErrNoHeadImage 错误.
func (info *OauthUserInfo) HeadImageSize() (size int, err error) {
	HeadImageURL := info.HeadImageURL
	if HeadImageURL == "" {
		err = ErrNoHeadImage
		return
	}

	lastSlashIndex := strings.LastIndex(HeadImageURL, "/")
	if lastSlashIndex == -1 {
		err = fmt.Errorf("invalid HeadImageURL: %s", HeadImageURL)
		return
	}
	HeadImageIndex := lastSlashIndex + 1
	if HeadImageIndex == len(HeadImageURL) {
		err = fmt.Errorf("invalid HeadImageURL: %s", HeadImageURL)
		return
	}

	sizeStr := HeadImageURL[HeadImageIndex:]

	size64, err := strconv.ParseUint(sizeStr, 10, 64)
	if err != nil {
		err = fmt.Errorf("invalid HeadImageURL: %s", HeadImageURL)
		return
	}

	if size64 == 0 {
		size64 = 640
	}
	size = int(size64)
	return
}

// 获取用户信息(需scope为 snsapi_userinfo).
//  NOTE:
//  1. Client 需要指定 OAuth2Config, OAuth2Token
//  2. lang 可能的取值是 zh_CN, zh_TW, en, 如果留空 "" 则默认为 zh_CN.
func (clt *Client) OauthUserInfo(openId string, accessToken string, lang string) (info *UserInfo, err error) {
	if lang == "" {
		lang = Language_zh_CN
	}
	if lang != Language_zh_CN && lang != Language_en && lang != Language_zh_TW {
		err = errors.New("错误的 lang 参数")
		return
	}

	_url := "https://api.weixin.qq.com/sns/userinfo" +
		"?access_token=" + url.QueryEscape(accessToken) +
		"&openid=" + url.QueryEscape(openId) +
		"&lang=" + url.QueryEscape(lang)
	httpResp, err := clt.httpClient.Get(_url)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http.Status: %s", httpResp.Status)
		return
	}

	var result struct {
		Error
		UserInfo
	}

	if err = json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	info = &result.UserInfo
	return

}

// 构造请求用户授权获取code的地址.
//  appId:       公众号的唯一标识
//  redirectURL: 授权后重定向的回调链接地址
//               如果用户同意授权，页面将跳转至 redirect_uri/?code=CODE&state=STATE。
//               若用户禁止授权，则重定向后不会带上code参数，仅会带上state参数redirect_uri?state=STATE
//  scope:       应用授权作用域，
//               snsapi_base （不弹出授权页面，直接跳转，只能获取用户openid），
//               snsapi_userinfo （弹出授权页面，可通过openid拿到昵称、性别、所在地。
//               并且，即使在未关注的情况下，只要用户授权，也能获取其信息）
//  state:       重定向后会带上state参数，开发者可以填写a-zA-Z0-9的参数值，最多128字节
func OAuthCodeURL(appId, redirectURL, scope string) (authUrl string, state string) {
	state = crypto.GetRandomKey()
	authUrl = "https://open.weixin.qq.com/connect/oauth2/authorize" +
		"?appid=" + url.QueryEscape(appId) +
		"&redirect_uri=" + url.QueryEscape(redirectURL) +
		"&response_type=code&scope=" + url.QueryEscape(scope) +
		"&state=" + url.QueryEscape(state) +
		"#wechat_redirect"
	return
}

// 用户相关的 oauth2 token 信息
type OAuth2Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`

	OpenId    string   `json:"openid"`
	Scopes    []string `json:"scope"` // 用户授权的作用域
	UnionId   string   `json:"unionid"`
	ExpiresAt int64    `json:"-"` // 过期时间, unixtime, 分布式系统要求时间同步, 建议使用 NTP
}

// 判断授权的 OAuth2Token.AccessToken 是否过期, 过期返回 true, 否则返回 false
func (token *OAuth2Token) accessTokenExpired() bool {
	return time.Now().Unix() >= token.ExpiresAt
}

// 通过code换取网页授权access_token.
//  NOTE:
//  1. Client 需要指定 OAuth2Config
//  2. 如果指定了 OAuth2Token, 则会更新这个 OAuth2Token, 同时返回的也是指定的 OAuth2Token;
//     否则会重新分配一个 OAuth2Token.
func (clt *Client) GetOauthToken(code string) (token *OAuth2Token, err error) {
	token = new(OAuth2Token)

	_url := "https://api.weixin.qq.com/sns/oauth2/access_token" +
		"?appid=" + url.QueryEscape(clt.appId) +
		"&secret=" + url.QueryEscape(clt.appSecret) +
		"&code=" + url.QueryEscape(code) +
		"&grant_type=authorization_code"
	token, err = clt.updateOauthToken(_url)
	return
}

// 刷新access_token（如果需要）.
//  NOTE: Client 需要指定 OAuth2Config, OAuth2Token
func (clt *Client) OauthTokenRefresh(refreshToken string) (token *OAuth2Token, err error) {

	_url := "https://api.weixin.qq.com/sns/oauth2/refresh_token" +
		"?appid=" + url.QueryEscape(clt.appId) +
		"&grant_type=refresh_token&refresh_token=" + url.QueryEscape(refreshToken)
	token, err = clt.updateOauthToken(_url)
	return
}

// 检验授权凭证（access_token）是否有效.
//  NOTE:
//  1. Client 需要指定 OAuth2Token
//  2. 先判断 err 然后再判断 valid
func (clt *Client) CheckOauthAccessTokenValid(accessToken string, openId string) (valid bool, err error) {

	_url := "https://api.weixin.qq.com/sns/auth?access_token=" + url.QueryEscape(accessToken) +
		"&openid=" + url.QueryEscape(openId)
	httpResp, err := clt.httpClient.Get(_url)
	if err != nil {
		return
	}

	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http.Status: %s", httpResp.Status)
		return
	}

	var result Error
	if err = json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
		return
	}

	switch result.ErrCode {
	case ErrCodeOK:
		valid = true
		return
	case 40001:
		return
	default:
		err = &result
		return
	}
}

// 从服务器获取新的 token 更新 tk
func (clt *Client) updateOauthToken(url string) (tk *OAuth2Token, err error) {
	tk = new(OAuth2Token)

	httpResp, err := clt.httpClient.Get(url)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http.Status: %s", httpResp.Status)
		return
	}

	var result struct {
		Error
		AccessToken  string `json:"access_token"`  // 网页授权接口调用凭证,注意：此access_token与基础支持的access_token不同
		RefreshToken string `json:"refresh_token"` // 用户刷新access_token
		ExpiresIn    int64  `json:"expires_in"`    // access_token接口调用凭证超时时间，单位（秒）
		OpenId       string `json:"openid"`        // 用户唯一标识，请注意，在未关注公众号时，用户访问公众号的网页，也会产生一个用户和公众号唯一的OpenID
		Scope        string `json:"scope"`         // 用户授权的作用域，使用逗号（,）分隔
		UnionId      string `json:"unionid"`
	}

	if err = json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}

	// 由于网络的延时, 分布式服务器之间的时间可能不是绝对同步, access_token 过期时间留了一个缓冲区;
	switch {
	case result.ExpiresIn > 60*60:
		result.ExpiresIn -= 60 * 20
	case result.ExpiresIn > 60*30:
		result.ExpiresIn -= 60 * 10
	case result.ExpiresIn > 60*15:
		result.ExpiresIn -= 60 * 5
	case result.ExpiresIn > 60*5:
		result.ExpiresIn -= 60
	case result.ExpiresIn > 60:
		result.ExpiresIn -= 20
	case result.ExpiresIn > 0:
	default:
		err = fmt.Errorf("invalid expires_in: %d", result.ExpiresIn)
		return
	}

	tk.AccessToken = result.AccessToken
	if result.RefreshToken != "" {
		tk.RefreshToken = result.RefreshToken
	}
	tk.ExpiresIn = result.ExpiresIn
	tk.ExpiresAt = time.Now().Unix() + result.ExpiresIn

	tk.OpenId = result.OpenId
	tk.UnionId = result.UnionId

	strs := strings.Split(result.Scope, ",")
	tk.Scopes = make([]string, 0, len(strs))
	for _, str := range strs {
		str = strings.TrimSpace(str)
		if str == "" {
			continue
		}
		tk.Scopes = append(tk.Scopes, str)
	}
	return
}
