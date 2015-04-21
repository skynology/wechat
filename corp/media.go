// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package corp

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
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
func (clt *Client) GetMediaDownloadURL(mediaId string) string {
	token, err := clt.Token()
	if err != nil {
		return ""
	}
	uri := "https://qyapi.weixin.qq.com/cgi-bin/media/get?media_id=" + url.QueryEscape(mediaId) +
		"&access_token=" + url.QueryEscape(token)

	return uri
}

func (clt *Client) UploadMediaFromReader(mediaType, filename string, reader io.Reader) (info *MediaInfo, err error) {
	var result struct {
		Error
		MediaInfo
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/media/upload?type=" +
		url.QueryEscape(mediaType) + "&access_token="
	if err = clt.UploadFromReader(incompleteURL, filename, reader, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	info = &result.MediaInfo
	return
}

// 下载多媒体到 io.Writer.
func (clt *Client) DownloadMediaToWriter(mediaId string, writer io.Writer) (err error) {
	token, err := clt.Token()
	if err != nil {
		return
	}

	hasRetried := false
RETRY:
	finalURL := "https://qyapi.weixin.qq.com/cgi-bin/media/get?media_id=" + url.QueryEscape(mediaId) +
		"&access_token=" + url.QueryEscape(token)

	httpResp, err := clt.httpClient.Get(finalURL)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("http.Status: %s", httpResp.Status)
	}

	ContentType, _, _ := mime.ParseMediaType(httpResp.Header.Get("Content-Type"))
	if ContentType != "text/plain" && ContentType != "application/json" {
		// 返回的是媒体流
		_, err = io.Copy(writer, httpResp.Body)
		return
	}

	// 返回的是错误信息
	var result Error
	if err = json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
		return
	}

	switch result.ErrCode {
	case ErrCodeOK:
		return // 基本不会出现
	case ErrCodeInvalidCredential, ErrCodeTimeout: // 失效(过期)重试一次
		if !hasRetried {
			hasRetried = true

			if token, err = clt.RefreshToken(); err != nil {
				return
			}
			goto RETRY
		}
		fallthrough
	default:
		err = &result
		return
	}
}
