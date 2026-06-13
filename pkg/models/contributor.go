package models

// Contributor 表示 NPM 包的贡献者信息
//
// 包含了贡献者的基本信息，与 Author 和 Maintainer 结构一致。
//
// 字段说明:
//   - Name: 贡献者名称
//   - Email: 贡献者电子邮件地址
//   - Url: 贡献者相关网站链接
type Contributor struct {
	Name  string `json:"name"`  // 贡献者名称
	Email string `json:"email"` // 贡献者电子邮件地址
	Url   string `json:"url"`   // 贡献者相关网站链接
}
