package models

import "encoding/json"

// Hook NPM Webhook
//
// NPM Webhook 允许在包发布、更新等事件发生时通知外部服务。
type Hook struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Name     string   `json:"name"`
	Endpoint string   `json:"endpoint"`
	Secret   string   `json:"secret,omitempty"`
	Created  string   `json:"created"`
	Updated  string   `json:"updated"`
	Events   []string `json:"events"`
	Package  string   `json:"package,omitempty"`
	Active   bool     `json:"active"`
	Deleted  bool     `json:"deleted,omitempty"`
}

// ToJsonString 将 Hook 对象转换为 JSON 字符串
func (x *Hook) ToJsonString() string {
	bytes, _ := json.Marshal(x)
	return string(bytes)
}

// HookCreation 创建 Webhook 的请求
type HookCreation struct {
	Name     string   `json:"name"`
	Endpoint string   `json:"endpoint"`
	Secret   string   `json:"secret,omitempty"`
	Events   []string `json:"events"`
	Package  string   `json:"package,omitempty"`
	Active   bool     `json:"active"`
}

// ToJsonString 将 HookCreation 对象转换为 JSON 字符串
func (x *HookCreation) ToJsonString() string {
	bytes, _ := json.Marshal(x)
	return string(bytes)
}

// HookUpdate 更新 Webhook 的请求
type HookUpdate struct {
	Endpoint string   `json:"endpoint,omitempty"`
	Secret   string   `json:"secret,omitempty"`
	Events   []string `json:"events,omitempty"`
	Active   *bool    `json:"active,omitempty"`
}

// ToJsonString 将 HookUpdate 对象转换为 JSON 字符串
func (x *HookUpdate) ToJsonString() string {
	bytes, _ := json.Marshal(x)
	return string(bytes)
}

// HookListOptions Webhook 列表查询参数
type HookListOptions struct {
	Package string `json:"package,omitempty"`
	Page    int    `json:"page,omitempty"`
	PerPage int    `json:"per_page,omitempty"`
}

// ToJsonString 将 HookListOptions 对象转换为 JSON 字符串
func (x *HookListOptions) ToJsonString() string {
	bytes, _ := json.Marshal(x)
	return string(bytes)
}
