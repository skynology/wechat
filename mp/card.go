// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     gaowenbin(gaowenbinmarr@gmail.com), chanxuehong(chanxuehong@gmail.com)

package mp

import (
	"errors"
	"fmt"
)

const (
	// 卡券类型
	CardTypeGeneralCoupon = "GENERAL_COUPON" // 通用券
	CardTypeGroupon       = "GROUPON"        // 团购券
	CardTypeGift          = "GIFT"           // 礼品券
	CardTypeCash          = "CASH"           // 代金券
	CardTypeDiscount      = "DISCOUNT"       // 折扣券
	CardTypeMemberCard    = "MEMBER_CARD"    // 会员卡
	CardTypeScenicTicket  = "SCENIC_TICKET"  // 景点门票
	CardTypeMovieTicket   = "MOVIE_TICKET"   // 电影票
	CardTypeBoardingPass  = "BOARDING_PASS"  // 飞机票
	CardTypeLuckyMoney    = "LUCKY_MONEY"    // 红包
	CardTypeMeetingTicket = "MEETING_TICKET" // 会议门票
)

// 卡卷数据结构
type Card struct {
	CardType string `json:"card_type,omitempty"`

	GeneralCoupon *GeneralCoupon `json:"general_coupon,omitempty"`
	Groupon       *Groupon       `json:"groupon,omitempty"`
	Gift          *Gift          `json:"gift,omitempty"`
	Cash          *Cash          `json:"cash,omitempty"`
	Discount      *Discount      `json:"discount,omitempty"`
	MemberCard    *MemberCard    `json:"member_card,omitempty"`
	ScenicTicket  *ScenicTicket  `json:"scenic_ticket,omitempty"`
	MovieTicket   *MovieTicket   `json:"movie_ticket,omitempty"`
	BoardingPass  *BoardingPass  `json:"boarding_pass,omitempty"`
	LuckyMoney    *LuckyMoney    `json:"lucky_money,omitempty"`
	MeetingTicket *MeetingTicket `json:"meeting_ticket,omitempty"`
}

// 通用券
type GeneralCoupon struct {
	BaseInfo      *CardBaseInfo `json:"base_info,omitempty"`
	DefaultDetail string        `json:"default_detail,omitempty"` // 描述文本
}

// 团购券
type Groupon struct {
	BaseInfo   *CardBaseInfo `json:"base_info,omitempty"`
	DealDetail string        `json:"deal_detail,omitempty"` // 团购券专用，团购详情
}

// 礼品券
type Gift struct {
	BaseInfo *CardBaseInfo `json:"base_info,omitempty"`
	Gift     string        `json:"gift,omitempty"` // 礼品券专用，表示礼品名字
}

// 代金券
type Cash struct {
	BaseInfo   *CardBaseInfo `json:"base_info,omitempty"`
	LeastCost  *int          `json:"least_cost,omitempty"`  // 代金券专用，表示起用金额（单位为分）
	ReduceCost *int          `json:"reduce_cost,omitempty"` // 代金券专用，表示减免金额（单位为分）
}

// 折扣券
type Discount struct {
	BaseInfo *CardBaseInfo `json:"base_info,omitempty"`
	Discount *int          `json:"discount,omitempty"` // 折扣券专用，表示打折额度（百分比）。填30 就是七折。
}

// 会员卡
type MemberCard struct {
	BaseInfo *CardBaseInfo `json:"base_info,omitempty"`

	// 是否支持积分，填写true 或false，如填写true，积分相关字段均为必填。
	// 填写false，积分字段无需填写。储值字段处理方式相同。
	SupplyBonus       *bool  `json:"supply_bonus,omitempty"`
	SupplyBalance     *bool  `json:"supply_balance,omitempty"`    // 是否支持储值
	BonusClearedRules string `json:"bonus_cleared,omitempty"`     // 积分清零规则
	BonusRules        string `json:"bonus_rules,omitempty"`       // 积分规则
	BalanceRules      string `json:"balance_rules,omitempty"`     // 储值说明
	Prerogative       string `json:"prerogative,omitempty"`       // 特权说明
	BindOldCardURL    string `json:"bind_old_card_url,omitempty"` // 绑定旧卡的url，与“activate_url”字段二选一必填。
	ActivateURL       string `json:"activate_url,omitempty"`      // 激活会员卡的url，与“bind_old_card_url”字段二选一必填。
	NeedPushOnView    *bool  `json:"need_push_on_view,omitempty"` // true 为用户点击进入会员卡时是否推送事件。
}

// 景点门票
type ScenicTicket struct {
	BaseInfo    *CardBaseInfo `json:"base_info,omitempty"`
	TicketClass string        `json:"ticket_class,omitempty"` // 票类型，例如平日全票，套票等
	GuideURL    string        `json:"guide_url,omitempty"`    // 导览图url
}

// 电影票
type MovieTicket struct {
	BaseInfo *CardBaseInfo `json:"base_info,omitempty"`
	Detail   string        `json:"detail,omitempty"` // 电影票详情
}

// 飞机票
type BoardingPass struct {
	BaseInfo *CardBaseInfo `json:"base_info,omitempty"`

	From          string `json:"from,omitempty"`           // 起点，上限为18 个汉字
	To            string `json:"to,omitempty"`             // 终点，上限为18 个汉字
	Flight        string `json:"flight,omitempty"`         // 航班
	DepartureTime int64  `json:"departure_time,omitempty"` // 起飞时间。Unix 时间戳格式
	LandingTime   int64  `json:"landing_time,omitempty"`   // 降落时间。Unix 时间戳格式
	CheckinURL    string `json:"check_in_url,omitempty"`   // 在线值机的链接
	Gate          string `json:"gate,omitempty"`           // 登机口。如发生登机口变更，建议商家实时调用该接口变更
	BoardingTime  int64  `json:"boarding_time,omitempty"`  // 登机时间，只显示“时分”不显示日期，按时间戳格式填写。如发生登机时间变更，建议商家实时调用该接口变更。
	AirModel      string `json:"air_model,omitempty"`      // 机型，上限为8 个汉字
}

// 红包
type LuckyMoney struct {
	BaseInfo *CardBaseInfo `json:"base_info,omitempty"`
}

// 会议门票
type MeetingTicket struct {
	BaseInfo      *CardBaseInfo `json:"base_info,omitempty"`
	MeetingDetail string        `json:"meeting_detail,omitempty"` // 会议详情
	MapURL        string        `json:"map_url,omitempty"`        // 会议导览图
}

// base_info ===================================================================

const (
	// 卡券code码展示类型
	CodeTypeText        = "CODE_TYPE_TEXT"         // 文本
	CodeTypeBarCode     = "CODE_TYPE_BARCODE"      // 一维码
	CodeTypeQRCode      = "CODE_TYPE_QRCODE"       // 二维码
	CodeTypeOnlyBarCode = "CODE_TYPE_ONLY_BARCODE" // 一维码无code 显示
	CodeTypeOnlyQRCode  = "CODE_TYPE_ONLY_QRCODE"  // 二维码无code 显示
)

const (
	// 卡卷的状态
	CardStatusNotVerify    = "CARD_STATUS_NOT_VERIFY"    // 待审核
	CardStatusVerifyFail   = "CARD_STATUS_VERIFY_FALL"   // 审核失败
	CardStatusVerifyOk     = "CARD_STATUS_VERIFY_OK"     // 通过审核
	CardStatusUserDelete   = "CARD_STATUS_USER_DELETE"   // 卡券被用户删除
	CardStatusUserDispatch = "CARD_STATUS_USER_DISPATCH" // 在公众平台投放过的卡券
)

// 基本的卡券数据，所有卡券通用
type CardBaseInfo struct {
	CardId string `json:"id,omitempty"`     // 查询的时候有返回
	Status string `json:"status,omitempty"` // 查询的时候有返回

	LogoURL     string `json:"logo_url,omitempty"`    // 卡券的商户logo，尺寸为300*300。
	CodeType    string `json:"code_type,omitempty"`   // code 码展示类型
	BrandName   string `json:"brand_name,omitempty"`  // 商户名字,字数上限为12 个汉字。（填写直接提供服务的商户名， 第三方商户名填写在source 字段）
	Title       string `json:"title,omitempty"`       // 券名，字数上限为9 个汉字。(建议涵盖卡券属性、服务及金额)
	SubTitle    string `json:"sub_title,omitempty"`   // 券名的副标题，字数上限为18个汉字。
	Color       string `json:"color,omitempty"`       // 券颜色。按色彩规范标注填写Color010-Color100
	Notice      string `json:"notice,omitempty"`      // 使用提醒，字数上限为9 个汉字。（一句话描述，展示在首页，示例：请出示二维码核销卡券）
	Description string `json:"description,omitempty"` // 使用说明。长文本描述，可以分行，上限为1000 个汉字。

	DateInfo *DateInfo `json:"date_info,omitempty"` // 有效日期
	SKU      *SKU      `json:"sku,omitempty"`       // 商品信息

	LocationIdList       []int64 `json:"location_id_list,omitempty"`        // 门店地址ID
	UseCustomCode        *bool   `json:"use_custom_code,omitempty"`         // 是否自定义code 码。
	BindOpenId           *bool   `json:"bind_openid,omitempty"`             // 是否指定用户领取，填写true或false。不填代表默认为否。
	CanShare             *bool   `json:"can_share,omitempty"`               // 领取卡券原生页面是否可分享，填写true 或false，true 代表可分享。默认可分享。
	CanGiveFriend        *bool   `json:"can_give_friend,omitempty"`         // 卡券是否可转赠，填写true 或false,true 代表可转赠。默认可转赠。
	UseLimit             *int    `json:"use_limit,omitempty"`               // 每人使用次数限制。
	GetLimit             *int    `json:"get_limit,omitempty"`               // 每人最大领取次数，不填写默认等于quantity。
	ServicePhone         string  `json:"service_phone,omitempty"`           // 客服电话
	Source               string  `json:"source,omitempty"`                  // 第三方来源名，如携程
	CustomURLName        string  `json:"custom_url_name,omitempty"`         // 商户自定义入口名称，与custom_url 字段共同使用，长度限制在5 个汉字内。
	CustomURL            string  `json:"custom_url,omitempty"`              // 商户自定义入口跳转外链的地址链接,跳转页面内容需与自定义cell 名称保持匹配。
	CustomURLSubTitle    string  `json:"custom_url_sub_title,omitempty"`    // 显示在入口右侧的tips，长度限制在6 个汉字内。
	PromotionURLName     string  `json:"promotion_url_name,omitempty"`      // 营销场景的自定义入口
	PromotionURL         string  `json:"promotion_url,omitempty"`           // 入口跳转外链的地址链接
	PromotionURLSubTitle string  `json:"promotion_url_sub_title,omitempty"` // 显示在入口右侧的tips，长度限制在6 个汉字内。
}

type DateInfo struct {
	// 使用时间的类型1：固定日期区间，2：固定时长（自领取后按天算）
	Type int `json:"type"`
	// 固定日期区间专用，表示起用时间。从1970 年1 月1 日00:00:00至起用时间的秒数，最终需转换为字符串形态传入，下同。（单位为秒）
	BeginTimestamp int64 `json:"begin_timestamp,omitempty"`
	// 固定日期区间专用，表示结束时间。（单位为秒）
	EndTimestamp int64 `json:"end_timestamp,omitempty"`
	// 固定时长专用，表示自领取后多少天内有效。（单位为天）领取后当天有效填写0。
	FixedTerm int `json:"fixed_term,omitempty"`
	// 固定时长专用，表示自领取后多少天开始生效。（单位为天）
	FixedBeginTerm int `json:"fixed_begin_term,omitempty"`
}

type SKU struct {
	Quantity int `json:"quantity,omitempty"` // 上架的数量。（不支持填写0或无限大）
}
type BoardingPassCheckinParameters struct {
	Code   string `json:"code"`              // 必须; 飞机票的序列号
	CardId string `json:"card_id,omitempty"` // 可选; 需办理值机的机票card_id。自定义code 的飞机票为必填

	PassengerName string `json:"passenger_name,omitempty"` // 必须; 乘客姓名，上限为15 个汉字。
	Class         string `json:"class,omitempty"`          // 必须; 舱等，如头等舱等，上限为5 个汉字。
	Seat          string `json:"seat,omitempty"`           // 可选; 乘客座位号。
	ETKT_NBR      string `json:"etkt_bnr,omitempty"`       // 必须; 电子客票号，上限为14 个数字
	QRCodeData    string `json:"qrcode_data,omitempty"`    // 可选; 二维码数据。乘客用于值机的二维码字符串，微信会通过此数据为用户生成值机用的二维码。
	IsCancel      *bool  `json:"is_cancel,omitempty"`      // 可选; 是否取消值机。填写true 或false。true 代表取消，如填写true 上述字段（如calss 等）均不做判断，机票返回未值机状态，乘客可重新值机。默认填写false
}

// 在线值机接口.
//  领取电影票后通过调用“更新电影票”接口update 电影信息及用户选座信息
func (clt *Client) BoardingPassCheckin(para *BoardingPassCheckinParameters) (err error) {
	if para == nil {
		return errors.New("nil BoardingPassCheckinParameters")
	}

	var result Error

	incompleteURL := "https://api.weixin.qq.com/card/boardingpass/checkin?access_token="
	if err = clt.PostJSON(incompleteURL, para, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result
		return
	}
	return
}

// 创建卡券接口.
//  Card 需要设置哪些字段请参考微信官方文档.
func (clt *Client) CardCreate(card *Card) (cardId string, err error) {
	if card == nil {
		err = errors.New("nil card")
		return
	}

	var request = struct {
		*Card `json:"card,omitempty"`
	}{
		Card: card,
	}

	var result struct {
		Error
		CardId string `json:"card_id"`
	}

	incompleteURL := "https://api.weixin.qq.com/card/create?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	cardId = result.CardId
	return
}

// 查询卡券详情.
//  返回的 Card 有哪些字段请参考微信官方文档.
func (clt *Client) CardGet(cardId string) (card *Card, err error) {
	var request = struct {
		CardId string `json:"card_id"`
	}{
		CardId: cardId,
	}

	var result struct {
		Error
		Card `json:"card"`
	}

	incompleteURL := "https://api.weixin.qq.com/card/get?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	card = &result.Card
	return
}

// 更改卡券信息接口.
//  支持更新部分通用字段及特殊卡券（会员卡、飞机票、电影票、红包）中特定字段的信息，请参考微信官方文档.。
//  注：更改卡券的部分字段后会重新提交审核，详情见字段说明。
func (clt *Client) CardUpdate(cardId string, card *Card) (err error) {
	if card == nil {
		return errors.New("nil Card")
	}
	card.CardType = "" // NOTE

	var request = struct {
		CardId string `json:"card_id"`
		*Card
	}{
		CardId: cardId,
		Card:   card,
	}

	var result Error

	incompleteURL := "https://api.weixin.qq.com/card/update?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}
	if result.ErrCode != ErrCodeOK {
		err = &result
		return
	}
	return
}

// 删除卡券
//  删除卡券接口允许商户删除任意一类卡券。删除卡券后，该卡券对应已生成的领取用二维码、添加到卡包JS API 均会失效。
//  注意：如用户在商家删除卡券前已领取一张或多张该卡券依旧有效。即删除卡券不能删除已被用户领取，保存在微信客户端中的卡券。
func (clt *Client) CardDelete(cardId string) (err error) {
	var request = struct {
		CardId string `json:"card_id"`
	}{
		CardId: cardId,
	}

	var result Error

	incompleteURL := "https://api.weixin.qq.com/card/delete?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result
		return
	}
	return
}

// 批量查询卡列表.
//  offset: 查询卡列表的起始偏移量，从0 开始，即offset: 5 是指从从列表里的第六个开始读取。
//  count : 需要查询的卡片的数量（数量最大50）
func (clt *Client) CardBatchGet(offset, count int) (cardIdList []string, totalNum int, err error) {
	if offset < 0 {
		err = fmt.Errorf("invalid offset: %d", offset)
		return
	}
	if count < 0 {
		err = fmt.Errorf("invalid count: %d", count)
		return
	}

	var request = struct {
		Offset int `json:"offset"`
		Count  int `json:"count"`
	}{
		Offset: offset,
		Count:  count,
	}

	var result struct {
		Error
		CardIdList []string `json:"card_id_list"`
		TotalNum   int      `json:"total_num"`
	}

	incompleteURL := "https://api.weixin.qq.com/card/batchget?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	cardIdList = result.CardIdList
	totalNum = result.TotalNum
	return
}

// 库存修改接口.
// cardId:      卡券ID
// increaseNum: 增加库存数量, 可以为负数
func (clt *Client) CardModifyStock(cardId string, increaseNum int) (err error) {
	var request struct {
		CardId             string `json:"card_id"`
		IncreaseStockValue int    `json:"increase_stock_value,omitempty"`
		ReduceStockValue   int    `json:"reduce_stock_value,omitempty"`
	}
	request.CardId = cardId
	switch {
	case increaseNum > 0:
		request.IncreaseStockValue = increaseNum
	case increaseNum < 0:
		request.ReduceStockValue = -increaseNum
	default: // increaseNum == 0
		return
	}

	var result Error

	incompleteURL := "https://api.weixin.qq.com/card/modifystock?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result
		return
	}
	return
}

type MeetingTicketUpdateUserParameters struct {
	Code   string `json:"code"`              // 必须; 用户的门票唯一序列号
	CardId string `json:"card_id,omitempty"` // 可选; 要更新门票序列号所述的card_id ， 生成券时use_custom_code 填写true 时必填。

	Zone       string `json:"zone,omitempty"`        // 可选; 区域
	Entrance   string `json:"entrance,omitempty"`    // 可选; 入口
	SeatNumber string `json:"seat_number,omitempty"` // 可选; 座位号
	BeginTime  int64  `json:"begin_time,omitempty"`  // 开场时间
	EndTime    int64  `json:"end_time,omitempty"`    // 结束时间
}

// 更新电影票.
//  领取电影票后通过调用“更新电影票”接口update 电影信息及用户选座信息
func (clt *Client) MeetingTicketUpdateUser(para *MeetingTicketUpdateUserParameters) (err error) {
	if para == nil {
		return errors.New("nil MeetingTicketUpdateUserParameters")
	}

	var result Error

	incompleteURL := "https://api.weixin.qq.com/card/meetingticket/updateuser?access_token="
	if err = clt.PostJSON(incompleteURL, para, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result
		return
	}
	return
}

// 创建二维码的参数
type CardQRCodeInfo struct {
	CardId        string `json:"card_id"`                  // 必须; 卡券ID
	Code          string `json:"code,omitempty"`           // 可选; 指定卡券code码，只能被领一次。use_custom_code 字段为true 的卡券必须填写，非自定义code 不必填写。
	OpenId        string `json:"openid,omitempty"`         // 可选; 指定领取者的openid，只有该用户能领取。bind_openid字段为true 的卡券必须填写，非自定义openid 不必填写。
	ExpireSeconds *int   `json:"expire_seconds,omitempty"` // 可选; 指定二维码的有效时间，范围是60 ~ 1800 秒。不填默认为永久有效。
	IsUniqueCode  *bool  `json:"is_unique_code,omitempty"` // 可选; 指定下发二维码，生成的二维码随机分配一个code，领取后不可再次扫描。填写true 或false。默认false。
	Balance       *int   `json:"balance,omitempty"`        // 可选; 红包余额，以分为单位。红包类型必填（LUCKY_MONEY），其他卡券类型不填。
	OuterId       *int64 `json:"outer_id,omitempty"`       // 可选; 领取场景值，用于领取渠道的数据统计，默认值为0，字段类型为整型。用户领取卡券后触发的事件推送中会带上此自定义场景值。
}

// 卡券投放, 创建二维码.
//  创建卡券后，商户可通过接口生成一张卡券二维码供用户扫码后添加卡券到卡包。
func (clt *Client) CardQRCodeCreate(qrcodeInfo *CardQRCodeInfo) (ticket string, err error) {
	if qrcodeInfo == nil {
		err = errors.New("nil CardQRCodeInfo")
		return
	}

	var request struct {
		ActionName string `json:"action_name"`
		ActionInfo struct {
			Card *CardQRCodeInfo `json:"card,omitempty"`
		} `json:"action_info"`
	}
	request.ActionName = "QR_CARD"
	request.ActionInfo.Card = qrcodeInfo

	var result struct {
		Error
		Ticket string `json:"ticket"`
	}

	incompleteURL := "https://api.weixin.qq.com/card/qrcode/create?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	ticket = result.Ticket
	return
}

type TestWhiteListSetParameters struct {
	OpenIdList   []string `json:"openid,omitempty"`   // 测试的openid 列表
	UserNameList []string `json:"username,omitempty"` // 测试的微信号列表
}

// 设置测试用户白名单.
//  由于卡券有审核要求，为方便公众号调试，可以设置一些测试帐号，这些帐号可领取未通过审核的卡券，体验整个流程。
//  注：同时支持“openid”、“username”两种字段设置白名单，总数上限为10 个。
func (clt *Client) SetCardTestWhiteList(para *TestWhiteListSetParameters) (err error) {
	if para == nil {
		return errors.New("nil TestWhiteListSetParameters")
	}

	var result Error

	incompleteURL := "https://api.weixin.qq.com/card/testwhitelist/set?access_token="
	if err = clt.PostJSON(incompleteURL, para, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result
		return
	}
	return
}
