// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package corp

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"strconv"
	"strings"

	"github.com/skynology/wechat/util"
)

type RequestParam struct {
	AppId     string
	Random    []byte
	Timestamp int64
	Token     string
	AESKey    [32]byte
	Nonce     string
}

// 回复消息的 http body
type ResponseHttpBody struct {
	XMLName      struct{} `xml:"xml" json:"-"`
	EncryptedMsg string   `xml:"Encrypt"`
	MsgSignature string   `xml:"MsgSignature"`
	TimeStamp    int64    `xml:"TimeStamp"`
	Nonce        string   `xml:"Nonce"`
}

func GetResponse(r *RequestParam, msg interface{}) (result ResponseHttpBody, err error) {
	if r == nil {
		err = errors.New("nil Request Parmas")
		return
	}
	if msg == nil {
		err = errors.New("nil message")
		return
	}

	MsgRawXML, err := xml.Marshal(msg)
	if err != nil {
		return
	}

	EncryptedMsg := util.AESEncryptMsg(r.Random, MsgRawXML, r.AppId, r.AESKey)
	base64EncryptedMsg := base64.StdEncoding.EncodeToString(EncryptedMsg)

	responseHttpBody := ResponseHttpBody{
		EncryptedMsg: base64EncryptedMsg,
		TimeStamp:    r.Timestamp,
		Nonce:        r.Nonce,
	}

	TimestampStr := strconv.FormatInt(responseHttpBody.TimeStamp, 10)
	responseHttpBody.MsgSignature = util.MsgSign(r.Token, TimestampStr,
		responseHttpBody.Nonce, responseHttpBody.EncryptedMsg)

	result = responseHttpBody
	return
}

// 用 '|' 连接 a 的各个元素
func JoinString(a []string) string {
	return strings.Join(a, "|")
}

// 用 '|' 连接 a 的各个元素的十进制字符串
func JoinInt64(a []int64) string {
	switch len(a) {
	case 0:
		return ""
	case 1:
		return strconv.FormatInt(a[0], 10)
	default:
		strs := make([]string, len(a))
		for i, n := range a {
			strs[i] = strconv.FormatInt(n, 10)
		}
		return strings.Join(strs, "|")
	}
}

// 用 '|' 分离 str
func SplitString(str string) []string {
	return strings.Split(str, "|")
}

// 用 '|' 分离 str, 然后将分离后的字符串都转换为整数
//  NOTE: 要求 str 都是整数合并的, 否则会出错
func SplitInt64(str string) (dst []int64, err error) {
	strs := strings.Split(str, "|")

	dst = make([]int64, len(strs))
	for i, str := range strs {
		dst[i], err = strconv.ParseInt(str, 10, 64)
		if err != nil {
			return
		}
	}
	return
}
