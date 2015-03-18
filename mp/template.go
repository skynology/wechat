package mp

import (
	"encoding/json"
	"errors"
)

type TemplateMessage struct {
	ToUser     string `json:"touser"`             // 必须, 接受者OpenID
	TemplateId string `json:"template_id"`        // 必须, 模版ID
	URL        string `json:"url,omitempty"`      // 可选, 用户点击后跳转的URL，该URL必须处于开发者在公众平台网站中设置的域中
	TopColor   string `json:"topcolor,omitempty"` // 可选, 整个消息的颜色, 可以不设置

	// 必须, JSON 格式的 []byte, 满足特定的模板需求
	RawJSONData json.RawMessage `json:"data"`
}

// 设置所属行业.
//  目前 industryId 的个数只能为 2.
func (clt *Client) SetIndustry(industryId ...int64) (err error) {
	if len(industryId) < 2 {
		return errors.New("industryId 的个数不能小于 2")
	}

	var request = struct {
		IndustryId1 int64 `json:"industry_id1"`
		IndustryId2 int64 `json:"industry_id2"`
	}{
		IndustryId1: industryId[0],
		IndustryId2: industryId[1],
	}

	var result Error

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/template/api_set_industry?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result
		return
	}
	return
}

// 从行业模板库选择模板添加到账号后台, 并返回模板id.
//  templateIdShort: 模板库中模板的编号，有“TM**”和“OPENTMTM**”等形式.
func (clt *Client) AddTemplate(templateIdShort string) (templateId string, err error) {
	var request = struct {
		TemplateIdShort string `json:"template_id_short"`
	}{
		TemplateIdShort: templateIdShort,
	}

	var result struct {
		Error
		TemplateId string `json:"template_id"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/template/api_add_template?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	templateId = result.TemplateId
	return
}

// 发送模板消息
func (clt *Client) SendTemplateMessage(msg *TemplateMessage) (msgid int64, err error) {
	if msg == nil {
		err = errors.New("nil TemplateMessage")
		return
	}

	var result struct {
		Error
		MsgId int64 `json:"msgid"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token="
	if err = clt.PostJSON(incompleteURL, msg, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	msgid = result.MsgId
	return
}
