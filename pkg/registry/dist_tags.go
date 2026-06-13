package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
	return x.GetDistTagsAbbreviated(ctx, packageName)
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
	targetUrl := fmt.Sprintf("%s/-/package/%s/dist-tags", x.options.RegistryURL, encodePackageName(packageName))
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get dist-tags for '%s': %w", packageName, err)
	}
	return unmarshalJson[map[string]string](bytes)
}

// GetDistTag 获取指定包的特定分发标签
//
// 使用 dist-tags 专用端点获取单个标签的版本号。
// 如果标签不存在，Registry 会返回 404 错误。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 包名称
//   - tag: 标签名称，如 "latest", "next", "beta" 等
//
// 返回值:
//   - string: 标签对应的版本号
//   - error: 如果请求失败或标签不存在则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	version, err := registry.GetDistTag(ctx, "react", "latest")
//	if err != nil {
//	    // 处理错误
//	}
//	fmt.Println("Latest:", version)
func (x *Registry) GetDistTag(ctx context.Context, packageName, tag string) (string, error) {
	targetUrl := fmt.Sprintf("%s/-/package/%s/dist-tags/%s", x.options.RegistryURL, encodePackageName(packageName), tag)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return "", fmt.Errorf("failed to get dist-tag '%s' for '%s': %w", tag, packageName, err)
	}

	// 尝试直接解析为 string（标准 NPM Registry 格式）
	var version string
	if err := json.Unmarshal(bytes, &version); err == nil {
		return version, nil
	}

	// 解析为 object 以判断响应类型
	var obj map[string]interface{}
	if err := json.Unmarshal(bytes, &obj); err != nil {
		return "", fmt.Errorf("failed to parse dist-tag response for '%s/%s': %w", packageName, tag, err)
	}

	// 检查是否为错误响应（如 Verdaccio 的 {"error": "File not found"}）
	if errMsg, ok := obj["error"]; ok {
		// 单标签端点不支持，回退到获取全部标签再提取
		allTags, fallbackErr := x.GetDistTagsAbbreviated(ctx, packageName)
		if fallbackErr != nil {
			return "", fmt.Errorf("dist-tag '%s' not found for '%s': %v", tag, packageName, errMsg)
		}
		if v, ok := allTags[tag]; ok {
			return v, nil
		}
		return "", fmt.Errorf("dist-tag '%s' does not exist for '%s'", tag, packageName)
	}

	// 某些私有仓库返回 {"tag": "version"} 格式
	if v, ok := obj[tag]; ok {
		if s, ok := v.(string); ok {
			return s, nil
		}
	}

	// 如果只有一个值，返回它
	for _, v := range obj {
		if s, ok := v.(string); ok {
			return s, nil
		}
	}

	return "", fmt.Errorf("failed to parse dist-tag response for '%s/%s': unexpected format", packageName, tag)
}

// SetDistTag 设置指定包的分发标签（指向特定版本）
//
// 需要认证Token。将指定的标签指向给定的版本号。
// 如果标签已存在，则更新；如果不存在，则创建。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 包名称
//   - tag: 标签名称，如 "latest", "next", "beta" 等
//   - version: 标签要指向的版本号
//
// 返回值:
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	options := NewOptions().SetToken("npm_xxxxx")
//	registry := NewRegistry(options)
//	ctx := context.Background()
//	err := registry.SetDistTag(ctx, "my-package", "next", "2.0.0-rc.1")
//	if err != nil {
//	    // 处理错误
//	}
func (x *Registry) SetDistTag(ctx context.Context, packageName, tag, version string) error {
	if err := x.requireToken(); err != nil {
		return err
	}

	targetUrl := fmt.Sprintf("%s/-/package/%s/dist-tags/%s", x.options.RegistryURL, encodePackageName(packageName), tag)
	// NPM Registry 期望 Body 为 JSON 字符串，如 "1.0.0"
	// 使用 sendJSON 以确保设置 Content-Type: application/json
	_, err := x.sendJSON(ctx, http.MethodPut, targetUrl, version, http.StatusOK, http.StatusCreated)
	if err != nil {
		return fmt.Errorf("failed to set dist-tag '%s' to '%s' for '%s': %w", tag, version, packageName, err)
	}
	return nil
}

// DeleteDistTag 删除指定包的分发标签
//
// 需要认证Token。删除指定的标签，不影响版本本身。
// 注意：删除 "latest" 标签可能会导致问题，谨慎操作。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 包名称
//   - tag: 要删除的标签名称
//
// 返回值:
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	options := NewOptions().SetToken("npm_xxxxx")
//	registry := NewRegistry(options)
//	ctx := context.Background()
//	err := registry.DeleteDistTag(ctx, "my-package", "beta")
//	if err != nil {
//	    // 处理错误
//	}
func (x *Registry) DeleteDistTag(ctx context.Context, packageName, tag string) error {
	if err := x.requireToken(); err != nil {
		return err
	}

	targetUrl := fmt.Sprintf("%s/-/package/%s/dist-tags/%s", x.options.RegistryURL, encodePackageName(packageName), tag)
	_, err := x.sendRequest(ctx, http.MethodDelete, targetUrl, nil, http.StatusOK, http.StatusNoContent, http.StatusCreated)
	if err != nil {
		return fmt.Errorf("failed to delete dist-tag '%s' for '%s': %w", tag, packageName, err)
	}
	return nil
}

// SetDistTags 批量设置指定包的分发标签
//
// 需要认证Token。使用 POST 方法合并更新多个标签，
// 不会影响未在请求中包含的已有标签。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 包名称
//   - tags: 标签到版本号的映射，如 {"next": "2.0.0-rc.1", "beta": "1.9.0-beta.3"}
//
// 返回值:
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	options := NewOptions().SetToken("npm_xxxxx")
//	registry := NewRegistry(options)
//	ctx := context.Background()
//	err := registry.SetDistTags(ctx, "my-package", map[string]string{
//	    "next": "2.0.0-rc.1",
//	    "beta": "1.9.0-beta.3",
//	})
//	if err != nil {
//	    // 处理错误
//	}
func (x *Registry) SetDistTags(ctx context.Context, packageName string, tags map[string]string) error {
	if err := x.requireToken(); err != nil {
		return err
	}

	targetUrl := fmt.Sprintf("%s/-/package/%s/dist-tags", x.options.RegistryURL, encodePackageName(packageName))
	_, err := x.sendJSON(ctx, http.MethodPost, targetUrl, tags, http.StatusOK, http.StatusCreated)
	if err != nil {
		return fmt.Errorf("failed to set dist-tags for '%s': %w", packageName, err)
	}
	return nil
}
