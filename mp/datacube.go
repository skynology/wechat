// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package mp

import (
	"errors"
)

// 获取统计数据通用的请求结构.
type DataCubeParam struct {
	// 获取数据的起始日期, YYYY-MM-DD 格式.
	// begin_date 和 end_date 的差值需小于"最大时间跨度"(比如最大时间跨度为1时,
	// begin_date 和 end_date 的差值只能为0, 才能小于1), 否则会报错
	BeginDate string `json:"begin_date,omitempty"`

	// 获取数据的结束日期, YYYY-MM-DD 格式.
	// end_date 允许设置的最大值为昨日
	EndDate string `json:"end_date,omitempty"`
}

type ArticleBaseData struct {
	IntPageReadUser  int `json:"int_page_read_user"`  // 图文页（点击群发图文卡片进入的页面）的阅读人数
	IntPageReadCount int `json:"int_page_read_count"` // 图文页的阅读次数
	OriPageReadUser  int `json:"ori_page_read_user"`  // 原文页（点击图文页“阅读原文”进入的页面）的阅读人数，无原文页时此处数据为0
	OriPageReadCount int `json:"ori_page_read_count"` // 原文页的阅读次数
	ShareUser        int `json:"share_user"`          // 分享的人数
	ShareCount       int `json:"share_count"`         // 分享的次数
	AddToFavUser     int `json:"add_to_fav_user"`     // 收藏的人数
	AddToFavCount    int `json:"add_to_fav_count"`    // 收藏的次数
}

// 图文群发每日数据
type ArticleSummaryData struct {
	RefDate string `json:"ref_date"` // 数据的日期, YYYY-MM-DD 格式

	// 这里的msgid实际上是由msgid（图文消息id）和index（消息次序索引）组成，
	// 例如12003_3， 其中12003是msgid，即一次群发的id消息的；
	// 3为index，假设该次群发的图文消息共5个文章（因为可能为多图文）， 3表示5个中的第3个
	MsgId string `json:"msgid"`
	Title string `json:"title"` // 图文消息的标题
	ArticleBaseData
}

// 获取图文群发每日数据.
func (clt *Client) GetArticleSummary(param *DataCubeParam) (list []ArticleSummaryData, err error) {
	if param == nil {
		err = errors.New("nil DataCubeParam")
		return
	}

	var result struct {
		Error
		List []ArticleSummaryData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getarticlesummary?access_token="
	if err = clt.PostJSON(incompleteURL, param, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}

// 图文群发总数据
type ArticleTotalData struct {
	RefDate string `json:"ref_date"` // 数据的日期, YYYY-MM-DD 格式
	MsgId   string `json:"msgid"`    // 同 ArticleSummaryData.MsgId
	Title   string `json:"title"`
	Details []struct {
		StatDate   string `json:"stat_date"`   // 统计的日期，在getarticletotal接口中，ref_date指的是文章群发出日期， 而stat_date是数据统计日期
		TargetUser int    `json:"target_user"` // 送达人数，一般约等于总粉丝数（需排除黑名单或其他异常情况下无法收到消息的粉丝）
		ArticleBaseData
	} `json:"details"`
}

// 获取图文群发总数据.
func (clt *Client) GetArticleTotal(param *DataCubeParam) (list []ArticleTotalData, err error) {
	if param == nil {
		err = errors.New("nil DataCubeParam")
		return
	}

	var result struct {
		Error
		List []ArticleTotalData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getarticletotal?access_token="
	if err = clt.PostJSON(incompleteURL, param, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}

// 图文统计数据
type UserReadData struct {
	RefDate    string `json:"ref_date"` // 数据的日期, YYYY-MM-DD 格式
	UserSource int    `json:"user_source"`
	ArticleBaseData
}

// 获取图文统计数据.
func (clt *Client) GetUserRead(param *DataCubeParam) (list []UserReadData, err error) {
	if param == nil {
		err = errors.New("nil DataCubeParam")
		return
	}

	var result struct {
		Error
		List []UserReadData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getuserread?access_token="
	if err = clt.PostJSON(incompleteURL, param, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}

// 图文统计分时数据
type UserReadHourData struct {
	UserReadData
	RefHour int `json:"ref_hour"` // 数据的小时，包括从000到2300，分别代表的是[000,100)到[2300,2400)，即每日的第1小时和最后1小时
}

// 获取图文统计分时数据.
func (clt *Client) GetUserReadHour(param *DataCubeParam) (list []UserReadHourData, err error) {
	if param == nil {
		err = errors.New("nil DataCubeParam")
		return
	}

	var result struct {
		Error
		List []UserReadHourData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getuserreadhour?access_token="
	if err = clt.PostJSON(incompleteURL, param, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}

// 图文分享转发数据
type UserShareData struct {
	RefDate    string `json:"ref_date"`    // 数据的日期, YYYY-MM-DD 格式
	ShareScene int    `json:"share_scene"` // 分享的场景, 1代表好友转发 2代表朋友圈 3代表腾讯微博 255代表其他
	ShareUser  int    `json:"share_user"`  // 分享的人数
	ShareCount int    `json:"share_count"` // 分享的次数
}

// 获取图文分享转发数据.
func (clt *Client) GetUserShare(param *DataCubeParam) (list []UserShareData, err error) {
	if param == nil {
		err = errors.New("nil DataCubeParam")
		return
	}

	var result struct {
		Error
		List []UserShareData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getusershare?access_token="
	if err = clt.PostJSON(incompleteURL, param, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}

// 图文分享转发分时数据
type UserShareHourData struct {
	UserShareData
	RefHour int `json:"ref_hour"` // 数据的小时，包括从000到2300，分别代表的是[000,100)到[2300,2400)，即每日的第1小时和最后1小时
}

// 获取图文分享转发分时数据.
func (clt *Client) GetUserShareHour(param *DataCubeParam) (list []UserShareHourData, err error) {
	if param == nil {
		err = errors.New("nil DataCubeParam")
		return
	}

	var result struct {
		Error
		List []UserShareHourData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getusersharehour?access_token="
	if err = clt.PostJSON(incompleteURL, param, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}

// 接口分析数据
type InterfaceSummaryData struct {
	RefDate       string `json:"ref_date"`        // 数据的日期, YYYY-MM-DD 格式
	CallbackCount int    `json:"callback_count"`  // 通过服务器配置地址获得消息后，被动回复用户消息的次数
	FailCount     int    `json:"fail_count"`      // 上述动作的失败次数
	TotalTimeCost int    `json:"total_time_cost"` // 总耗时，除以callback_count即为平均耗时
	MaxTimeCost   int    `json:"max_time_cost"`   // 最大耗时
}

// 获取接口分析数据.
func (clt *Client) GetInterfaceSummary(param *DataCubeParam) (list []InterfaceSummaryData, err error) {
	if param == nil {
		err = errors.New("nil DataCubeParam")
		return
	}

	var result struct {
		Error
		List []InterfaceSummaryData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getinterfacesummary?access_token="
	if err = clt.PostJSON(incompleteURL, param, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}

type InterfaceSummaryHourData struct {
	InterfaceSummaryData
	RefHour *int `json:"ref_hour,omitempty"` // 数据的小时，包括从000到2300，分别代表的是[000,100)到[2300,2400)，即每日的第1小时和最后1小时
}

// 获取接口分析分时数据.
func (clt *Client) GetInterfaceSummaryHour(param *DataCubeParam) (list []InterfaceSummaryHourData, err error) {
	if param == nil {
		err = errors.New("nil DataCubeParam")
		return
	}

	var result struct {
		Error
		List []InterfaceSummaryHourData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getinterfacesummaryhour?access_token="
	if err = clt.PostJSON(incompleteURL, param, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}

// 消息发送概况数据
type UpstreamMsgData struct {
	RefDate string `json:"ref_date"` // 数据的日期, YYYY-MM-DD 格式

	// 消息类型，代表含义如下：
	// 1代表文字
	// 2代表图片
	// 3代表语音
	// 4代表视频
	// 6代表第三方应用消息（链接消息）
	MsgType  int `json:"msg_type"`
	MsgUser  int `json:"msg_user"`  // 上行发送了（向公众号发送了）消息的用户数
	MsgCount int `json:"msg_count"` // 上行发送了消息的消息总数
}

// 获取消息发送概况数据.
func (clt *Client) GetUpstreamMsg(param *DataCubeParam) (list []UpstreamMsgData, err error) {
	if param == nil {
		err = errors.New("nil DataCubeParam")
		return
	}

	var result struct {
		Error
		List []UpstreamMsgData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getupstreammsg?access_token="
	if err = clt.PostJSON(incompleteURL, param, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}

// 消息分送分时数据
type UpstreamMsgHourData struct {
	UpstreamMsgData
	RefHour int `json:"ref_hour"` // 数据的小时，包括从000到2300，分别代表的是[000,100)到[2300,2400)，即每日的第1小时和最后1小时
}

// 获取消息分送分时数据.
func (clt *Client) GetUpstreamMsgHour(param *DataCubeParam) (list []UpstreamMsgHourData, err error) {
	if param == nil {
		err = errors.New("nil DataCubeParam")
		return
	}

	var result struct {
		Error
		List []UpstreamMsgHourData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getupstreammsghour?access_token="
	if err = clt.PostJSON(incompleteURL, param, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}

// 消息发送周数据
type UpstreamMsgWeekData struct {
	UpstreamMsgData
}

// 获取消息发送周数据.
func (clt *Client) GetUpstreamMsgWeek(param *DataCubeParam) (list []UpstreamMsgWeekData, err error) {
	if param == nil {
		err = errors.New("nil DataCubeParam")
		return
	}

	var result struct {
		Error
		List []UpstreamMsgWeekData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getupstreammsgweek?access_token="
	if err = clt.PostJSON(incompleteURL, param, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}

// 消息发送月数据
type UpstreamMsgMonthData struct {
	UpstreamMsgData
}

// 获取消息发送月数据.
func (clt *Client) GetUpstreamMsgMonth(param *DataCubeParam) (list []UpstreamMsgMonthData, err error) {
	if param == nil {
		err = errors.New("nil DataCubeParam")
		return
	}

	var result struct {
		Error
		List []UpstreamMsgMonthData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getupstreammsgmonth?access_token="
	if err = clt.PostJSON(incompleteURL, param, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}

// 消息发送分布数据
type UpstreamMsgDistData struct {
	RefDate       string `json:"ref_date"`       // 数据的日期, YYYY-MM-DD 格式
	CountInterval int    `json:"count_interval"` // 当日发送消息量分布的区间，0代表 “0”，1代表“1-5”，2代表“6-10”，3代表“10次以上”
	MsgUser       int    `json:"msg_user"`       // 上行发送了（向公众号发送了）消息的用户数
}

// 获取消息发送分布数据.
func (clt *Client) GetUpstreamMsgDist(param *DataCubeParam) (list []UpstreamMsgDistData, err error) {
	if param == nil {
		err = errors.New("nil DataCubeParam")
		return
	}

	var result struct {
		Error
		List []UpstreamMsgDistData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getupstreammsgdist?access_token="
	if err = clt.PostJSON(incompleteURL, param, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}

// 消息发送分布周数据
type UpstreamMsgDistWeekData struct {
	UpstreamMsgDistData
}

// 获取消息发送分布周数据.
func (clt *Client) GetUpstreamMsgDistWeek(param *DataCubeParam) (list []UpstreamMsgDistWeekData, err error) {
	if param == nil {
		err = errors.New("nil Request")
		return
	}

	var result struct {
		Error
		List []UpstreamMsgDistWeekData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getupstreammsgdistweek?access_token="
	if err = clt.PostJSON(incompleteURL, param, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}

// 消息发送分布月数据
type UpstreamMsgDistMonthData struct {
	UpstreamMsgDistData
}

// 获取消息发送分布月数据.
func (clt *Client) GetUpstreamMsgDistMonth(param *DataCubeParam) (list []UpstreamMsgDistMonthData, err error) {
	if param == nil {
		err = errors.New("nil DataCubeParam")
		return
	}

	var result struct {
		Error
		List []UpstreamMsgDistMonthData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getupstreammsgdistmonth?access_token="
	if err = clt.PostJSON(incompleteURL, param, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}

// 用户增减数据
type UserSummaryData struct {
	RefDate string `json:"ref_date"` // 数据的日期, YYYY-MM-DD 格式

	// 用户的渠道，数值代表的含义如下：
	// 0  代表其他
	// 30 代表扫二维码
	// 17 代表名片分享
	// 35 代表搜号码（即微信添加朋友页的搜索）
	// 39 代表查询微信公众帐号
	// 43 代表图文页右上角菜单
	UserSource int `json:"user_source"`

	NewUser    int `json:"new_user"`    // 新增的用户数量
	CancelUser int `json:"cancel_user"` // 取消关注的用户数量，new_user减去cancel_user即为净增用户数量
}

// 获取用户增减数据.
func (clt *Client) GetUserSummary(param *DataCubeParam) (list []UserSummaryData, err error) {
	if param == nil {
		err = errors.New("nil DataCubeParam")
		return
	}

	var result struct {
		Error
		List []UserSummaryData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getusersummary?access_token="
	if err = clt.PostJSON(incompleteURL, param, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}

// 累计用户数据
type UserCumulateData struct {
	RefDate      string `json:"ref_date"`      // 数据的日期, YYYY-MM-DD 格式
	CumulateUser int    `json:"cumulate_user"` // 总用户量
}

// 获取累计用户数据.
func (clt *Client) GetUserCumulate(param *DataCubeParam) (list []UserCumulateData, err error) {
	if param == nil {
		err = errors.New("nil DataCubeParam")
		return
	}

	var result struct {
		Error
		List []UserCumulateData `json:"list"`
	}

	incompleteURL := "https://api.weixin.qq.com/datacube/getusercumulate?access_token="
	if err = clt.PostJSON(incompleteURL, param, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	list = result.List
	return
}
