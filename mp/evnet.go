package mp

import (
	"fmt"
	"strings"
)

const (
	EventTypeClick = "CLICK" // 点击菜单拉取消息时的事件推送
	EventTypeView  = "VIEW"  // 点击菜单跳转链接时的事件推送

	EventTypeSubscribe   = "subscribe"   // 订阅, 包括点击订阅和扫描二维码
	EventTypeUnsubscribe = "unsubscribe" // 取消订阅
	EventTypeScan        = "SCAN"        // 已经订阅的用户扫描二维码事件
	EventTypeLocation    = "LOCATION"    // 上报地理位置事件

	// 请注意, 下面的事件仅支持微信iPhone5.4.1以上版本, 和Android5.4以上版本的微信用户,
	// 旧版本微信用户点击后将没有回应, 开发者也不能正常接收到事件推送.
	EventTypeScanCodePush    = "scancode_push"      // scancode_push：扫码推事件的事件推送
	EventTypeScanCodeWaitMsg = "scancode_waitmsg"   // scancode_waitmsg：扫码推事件且弹出“消息接收中”提示框的事件推送
	EventTypePicSysPhoto     = "pic_sysphoto"       // pic_sysphoto：弹出系统拍照发图的事件推送
	EventTypePicPhotoOrAlbum = "pic_photo_or_album" // pic_photo_or_album：弹出拍照或者相册发图的事件推送
	EventTypePicWeixin       = "pic_weixin"         // pic_weixin：弹出微信相册发图器的事件推送
	EventTypeLocationSelect  = "location_select"    // location_select：弹出地理位置选择器的事件推送
)

// 关注事件(普通关注)
type SubscribeEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	Event string `xml:"Event" json:"Event"` // 事件类型，subscribe(订阅)
}

func GetSubscribeEvent(msg *MixedMessage) *SubscribeEvent {
	return &SubscribeEvent{
		CommonMessageHeader: msg.CommonMessageHeader,
		Event:               msg.Event,
	}
}

// 取消关注
type UnsubscribeEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	Event string `xml:"Event" json:"Event"` // 事件类型，unsubscribe(取消订阅)
}

func GetUnsubscribeEvent(msg *MixedMessage) *UnsubscribeEvent {
	return &UnsubscribeEvent{
		CommonMessageHeader: msg.CommonMessageHeader,
		Event:               msg.Event,
	}
}

// 用户未关注时，扫描带参数二维码进行关注后的事件推送
type SubscribeByScanEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	Event    string `xml:"Event"    json:"Event"`    // 事件类型，subscribe
	EventKey string `xml:"EventKey" json:"EventKey"` // 事件KEY值，qrscene_为前缀，后面为二维码的参数值
	Ticket   string `xml:"Ticket"   json:"Ticket"`   // 二维码的ticket，可用来换取二维码图片
}

// 获取二维码参数
func (event *SubscribeByScanEvent) Scene() (scene string, err error) {
	const prefix = "qrscene_"
	if !strings.HasPrefix(event.EventKey, prefix) {
		err = fmt.Errorf("EventKey 应该以 %q 为前缀: %q", prefix, event.EventKey)
		return
	}
	scene = event.EventKey[len(prefix):]
	return
}

func GetSubscribeByScanEvent(msg *MixedMessage) *SubscribeByScanEvent {
	return &SubscribeByScanEvent{
		CommonMessageHeader: msg.CommonMessageHeader,
		Event:               msg.Event,
		EventKey:            msg.EventKey,
		Ticket:              msg.Ticket,
	}
}

// 用户已关注时，扫描带参数二维码的事件推送
type ScanEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	Event    string `xml:"Event"    json:"Event"`    // 事件类型，SCAN
	EventKey string `xml:"EventKey" json:"EventKey"` // 事件KEY值，是一个32位无符号整数，即创建二维码时的二维码scene_id
	Ticket   string `xml:"Ticket"   json:"Ticket"`   // 二维码的ticket，可用来换取二维码图片
}

func GetScanEvent(msg *MixedMessage) *ScanEvent {
	return &ScanEvent{
		CommonMessageHeader: msg.CommonMessageHeader,
		Event:               msg.Event,
		EventKey:            msg.EventKey,
		Ticket:              msg.Ticket,
	}
}

// 上报地理位置事件
type LocationEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	Event     string  `xml:"Event"     json:"Event"`     // 事件类型，LOCATION
	Latitude  float64 `xml:"Latitude"  json:"Latitude"`  // 地理位置纬度
	Longitude float64 `xml:"Longitude" json:"Longitude"` // 地理位置经度
	Precision float64 `xml:"Precision" json:"Precision"` // 地理位置精度
}

func GetLocationEvent(msg *MixedMessage) *LocationEvent {
	return &LocationEvent{
		CommonMessageHeader: msg.CommonMessageHeader,
		Event:               msg.Event,
		Latitude:            msg.Latitude,
		Longitude:           msg.Longitude,
		Precision:           msg.Precision,
	}
}

// 点击菜单拉取消息时的事件推送
type ClickEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	Event    string `xml:"Event"    json:"Event"`    // 事件类型, CLICK
	EventKey string `xml:"EventKey" json:"EventKey"` // 事件KEY值, 与自定义菜单接口中KEY值对应
}

func GetClickEvent(msg *MixedMessage) *ClickEvent {
	return &ClickEvent{
		CommonMessageHeader: msg.CommonMessageHeader,
		Event:               msg.Event,
		EventKey:            msg.EventKey,
	}
}

// 点击菜单跳转链接时的事件推送
type ViewEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	Event    string `xml:"Event"    json:"Event"`    // 事件类型, VIEW
	EventKey string `xml:"EventKey" json:"EventKey"` // 事件KEY值, 设置的跳转URL
}

func GetViewEvent(msg *MixedMessage) *ViewEvent {
	return &ViewEvent{
		CommonMessageHeader: msg.CommonMessageHeader,
		Event:               msg.Event,
		EventKey:            msg.EventKey,
	}
}

// scancode_push：扫码推事件的事件推送
type ScanCodePushEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	Event    string `xml:"Event"    json:"Event"`    // 事件类型, scancode_push
	EventKey string `xml:"EventKey" json:"EventKey"` // 事件KEY值, 由开发者在创建菜单时设定

	ScanCodeInfo struct {
		ScanType   string `xml:"ScanType"   json:"ScanType"`   // 扫描类型, 一般是qrcode
		ScanResult string `xml:"ScanResult" json:"ScanResult"` // 扫描结果, 即二维码对应的字符串信息
	} `xml:"ScanCodeInfo" json:"ScanCodeInfo"`
}

func GetScanCodePushEvent(msg *MixedMessage) *ScanCodePushEvent {
	return &ScanCodePushEvent{
		CommonMessageHeader: msg.CommonMessageHeader,
		Event:               msg.Event,
		EventKey:            msg.EventKey,
		ScanCodeInfo:        msg.ScanCodeInfo,
	}
}

// scancode_waitmsg：扫码推事件且弹出“消息接收中”提示框的事件推送
type ScanCodeWaitMsgEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	Event    string `xml:"Event"    json:"Event"`    // 事件类型, scancode_waitmsg
	EventKey string `xml:"EventKey" json:"EventKey"` // 事件KEY值, 由开发者在创建菜单时设定

	ScanCodeInfo struct {
		ScanType   string `xml:"ScanType"   json:"ScanType"`   // 扫描类型, 一般是qrcode
		ScanResult string `xml:"ScanResult" json:"ScanResult"` // 扫描结果, 即二维码对应的字符串信息
	} `xml:"ScanCodeInfo" json:"ScanCodeInfo"`
}

func GetScanCodeWaitMsgEvent(msg *MixedMessage) *ScanCodeWaitMsgEvent {
	return &ScanCodeWaitMsgEvent{
		CommonMessageHeader: msg.CommonMessageHeader,
		Event:               msg.Event,
		EventKey:            msg.EventKey,
		ScanCodeInfo:        msg.ScanCodeInfo,
	}
}

// pic_sysphoto：弹出系统拍照发图的事件推送
type PicSysPhotoEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	Event    string `xml:"Event"    json:"Event"`    // 事件类型, pic_sysphoto
	EventKey string `xml:"EventKey" json:"EventKey"` // 事件KEY值, 由开发者在创建菜单时设定

	SendPicsInfo struct {
		Count   int `xml:"Count" json:"Count"` // 发送的图片数量
		PicList []struct {
			PicMD5Sum string `xml:"PicMd5Sum" json:"PicMd5Sum"` // 图片的MD5值, 开发者若需要, 可用于验证接收到图片
		} `xml:"PicList>item,omitempty" json:"PicList,omitempty"` // 图片列表
	} `xml:"SendPicsInfo" json:"SendPicsInfo"` // 发送的图片信息
}

func GetPicSysPhotoEvent(msg *MixedMessage) *PicSysPhotoEvent {
	return &PicSysPhotoEvent{
		CommonMessageHeader: msg.CommonMessageHeader,
		Event:               msg.Event,
		EventKey:            msg.EventKey,
		SendPicsInfo:        msg.SendPicsInfo,
	}
}

// pic_photo_or_album：弹出拍照或者相册发图的事件推送
type PicPhotoOrAlbumEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	Event    string `xml:"Event"    json:"Event"`    // 事件类型, pic_photo_or_album
	EventKey string `xml:"EventKey" json:"EventKey"` // 事件KEY值, 由开发者在创建菜单时设定

	SendPicsInfo struct {
		Count   int `xml:"Count" json:"Count"` // 发送的图片数量
		PicList []struct {
			PicMD5Sum string `xml:"PicMd5Sum" json:"PicMd5Sum"` // 图片的MD5值, 开发者若需要, 可用于验证接收到图片
		} `xml:"PicList>item,omitempty" json:"PicList,omitempty"` // 图片列表
	} `xml:"SendPicsInfo" json:"SendPicsInfo"` // 发送的图片信息
}

func GetPicPhotoOrAlbumEvent(msg *MixedMessage) *PicPhotoOrAlbumEvent {
	return &PicPhotoOrAlbumEvent{
		CommonMessageHeader: msg.CommonMessageHeader,
		Event:               msg.Event,
		EventKey:            msg.EventKey,
		SendPicsInfo:        msg.SendPicsInfo,
	}
}

// pic_weixin：弹出微信相册发图器的事件推送
type PicWeixinEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	Event    string `xml:"Event"    json:"Event"`    // 事件类型, pic_weixin
	EventKey string `xml:"EventKey" json:"EventKey"` // 事件KEY值, 由开发者在创建菜单时设定

	SendPicsInfo struct {
		Count   int `xml:"Count" json:"Count"` // 发送的图片数量
		PicList []struct {
			PicMD5Sum string `xml:"PicMd5Sum" json:"PicMd5Sum"` // 图片的MD5值, 开发者若需要, 可用于验证接收到图片
		} `xml:"PicList>item,omitempty" json:"PicList,omitempty"` // 图片列表
	} `xml:"SendPicsInfo" json:"SendPicsInfo"` // 发送的图片信息
}

func GetPicWeixinEvent(msg *MixedMessage) *PicWeixinEvent {
	return &PicWeixinEvent{
		CommonMessageHeader: msg.CommonMessageHeader,
		Event:               msg.Event,
		EventKey:            msg.EventKey,
		SendPicsInfo:        msg.SendPicsInfo,
	}
}

// location_select：弹出地理位置选择器的事件推送
type LocationSelectEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	Event    string `xml:"Event"    json:"Event"`    // 事件类型, location_select
	EventKey string `xml:"EventKey" json:"EventKey"` // 事件KEY值, 由开发者在创建菜单时设定

	SendLocationInfo struct {
		LocationX float64 `xml:"Location_X" json:"Location_X"` // 地理位置纬度
		LocationY float64 `xml:"Location_Y" json:"Location_Y"` // 地理位置经度
		Scale     int     `xml:"Scale"      json:"Scale"`      // 精度, 可理解为精度或者比例尺、越精细的话 scale越高
		Label     string  `xml:"Label"      json:"Label"`      // 地理位置的字符串信息
		Poiname   string  `xml:"Poiname"    json:"Poiname"`    // 朋友圈POI的名字, 可能为空
	} `xml:"SendLocationInfo" json:"SendLocationInfo"` // 发送的位置信息
}

func GetLocationSelectEvent(msg *MixedMessage) *LocationSelectEvent {
	return &LocationSelectEvent{
		CommonMessageHeader: msg.CommonMessageHeader,
		Event:               msg.Event,
		EventKey:            msg.EventKey,
		SendLocationInfo:    msg.SendLocationInfo,
	}
}
