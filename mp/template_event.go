package mp

const (
	EventTypeTemplateSendJobFinish = "TEMPLATESENDJOBFINISH"
)

const (
	TemplateSendStatusSuccess            = "success"               // 送达成功时
	TemplateSendStatusFailedUserBlock    = "failed:user block"     // 送达由于用户拒收（用户设置拒绝接收公众号消息）而失败
	TemplateSendStatusFailedSystemFailed = "failed: system failed" // 送达由于其他原因失败
)

// 在模版消息发送任务完成后，微信服务器会将是否送达成功作为通知，发送到开发者中心中填写的服务器配置地址中。
type TemplateSendJobFinishEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	Event  string `xml:"Event"  json:"Event"` // 事件信息，此处为 TEMPLATESENDJOBFINISH
	MsgId  int64  `xml:"MsgId"  json:"MsgId"` // 模板消息ID
	Status string `xml:"Status" json:"Status"`
}

func GetTemplateSendJobFinishEvent(msg *MixedMessage) *TemplateSendJobFinishEvent {
	return &TemplateSendJobFinishEvent{
		CommonMessageHeader: msg.CommonMessageHeader,
		Event:               msg.Event,
		MsgId:               msg.MsgID, // NOTE
		Status:              msg.Status,
	}
}
