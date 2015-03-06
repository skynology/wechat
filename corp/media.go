// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package corp

import (
	"net/url"
)

const (
	MediaTypeImage = "image"
	MediaTypeVoice = "voice"
	MediaTypeVideo = "video"
	MediaTypeFile  = "file"
)

type MediaInfo struct {
	MediaType string `json:"type"`       // 媒体文件类型，分别有图片（image）、语音（voice）、视频（video）,普通文件(file)
	MediaId   string `json:"media_id"`   // 媒体文件上传后获取的唯一标识
	CreatedAt int64  `json:"created_at"` // 媒体文件上传时间戳
}

// 获取上media下载URL, 用于保存到文件服务器
func (clt *Client) GetMediaDownloadURL(mediaId string) (uri string, err error) {
	token, err := clt.Token()
	if err != nil {
		return
	}
	uri = "https://qyapi.weixin.qq.com/cgi-bin/media/get?media_id=" + url.QueryEscape(mediaId) +
		"&access_token=" + url.QueryEscape(token)

	return
}
