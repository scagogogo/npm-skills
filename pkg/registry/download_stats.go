package registry

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/scagogogo/npm-crawler/pkg/models"
)

// downloadStatsBaseURL 是 NPM 下载统计 API 的基础 URL
const downloadStatsBaseURL = "https://api.npmjs.org/downloads"

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
//	if err != nil {
//	    // 处理错误
//	}
//	fmt.Println("下载次数:", stats.Downloads)
func (x *Registry) GetDownloadStats(ctx context.Context, packageName, period string) (*models.DownloadStats, error) {
	targetUrl := fmt.Sprintf("%s/point/%s/%s", downloadStatsBaseURL, period, packageName)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, err
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
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	stats, err := registry.GetDownloadRangeStats(ctx, "react", "last-week")
//	if err != nil {
//	    // 处理错误
//	}
//	for _, day := range stats.Downloads {
//	    fmt.Printf("%s: %d\n", day.Day, day.Downloads)
//	}
func (x *Registry) GetDownloadRangeStats(ctx context.Context, packageName, period string) (*models.DownloadRangeStats, error) {
	targetUrl := fmt.Sprintf("%s/range/%s/%s", downloadStatsBaseURL, period, packageName)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, err
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
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	stats, err := registry.GetDownloadStatsByDateRange(ctx, "react", "2024-01-01", "2024-01-31")
func (x *Registry) GetDownloadStatsByDateRange(ctx context.Context, packageName, start, end string) (*models.DownloadStats, error) {
	targetUrl := fmt.Sprintf("%s/point/%s:%s/%s", downloadStatsBaseURL, start, end, packageName)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, err
	}
	return unmarshalJson[*models.DownloadStats](bytes)
}

// GetBulkDownloadStats 批量获取多个包的下载统计（最多 128 个包）
//
// 适用于比较多个包的下载量，一次请求替代多次单独查询。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageNames: 包名称切片，最多 128 个
//   - period: 统计周期，"last-day"、"last-week"、"last-month"
//
// 返回值:
//   - map[string]*models.DownloadStats: 包名到下载统计的映射
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	stats, err := registry.GetBulkDownloadStats(ctx, []string{"react", "vue", "angular"}, "last-week")
func (x *Registry) GetBulkDownloadStats(ctx context.Context, packageNames []string, period string) (map[string]*models.DownloadStats, error) {
	if len(packageNames) == 0 {
		return nil, fmt.Errorf("packageNames must not be empty")
	}
	if len(packageNames) > 128 {
		return nil, fmt.Errorf("packageNames must not exceed 128, got %d", len(packageNames))
	}
	escaped := url.QueryEscape(strings.Join(packageNames, ","))
	targetUrl := fmt.Sprintf("%s/point/%s/%s", downloadStatsBaseURL, period, escaped)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, err
	}
	return unmarshalJson[map[string]*models.DownloadStats](bytes)
}

// GetBulkDownloadRangeStats 批量获取多个包的每日下载统计（最多 128 个包）
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - packageNames: 包名称切片，最多 128 个
//   - period: 统计周期
//
// 返回值:
//   - map[string]*models.DownloadRangeStats: 包名到区间统计的映射
//   - error: 如果请求失败则返回错误
func (x *Registry) GetBulkDownloadRangeStats(ctx context.Context, packageNames []string, period string) (map[string]*models.DownloadRangeStats, error) {
	if len(packageNames) == 0 {
		return nil, fmt.Errorf("packageNames must not be empty")
	}
	if len(packageNames) > 128 {
		return nil, fmt.Errorf("packageNames must not exceed 128, got %d", len(packageNames))
	}
	escaped := url.QueryEscape(strings.Join(packageNames, ","))
	targetUrl := fmt.Sprintf("%s/range/%s/%s", downloadStatsBaseURL, period, escaped)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, err
	}
	return unmarshalJson[map[string]*models.DownloadRangeStats](bytes)
}
