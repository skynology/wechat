package corp

import "errors"

type Chat struct {
	ChatId  string   `json:"chatid"`
	Name    string   `json:"name"`
	OwnerId string   `json:"owner"`
	Users   []string `json:"userlist"`
}

// 创建会话
func (clt *Client) CreateChat(chat Chat) (err error) {
	var result Error

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/chat/create?access_token="
	if err = clt.PostJSON(incompleteURL, chat, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result
		return
	}
	return
}

// 获取会话
func (clt *Client) GetChat(id string) (chat Chat, err error) {
	var result struct {
		Error
		Chat Chat `json:"chat_info"`
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/chat/get?chatid=" +
		id + "&access_token="
	if err = clt.GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	chat = result.Chat
	return
}

// 退出会话
func (clt *Client) QuitChat(chatid string, userId string) (err error) {
	var result Error

	data := map[string]interface{}{
		"chatid":  chatid,
		"op_user": userId,
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/chat/quit?access_token="
	if err = clt.PostJSON(incompleteURL, data, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result
		return
	}
	return
}

type ChatRecevier struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

func (c *ChatRecevier) CheckValid() error {
	if c.Type != "single" && c.Type != "group" {
		return errors.New("the type should be 'single' or 'group' ")
	}
	return nil
}

// 清楚未读会话状态
func (clt *Client) CleanChatNotify(userId string, receiver ChatRecevier) (err error) {
	var result Error

	data := map[string]interface{}{
		"op_user": userId,
		"chat":    receiver,
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/chat/quit?access_token="
	if err = clt.PostJSON(incompleteURL, data, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result
		return
	}
	return
}

// 设置成员新消息免打扰
// isMute = true: 免打扰
// isMute = false: 解封免打扰
func (clt *Client) SetChatMute(userIds []string, isMute bool) (invalidusers []string, err error) {
	var result struct {
		Error
		InvalidUser []string `json:"invaliduser"`
	}

	status := 1
	if !isMute {
		status = 0
	}

	list := make([]map[string]interface{}, 0)

	for _, c := range userIds {
		list = append(list, map[string]interface{}{
			"userid": c,
			"status": status,
		})
	}

	data := map[string]interface{}{
		"user_mute_list": list,
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/chat/setmute?access_token="
	if err = clt.PostJSON(incompleteURL, data, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	invalidusers = result.InvalidUser
	return
}

type CommonChatHeader struct {
	Recevier ChatRecevier `json:"receiver"`
	Sender   string       `json:"sender"`
	MsgType  string       `json:"msgtype"`
}

type ChatText struct {
	CommonChatHeader
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
}
type ChatImage struct {
	CommonChatHeader

	Image struct {
		MediaId string `json:"media_id"` // 图片媒体文件id，可以调用上传媒体文件接口获取
	} `json:"image"`
}

type ChatFile struct {
	CommonChatHeader

	File struct {
		MediaId string `json:"media_id"` // 媒体文件id，可以调用上传媒体文件接口获取
	} `json:"file"`
}

func (clt *Client) SendChatText(msg *ChatText) (err error) {
	if msg == nil {
		err = errors.New("nil msg")
		return
	}
	return clt.SendChat(msg)
}

func (clt *Client) SendChatImage(msg *ChatImage) (err error) {
	if msg == nil {
		err = errors.New("nil msg")
		return
	}
	return clt.SendChat(msg)
}

func (clt *Client) SendChat(msg interface{}) (err error) {
	var result struct {
		Error
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/chat/send?access_token="
	if err = clt.PostJSON(incompleteURL, msg, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}

	return
}
