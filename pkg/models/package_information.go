package models

import (
	"encoding/json"
)

// Package 表示一个 NPM 包的完整信息结构
//
// 包含了从 NPM Registry 获取的包的所有元数据信息，包括基本信息、
// 版本历史、维护者、分发标签等多种信息。
//
// 主要字段说明:
//   - Name: 包名称
//   - Description: 包描述
//   - Versions: 包含所有版本信息的映射，键为版本号
//   - Maintainers: 维护者列表
//   - DistTags: 分发标签信息，如 "latest"、"next" 等
//   - Time: 各版本发布时间信息
type Package struct {
	ID             string                 `json:"_id"`
	Rev            string                 `json:"_rev"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	DistTags       map[string]string      `json:"dist-tags"`
	Versions       map[string]Version     `json:"versions"`
	Maintainers    []Maintainer           `json:"maintainers"`
	Time           map[string]string      `json:"time"`
	Repository     Repository             `json:"repository"`
	ReadMe         string                 `json:"readme"`
	ReadMeFilename string                 `json:"readmeFilename"`
	Homepage       string                 `json:"homepage"`
	Bugs           *Bugs                  `json:"bugs"`
	License        string                 `json:"license"`
	Users          map[string]bool        `json:"users"`
	Keywords       []string               `json:"keywords"`
	Author         Author                 `json:"author"`
	Contributors   []Contributor          `json:"contributors"`
	Deprecated     interface{}            `json:"deprecated"` // 弃用说明，可以是 string、bool 或 nil
	Funding        interface{}            `json:"funding"`    // 资金赞助信息
	Attachments    map[string]Attachment  `json:"_attachments"` // 发布时包含的附件（tarball 数据）
}

// Attachment represents a package attachment (tarball data) in the _attachments field.
type Attachment struct {
	ContentType string `json:"content_type"`
	Data        string `json:"data"`   // Base64-encoded tarball data
	Length      int    `json:"length"` // Size in bytes
}

// ToJsonString 将 Package 对象转换为 JSON 字符串
//
// 此方法将包的完整信息序列化为 JSON 格式的字符串，方便存储或传输。
// 如果序列化过程中发生错误，将返回包含错误信息的字符串。
//
// 返回值:
//   - string: JSON 格式的字符串表示
//
// 示例:
//
//	pkg := GetPackageInformation(...)
//	jsonStr := pkg.ToJsonString()
//	fmt.Println(jsonStr)
func (x *Package) ToJsonString() string {
	bytes, _ := json.Marshal(x)
	return string(bytes)
}

// IsDeprecated returns whether this package is deprecated.
// Handles the fact that Deprecated can be a string, bool, or nil.
func (x *Package) IsDeprecated() bool {
	if x.Deprecated == nil {
		return false
	}
	switch d := x.Deprecated.(type) {
	case bool:
		return d
	case string:
		return d != ""
	default:
		return false
	}
}

// DeprecatedMessage returns the deprecation message if this package is deprecated.
// Returns empty string if not deprecated.
func (x *Package) DeprecatedMessage() string {
	if x.Deprecated == nil {
		return ""
	}
	switch d := x.Deprecated.(type) {
	case bool:
		if d {
			return "this package has been deprecated"
		}
		return ""
	case string:
		return d
	default:
		return ""
	}
}

// Maintainer 表示 NPM 包的维护者信息
//
// 字段说明:
//   - Name: 维护者名称
//   - Email: 维护者电子邮件地址
//   - Url: 维护者相关网站链接
type Maintainer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Url   string `json:"url"`
}

// Repository 类型定义在 repository.go 文件中
// 表示 NPM 包的代码仓库信息
//
// 字段说明:
//   - Type: 仓库类型，通常为 "git"
//   - Url: 仓库 URL 地址
//   - Directory: 包在仓库中的目录位置（对于 monorepo 项目）

// Author 类型定义在 author.go 文件中
// 表示 NPM 包的作者信息
//
// 字段说明:
//   - Name: 作者名称
//   - Email: 作者电子邮件地址
//   - Url: 作者相关网站链接

// Contributor 类型定义在 contributor.go 文件中
// 表示 NPM 包的贡献者信息
//
// 字段说明:
//   - Name: 贡献者名称
//   - Email: 贡献者电子邮件地址
//   - Url: 贡献者相关网站链接
