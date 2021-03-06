// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package corp

import (
	"errors"
	"strconv"
)

type DepartmentCreateParameters struct {
	Name     string `json:"name,omitempty"`  // 部门名称。长度限制为1~64个字符
	ParentId int64  `json:"parentid"`        // 父亲部门id。根部门id为1
	Order    int64  `json:"order,omitempty"` // 在父部门中的次序。从1开始，数字越大排序越靠后
	Id       int64  `json:"id,omitempty"`    // 部门ID。用指定部门ID新建部门，不指定此参数时，则自动生成
}

func (para *DepartmentCreateParameters) SetOrder(order int64) {
	para.Order = order
}

func (para *DepartmentCreateParameters) SetId(id int64) {
	para.Id = id
}

// 创建部门
func (clt *Client) DepartmentCreate(para *DepartmentCreateParameters) (id int64, err error) {
	if para == nil {
		err = errors.New("nil parameters")
		return
	}

	var result struct {
		Error
		Id int64 `json:"id"`
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/department/create?access_token="
	if err = clt.PostJSON(incompleteURL, para, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	id = result.Id
	return
}

type DepartmentUpdateParameters struct {
	Id       int64  `json:"id"`                 // 部门id
	Name     string `json:"name,omitempty"`     // 更新的部门名称。长度限制为1~64个字符。修改部门名称时指定该参数
	ParentId int64  `json:"parentid,omitempty"` // 父亲部门id。根部门id为1
	Order    int64  `json:"order,omitempty"`    // 在父部门中的次序。从1开始，数字越大排序越靠后，当数字大于该层部门数时表示移动到最末尾。
}

func (para *DepartmentUpdateParameters) SetParentId(parentId int64) {
	para.ParentId = parentId
}

func (para *DepartmentUpdateParameters) SetOrder(order int64) {
	para.Order = order
}

// 更新部门
func (clt *Client) DepartmentUpdate(para *DepartmentUpdateParameters) (err error) {
	if para == nil {
		err = errors.New("nil parameters")
		return
	}

	var result Error

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/department/update?access_token="
	if err = clt.PostJSON(incompleteURL, para, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result
		return
	}
	return
}

// 删除部门
func (clt *Client) DepartmentDelete(id int64) (err error) {
	var result Error

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/department/delete?id=" +
		strconv.FormatInt(id, 10) + "&access_token="
	if err = clt.GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result
		return
	}
	return
}

type Department struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	ParentId int64  `json:"parentid"`
}

// 获取 rootId 部门的子部门
func (clt *Client) DepartmentList(rootId int64) (departments []Department, err error) {
	var result struct {
		Error
		Departments []Department `json:"department"`
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/department/list?id=" +
		strconv.FormatInt(rootId, 10) + "&access_token="
	if err = clt.GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	departments = result.Departments
	return
}
