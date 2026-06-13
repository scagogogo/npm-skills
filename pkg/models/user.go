package models

import "encoding/json"

// LoginResult 登录结果
//
// NPM Registry 登录或注册成功后返回的结果，包含认证 Token。
// 注意：不同 Registry 的 ok 字段格式可能不同：
// - 官方 NPM Registry: ok 为 true（布尔值）
// - Verdaccio 等私有仓库: ok 为描述字符串（如 "user 'xxx' created"）
//
// Token 值在序列化时会自动掩码，只显示前4位和后4位字符。
type LoginResult struct {
	ID    string `json:"id"`
	Rev   string `json:"rev"`
	Token string `json:"token"`
	Ok    OkBool `json:"ok"`
}

// MarshalJSON 自定义 JSON 序列化，对 Token 字段进行掩码处理
func (l LoginResult) MarshalJSON() ([]byte, error) {
	type Alias LoginResult
	return json.Marshal(&struct {
		Token string `json:"token"`
		Alias
	}{
		Token: maskLoginToken(l.Token),
		Alias: Alias(l),
	})
}

// maskLoginToken 对登录 token 进行掩码
func maskLoginToken(token string) string {
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + "..." + token[len(token)-4:]
}

// OkBool 是一个兼容布尔值和字符串的 JSON 解析类型
// 某些 Registry 返回 ok: true，某些返回 ok: "user created"
type OkBool struct {
	value  bool
	str    string
	isBool bool
}

// Bool 返回布尔值
func (o *OkBool) Bool() bool {
	return o.value
}

// String 返回字符串值
func (o *OkBool) String() string {
	if o.isBool {
		if o.value {
			return "true"
		}
		return "false"
	}
	return o.str
}

// UnmarshalJSON 实现 json.Unmarshaler 接口，兼容布尔值和字符串
func (o *OkBool) UnmarshalJSON(data []byte) error {
	// 尝试解析为布尔值
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		o.value = b
		o.isBool = true
		return nil
	}
	// 尝试解析为字符串
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	o.str = s
	o.isBool = false
	return nil
}

// MarshalJSON 实现 json.Marshaler 接口
func (o OkBool) MarshalJSON() ([]byte, error) {
	if o.isBool {
		return json.Marshal(o.value)
	}
	return json.Marshal(o.str)
}

// ToJsonString 将 LoginResult 对象转换为 JSON 字符串
func (x *LoginResult) ToJsonString() string {
	bytes, _ := json.Marshal(x)
	return string(bytes)
}

// UserCreation 创建用户请求
//
// 用于注册新用户时提交的信息。
type UserCreation struct {
	ID       string `json:"_id"`     // 固定格式: "org.couchdb.user:<name>"
	Name     string `json:"name"`    // 用户名
	Password string `json:"password"` // 密码
	Email    string `json:"email"`    // 邮箱地址
	Type     string `json:"type"`     // 固定值: "user"
}

// ToJsonString 将 UserCreation 对象转换为 JSON 字符串
func (x *UserCreation) ToJsonString() string {
	bytes, _ := json.Marshal(x)
	return string(bytes)
}

// UserProfile 用户信息
//
// 从 Registry 获取的用户信息。
type UserProfile struct {
	ID            string `json:"_id"`               // 格式: "org.couchdb.user:<name>"
	Rev           string `json:"_rev"`              // CouchDB 文档修订版本
	Name          string `json:"name"`              // 用户名
	Email         string `json:"email"`             // 邮箱地址
	Type          string `json:"type"`              // 类型，通常为 "user"
	EmailVerified bool   `json:"email_verified"`    // 邮箱是否已验证
	Avatar        string `json:"avatar,omitempty"`  // 头像 URL
	GitHub        string `json:"github,omitempty"`  // 关联的 GitHub 账户
	Created       string `json:"created,omitempty"` // 账户创建时间
	Updated       string `json:"updated,omitempty"` // 账户更新时间
}

// ToJsonString 将 UserProfile 对象转换为 JSON 字符串
func (x *UserProfile) ToJsonString() string {
	bytes, _ := json.Marshal(x)
	return string(bytes)
}
