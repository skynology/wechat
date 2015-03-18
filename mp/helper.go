// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package mp

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"strconv"

	"github.com/skynology/wechat/util"
)

// 安全模式回复消息的 http body
type ResponseHttpBody struct {
	XMLName      struct{} `xml:"xml" json:"-"`
	EncryptedMsg string   `xml:"Encrypt"`
	MsgSignature string   `xml:"MsgSignature"`
	TimeStamp    int64    `xml:"TimeStamp"`
	Nonce        string   `xml:"Nonce"`
}

type RequestParam struct {
	AppId     string
	Random    []byte
	Timestamp int64
	Token     string
	AESKey    [32]byte
	Nonce     string
}

// 回复消息给微信服务器(安全模式).
//  要求 msg 是有效的消息数据结构(经过 encoding/xml marshal 后符合消息的格式);
//  如果有必要可以修改 Request 里面的某些值, 比如 TimeStamp.
func GetAESResponse(r *RequestParam, msg interface{}) (result ResponseHttpBody, err error) {
	if r == nil {
		err = errors.New("nil request params")
		return
	}
	if msg == nil {
		err = errors.New("nil message")
		return
	}

	rawXMLMsg, err := xml.Marshal(msg)
	if err != nil {
		return
	}

	encryptedMsg := util.AESEncryptMsg(r.Random, rawXMLMsg, r.AppId, r.AESKey)
	base64EncryptedMsg := base64.StdEncoding.EncodeToString(encryptedMsg)

	responseHttpBody := ResponseHttpBody{
		EncryptedMsg: base64EncryptedMsg,
		TimeStamp:    r.Timestamp,
		Nonce:        r.Nonce,
	}

	timestampStr := strconv.FormatInt(responseHttpBody.TimeStamp, 10)
	responseHttpBody.MsgSignature = util.MsgSign(r.Token, timestampStr,
		responseHttpBody.Nonce, responseHttpBody.EncryptedMsg)

	result = responseHttpBody
	return
}
