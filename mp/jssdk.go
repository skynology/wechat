package mp

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/skynology/wechat/util"
)

// 微信 js-sdk wx.config 的参数签名.
func WXConfigSign(jsapiTicket, url string) (signature string, timestamp string, noncestr string) {
	noncestr = util.RandString(15)
	timestamp = fmt.Sprintf("%v", time.Now().Unix())

	n := len("jsapi_ticket=") + len(jsapiTicket) +
		len("&noncestr=") + len(noncestr) +
		len("&timestamp=") + len(timestamp) +
		len("&url=") + len(url)

	buf := make([]byte, 0, n)

	buf = append(buf, "jsapi_ticket="...)
	buf = append(buf, jsapiTicket...)
	buf = append(buf, "&noncestr="...)
	buf = append(buf, noncestr...)
	buf = append(buf, "&timestamp="...)
	buf = append(buf, timestamp...)
	buf = append(buf, "&url="...)
	buf = append(buf, url...)

	// fmt.Println("wx sign is:", string(buf))

	hashsum := sha1.Sum(buf)
	signature = hex.EncodeToString(hashsum[:])

	return
}
