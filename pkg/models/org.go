package models

import "encoding/json"

// Organization NPM组织信息
//
// NPM 组织用于管理团队和包的访问权限，适合企业和团队协作。
type Organization struct {
	Name  string `json:"name"`
	Scope string `json:"scope,omitempty"`
}

// ToJsonString 将 Organization 对象转换为 JSON 字符串
func (x *Organization) ToJsonString() string {
	bytes, _ := json.Marshal(x)
	return string(bytes)
}

// Team NPM团队信息
//
// 团队隶属于组织，用于更细粒度地管理包的访问权限。
type Team struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name,omitempty"`
	Description string `json:"description,omitempty"`
}

// ToJsonString 将 Team 对象转换为 JSON 字符串
func (x *Team) ToJsonString() string {
	bytes, _ := json.Marshal(x)
	return string(bytes)
}

// OrgCreation 创建组织的请求
type OrgCreation struct {
	Name string `json:"name"`
}

// ToJsonString 将 OrgCreation 对象转换为 JSON 字符串
func (x *OrgCreation) ToJsonString() string {
	bytes, _ := json.Marshal(x)
	return string(bytes)
}

// TeamCreation 创建团队的请求
type TeamCreation struct {
	Name string `json:"name"`
}

// ToJsonString 将 TeamCreation 对象转换为 JSON 字符串
func (x *TeamCreation) ToJsonString() string {
	bytes, _ := json.Marshal(x)
	return string(bytes)
}
