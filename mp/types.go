package mp

// 微信服务器推送过来的消息(事件)通用的消息头
type CommonMessageHeader struct {
	ToUserName   string `xml:"ToUserName"   json:"ToUserName"`
	FromUserName string `xml:"FromUserName" json:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"   json:"CreateTime"`
	MsgType      string `xml:"MsgType"      json:"MsgType"`
}

// 微信服务器推送过来的消息(事件)的合集.
type MixedMessage struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	// fuck, MsgId != MsgID
	MsgId int64 `xml:"MsgId" json:"MsgId"`
	MsgID int64 `xml:"MsgID" json:"MsgID"`

	Content      string  `xml:"Content"      json:"Content"`
	MediaId      string  `xml:"MediaId"      json:"MediaId"`
	PicURL       string  `xml:"PicUrl"       json:"PicUrl"`
	Format       string  `xml:"Format"       json:"Format"`
	Recognition  string  `xml:"Recognition"  json:"Recognition"`
	ThumbMediaId string  `xml:"ThumbMediaId" json:"ThumbMediaId"`
	LocationX    float64 `xml:"Location_X"   json:"Location_X"`
	LocationY    float64 `xml:"Location_Y"   json:"Location_Y"`
	Scale        int     `xml:"Scale"        json:"Scale"`
	Label        string  `xml:"Label"        json:"Label"`
	Title        string  `xml:"Title"        json:"Title"`
	Description  string  `xml:"Description"  json:"Description"`
	URL          string  `xml:"Url"          json:"Url"`

	Event    string `xml:"Event"    json:"Event"`
	EventKey string `xml:"EventKey" json:"EventKey"`

	ScanCodeInfo struct {
		ScanType   string `xml:"ScanType"   json:"ScanType"`
		ScanResult string `xml:"ScanResult" json:"ScanResult"`
	} `xml:"ScanCodeInfo" json:"ScanCodeInfo"`

	SendPicsInfo struct {
		Count   int `xml:"Count" json:"Count"`
		PicList []struct {
			PicMD5Sum string `xml:"PicMd5Sum" json:"PicMd5Sum"`
		} `xml:"PicList>item,omitempty" json:"PicList,omitempty"`
	} `xml:"SendPicsInfo" json:"SendPicsInfo"`

	SendLocationInfo struct {
		LocationX float64 `xml:"Location_X" json:"Location_X"`
		LocationY float64 `xml:"Location_Y" json:"Location_Y"`
		Scale     int     `xml:"Scale"      json:"Scale"`
		Label     string  `xml:"Label"      json:"Label"`
		Poiname   string  `xml:"Poiname"    json:"Poiname"`
	} `xml:"SendLocationInfo" json:"SendLocationInfo"`

	Ticket      string  `xml:"Ticket"      json:"Ticket"`
	Latitude    float64 `xml:"Latitude"    json:"Latitude"`
	Longitude   float64 `xml:"Longitude"   json:"Longitude"`
	Precision   float64 `xml:"Precision"   json:"Precision"`
	Status      string  `xml:"Status"      json:"Status"`
	TotalCount  int     `xml:"TotalCount"  json:"TotalCount"`
	FilterCount int     `xml:"FilterCount" json:"FilterCount"`
	SentCount   int     `xml:"SentCount"   json:"SentCount"`
	ErrorCount  int     `xml:"ErrorCount"  json:"ErrorCount"`
	OrderId     string  `xml:"OrderId"     json:"OrderId"`
	OrderStatus int     `xml:"OrderStatus" json:"OrderStatus"`
	ProductId   string  `xml:"ProductId"   json:"ProductId"`
	SKUInfo     string  `xml:"SkuInfo"     json:"SkuInfo"`

	// card
	CardId         string `xml:"CardId"         json:"CardId"`
	IsGiveByFriend int    `xml:"IsGiveByFriend" json:"IsGiveByFriend"`
	FriendUserName string `xml:"FriendUserName" json:"FriendUserName"`
	UserCardCode   string `xml:"UserCardCode"   json:"UserCardCode"`
	OuterId        int64  `xml:"OuterId"        json:"OuterId"`

	// poi
	UniqId string `xml:"UniqId" json:"UniqId"`
	PoiId  string `xml:"PoiId"  json:"PoiId"`
	Result string `xml:"Result" json:"Result"`
	Msg    string `xml:"Msg"    json:"Msg"`
}
