package mp

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
)

const (
	TemporaryQRCodeExpireSecondsLimit = 1800   // 临时二维码 expire_seconds 限制
	PermanentQRCodeSceneIdLimit       = 100000 // 永久二维码 scene_id 限制
)

// 永久二维码
type PermanentQRCode struct {
	// 下面两个字段同时只有一个有效
	SceneId     int    `json:"scene_id"`  // 场景值
	SceneString string `json:"scene_str"` // 场景值ID（字符串形式的ID），字符串类型，长度限制为1到64

	Ticket string `json:"ticket"` // 获取的二维码ticket，凭借此ticket可以在有效时间内换取二维码。
	URL    string `json:"url"`    // 二维码图片解析后的地址，开发者可根据该地址自行生成需要的二维码图片
}

// 二维码图片的URL, 可以GET此URL下载二维码或者在线显示此二维码.
func (qrcode *PermanentQRCode) PicURL() string {
	return "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + url.QueryEscape(qrcode.Ticket)
}

// 临时二维码
type TemporaryQRCode struct {
	PermanentQRCode
	ExpiresIn int `json:"expire_seconds"` // 二维码的有效时间，以秒为单位。
}

// 二维码图片的URL, 可以GET此URL下载二维码或者在线显示此二维码.
func (clt *Client) QRCodeURL(ticket string) string {
	return "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + url.QueryEscape(ticket)
}

// 创建临时二维码
//  SceneId:       场景值ID，为32位非0整型
//  ExpireSeconds: 二维码的有效时间，以秒为单位。
func (clt *Client) CreateTemporaryQRCode(SceneId int, ExpireSeconds int) (qrcode *TemporaryQRCode, err error) {
	var request struct {
		ExpireSeconds int    `json:"expire_seconds"`
		ActionName    string `json:"action_name"`
		ActionInfo    struct {
			Scene struct {
				SceneId int `json:"scene_id"`
			} `json:"scene"`
		} `json:"action_info"`
	}
	request.ExpireSeconds = ExpireSeconds
	request.ActionName = "QR_SCENE"
	request.ActionInfo.Scene.SceneId = SceneId

	// fmt.Println("qrcode:", request)

	var result struct {
		Error
		TemporaryQRCode
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	result.TemporaryQRCode.SceneId = SceneId
	qrcode = &result.TemporaryQRCode
	return
}

// 创建永久二维码
//  SceneId: 场景值ID，最大值为100000（目前参数只支持1--100000）
func (clt *Client) CreatePermanentQRCode(SceneId int) (qrcode *PermanentQRCode, err error) {
	var request struct {
		ActionName string `json:"action_name"`
		ActionInfo struct {
			Scene struct {
				SceneId int `json:"scene_id"`
			} `json:"scene"`
		} `json:"action_info"`
	}
	request.ActionName = "QR_LIMIT_SCENE"
	request.ActionInfo.Scene.SceneId = SceneId

	var result struct {
		Error
		PermanentQRCode
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	result.PermanentQRCode.SceneId = SceneId
	qrcode = &result.PermanentQRCode
	return
}

// 创建永久二维码
//  SceneString: 场景值ID（字符串形式的ID），字符串类型，长度限制为1到64
func (clt *Client) CreatePermanentQRCodeWithSceneString(SceneString string) (qrcode *PermanentQRCode, err error) {
	var request struct {
		ActionName string `json:"action_name"`
		ActionInfo struct {
			Scene struct {
				SceneString string `json:"scene_str"`
			} `json:"scene"`
		} `json:"action_info"`
	}
	request.ActionName = "QR_LIMIT_SCENE"
	request.ActionInfo.Scene.SceneString = SceneString

	var result struct {
		Error
		PermanentQRCode
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	result.PermanentQRCode.SceneString = SceneString
	qrcode = &result.PermanentQRCode
	return
}

// 通过ticket换取二维码, 写入到 writer.
//  NOTE: 调用者保证所有参数有效.
func qrcodeDownloadToWriter(ticket string, writer io.Writer, httpClient *http.Client) (err error) {
	qrcodeURL := "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + url.QueryEscape(ticket)
	httpResp, err := httpClient.Get(qrcodeURL)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode == http.StatusOK {
		_, err = io.Copy(writer, httpResp.Body)
		return
	}
	return errors.New("下载二维码出错, ticket: " + ticket)
}

// 通过ticket换取二维码, 写入到 writer.
//  如果 httpClient == nil 则默认用 http.DefaultClient.
func QRCodeDownloadToWriter(ticket string, writer io.Writer, httpClient *http.Client) (err error) {
	if writer == nil {
		return errors.New("nil writer")
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return qrcodeDownloadToWriter(ticket, writer, httpClient)
}

// 通过ticket换取二维码, 写入到 writer.
func (clt *Client) QRCodeDownloadToWriter(ticket string, writer io.Writer) (err error) {
	if writer == nil {
		return errors.New("nil writer")
	}
	if clt.httpClient == nil {
		clt.httpClient = http.DefaultClient
	}
	return qrcodeDownloadToWriter(ticket, writer, clt.httpClient)
}

// 通过ticket换取二维码, 写入到 filepath 路径的文件.
//  如果 httpClient == nil 则默认用 http.DefaultClient
func QRCodeDownload(ticket, filepath string, httpClient *http.Client) (err error) {
	file, err := os.Create(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return qrcodeDownloadToWriter(ticket, file, httpClient)
}

// 通过ticket换取二维码, 写入到 filepath 路径的文件.
func (clt *Client) QRCodeDownload(ticket, filepath string) (err error) {
	file, err := os.Create(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	if clt.httpClient == nil {
		clt.httpClient = http.DefaultClient
	}
	return qrcodeDownloadToWriter(ticket, file, clt.httpClient)
}

// 将一条长链接转成短链接.
//  主要使用场景：
//  开发者用于生成二维码的原链接（商品、支付二维码等）太长导致扫码速度和成功率下降，
//  将原长链接通过此接口转成短链接再生成二维码将大大提升扫码速度和成功率。
func (clt *Client) ShortURL(LongURL string) (ShortURL string, err error) {
	var request = struct {
		Action  string `json:"action"`
		LongURL string `json:"long_url"`
	}{
		Action:  "long2short",
		LongURL: LongURL,
	}

	var result struct {
		Error
		ShortURL string `json:"short_url"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/shorturl?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	ShortURL = result.ShortURL
	return
}
