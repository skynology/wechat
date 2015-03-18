package mp

// 发送对象类型
type MassMessageToType string

const (
	//MassMessageToAll      MassMessageToType = "all"
	MassMessageToGroup    MassMessageToType = "group"
	MassMessageToUserList MassMessageToType = "userlist"
	MassMessageToPreview  MassMessageToType = "preview"
)

const (
	MassMsgTypeText  = "text"
	MassMsgTypeImage = "image"
	MassMsgTypeVoice = "voice"
	MassMsgTypeVideo = "mpvideo"
	MassMsgTypeNews  = "mpnews"
)

type CommonMassMessageHeader struct {
	Filter struct {
		GroupId int64 `json:"group_id,omitempty"`
		IsToAll bool  `json:"is_to_all"`
	} `json:"filter"`
	ToUserList interface{} `json:"touser,omitempty"`
	MsgType    string      `json:"msgtype"`
}

type MassText struct {
	CommonMassMessageHeader
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
}

func NewMassText(groupId int64, content string) *MassText {
	var msg MassText
	msg.MsgType = MassMsgTypeText
	msg.Filter.GroupId = groupId
	msg.Text.Content = content
	return &msg
}

type MassImage struct {
	CommonMassMessageHeader
	Image struct {
		MediaId string `json:"media_id"`
	} `json:"image"`
}

func NewMassImage(groupId int64, mediaId string) *MassImage {
	var msg MassImage
	msg.MsgType = MassMsgTypeImage
	msg.Filter.GroupId = groupId
	msg.Image.MediaId = mediaId
	return &msg
}

type MassVoice struct {
	CommonMassMessageHeader
	Voice struct {
		MediaId string `json:"media_id"`
	} `json:"voice"`
}

func NewMassVoice(groupId int64, mediaId string) *MassVoice {
	var msg MassVoice
	msg.MsgType = MassMsgTypeVoice
	msg.Filter.GroupId = groupId
	msg.Voice.MediaId = mediaId
	return &msg
}

type MassVideo struct {
	CommonMassMessageHeader
	Video struct {
		MediaId string `json:"media_id"`
	} `json:"mpvideo"`
}

// 新建视频消息
//  NOTE: mediaId 应该通过 media.Client.CreateVideo 得到
func NewMassVideo(groupId int64, mediaId string) *MassVideo {
	var msg MassVideo
	msg.MsgType = MassMsgTypeVideo
	msg.Filter.GroupId = groupId
	msg.Video.MediaId = mediaId
	return &msg
}

// 图文消息
type MassNews struct {
	CommonMassMessageHeader
	News struct {
		MediaId string `json:"media_id"`
	} `json:"mpnews"`
}

// 新建图文消息
//  NOTE: mediaId 应该通过 media.Client.CreateNews 得到
func NewMassNews(groupId int64, mediaId string) *MassNews {
	var msg MassNews
	msg.MsgType = MassMsgTypeNews
	msg.Filter.GroupId = groupId
	msg.News.MediaId = mediaId
	return &msg
}

func (clt *Client) SendMassMassage(filterType MassMessageToType, msg interface{}) (msgid int64, err error) {
	var result struct {
		Error
		Type  string `json:"type"`
		MsgId int64  `json:"msg_id"`
	}

	urlPrefix := "https://api.weixin.qq.com/cgi-bin/message/mass/"
	incompleteURL := urlPrefix + "sendall?access_token="
	if filterType == MassMessageToPreview {
		incompleteURL = urlPrefix + "preview?access_token="
	} else if filterType == MassMessageToUserList {
		incompleteURL = urlPrefix + "send?access_token="
	}

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

func (clt *Client) DeleteMassMassage(msgId int64) (err error) {
	var result struct {
		Error
	}
	msg := struct {
		MsgId int64 `json:"msg_id"`
	}{
		MsgId: msgId,
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/message/mass/delete?access_token="
	if err = clt.PostJSON(incompleteURL, msg, &result); err != nil {
		return
	}
	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	return
}
