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
	Module      string      `json:"module"`      // ES Module 入口文件
	Types       string      `json:"types"`       // TypeScript 类型声明入口
	Typings     string      `json:"typings"`     // TypeScript 类型声明入口（别名）
	Type        string      `json:"type"`        // 模块类型: "commonjs" 或 "module"
	Exports     interface{} `json:"exports"`     // 现代包导出映射（可以是字符串或条件导出对象）
	Bin         interface{} `json:"bin"`         // 可执行文件映射（可以是字符串或 map[string]string）
	Scripts     Script      `json:"scripts"`     // 脚本命令定义（支持任意 npm script key）
	Repository  *Repository `json:"repository"`  // 代码仓库信息
	Keywords    []string    `json:"keywords"`    // 关键词列表
	Author      *User       `json:"author"`      // 作者信息
	License     string      `json:"license"`     // 许可证类型
	Bugs        *Bugs       `json:"bugs"`        // 问题跟踪链接
	Homepage    string      `json:"homepage"`    // 项目主页
	Funding     interface{} `json:"funding"`     // 资金赞助信息（可以是字符串、对象或数组）

	// 依赖关系，key是依赖的包，value是版本约束
	Dependencies         map[string]string `json:"dependencies"`         // 运行时依赖
	DevDependencies      map[string]string `json:"devDependencies"`      // 开发时依赖
	PeerDependencies     map[string]string `json:"peerDependencies"`     // 对等依赖
	OptionalDependencies map[string]string `json:"optionalDependencies"` // 可选依赖
	BundledDependencies  interface{}       `json:"bundledDependencies"`  // 捆绑依赖（可以是 bool 或 []string）
	BundleDependencies   interface{}       `json:"bundleDependencies"`   // 捆绑依赖（别名）

	// 平台和引擎约束
	Engines map[string]string `json:"engines"` // 引擎版本约束，如 {"node": ">=14", "npm": ">=6"}
	OS      []string          `json:"os"`      // 操作系统兼容性约束，如 ["linux", "darwin", "!win32"]
	CPU     []string          `json:"cpu"`     // CPU 架构兼容性约束，如 ["x64", "arm64"]

	// 元数据字段
	PeerDependenciesMeta     map[string]PeerDependencyMeta `json:"peerDependenciesMeta"`     // 对等依赖的元数据
	OptionalDependenciesMeta map[string]interface{}        `json:"optionalDependenciesMeta"` // 可选依赖的元数据
	Workspaces               interface{}                   `json:"workspaces"`               // 工作区配置（可以是字符串或字符串数组）

	ID          string  `json:"_id"`         // 包ID，通常为 "name@version"
	Dist        *Dist   `json:"dist"`        // 分发信息，包含下载URL和校验和
	From        string  `json:"_from"`       // 包的来源
	NpmVersion  string  `json:"_npmVersion"` // 发布时使用的 npm 版本
	NpmUser     *User   `json:"_npmUser"`    // 发布包的用户信息
	Maintainers []*User `json:"maintainers"` // 维护者列表

	// 目录结构信息
	Directories Directories `json:"directories"`

	Deprecated interface{} `json:"deprecated"` // 弃用说明，string 或 bool 类型

	HasShrinkwrap bool   `json:"_hasShrinkwrap"` // 是否包含 shrinkwrap
	NodeVersion   string `json:"_nodeVersion"`   // 发布时的 Node 版本
}

// Directories 表示 NPM 包的目录结构
type Directories struct {
	Man string `json:"man"` // man 手册页目录
	Lib string `json:"lib"` // 库文件目录
	Bin string `json:"bin"` // 可执行文件目录
}

// PeerDependencyMeta 表示对等依赖的元数据信息
type PeerDependencyMeta struct {
	Optional bool `json:"optional"` // 标记此对等依赖是否为可选
}

// IsDeprecated returns whether this version is deprecated.
// Handles the fact that Deprecated can be a string, bool, or nil.
func (v *Version) IsDeprecated() bool {
	if v.Deprecated == nil {
		return false
	}
	switch d := v.Deprecated.(type) {
	case bool:
		return d
	case string:
		return d != ""
	default:
		return false
	}
}

// DeprecatedMessage returns the deprecation message if this version is deprecated.
// Returns empty string if not deprecated.
func (v *Version) DeprecatedMessage() string {
	if v.Deprecated == nil {
		return ""
	}
	switch d := v.Deprecated.(type) {
	case bool:
		if d {
			return "this version has been deprecated"
		}
		return ""
	case string:
		return d
	default:
		return ""
	}
}
