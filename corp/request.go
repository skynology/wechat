// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package corp

type ReqText struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	MsgId   int64  `xml:"MsgId"   json:"MsgId"`   // 消息id，64位整型
	Content string `xml:"Content" json:"Content"` // 文本消息内容
}

func GetText(msg *MixedMessage) *ReqText {
	return &ReqText{
		CommonMessageHeader: msg.CommonMessageHeader,
		MsgId:               msg.MsgId,
		Content:             msg.Content,
	}
}

type ReqImage struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	MsgId   int64  `xml:"MsgId"   json:"MsgId"`   // 消息id，64位整型
	MediaId string `xml:"MediaId" json:"MediaId"` // 图片媒体文件id，可以调用获取媒体文件接口拉取数据
	PicURL  string `xml:"PicUrl"  json:"PicUrl"`  // 图片链接
}

func GetImage(msg *MixedMessage) *ReqImage {
	return &ReqImage{
		CommonMessageHeader: msg.CommonMessageHeader,
		MsgId:               msg.MsgId,
		MediaId:             msg.MediaId,
		PicURL:              msg.PicURL,
	}
}

type ReqVoice struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	MsgId   int64  `xml:"MsgId"   json:"MsgId"`   // 消息id，64位整型
	MediaId string `xml:"MediaId" json:"MediaId"` // 语音媒体文件id，可以调用获取媒体文件接口拉取数据
	Format  string `xml:"Format"  json:"Format"`  // 语音格式，如amr，speex等
}

func GetVoice(msg *MixedMessage) *ReqVoice {
	return &ReqVoice{
		CommonMessageHeader: msg.CommonMessageHeader,
		MsgId:               msg.MsgId,
		MediaId:             msg.MediaId,
		Format:              msg.Format,
	}
}

type ReqVideo struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	MsgId        int64  `xml:"MsgId"        json:"MsgId"`        // 消息id，64位整型
	MediaId      string `xml:"MediaId"      json:"MediaId"`      // 视频媒体文件id，可以调用获取媒体文件接口拉取数据
	ThumbMediaId string `xml:"ThumbMediaId" json:"ThumbMediaId"` // 视频消息缩略图的媒体id，可以调用获取媒体文件接口拉取数据
}

func GetVideo(msg *MixedMessage) *ReqVideo {
	return &ReqVideo{
		CommonMessageHeader: msg.CommonMessageHeader,
		MsgId:               msg.MsgId,
		MediaId:             msg.MediaId,
		ThumbMediaId:        msg.ThumbMediaId,
	}
}

type ReqLocation struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	MsgId     int64   `xml:"MsgId"      json:"MsgId"`      // 消息id，64位整型
	LocationX float64 `xml:"Location_X" json:"Location_X"` // 地理位置纬度
	LocationY float64 `xml:"Location_Y" json:"Location_Y"` // 地理位置经度
	Scale     int     `xml:"Scale"      json:"Scale"`      // 地图缩放大小
	Label     string  `xml:"Label"      json:"Label"`      // 地理位置信息
}

func GetLocation(msg *MixedMessage) *ReqLocation {
	return &ReqLocation{
		CommonMessageHeader: msg.CommonMessageHeader,
		MsgId:               msg.MsgId,
		LocationX:           msg.LocationX,
		LocationY:           msg.LocationY,
		Scale:               msg.Scale,
		Label:               msg.Label,
	}
}
