package models

// Repository 表示 NPM 包的代码仓库信息
//
// 包含仓库的类型、URL 和子目录信息。
// 对于 monorepo 项目，Directory 字段指定包在仓库中的位置。
//
// 主要字段说明:
//   - Type: 仓库类型，通常为 "git"
//   - URL: 仓库 URL 地址
//   - Directory: 包在仓库中的子目录（monorepo 项目常用）
type Repository struct {
	Type      string `json:"type"`      // 仓库类型，如 "git"
	URL       string `json:"url"`       // 仓库 URL 地址
	Directory string `json:"directory"` // 包在仓库中的目录位置（对于 monorepo 项目，如 "packages/core"）
}
