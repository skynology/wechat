package mp

import (
	"errors"
	"fmt"
)

// 文本消息
type ResText struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	Content string `xml:"Content" json:"Content"` // 回复的消息内容, 支持换行符
}

// 新建文本消息
//  NOTE: content 支持换行符
func NewResText(to, from string, timestamp int64, content string) (text *ResText) {
	return &ResText{
		CommonMessageHeader: CommonMessageHeader{
			ToUserName:   to,
			FromUserName: from,
			CreateTime:   timestamp,
			MsgType:      MsgTypeText,
		},
		Content: content,
	}
}

// 图片消息
type ResImage struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	Image struct {
		MediaId string `xml:"MediaId" json:"MediaId"` // MediaId 通过上传多媒体文件得到
	} `xml:"Image" json:"Image"`
}

// 新建图片消息
//  MediaId 通过上传多媒体文件得到
func NewResImage(to, from string, timestamp int64, mediaId string) (image *ResImage) {
	image = &ResImage{
		CommonMessageHeader: CommonMessageHeader{
			ToUserName:   to,
			FromUserName: from,
			CreateTime:   timestamp,
			MsgType:      MsgTypeImage,
		},
	}
	image.Image.MediaId = mediaId
	return
}

// 语音消息
type ResVoice struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	Voice struct {
		MediaId string `xml:"MediaId" json:"MediaId"` // MediaId 通过上传多媒体文件得到
	} `xml:"Voice" json:"Voice"`
}

// 新建语音消息
//  MediaId 通过上传多媒体文件得到
func NewResVoice(to, from string, timestamp int64, mediaId string) (voice *ResVoice) {
	voice = &ResVoice{
		CommonMessageHeader: CommonMessageHeader{
			ToUserName:   to,
			FromUserName: from,
			CreateTime:   timestamp,
			MsgType:      MsgTypeVoice,
		},
	}
	voice.Voice.MediaId = mediaId
	return
}

// 视频消息
type ResVideo struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	Video struct {
		MediaId     string `xml:"MediaId"               json:"MediaId"`               // MediaId 通过上传多媒体文件得到
		Title       string `xml:"Title,omitempty"       json:"Title,omitempty"`       // 视频消息的标题
		Description string `xml:"Description,omitempty" json:"Description,omitempty"` // 视频消息的描述
	} `xml:"Video" json:"Video"`
}

// 新建视频消息
//  MediaId 通过上传多媒体文件得到
//  title, description 可以为 ""
func NewResVideo(to, from string, timestamp int64, mediaId, title, description string) (video *ResVideo) {
	video = &ResVideo{
		CommonMessageHeader: CommonMessageHeader{
			ToUserName:   to,
			FromUserName: from,
			CreateTime:   timestamp,
			MsgType:      MsgTypeVideo,
		},
	}
	video.Video.MediaId = mediaId
	video.Video.Title = title
	video.Video.Description = description
	return
}

// 音乐消息
type ResMusic struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	Music struct {
		Title        string `xml:"Title,omitempty"        json:"Title,omitempty"`        // 音乐标题
		Description  string `xml:"Description,omitempty"  json:"Description,omitempty"`  // 音乐描述
		MusicURL     string `xml:"MusicUrl,omitempty"     json:"MusicUrl,omitempty"`     // 音乐链接
		HQMusicURL   string `xml:"HQMusicUrl,omitempty"   json:"HQMusicUrl,omitempty"`   // 高质量音乐链接, WIFI环境优先使用该链接播放音乐
		ThumbMediaId string `xml:"ThumbMediaId,omitempty" json:"ThumbMediaId,omitempty"` // 缩略图的媒体id, 通过上传多媒体文件得到
	} `xml:"Music" json:"Music"`
}

// 新建音乐消息
//  thumbMediaId 通过上传多媒体文件得到
//  title, description 可以为 ""
func NewResMusic(to, from string, timestamp int64, thumbMediaId, musicURL,
	HQMusicURL, title, description string) (music *ResMusic) {

	music = &ResMusic{
		CommonMessageHeader: CommonMessageHeader{
			ToUserName:   to,
			FromUserName: from,
			CreateTime:   timestamp,
			MsgType:      MsgTypeMusic,
		},
	}
	music.Music.Title = title
	music.Music.Description = description
	music.Music.MusicURL = musicURL
	music.Music.HQMusicURL = HQMusicURL
	music.Music.ThumbMediaId = thumbMediaId
	return
}

// 图文消息里的 Article
type ResArticle struct {
	Title       string `xml:"Title,omitempty"       json:"Title,omitempty"`       // 图文消息标题
	Description string `xml:"Description,omitempty" json:"Description,omitempty"` // 图文消息描述
	PicURL      string `xml:"PicUrl,omitempty"      json:"PicUrl,omitempty"`      // 图片链接, 支持JPG, PNG格式, 较好的效果为大图360*200, 小图200*200
	URL         string `xml:"Url,omitempty"         json:"Url,omitempty"`         // 点击图文消息跳转链接
}

// 图文消息
type ResNews struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	ArticleCount int       `xml:"ArticleCount"            json:"ArticleCount"`       // 图文消息个数, 限制为10条以内
	Articles     []Article `xml:"Articles>item,omitempty" json:"Articles,omitempty"` // 多条图文消息信息, 默认第一个item为大图, 注意, 如果图文数超过10, 则将会无响应
}

// NOTE: articles 的长度不能超过 NewsArticleCountLimit
func NewResNews(to, from string, timestamp int64, articles []Article) (news *ResNews) {
	news = &ResNews{
		CommonMessageHeader: CommonMessageHeader{
			ToUserName:   to,
			FromUserName: from,
			CreateTime:   timestamp,
			MsgType:      MsgTypeNews,
		},
	}
	news.Articles = articles
	news.ArticleCount = len(articles)
	return
}

// 检查 News 是否有效，有效返回 nil，否则返回错误信息
func (news *ResNews) CheckValid() (err error) {
	n := len(news.Articles)
	if n != news.ArticleCount {
		err = fmt.Errorf("图文消息的 ArticleCount == %d, 实际文章个数为 %d", news.ArticleCount, n)
		return
	}
	if n <= 0 {
		err = errors.New("图文消息里没有文章")
		return
	}
	if n > NewsArticleCountLimit {
		err = fmt.Errorf("图文消息的文章个数不能超过 %d, 现在为 %d", NewsArticleCountLimit, n)
		return
	}
	return
}

type TransInfo struct {
	KfAccount string `xml:"KfAccount" json:"KfAccount"`
}

// 将消息转发到多客服
type TransferToCustomerService struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	*TransInfo `xml:"TransInfo,omitempty" json:"TransInfo,omitempty"`
}

// 如果不指定客服则 kfAccount 留空.
func NewTransferToCustomerService(to, from string, timestamp int64, kfAccount string) (msg *TransferToCustomerService) {
	msg = &TransferToCustomerService{
		CommonMessageHeader: CommonMessageHeader{
			ToUserName:   to,
			FromUserName: from,
			CreateTime:   timestamp,
			MsgType:      MsgTypeTransferCustomerService,
		},
	}

	if kfAccount != "" {
		msg.TransInfo = &TransInfo{
			KfAccount: kfAccount,
		}
	}
	return
}
