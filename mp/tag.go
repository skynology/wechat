package mp

import "errors"

const TagCountLimit = 100 // 一个公众账号，最多支持创建100个标签

// 用户标签
type Tag struct {
	Id        int64  `json:"id"`    // 标签id, 由微信分配
	Name      string `json:"name"`  // 标签名字, UTF8编码
	UserCount int    `json:"count"` // 标签内用户数量
}

// 创建标签.
//  name: 标签名字（30个字符以内）.
func (clt *Client) CreateTag(name string) (tag *Tag, err error) {
	if name == "" {
		err = errors.New(`name == ""`)
		return
	}

	var request struct {
		Tag struct {
			Name string `json:"name"`
		} `json:"tag"`
	}
	request.Tag.Name = name

	var result struct {
		Error
		Tag `json:"tag"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/tags/create?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}

	tag = &result.Tag
	return
}

// 查询所有标签.
func (clt *Client) ListTag() (tags []Tag, err error) {
	var result = struct {
		Error
		Tags []Tag `json:"tags"`
	}{
		Tags: make([]Tag, 0),
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/tags/get?access_token="
	if err = clt.GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	tags = result.Tags
	return
}

// 修改标签名.
//  name: 标签名字（30个字符以内）.
func (clt *Client) TagRename(tagId int64, newName string) (err error) {
	if newName == "" {
		return errors.New(`newName == ""`)
	}

	var request struct {
		Tag struct {
			Id   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"tag"`
	}
	request.Tag.Id = tagId
	request.Tag.Name = newName

	var result Error

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/tags/update?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result
		return
	}
	return
}

// 删除标签.
//  name: 标签名字（30个字符以内）.
func (clt *Client) TagDelete(tagId int64) (err error) {

	var request struct {
		TagId int64 `json:"tagid"`
	}
	request.TagId = tagId

	var result Error

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/tags/delete?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result
		return
	}
	return
}

// 标签下的用户openid列表
type TagUserListResult struct {
	GotCount int `json:"count"` // 拉取的OPENID个数，最大值为10000

	Data struct {
		OpenId []string `json:"openid,omitempty"`
	} `json:"data"` // 列表数据，OPENID的列表

	// 拉取列表的后一个用户的OPENID, 如果 next_openid == "" 则表示没有了用户数据
	NextOpenId string `json:"next_openid"`
}

// 获取关注者列表, 每次最多能获取 10000 个用户, 如果 beginOpenId == "" 则表示从头获取
func (clt *Client) TagUserList(tagId int64, beginOpenId string) (data *TagUserListResult, err error) {
	var result struct {
		Error
		TagUserListResult
	}

	var incompleteURL string
	incompleteURL = "https://api.weixin.qq.com/cgi-bin/user/tag/get?access_token="

	var request struct {
		TagId      int64  `json:"tagid"`
		NextOpenId string `json:"next_openid"`
	}
	request.TagId = tagId
	request.NextOpenId = beginOpenId

	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	data = &result.TagUserListResult
	return
}

// 批量移动用户标签.
func (clt *Client) UserBatchMoveToTag(openIdList []string, toTagId int64) (err error) {
	var request = struct {
		OpenIdList []string `json:"openid_list"`
		ToTagId    int64    `json:"tagid"`
	}{
		OpenIdList: openIdList,
		ToTagId:    toTagId,
	}

	var result Error

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/tags/members/batchtagging?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result
		return
	}
	return
}

// 批量删除用户标签.
func (clt *Client) UserBatchRemoveFromTag(openIdList []string, toTagId int64) (err error) {
	var request = struct {
		OpenIdList []string `json:"openid_list"`
		ToTagId    int64    `json:"tagid"`
	}{
		OpenIdList: openIdList,
		ToTagId:    toTagId,
	}

	var result Error

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/tags/members/batchuntagging?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result
		return
	}
	return
}

// 获取用户标签列表
func (clt *Client) UserTagIdList(openid string) (tagIds []int64, err error) {
	var request = struct {
		OpenId string `json:openid"`
	}{
		OpenId: openid,
	}

	var result struct {
		Error
		TagIds []int64 `json:"tagid_list"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/tags/getidlist?access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}

	tagIds = result.TagIds
	return
}
