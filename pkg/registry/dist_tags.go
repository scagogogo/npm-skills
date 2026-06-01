package registry

import (
	"context"
	"fmt"
)

// GetDistTags 获取指定包的所有分发标签（dist-tags）
//
// Dist-tags 是 NPM 的版本别名机制，最常见的有 "latest"（最新稳定版）、
// "next"（下一个版本）、"beta" 等。此方法返回包的所有 dist-tags 映射。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 要查询的包名称
//
// 返回值:
//   - map[string]string: 标签名到版本号的映射，如 {"latest": "18.2.0", "next": "19.0.0-rc.1"}
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	tags, err := registry.GetDistTags(ctx, "react")
//	if err != nil {
//	    // 处理错误
//	}
//	fmt.Println("Latest:", tags["latest"])
//	fmt.Println("Next:", tags["next"])
func (x *Registry) GetDistTags(ctx context.Context, packageName string) (map[string]string, error) {
	// 使用已有的 GetPackageInformation 方法提取 dist-tags
	pkg, err := x.GetPackageInformation(ctx, packageName)
	if err != nil {
		return nil, fmt.Errorf("failed to get dist-tags for '%s': %w", packageName, err)
	}
	return pkg.DistTags, nil
}

// GetDistTagsAbbreviated 获取指定包的 dist-tags（使用精简 API）
//
// 使用 NPM 的 dist-tags 专用端点，只返回 dist-tags 数据，
// 不获取完整包信息，速度更快。
//
// 参数:
//   - ctx: 上下文
//   - packageName: 包名称（scoped 包需使用 URL 编码，如 @nestjs/core → @nestjs%2Fcore）
//
// 返回值:
//   - map[string]string: 标签名到版本号的映射
//   - error: 如果请求失败则返回错误
func (x *Registry) GetDistTagsAbbreviated(ctx context.Context, packageName string) (map[string]string, error) {
	targetUrl := fmt.Sprintf("%s/-/package/%s/dist-tags", x.options.RegistryURL, packageName)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get dist-tags for '%s': %w", packageName, err)
	}
	return unmarshalJson[map[string]string](bytes)
}