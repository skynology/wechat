// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package mp

import (
	"errors"
	"fmt"
)

type PoiAddParameters struct {
	BaseInfo struct {
		Sid          string   `json:"sid,omitempty"`           // 可选, 商户自己的id，用于后续审核通过收到poi_id 的通知时，做对应关系。请商户自己保证唯一识别性
		BusinessName string   `json:"business_name,omitempty"` // 必须, 门店名称（仅为商户名，如：国美、麦当劳，不应包含地区、店号等信息，错误示例：北京国美）
		BranchName   string   `json:"branch_name,omitempty"`   // 可选, 分店名称（不应包含地区信息、不应与门店名重复，错误示例：北京王府井店）
		Province     string   `json:"province,omitempty"`      // 必须, 门店所在的省份（直辖市填城市名,如：北京市）
		City         string   `json:"city,omitempty"`          // 必须, 门店所在的城市
		District     string   `json:"district,omitempty"`      // 可选, 门店所在地区
		Address      string   `json:"address,omitempty"`       // 必须, 门店所在的详细街道地址（不要填写省市信息）
		Telephone    string   `json:"telephone,omitempty"`     // 必须, 门店的电话（纯数字，区号、分机号均由“-”隔开）
		Categories   []string `json:"categories,omitempty"`    // 必须, 门店的类型（详细分类参见分类附表，不同级分类用“,”隔开，如：美食，川菜，火锅）
		OffsetType   int      `json:"offset_type"`             // 必须, 坐标类型，1 为火星坐标（目前只能选1）
		Longitude    float64  `json:"longitude"`               // 必须, 门店所在地理位置的经度
		Latitude     float64  `json:"latitude"`                // 必须, 门店所在地理位置的纬度（经纬度均为火星坐标，最好选用腾讯地图标记的坐标）
		PhotoList    []string `json:"photo_list,omitempty"`    // 必须, 图片列表，url 形式，可以有多张图片，尺寸为640*340px。必须为上一接口生成的url
		Recommend    string   `json:"recommend,omitempty"`     // 可选, 推荐品，餐厅可为推荐菜；酒店为推荐套房；景点为推荐游玩景点等，针对自己行业的推荐内容
		Special      string   `json:"special,omitempty"`       // 必须, 特色服务，如免费wifi，免费停车，送货上门等商户能提供的特色功能或服务
		Introduction string   `json:"introduction,omitempty"`  // 可选, 商户简介，主要介绍商户信息等
		OpenTime     string   `json:"open_time,omitempty"`     // 必须, 营业时间，24 小时制表示，用“-”连接，如8:00-20:00
		AvgPrice     int      `json:"avg_price,omitempty"`     // 可选, 人均价格，大于0 的整数
	} `json:"base_info"`
}

// 创建门店.
func (clt *Client) PoiAdd(para *PoiAddParameters) (err error) {
	if para == nil {
		return errors.New("nil PoiAddParameters")
	}

	var request = struct {
		*PoiAddParameters `json:"business,omitempty"`
	}{
		PoiAddParameters: para,
	}

	var result Error

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/poi/addpoi?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result
		return
	}
	return
}

type PoiBrief struct {
	BaseInfo struct {
		PoiId          string `json:"poi_id,omitempty"`        // Poi 的id, 只有审核通过后才有
		AvailableState int    `json:"available_state"`         // 门店是否可用状态。1 表示系统错误、2 表示审核中、3 审核通过、4 审核驳回。当该字段为1、2、4 状态时，poi_id 为空
		Sid            string `json:"sid,omitempty"`           // 商户自己的id，用于后续审核通过收到poi_id 的通知时，做对应关系。请商户自己保证唯一识别性
		BusinessName   string `json:"business_name,omitempty"` // 门店名称（仅为商户名，如：国美、麦当劳，不应包含地区、店号等信息，错误示例：北京国美）
		BranchName     string `json:"branch_name,omitempty"`   // 分店名称（不应包含地区信息、不应与门店名重复，错误示例：北京王府井店）
		Address        string `json:"address,omitempty"`       // 门店所在的详细街道地址（不要填写省市信息）
	} `json:"base_info"`
}

// 查询门店列表.
//  begin: 开始位置，0 即为从第一条开始查询
//  limit: 返回数据条数，最大允许50，默认为20
func (clt *Client) PoiList(begin, limit int) (list []PoiBrief, totalCount int, err error) {
	if begin < 0 {
		err = fmt.Errorf("invalid begin: %d", begin)
		return
	}
	if limit < 0 {
		err = fmt.Errorf("invalid limit: %d", limit)
		return
	}

	var request = struct {
		Begin int `json:"begin"`
		Limit int `json:"limit,omitempty"`
	}{
		Begin: begin,
		Limit: limit,
	}

	var result struct {
		Error
		PoiList    []PoiBrief `json:"business_list"`
		TotalCount int        `json:"total_count"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/poi/getpoilist?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.PoiList
	totalCount = result.TotalCount
	return
}

type Poi struct {
	BaseInfo struct {
		PoiId          string `json:"poi_id,omitempty"` // Poi 的id, 只有审核通过后才有
		AvailableState int    `json:"available_state"`  // 门店是否可用状态。1 表示系统错误、2 表示审核中、3 审核通过、4 审核驳回。当该字段为1、2、4 状态时，poi_id 为空
		UpdateStatus   int    `json:"update_status"`    // 扩展字段是否正在更新中。1 表示扩展字段正在更新中，尚未生效，不允许再次更新； 0 表示扩展字段没有在更新中或更新已生效，可以再次更新

		Sid          string   `json:"sid,omitempty"`           // 商户自己的id，用于后续审核通过收到poi_id 的通知时，做对应关系。请商户自己保证唯一识别性
		BusinessName string   `json:"business_name,omitempty"` // 门店名称（仅为商户名，如：国美、麦当劳，不应包含地区、店号等信息，错误示例：北京国美）
		BranchName   string   `json:"branch_name,omitempty"`   // 分店名称（不应包含地区信息、不应与门店名重复，错误示例：北京王府井店）
		Province     string   `json:"province,omitempty"`      // 门店所在的省份（直辖市填城市名,如：北京市）
		City         string   `json:"city,omitempty"`          // 门店所在的城市
		District     string   `json:"district,omitempty"`      // 门店所在地区
		Address      string   `json:"address,omitempty"`       // 门店所在的详细街道地址（不要填写省市信息）
		Telephone    string   `json:"telephone,omitempty"`     // 门店的电话（纯数字，区号、分机号均由“-”隔开）
		Categories   []string `json:"categories,omitempty"`    // 门店的类型（详细分类参见分类附表，不同级分类用“,”隔开，如：美食，川菜，火锅）
		OffsetType   int      `json:"offset_type"`             // 坐标类型，1 为火星坐标（目前只能选1）
		Longitude    float64  `json:"longitude"`               // 门店所在地理位置的经度
		Latitude     float64  `json:"latitude"`                // 门店所在地理位置的纬度（经纬度均为火星坐标，最好选用腾讯地图标记的坐标）
		PhotoList    []string `json:"photo_list,omitempty"`    // 图片列表，url 形式，可以有多张图片，尺寸为640*340px。必须为上一接口生成的url
		Recommend    string   `json:"recommend,omitempty"`     // 推荐品，餐厅可为推荐菜；酒店为推荐套房；景点为推荐游玩景点等，针对自己行业的推荐内容
		Special      string   `json:"special,omitempty"`       // 特色服务，如免费wifi，免费停车，送货上门等商户能提供的特色功能或服务
		Introduction string   `json:"introduction,omitempty"`  // 商户简介，主要介绍商户信息等
		OpenTime     string   `json:"open_time,omitempty"`     // 营业时间，24 小时制表示，用“-”连接，如8:00-20:00
		AvgPrice     int      `json:"avg_price,omitempty"`     // 人均价格，大于0 的整数
	} `json:"base_info"`
}

// 查询门店信息.
func (clt *Client) PoiGet(poiId string) (poi *Poi, err error) {
	var request = struct {
		PoiId string `json:"poi_id"`
	}{
		PoiId: poiId,
	}

	var result struct {
		Error
		Poi `json:"business"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/poi/getpoi?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	if result.Poi.BaseInfo.PoiId == "" {
		result.Poi.BaseInfo.PoiId = poiId
	}
	poi = &result.Poi
	return
}

type PoiUpdateParameters struct {
	BaseInfo struct {
		PoiId string `json:"poi_id,omitempty"`

		Telephone    string   `json:"telephone,omitempty"`    // 必须, 门店的电话（纯数字，区号、分机号均由“-”隔开）
		PhotoList    []string `json:"photo_list,omitempty"`   // 必须, 图片列表，url 形式，可以有多张图片，尺寸为640*340px。必须为上一接口生成的url
		Recommend    string   `json:"recommend,omitempty"`    // 可选, 推荐品，餐厅可为推荐菜；酒店为推荐套房；景点为推荐游玩景点等，针对自己行业的推荐内容
		Special      string   `json:"special,omitempty"`      // 必须, 特色服务，如免费wifi，免费停车，送货上门等商户能提供的特色功能或服务
		Introduction string   `json:"introduction,omitempty"` // 可选, 商户简介，主要介绍商户信息等
		OpenTime     string   `json:"open_time,omitempty"`    // 必须, 营业时间，24 小时制表示，用“-”连接，如8:00-20:00
		AvgPrice     int      `json:"avg_price,omitempty"`    // 可选, 人均价格，大于0 的整数
	} `json:"base_info"`
}

// 修改门店服务信息.
//  商户可以通过该接口，修改门店的服务信息，包括：图片列表、营业时间、推荐、特色服务、简
//  介、人均价格、电话7 个字段。目前基础字段包括（名称、坐标、地址等不可修改）
func (clt *Client) PoiUpdate(para *PoiUpdateParameters) (err error) {
	if para == nil {
		return errors.New("nil PoiUpdateParameters")
	}

	var request = struct {
		*PoiUpdateParameters `json:"business,omitempty"`
	}{
		PoiUpdateParameters: para,
	}

	var result Error

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/poi/updatepoi?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result
		return
	}
	return
}

// 删除门店.
func (clt *Client) PoiDelete(poiId string) (err error) {
	var request = struct {
		PoiId string `json:"poi_id"`
	}{
		PoiId: poiId,
	}

	var result Error

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/poi/delpoi?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result
		return
	}
	return
}

const (
	EventTypePoiCheckNotify = "poi_check_notify" // Poi 审核结果事件推送
)

// Poi 审核结果事件推送
type PoiCheckNotifyEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonMessageHeader

	Event  string `xml:"Event"  json:"Event"`  // 事件类型, poi_check_notify
	UniqId string `xml:"UniqId" json:"UniqId"` // 商户自己内部ID，即字段中的sid
	PoiId  string `xml:"PoiId"  json:"PoiId"`  // 微信的门店ID，微信内门店唯一标示ID
	Result string `xml:"Result" json:"Result"` // 审核结果，成功succ 或失败fail
	Msg    string `xml:"Msg"    json:"Msg"`    // 成功的通知信息，或审核失败的驳回理由
}

func GetPoiCheckNotifyEvent(msg *MixedMessage) *PoiCheckNotifyEvent {
	return &PoiCheckNotifyEvent{
		CommonMessageHeader: msg.CommonMessageHeader,
		Event:               msg.Event,
		UniqId:              msg.UniqId,
		PoiId:               msg.PoiId,
		Result:              msg.Result,
		Msg:                 msg.Msg,
	}
}
