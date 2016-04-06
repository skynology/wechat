// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package pay

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// 微信支付签名.
//  parameters: 待签名的参数
//  apiKey:     API密钥
//  fn:         func() hash.Hash, 如果 fn == nil 则默认用 md5.New
func (cli *Client) Sign(data interface{}) string {
	parameters := convertStructToMap(data)
	//	fmt.Println("sign param:", parameters)

	ks := make([]string, 0, len(parameters))
	for k := range parameters {
		if k == "sign" {
			continue
		}
		ks = append(ks, k)
	}
	sort.Strings(ks)

	h := md5.New()
	signature := make([]byte, h.Size()*2)

	for _, k := range ks {
		v := parameters[k]
		if v == "" {
			continue
		}
		h.Write([]byte(k))
		h.Write([]byte{'='})
		h.Write([]byte(v))
		h.Write([]byte{'&'})
	}
	h.Write([]byte("key="))
	h.Write([]byte(cli.apiKey))

	hex.Encode(signature, h.Sum(nil))
	return string(bytes.ToUpper(signature))
}

// 把支付相关函数所用到的struct参数转为map, 以便排序并加sign
func convertStructToMap(data interface{}) map[string]string {
	if check, ok := data.(map[string]string); ok {
		return check
	}

	result := map[string]string{}
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		key, _ := parseTag(field.Tag.Get("xml"))
		if key == "xml" || key == "-" {
			continue
		}
		// 过滤xml跟节点
		if key == "" {
			key = field.Name
		}
		value := fmt.Sprintf("%v", v.FieldByName(field.Name).Interface())
		if value == "" {
			continue
		}
		result[key] = value
	}
	return result
}

// tagOptions contains a slice of tag options
type tagOptions []string

// Has returns true if the given optiton is available in tagOptions
func (t tagOptions) Has(opt string) bool {
	for _, tagOpt := range t {
		if tagOpt == opt {
			return true
		}
	}

	return false
}

// parseTag splits a struct field's tag into its name and a list of options
// which comes after a name. A tag is in the form of: "name,option1,option2".
// The name can be neglectected.
func parseTag(tag string) (string, tagOptions) {
	// tag is one of followings:
	// ""
	// "name"
	// "name,opt"
	// "name,opt,opt2"
	// ",opt"

	res := strings.Split(tag, ",")
	return res[0], res[1:]
}
