package corp

import "strconv"

type AgentUserInfo struct {
	UserId string `json:"userid"`
	Status string `json:"status"`
}

// 获取企业号应用
type AgentParameters struct {
	AgentId       int64  `json:"agentid"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	SquareLogoUrl string `json:"square_logo_url"`
	RoundLogoUrl  string `json:"round_logo_url"`

	AllowUserInfos struct {
		User []AgentUserInfo `json:"user"`
	} `json:"allow_userinfos"`
	AllowPartys struct {
		PartyId []int64 `json:"partyid"`
	} `json:"allow_partys"`
	AllowTags struct {
		TagId []int64 `json:"tagid"`
	} `json:"allow_tags"`
	Close              int64  `json:"close"`
	RedirectDomain     string `json:"redirect_domain"`
	ReportLocationFlag int64  `json:"report_location_flag"`
	Isreportuser       int64  `json:"isreportuser"`
	Isreportenter      int64  `json:"isreportenter"`
}

// 获取企业号应用
func (clt *Client) GetAgent(agentId int64) (agent AgentParameters, err error) {
	var result struct {
		Error
		AgentParameters
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/agent/get?agentid=" +
		strconv.FormatInt(agentId, 10) + "&access_token="
	if err = clt.GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	agent = result.AgentParameters
	return
}

type UpdateAgentParameters struct {
	AgentId            int64  `json:"agentid"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	LogoMediaId        string `json:"logo_mediaid"`
	RedirectDomain     string `json:"redirect_domain"`
	ReportLocationFlag int64  `json:"report_location_flag"`
	Isreportuser       int64  `json:"isreportuser"`
	Isreportenter      int64  `json:"isreportenter"`
}

// 设置企业号应用
func (clt *Client) SetAgent(data UpdateAgentParameters) (err error) {
	var result struct {
		Error
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/agent/set?access_token="
	if err = clt.PostJSON(incompleteURL, data, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	return
}

type AgentListItem struct {
	AgentId       int64  `json:"agentid"`
	Name          string `json:"name"`
	SquareLogoUrl string `json:"square_logo_url"`
	RoundLogoUrl  string `json:"round_logo_url"`
}
type AgentListParameters struct {
	AgentList []AgentListItem `json:"agentlist"`
}

// 获取应用概况列表说明
func (clt *Client) GetAgentList() (agent AgentListParameters, err error) {
	var result struct {
		Error
		AgentListParameters
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/agent/list?access_token="
	if err = clt.GetJSON(incompleteURL, &result); err != nil {
		return
	}

	if result.ErrCode != ErrCodeOK {
		err = &result.Error
		return
	}
	agent = result.AgentListParameters
	return
}
