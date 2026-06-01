package models

// Version 表示 NPM 包的特定版本信息
//
// 此结构包含了 NPM 包某个特定版本的详细信息，包括基本信息、
// 依赖关系、分发信息等。通常作为 Package 结构体中 Versions 字段的值。
//
// 主要字段说明:
//   - Name: 包名称
//   - Version: 版本号，如 "1.0.0"
//   - Description: 版本描述
//   - Dependencies: 运行时依赖，键为依赖包名，值为版本约束
//   - DevDependencies: 开发时依赖
//   - Dist: 分发信息，包含下载 URL 和校验和
type Version struct {
	Name        string      `json:"name"`        // 包名称
	Version     string      `json:"version"`     // 版本号，如 "1.0.0"
	Description string      `json:"description"` // 版本描述
	Main        string      `json:"main"`        // 主入口文件
	Scripts     Script      `json:"scripts"`      // 脚本命令定义（支持任意 npm script key）
	Repository  *Repository `json:"repository"`  // 代码仓库信息
	Keywords    []string    `json:"keywords"`    // 关键词列表
	Author      *User       `json:"author"`      // 作者信息
	License     string      `json:"license"`     // 许可证类型
	Bugs        *Bugs       `json:"bugs"`        // 问题跟踪链接
	Homepage    string      `json:"homepage"`    // 项目主页

	// 依赖关系，key是依赖的包，value是版本约束
	Dependencies    map[string]string `json:"dependencies"`    // 运行时依赖
	DevDependencies map[string]string `json:"devDependencies"` // 开发时依赖

	ID          string  `json:"_id"`         // 包ID，通常为 "name@version"
	Dist        *Dist   `json:"dist"`        // 分发信息，包含下载URL和校验和
	From        string  `json:"_from"`       // 包的来源
	NpmVersion  string  `json:"_npmVersion"` // 发布时使用的 npm 版本
	NpmUser     *User   `json:"_npmUser"`    // 发布包的用户信息
	Maintainers []*User `json:"maintainers"` // 维护者列表

	// 目录结构信息
	Directories struct {
	} `json:"directories"`

	Deprecated interface{} `json:"deprecated"` // 弃用说明，string 或 bool 类型
}

// Script 类型定义在其他文件中
// 表示 NPM 包的脚本命令定义
//
// 包含 NPM 包的各种生命周期脚本，如 install、test、build 等
// 这些脚本可以通过 npm run [script-name] 来执行

// Dist 类型定义在其他文件中
// 表示 NPM 包的分发信息
//
// 包含下载包所需的信息，如 tarball URL 和 校验和
// 主要字段:
//   - Shasum: 包的 SHA 校验和
//   - Tarball: 包的下载 URL
//   - Integrity: 完整性校验值，通常为 SRI 格式
//   - NpmSignature: NPM 签名信息

// User 类型定义在其他文件中
// 表示与 NPM 包相关的用户信息
//
// 可用于表示作者、维护者或发布者
// 主要字段:
//   - Name: 用户名
//   - Email: 电子邮件地址
//   - URL: 用户网站或主页

// Bugs 类型定义在其他文件中
// 表示 NPM 包的问题跟踪信息
//
// 通常包含问题跟踪系统的 URL
// 主要字段:
//   - URL: 问题跟踪系统 URL
