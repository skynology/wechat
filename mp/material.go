package mp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
)

type News []Article

// 删除永久素材.
func (clt *Client) DeleteMaterial(mediaId string) (err error) {
	var request = struct {
		MediaId string `json:"media_id"`
	}{
		MediaId: mediaId,
	}

	var result Error

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/material/del_material?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result
		return
	}
	return
}

// 公众号永久素材总数信息
type MaterialCountInfo struct {
	VoiceCount int `json:"voice_count"`
	VideoCount int `json:"video_count"`
	ImageCount int `json:"image_count"`
	NewsCount  int `json:"news_count"`
}

// 获取素材总数.
func (clt *Client) GetMaterialCount() (info *MaterialCountInfo, err error) {
	var result struct {
		Error
		MaterialCountInfo
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/material/get_materialcount?access_token="
	if err = clt.GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	info = &result.MaterialCountInfo
	return
}

type MaterialInfo struct {
	MediaId    string `json:"media_id"`    // 素材id
	Name       string `json:"name"`        // 文件名称
	UpdateTime int64  `json:"update_time"` // 最后更新时间
}

// 获取素材列表.
//
//  materialType: 素材的类型，图片（image）、视频（video）、语音 （voice）
//  offset:       从全部素材的该偏移位置开始返回，0表示从第一个素材 返回
//  count:        返回素材的数量，取值在1到20之间
//
//  TotalCount:   该类型的素材的总数
//  ItemCount:    本次调用获取的素材的数量
//  Items:        本次调用获取的素材
func (clt *Client) BatchGetMaterial(materialType string, offset, count int) (TotalCount, ItemCount int, Items []MaterialInfo, err error) {
	var request = struct {
		MaterialType string `json:"type"`
		Offset       int    `json:"offset"`
		Count        int    `json:"count"`
	}{
		MaterialType: materialType,
		Offset:       offset,
		Count:        count,
	}

	var result struct {
		Error
		TotalCount int            `json:"total_count"`
		ItemCount  int            `json:"item_count"`
		Items      []MaterialInfo `json:"item"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/material/batchget_material?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	TotalCount = result.TotalCount
	ItemCount = result.ItemCount
	Items = result.Items
	return
}

// 新增永久图文素材.
func (clt *Client) AddMeterialNews(news News) (mediaId string, err error) {
	if len(news) == 0 {
		err = errors.New("图文素材是空的")
		return
	}
	if len(news) > NewsArticleCountLimit {
		err = fmt.Errorf("图文素材的文章个数不能超过 %d, 现在为 %d", NewsArticleCountLimit, len(news))
		return
	}

	var request = struct {
		Articles []Article `json:"articles,omitempty"`
	}{
		Articles: news,
	}

	var result struct {
		Error
		MediaId string `json:"media_id"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/material/add_news?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	mediaId = result.MediaId
	return
}

// 修改永久图文素材.
//  再次fuck微信開發組, 這個api是猜的!
func (clt *Client) UpdateMeterialNews(mediaId string, index int, news News) (err error) {
	var request = struct {
		MediaId  string `json:"media_id"`
		Index    int    `json:"index"`
		Articles News   `json:"articles,omitempty"`
	}{
		MediaId:  mediaId,
		Index:    index,
		Articles: news,
	}

	var result Error

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/material/update_news?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result
		return
	}
	return
}

// 获取永久图文素材.
func (clt *Client) GetMeterialNews(mediaId string) (news News, err error) {
	var request = struct {
		MediaId string `json:"media_id"`
	}{
		MediaId: mediaId,
	}

	var result struct {
		Error
		Articles []Article `json:"news_item"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/material/get_material?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	news = result.Articles
	return
}

// 获取永久素材时, 若是视频类型, 返回如下格式
type Video struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DownURL     string `json:"down_url"`
}

// 获取永久图文素材.
func (clt *Client) GetMeterialVideo(mediaId string) (video Video, err error) {
	var request = struct {
		MediaId string `json:"media_id"`
	}{
		MediaId: mediaId,
	}

	var result struct {
		Error
		Video
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/material/get_material?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	video = result.Video
	return
}

const (
	MaterialTypeImage = "image"
	MaterialTypeVoice = "voice"
	MaterialTypeVideo = "video"
	MaterialTypeThumb = "thumb"
	MaterialTypeNews  = "news"
)

// 上传多媒体图片
func (clt *Client) UploadImage(filepath string) (mediaId string, err error) {
	return clt.uploadMaterial(MaterialTypeImage, filepath)
}

// 上传多媒体缩略图
func (clt *Client) UploadThumb(filepath string) (mediaId string, err error) {
	return clt.uploadMaterial(MaterialTypeThumb, filepath)
}

// 上传多媒体语音
func (clt *Client) UploadVoice(filepath string) (mediaId string, err error) {
	return clt.uploadMaterial(MaterialTypeVoice, filepath)
}

// 上传多媒体
func (clt *Client) uploadMaterial(materialType, _filepath string) (mediaId string, err error) {
	file, err := os.Open(_filepath)
	if err != nil {
		return
	}
	defer file.Close()

	return clt.uploadMaterialFromReader(materialType, filepath.Base(_filepath), file)
}

// 上传多媒体图片
//  NOTE: 参数 filename 不是文件路径, 是指定 multipart/form-data 里面文件名称
func (clt *Client) UploadMeterialImageFromReader(filename string, reader io.Reader) (mediaId string, err error) {
	if filename == "" {
		err = errors.New("empty filename")
		return
	}
	if reader == nil {
		err = errors.New("nil reader")
		return
	}
	return clt.uploadMaterialFromReader(MaterialTypeImage, filename, reader)
}

// 上传多媒体缩略图
//  NOTE: 参数 filename 不是文件路径, 是指定 multipart/form-data 里面文件名称
func (clt *Client) UploadMeteriralThumbFromReader(filename string, reader io.Reader) (mediaId string, err error) {
	if filename == "" {
		err = errors.New("empty filename")
		return
	}
	if reader == nil {
		err = errors.New("nil reader")
		return
	}
	return clt.uploadMaterialFromReader(MaterialTypeThumb, filename, reader)
}

// 上传多媒体语音
//  NOTE: 参数 filename 不是文件路径, 是指定 multipart/form-data 里面文件名称
func (clt *Client) UploadMeterialVoiceFromReader(filename string, reader io.Reader) (mediaId string, err error) {
	if filename == "" {
		err = errors.New("empty filename")
		return
	}
	if reader == nil {
		err = errors.New("nil reader")
		return
	}
	return clt.uploadMaterialFromReader(MaterialTypeVoice, filename, reader)
}

// 上传多媒体缩视频
//  NOTE: 参数 filename 不是文件路径, 是指定 multipart/form-data 里面文件名称
func (clt *Client) UploadMeterialVideoFromReader(filename string, reader io.Reader, title, introduction string) (mediaId string, err error) {
	if filename == "" {
		err = errors.New("empty filename")
		return
	}
	if reader == nil {
		err = errors.New("nil reader")
		return
	}
	return clt.uploadVideoFromReader(filename, reader, title, introduction)
}
func (clt *Client) uploadMaterialFromReader(materialType, filename string, reader io.Reader) (mediaId string, err error) {
	var result struct {
		Error
		MediaId string `json:"media_id"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/material/add_material?type=" +
		url.QueryEscape(materialType) + "&access_token="
	if err = clt.UploadFromReader(incompleteURL, "media", filename, reader, "", nil, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	mediaId = result.MediaId
	return
}
func (clt *Client) uploadVideoFromReader(filename string, reader io.Reader,
	title, introduction string) (mediaId string, err error) {

	var desc = struct {
		Title        string `json:"title"`
		Introduction string `json:"introduction"`
	}{
		Title:        title,
		Introduction: introduction,
	}

	descBytes, err := json.Marshal(&desc)
	if err != nil {
		return
	}

	var result struct {
		Error
		MediaId string `json:"media_id"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/material/add_material?type=video&access_token="
	if err = clt.UploadFromReader(incompleteURL, "media", filename, reader, "description", descBytes, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	mediaId = result.MediaId
	return
}
