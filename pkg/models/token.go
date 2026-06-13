package models

import (
	"encoding/json"
	"time"
)

// Token NPM访问令牌
//
// 表示一个 NPM Registry 的访问令牌，可用于认证 API 请求。
// 注意：Token 值在序列化时会自动掩码，只显示前4位和后4位字符。
type Token struct {
	ID       string    `json:"id"`
	Token    string    `json:"token"`
	Key      string    `json:"key"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
	Readonly bool      `json:"readonly"`
	CIDR     []string  `json:"cidr_whitelist,omitempty"`
}

// maskToken 对 token 值进行掩码处理，只显示前4位和后4位
// 例如 "abcd12345678efgh" → "abcd...efgh"
func maskToken(token string) string {
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + "..." + token[len(token)-4:]
}

// MarshalJSON 自定义 JSON 序列化，对 Token 字段进行掩码处理
func (x Token) MarshalJSON() ([]byte, error) {
	type Alias Token
	return json.Marshal(&struct {
		Token string `json:"token"`
		Alias
	}{
		Token: maskToken(x.Token),
		Alias: Alias(x),
	})
}

// ToJsonString 将 Token 对象转换为 JSON 字符串（Token 值已掩码）
func (x *Token) ToJsonString() string {
	bytes, _ := json.Marshal(x)
	return string(bytes)
}

// TokenCreation 创建令牌的请求参数
//
// 用于创建新的 NPM 访问令牌时提供的信息。
// 需要提供当前用户的密码以确认身份。
type TokenCreation struct {
	Password string   `json:"password"`             // 当前用户的密码
	Readonly bool     `json:"readonly"`             // 是否为只读令牌
	CIDR     []string `json:"cidr_whitelist,omitempty"` // IP 白名单
}

// ToJsonString 将 TokenCreation 对象转换为 JSON 字符串
func (x *TokenCreation) ToJsonString() string {
	bytes, _ := json.Marshal(x)
	return string(bytes)
}