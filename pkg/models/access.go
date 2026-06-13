package models

import "encoding/json"

// PackageAccess 包的访问权限信息
//
// 包含包的访问级别设置，如公开(public)或私有(restricted)。
type PackageAccess struct {
	Package string            `json:"package"`
	Access  map[string]string `json:"access"` // 如 {"read": "public", "write": "restricted"}
}

// ToJsonString 将 PackageAccess 对象转换为 JSON 字符串
func (x *PackageAccess) ToJsonString() string {
	bytes, _ := json.Marshal(x)
	return string(bytes)
}

// PackageAccessUpdate 更新包访问权限的请求
//
// 用于设置包的访问级别，可选值为 "public" 或 "restricted"。
type PackageAccessUpdate struct {
	Access string `json:"access"` // "public" 或 "restricted"
}

// ToJsonString 将 PackageAccessUpdate 对象转换为 JSON 字符串
func (x *PackageAccessUpdate) ToJsonString() string {
	bytes, _ := json.Marshal(x)
	return string(bytes)
}

// Collaborator 包的协作者信息
//
// 表示有权访问包的用户或团队及其权限级别。
type Collaborator struct {
	Name        string `json:"name"`
	Email       string `json:"email,omitempty"`
	Permissions string `json:"permissions"` // "read" 或 "write"
}

// ToJsonString 将 Collaborator 对象转换为 JSON 字符串
func (x *Collaborator) ToJsonString() string {
	bytes, _ := json.Marshal(x)
	return string(bytes)
}

// Permission 权限级别
//
// 定义协作者对包的访问权限级别。
type Permission string

const (
	// PermissionRead 只读权限，允许协作者查看包信息
	PermissionRead Permission = "read"
	// PermissionWrite 写权限，允许协作者发布新版本和管理包
	PermissionWrite Permission = "write"
)
