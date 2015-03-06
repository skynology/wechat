package corp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

type Client struct {
	appId      string
	appSecret  string
	httpClient *http.Client
	tokenInfo  tokenInfo
	ticketInfo ticketInfo
}

func NewClient(appId string, appSecret string) *Client {
	return &Client{
		appId:      appId,
		appSecret:  appSecret,
		httpClient: http.DefaultClient,
	}
}

func (c *Client) SetToken(token tokenInfo) {
	c.tokenInfo = token
}
func (c *Client) SetTicket(ticket ticketInfo) {
	c.ticketInfo = ticket
}

type tokenInfo struct {
	Token     string `json:"access_token"`
	ExpiresIn int64  `json:"expires_in"`
}

func (c *Client) Token() (token string, err error) {
	if c.isValidToken() {
		token = c.tokenInfo.Token
		return
	}
	return c.RefreshToken()
}
func (c *Client) RefreshToken() (token string, err error) {
	tokenInfo, err := c.getToken()
	if err != nil {
		return
	}
	c.tokenInfo = tokenInfo
	token = c.tokenInfo.Token
	return
}

func (c *Client) isValidToken() bool {
	timeNowUnix := time.Now().Unix()

	if timeNowUnix+2 >= c.tokenInfo.ExpiresIn || c.tokenInfo.Token == "" {
		return false
	}

	return true
}
func (c *Client) getToken() (token tokenInfo, err error) {

	_url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s",
		url.QueryEscape(c.appId), url.QueryEscape(c.appSecret))

	if _url == "" {
		err = errors.New("invalid client type")
		return
	}

	httpResp, err := c.httpClient.Get(_url)
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
		tokenInfo
	}
	if err = json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	if result.ExpiresIn <= 0 {
		err = errors.New("invalid expires_in: " + strconv.FormatInt(result.ExpiresIn, 10))
		return
	}

	switch {
	case result.ExpiresIn > 60*60:
		result.ExpiresIn -= 60 * 10
	case result.ExpiresIn > 60*30:
		result.ExpiresIn -= 60 * 5
	case result.ExpiresIn > 60*5:
		result.ExpiresIn -= 60
	case result.ExpiresIn > 60:
		result.ExpiresIn -= 10
	case result.ExpiresIn > 0:
	default:
		err = errors.New("invalid expires_in: " + strconv.FormatInt(result.ExpiresIn, 10))
		return
	}
	token = result.tokenInfo
	return
}

type ticketInfo struct {
	Ticket    string `json:"ticket"`
	ExpiresIn int64  `json:"expires_in"` // 有效时间, seconds
}

func (c *Client) isValidTicket() bool {
	timeNowUnix := time.Now().Unix()
	if timeNowUnix+2 >= c.ticketInfo.ExpiresIn || c.ticketInfo.Ticket == "" {
		return false
	}
	return true
}
func (c *Client) Ticket() (ticket string, err error) {
	if c.isValidTicket() {
		ticket = c.ticketInfo.Ticket
		return
	}
	return c.RefreshTicket()
}
func (c *Client) RefreshTicket() (ticket string, err error) {
	ticketInfo, err := c.getTicket()
	if err != nil {
		return
	}
	c.ticketInfo = ticketInfo
	ticket = c.ticketInfo.Ticket
	return
}

// 从微信服务器获取 jsapi_ticket.
func (c *Client) getTicket() (ticket ticketInfo, err error) {
	var result struct {
		Error
		ticketInfo
	}
	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/get_jsapi_ticket?access_token="
	if err = c.GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}

	// 由于网络的延时, jsapi_ticket 过期时间留了一个缓冲区
	switch {
	case result.ExpiresIn > 60*60:
		result.ExpiresIn -= 60 * 10
	case result.ExpiresIn > 60*30:
		result.ExpiresIn -= 60 * 5
	case result.ExpiresIn > 60*5:
		result.ExpiresIn -= 60
	case result.ExpiresIn > 60:
		result.ExpiresIn -= 10
	case result.ExpiresIn > 0:
	default:
		err = errors.New("invalid expires_in: " + strconv.FormatInt(result.ExpiresIn, 10))
		return
	}

	ticket = result.ticketInfo
	return
}

// 用 encoding/json 把 request marshal 为 JSON, 放入 http 请求的 body 中,
// POST 到微信服务器, 然后将微信服务器返回的 JSON 用 encoding/json 解析到 response.
//
//  NOTE:
//  1. 一般不用调用这个方法, 请直接调用高层次的封装方法;
//  2. 最终的 URL == incompleteURL + access_token;
//  3. response 要求是 struct 的指针, 并且该 struct 拥有属性:
//     ErrCode int `json:"errcode"` (可以是直接属性, 也可以是匿名属性里的属性)
func (c *Client) PostJSON(incompleteURL string, request interface{}, response interface{}) (err error) {
	b, err := json.Marshal(request)
	if err != nil {
		return
	}

	b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
	b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
	b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)

	token, err := c.Token()
	if err != nil {
		return
	}

	hasRetried := false
RETRY:
	finalURL := incompleteURL + url.QueryEscape(token)

	httpResp, err := c.httpClient.Post(finalURL, "application/json; charset=utf-8", bytes.NewReader(b))
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("http.Status: %s", httpResp.Status)
	}

	if err = json.NewDecoder(httpResp.Body).Decode(response); err != nil {
		return
	}

	ErrCode := reflect.ValueOf(response).Elem().FieldByName("ErrCode").Int()

	switch ErrCode {
	case ErrCodeOK:
		return
	case ErrCodeTimeout, ErrCodeInvalidCredential:
		if !hasRetried {
			hasRetried = true

			if token, err = c.RefreshToken(); err != nil {
				return
			}
			goto RETRY
		}
		fallthrough
	default:
		return
	}
}

// GET 微信资源, 然后将微信服务器返回的 JSON 用 encoding/json 解析到 response.
//
//  NOTE:
//  1. 一般不用调用这个方法, 请直接调用高层次的封装方法;
//  2. 最终的 URL == incompleteURL + access_token;
//  3. response 要求是 struct 的指针, 并且该 struct 拥有属性:
//     ErrCode int `json:"errcode"` (可以是直接属性, 也可以是匿名属性里的属性)
func (c *Client) GetJSON(incompleteURL string, response interface{}) (err error) {
	token, err := c.Token()
	if err != nil {
		return
	}

	hasRetried := false
RETRY:
	finalURL := incompleteURL + url.QueryEscape(token)

	httpResp, err := c.httpClient.Get(finalURL)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("http.Status: %s", httpResp.Status)
	}

	if err = json.NewDecoder(httpResp.Body).Decode(response); err != nil {
		return
	}

	ErrCode := reflect.ValueOf(response).Elem().FieldByName("ErrCode").Int()

	switch ErrCode {
	case ErrCodeOK:
		return
	case ErrCodeTimeout, ErrCodeInvalidCredential:
		if !hasRetried {
			hasRetried = true

			if token, err = c.RefreshToken(); err != nil {
				return
			}
			goto RETRY
		}
		fallthrough
	default:
		return
	}
}
