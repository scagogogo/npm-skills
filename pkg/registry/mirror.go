package registry

// ------------------------------------------------- --------------------------------------------------------------------

const RegistryUrlYarn = "https://registry.yarnpkg.com"

// NewYarnRegistry 创建使用 Yarn 官方镜像源的 Registry 客户端
//
// Yarn 官方镜像源特点:
//   - 由 Yarn 团队维护的官方 NPM 镜像
//   - 通常提供稳定且与 npmjs.org 同步的包数据
//   - 适合在全球范围内使用，特别是 Yarn 用户
//
// 返回值:
//   - *Registry: 配置为使用 Yarn 镜像源的 Registry 客户端
//
// 使用示例:
//
//	registry := NewYarnRegistry()
//	ctx := context.Background()
//	pkg, err := registry.GetPackageInformation(ctx, "react")
func NewYarnRegistry() *Registry {
	return NewRegistry(NewOptions().SetRegistryURL(RegistryUrlYarn))
}

// ------------------------------------------------- --------------------------------------------------------------------

const RegistryUrlCnpm = "http://r.cnpmjs.org"

// NewCnpmRegistry 创建使用 CNPM 镜像源的 Registry 客户端
//
// CNPM 镜像源特点:
//   - 由中国社区维护的 NPM 镜像
//   - 针对中国大陆网络环境优化
//   - 同步频率较高，但可能略有延迟
//
// 返回值:
//   - *Registry: 配置为使用 CNPM 镜像源的 Registry 客户端
//
// 使用示例:
//
//	registry := NewCnpmRegistry()
//	ctx := context.Background()
//	pkg, err := registry.GetPackageInformation(ctx, "antd")
func NewCnpmRegistry() *Registry {
	return NewRegistry(NewOptions().SetRegistryURL(RegistryUrlCnpm))
}

// ------------------------------------------------- --------------------------------------------------------------------

const RegistryUrlTaoBao = "https://registry.npm.taobao.org"

// NewTaoBaoRegistry 创建使用淘宝 NPM 镜像源的 Registry 客户端
//
// 淘宝镜像源特点:
//   - 由阿里巴巴维护的国内最大 NPM 镜像
//   - 针对中国大陆网络环境优化，访问速度快
//   - 10分钟同步一次，覆盖绝大部分包
//   - 适合在中国大陆地区使用
//
// 返回值:
//   - *Registry: 配置为使用淘宝镜像源的 Registry 客户端
//
// 使用示例:
//
//	registry := NewTaoBaoRegistry()
//	ctx := context.Background()
//	pkg, err := registry.GetPackageInformation(ctx, "vue")
func NewTaoBaoRegistry() *Registry {
	return NewRegistry(NewOptions().SetRegistryURL(RegistryUrlTaoBao))
}

// ------------------------------------------------- --------------------------------------------------------------------

const RegistryUrlNpmjsCom = "https://skimdb.npmjs.com"

// NewNpmjsComRegistry 创建使用 npmjs.com 镜像源的 Registry 客户端
//
// npmjs.com 镜像源特点:
//   - 是 npm 的 CouchDB 数据镜像
//   - 包含较全的元数据信息
//   - 适合需要频繁查询包元数据但不下载包的场景
//
// 返回值:
//   - *Registry: 配置为使用 npmjs.com 镜像源的 Registry 客户端
//
// 使用示例:
//
//	registry := NewNpmjsComRegistry()
//	ctx := context.Background()
//	info, err := registry.GetRegistryInformation(ctx)
func NewNpmjsComRegistry() *Registry {
	return NewRegistry(NewOptions().SetRegistryURL(RegistryUrlNpmjsCom))
}

// ------------------------------------------------- --------------------------------------------------------------------

const RegistryUrlTencent = "http://mirrors.cloud.tencent.com/npm"

// NewTencentRegistry 创建使用腾讯云镜像源的 Registry 客户端
//
// 腾讯云镜像源特点:
//   - 由腾讯云团队维护
//   - 针对腾讯云服务和中国大陆网络环境优化
//   - 同步频率适中，稳定性良好
//   - 适合腾讯云用户和中国大陆地区使用
//
// 返回值:
//   - *Registry: 配置为使用腾讯云镜像源的 Registry 客户端
//
// 使用示例:
//
//	registry := NewTencentRegistry()
//	ctx := context.Background()
//	pkg, err := registry.GetPackageInformation(ctx, "tencent-cloud-sdk")
func NewTencentRegistry() *Registry {
	return NewRegistry(NewOptions().SetRegistryURL(RegistryUrlTencent))
}

// ------------------------------------------------- --------------------------------------------------------------------

const RegistryUrlNpmMirror = "https://registry.npmmirror.com"

// NewNpmMirrorRegistry 创建使用 NPM Mirror 镜像源的 Registry 客户端
//
// NPM Mirror 特点:
//   - 原淘宝 NPM 镜像的新域名
//   - 提供稳定、快速的包获取服务
//   - 针对中国大陆网络环境优化
//   - 适合在中国大陆地区使用
//
// 返回值:
//   - *Registry: 配置为使用 NPM Mirror 镜像源的 Registry 客户端
//
// 使用示例:
//
//	registry := NewNpmMirrorRegistry()
//	ctx := context.Background()
//	pkg, err := registry.GetPackageInformation(ctx, "react")
func NewNpmMirrorRegistry() *Registry {
	return NewRegistry(NewOptions().SetRegistryURL(RegistryUrlNpmMirror))
}

// ------------------------------------------------- --------------------------------------------------------------------

const RegistryUrlHuaWeiCloud = "https://mirrors.huaweicloud.com/repository/npm"

// NewHuaWeiCloudRegistry 创建使用华为云镜像源的 Registry 客户端
//
// 华为云镜像源特点:
//   - 由华为云团队维护
//   - 针对华为云服务和中国大陆网络环境优化
//   - 提供稳定的同步和访问速度
//   - 适合华为云用户和中国大陆地区使用
//
// 返回值:
//   - *Registry: 配置为使用华为云镜像源的 Registry 客户端
//
// 使用示例:
//
//	registry := NewHuaWeiCloudRegistry()
//	ctx := context.Background()
//	pkg, err := registry.GetPackageInformation(ctx, "huaweicloud-sdk")
func NewHuaWeiCloudRegistry() *Registry {
	return NewRegistry(NewOptions().SetRegistryURL(RegistryUrlHuaWeiCloud))
}

// MirrorEntry 表示一个镜像源的元数据信息
//
// 包含镜像源的名称、URL、地理区域和描述信息。
// 此结构体是 CLI 和 SDK 共享镜像元数据的标准数据结构。
//
// 主要字段说明:
//   - Name: 镜像源标识名，用于 --mirror flag
//   - URL: 镜像源的 Registry URL
//   - Region: 地理区域（"Global" 或 "China"）
//   - Description: 镜像源的可读描述
type MirrorEntry struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Region      string `json:"region"`
	Description string `json:"description"`
}

// ListMirrors 返回所有支持的镜像源列表
//
// 返回值:
//   - []MirrorEntry: 包含所有镜像源元数据的切片
//
// 使用示例:
//
//	mirrors := registry.ListMirrors()
//	for _, m := range mirrors {
//	    fmt.Printf("%s: %s (%s)\n", m.Name, m.URL, m.Region)
//	}
func ListMirrors() []MirrorEntry {
	return []MirrorEntry{
		{"official", DefaultRegistryURL, "Global", "NPM Official Registry"},
		{"taobao", RegistryUrlTaoBao, "China", "Taobao NPM Mirror (Alibaba)"},
		{"npm-mirror", RegistryUrlNpmMirror, "China", "NPM Mirror (new Taobao domain)"},
		{"huawei", RegistryUrlHuaWeiCloud, "China", "Huawei Cloud Mirror"},
		{"tencent", RegistryUrlTencent, "China", "Tencent Cloud Mirror"},
		{"cnpm", RegistryUrlCnpm, "China", "CNPM Mirror"},
		{"yarn", RegistryUrlYarn, "Global", "Yarn Official Mirror"},
		{"npmjscom", RegistryUrlNpmjsCom, "Global", "NPM CouchDB Mirror"},
	}
}

// ------------------------------------------------- --------------------------------------------------------------------
