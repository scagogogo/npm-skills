package registry

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/scagogogo/npm-skills/pkg/models"
)

// defaultDownloadStatsBaseURL 是 NPM 下载统计 API 的默认基础 URL
const defaultDownloadStatsBaseURL = "https://api.npmjs.org/downloads"

// ErrDownloadStatsNotAvailable 当私有仓库不支持下载统计 API 时返回此错误
var ErrDownloadStatsNotAvailable = fmt.Errorf("download stats API not available for this registry; set Options.DownloadStatsURL if your registry provides a compatible API")

// downloadStatsURL 返回下载统计 API 的基础 URL
//
// 对于私有仓库（非公共镜像），如果没有显式设置 DownloadStatsURL，
// 将返回空字符串，调用方应据此跳过下载统计请求。
func (x *Registry) downloadStatsURL() string {
	if x.options.DownloadStatsURL != "" {
		return x.options.DownloadStatsURL
	}
	if x.IsPrivateRegistry() {
		return ""
	}
	return defaultDownloadStatsBaseURL
}

// requireDownloadStatsURL 检查下载统计 API 是否可用
func (x *Registry) requireDownloadStatsURL() error {
	if x.downloadStatsURL() == "" {
		return ErrDownloadStatsNotAvailable
	}
	return nil
}

// GetDownloadStats 获取指定 NPM 包的下载统计信息（单个包、预定义周期）
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 要查询的包名称
//   - period: 统计周期，例如 "last-day", "last-week", "last-month"
//
// 返回值:
//   - *models.DownloadStats: 下载统计信息
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	stats, err := registry.GetDownloadStats(ctx, "react", "last-week")
func (x *Registry) GetDownloadStats(ctx context.Context, packageName, period string) (*models.DownloadStats, error) {
	if err := x.requireDownloadStatsURL(); err != nil {
		return nil, err
	}
	if err := requirePackageName(packageName); err != nil {
		return nil, err
	}
	targetUrl := fmt.Sprintf("%s/point/%s/%s", x.downloadStatsURL(), period, url.PathEscape(packageName))
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get download stats for '%s' in period '%s': %w", packageName, period, err)
	}
	return unmarshalJson[*models.DownloadStats](bytes)
}

// GetDownloadRangeStats 获取指定 NPM 包的每日下载统计（区间数据）
//
// 返回每日下载次数数组，适用于绘制下载趋势图。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 要查询的包名称
//   - period: 统计周期，"last-day"、"last-week"、"last-month"
//
// 返回值:
//   - *models.DownloadRangeStats: 包含每日下载数据的区间统计信息
//   - error: 如果请求失败则返回错误
func (x *Registry) GetDownloadRangeStats(ctx context.Context, packageName, period string) (*models.DownloadRangeStats, error) {
	if err := x.requireDownloadStatsURL(); err != nil {
		return nil, err
	}
	if err := requirePackageName(packageName); err != nil {
		return nil, err
	}
	targetUrl := fmt.Sprintf("%s/range/%s/%s", x.downloadStatsURL(), period, url.PathEscape(packageName))
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get download range stats for '%s' in period '%s': %w", packageName, period, err)
	}
	return unmarshalJson[*models.DownloadRangeStats](bytes)
}

// GetDownloadStatsByDateRange 获取指定日期范围的下载统计
//
// 使用自定义日期范围而非预定义周期。日期格式为 YYYY-MM-DD。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 要查询的包名称
//   - start: 开始日期，格式 YYYY-MM-DD
//   - end: 结束日期，格式 YYYY-MM-DD
//
// 返回值:
//   - *models.DownloadStats: 下载统计信息
//   - error: 如果请求失败则返回错误
func (x *Registry) GetDownloadStatsByDateRange(ctx context.Context, packageName, start, end string) (*models.DownloadStats, error) {
	if err := x.requireDownloadStatsURL(); err != nil {
		return nil, err
	}
	if err := requirePackageName(packageName); err != nil {
		return nil, err
	}
	targetUrl := fmt.Sprintf("%s/point/%s:%s/%s", x.downloadStatsURL(), start, end, url.PathEscape(packageName))
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get download stats for '%s' from %s to %s': %w", packageName, start, end, err)
	}
	return unmarshalJson[*models.DownloadStats](bytes)
}

// GetDownloadRangeStatsByDateRange 获取指定日期范围的每日下载统计（区间数据）
//
// 与 GetDownloadStatsByDateRange 不同，此方法返回每日的下载明细数据，
// 适用于绘制下载趋势图或进行时间序列分析。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageName: 要查询的包名称
//   - start: 开始日期，格式 YYYY-MM-DD
//   - end: 结束日期，格式 YYYY-MM-DD
//
// 返回值:
//   - *models.DownloadRangeStats: 包含每日下载数据的区间统计信息
//   - error: 如果请求失败则返回错误
func (x *Registry) GetDownloadRangeStatsByDateRange(ctx context.Context, packageName, start, end string) (*models.DownloadRangeStats, error) {
	if err := x.requireDownloadStatsURL(); err != nil {
		return nil, err
	}
	if err := requirePackageName(packageName); err != nil {
		return nil, err
	}
	targetUrl := fmt.Sprintf("%s/range/%s:%s/%s", x.downloadStatsURL(), start, end, url.PathEscape(packageName))
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get download range stats for '%s' from %s to %s: %w", packageName, start, end, err)
	}
	return unmarshalJson[*models.DownloadRangeStats](bytes)
}

// GetBulkDownloadStats 批量获取多个包的下载统计
//
// 自动将超过 128 个包的请求分批处理。
//
// 参数:
//   - ctx: 上下文
//   - packageNames: 包名称切片
//   - period: 统计周期，"last-day"、"last-week"、"last-month"
//
// 返回值:
//   - map[string]*models.DownloadStats: 包名到下载统计的映射
//   - error: 如果请求失败则返回错误
func (x *Registry) GetBulkDownloadStats(ctx context.Context, packageNames []string, period string) (map[string]*models.DownloadStats, error) {
	if err := x.requireDownloadStatsURL(); err != nil {
		return nil, err
	}
	if len(packageNames) == 0 {
		return nil, fmt.Errorf("packageNames must not be empty")
	}

	result := make(map[string]*models.DownloadStats)
	for i := 0; i < len(packageNames); i += 128 {
		end := i + 128
		if end > len(packageNames) {
			end = len(packageNames)
		}
		batch := packageNames[i:end]

		escaped := url.QueryEscape(strings.Join(batch, ","))
		targetUrl := fmt.Sprintf("%s/point/%s/%s", x.downloadStatsURL(), period, escaped)
		bytes, err := x.getBytes(ctx, targetUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to get bulk download stats in period '%s': %w", period, err)
		}
		batchResult, err := unmarshalJson[map[string]*models.DownloadStats](bytes)
		if err != nil {
			return nil, err
		}
		for k, v := range batchResult {
			result[k] = v
		}
	}
	return result, nil
}

// GetBulkDownloadRangeStats 批量获取多个包的每日下载统计
//
// 自动将超过 128 个包的请求分批处理。
//
// 参数:
//   - ctx: 上下文
//   - packageNames: 包名称切片
//   - period: 统计周期
//
// 返回值:
//   - map[string]*models.DownloadRangeStats: 包名到区间统计的映射
//   - error: 如果请求失败则返回错误
func (x *Registry) GetBulkDownloadRangeStats(ctx context.Context, packageNames []string, period string) (map[string]*models.DownloadRangeStats, error) {
	if err := x.requireDownloadStatsURL(); err != nil {
		return nil, err
	}
	if len(packageNames) == 0 {
		return nil, fmt.Errorf("packageNames must not be empty")
	}

	result := make(map[string]*models.DownloadRangeStats)
	for i := 0; i < len(packageNames); i += 128 {
		end := i + 128
		if end > len(packageNames) {
			end = len(packageNames)
		}
		batch := packageNames[i:end]

		escaped := url.QueryEscape(strings.Join(batch, ","))
		targetUrl := fmt.Sprintf("%s/range/%s/%s", x.downloadStatsURL(), period, escaped)
		bytes, err := x.getBytes(ctx, targetUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to get bulk download range stats in period '%s': %w", period, err)
		}
		batchResult, err := unmarshalJson[map[string]*models.DownloadRangeStats](bytes)
		if err != nil {
			return nil, err
		}
		for k, v := range batchResult {
			result[k] = v
		}
	}
	return result, nil
}

// GetBulkDownloadStatsByDateRange 批量获取多个包在自定义日期范围内的下载统计
//
// 自动将超过 128 个包的请求分批处理。
//
// 参数:
//   - ctx: 上下文
//   - packageNames: 包名称切片
//   - start: 开始日期，格式 YYYY-MM-DD
//   - end: 结束日期，格式 YYYY-MM-DD
//
// 返回值:
//   - map[string]*models.DownloadStats: 包名到下载统计的映射
//   - error: 如果请求失败则返回错误
func (x *Registry) GetBulkDownloadStatsByDateRange(ctx context.Context, packageNames []string, start, end string) (map[string]*models.DownloadStats, error) {
	if err := x.requireDownloadStatsURL(); err != nil {
		return nil, err
	}
	if len(packageNames) == 0 {
		return nil, fmt.Errorf("packageNames must not be empty")
	}

	result := make(map[string]*models.DownloadStats)
	for i := 0; i < len(packageNames); i += 128 {
		batchEnd := i + 128
		if batchEnd > len(packageNames) {
			batchEnd = len(packageNames)
		}
		batch := packageNames[i:batchEnd]

		escaped := url.QueryEscape(strings.Join(batch, ","))
		targetUrl := fmt.Sprintf("%s/point/%s:%s/%s", x.downloadStatsURL(), start, end, escaped)
		bytes, err := x.getBytes(ctx, targetUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to get bulk download stats from %s to %s: %w", start, end, err)
		}
		batchResult, err := unmarshalJson[map[string]*models.DownloadStats](bytes)
		if err != nil {
			return nil, err
		}
		for k, v := range batchResult {
			result[k] = v
		}
	}
	return result, nil
}

// GetBulkDownloadRangeStatsByDateRange 批量获取多个包在自定义日期范围内的每日下载统计
//
// 自动将超过 128 个包的请求分批处理。
//
// 参数:
//   - ctx: 上下文
//   - packageNames: 包名称切片
//   - start: 开始日期，格式 YYYY-MM-DD
//   - end: 结束日期，格式 YYYY-MM-DD
//
// 返回值:
//   - map[string]*models.DownloadRangeStats: 包名到区间统计的映射
//   - error: 如果请求失败则返回错误
func (x *Registry) GetBulkDownloadRangeStatsByDateRange(ctx context.Context, packageNames []string, start, end string) (map[string]*models.DownloadRangeStats, error) {
	if err := x.requireDownloadStatsURL(); err != nil {
		return nil, err
	}
	if len(packageNames) == 0 {
		return nil, fmt.Errorf("packageNames must not be empty")
	}

	result := make(map[string]*models.DownloadRangeStats)
	for i := 0; i < len(packageNames); i += 128 {
		batchEnd := i + 128
		if batchEnd > len(packageNames) {
			batchEnd = len(packageNames)
		}
		batch := packageNames[i:batchEnd]

		escaped := url.QueryEscape(strings.Join(batch, ","))
		targetUrl := fmt.Sprintf("%s/range/%s:%s/%s", x.downloadStatsURL(), start, end, escaped)
		bytes, err := x.getBytes(ctx, targetUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to get bulk download range stats from %s to %s: %w", start, end, err)
		}
		batchResult, err := unmarshalJson[map[string]*models.DownloadRangeStats](bytes)
		if err != nil {
			return nil, err
		}
		for k, v := range batchResult {
			result[k] = v
		}
	}
	return result, nil
}
