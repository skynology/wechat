package mp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
)

const (
	NewsArticleCountLimit = 10 // 图文消息里文章的个数限制
)

const (
	MediaTypeImage = "image"
	MediaTypeVoice = "voice"
	MediaTypeVideo = "video"
	MediaTypeThumb = "thumb"
	MediaTypeNews  = "news"
)

type MediaInfo struct {
	MediaType string `json:"type"`       // 图片（image）、语音（voice）、视频（video）、缩略图（thumb）和 图文消息（news）
	MediaId   string `json:"media_id"`   // 媒体的唯一标识
	CreatedAt int64  `json:"created_at"` // 媒体创建的时间戳
}

// 图文消息里的 Article
type Article struct {
	ThumbMediaId     string `json:"thumb_media_id"`               // 图文消息缩略图的 media_id, 可以在上传多媒体文件接口中获得
	Title            string `json:"title"`                        // 图文消息的标题
	Author           string `json:"author,omitempty"`             // 图文消息的作者
	Digest           string `json:"digest,omitempty"`             // 图文消息的摘要
	Content          string `json:"content"`                      // 图文消息页面的内容，支持HTML标签
	ContentSourceURL string `json:"content_source_url,omitempty"` // 在图文消息页面点击“阅读原文”后的页面
	ShowCoverPic     int    `json:"show_cover_pic"`               // 是否显示封面, 1为显示, 0为不显示
}

func (article *Article) SetShowCoverPic(b bool) {
	if b {
		article.ShowCoverPic = 1
	} else {
		article.ShowCoverPic = 0
	}
}

// 获取临时素材下载地址
func (clt *Client) GetMediaDownloadURL(mediaId string) string {
	token, err := clt.Token()
	if err != nil {
		return ""
	}
	finalURL := "http://file.api.weixin.qq.com/cgi-bin/media/get?media_id=" + url.QueryEscape(mediaId) +
		"&access_token=" + url.QueryEscape(token)

	return finalURL
}

// 下载多媒体到文件.
func (clt *Client) DownloadMedia(mediaId, filepath string) (err error) {
	file, err := os.Create(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	return clt.downloadMediaToWriter(mediaId, file)
}

// 下载多媒体到 io.Writer.
func (clt *Client) DownloadMediaToWriter(mediaId string, writer io.Writer) error {
	if writer == nil {
		return errors.New("nil writer")
	}
	return clt.downloadMediaToWriter(mediaId, writer)
}

// 下载多媒体到 io.Writer.
func (clt *Client) downloadMediaToWriter(mediaId string, writer io.Writer) (err error) {
	token, err := clt.Token()
	if err != nil {
		return
	}

	hasRetried := false
RETRY:
	finalURL := "http://file.api.weixin.qq.com/cgi-bin/media/get?media_id=" + url.QueryEscape(mediaId) +
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

// 根据上传的缩略图媒体创建图文消息素材.
//  articles 的长度不能大于 NewsArticleCountLimit.
func (clt *Client) CreateNews(articles []Article) (info *MediaInfo, err error) {
	if len(articles) == 0 {
		err = errors.New("图文消息是空的")
		return
	}
	if len(articles) > NewsArticleCountLimit {
		err = fmt.Errorf("图文消息的文章个数不能超过 %d, 现在为 %d", NewsArticleCountLimit, len(articles))
		return
	}

	var request = struct {
		Articles []Article `json:"articles,omitempty"`
	}{
		Articles: articles,
	}

	var result struct {
		Error
		MediaInfo
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/media/uploadnews?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	info = &result.MediaInfo
	return
}

// 根据上传的视频文件 media_id 创建视频媒体, 群发视频消息应该用这个函数得到的 media_id.
//  NOTE: title, description 可以为空.
func (clt *Client) CreateVideo(mediaId, title, description string) (info *MediaInfo, err error) {
	var request = struct {
		MediaId     string `json:"media_id"`
		Title       string `json:"title,omitempty"`
		Description string `json:"description,omitempty"`
	}{
		MediaId:     mediaId,
		Title:       title,
		Description: description,
	}

	var result struct {
		Error
		MediaInfo
	}

	incompleteURL := "https://file.api.weixin.qq.com/cgi-bin/media/uploadvideo?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	info = &result.MediaInfo
	return
}

// 上传多媒体图片
//  NOTE: 参数 filename 不是文件路径, 是指定 multipart/form-data 里面文件名称
func (clt *Client) UploadImageFromReader(filename string, reader io.Reader) (info *MediaInfo, err error) {
	if filename == "" {
		err = errors.New("empty filename")
		return
	}
	if reader == nil {
		err = errors.New("nil reader")
		return
	}
	return clt.uploadMediaFromReader(MediaTypeImage, filename, reader)
}

// 上传多媒体语音
//  NOTE: 参数 filename 不是文件路径, 是指定 multipart/form-data 里面文件名称
func (clt *Client) UploadVoiceFromReader(filename string, reader io.Reader) (info *MediaInfo, err error) {
	if filename == "" {
		err = errors.New("empty filename")
		return
	}
	if reader == nil {
		err = errors.New("nil reader")
		return
	}
	return clt.uploadMediaFromReader(MediaTypeVoice, filename, reader)
}

// 上传多媒体视频
//  NOTE: 参数 filename 不是文件路径, 是指定 multipart/form-data 里面文件名称
func (clt *Client) UploadVideoFromReader(filename string, reader io.Reader) (info *MediaInfo, err error) {
	if filename == "" {
		err = errors.New("empty filename")
		return
	}
	if reader == nil {
		err = errors.New("nil reader")
		return
	}
	return clt.uploadMediaFromReader(MediaTypeVideo, filename, reader)
}
func (clt *Client) uploadMediaFromReader(mediaType, filename string, reader io.Reader) (info *MediaInfo, err error) {
	var result struct {
		Error
		MediaInfo
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/media/upload?type=" +
		url.QueryEscape(mediaType) + "&access_token="
	if err = clt.UploadFromReader(incompleteURL, "media", filename, reader, "", nil, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	info = &result.MediaInfo
	return
}

// 上传多媒体缩略图
//  NOTE: 参数 filename 不是文件路径, 是指定 multipart/form-data 里面文件名称
func (clt *Client) UploadThumbFromReader(filename string, reader io.Reader) (info *MediaInfo, err error) {
	if filename == "" {
		err = errors.New("empty filename")
		return
	}
	if reader == nil {
		err = errors.New("nil reader")
		return
	}
	return clt.uploadThumbFromReader(filename, reader)
}

func (clt *Client) uploadThumbFromReader(filename string, reader io.Reader) (info *MediaInfo, err error) {
	var result struct {
		Error
		MediaType string `json:"type"`
		MediaId   string `json:"thumb_media_id"`
		CreatedAt int64  `json:"created_at"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/media/upload?type=thumb&access_token="
	if err = clt.UploadFromReader(incompleteURL, "media", filename, reader, "", nil, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	info = &MediaInfo{
		MediaType: result.MediaType,
		MediaId:   result.MediaId,
		CreatedAt: result.CreatedAt,
	}
	return
}
