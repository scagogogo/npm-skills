package models

import (
	"encoding/json"
)

// DownloadStats 表示 NPM 包的下载统计信息
//
// 此结构包含了从 NPM API 获取的包下载统计数据，包含下载次数、
// 统计周期、包名等信息。
//
// 主要字段说明:
//   - Downloads: 下载次数
//   - Start: 统计开始日期
//   - End: 统计结束日期
//   - Package: 包名称
type DownloadStats struct {
	Downloads int    `json:"downloads"` // 下载次数
	Start     string `json:"start"`     // 统计开始日期 (YYYY-MM-DD)
	End       string `json:"end"`       // 统计结束日期 (YYYY-MM-DD)
	Package   string `json:"package"`   // 包名称
}

// ToJsonString 将 DownloadStats 对象转换为 JSON 字符串
//
// 此方法将下载统计信息序列化为 JSON 格式的字符串，方便存储或传输。
// 如果序列化过程中发生错误，将返回包含错误信息的字符串。
//
// 返回值:
//   - string: JSON 格式的字符串表示
//
// 示例:
//
//	stats := &DownloadStats{...}
//	jsonStr := stats.ToJsonString()
//	fmt.Println(jsonStr)
func (ds *DownloadStats) ToJsonString() string {
	bytes, err := json.Marshal(ds)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

// DownloadRangeStats 表示时间范围内的下载统计信息
//
// 用于获取多个日期的下载统计数据。
//
// 主要字段说明:
//   - Start: 统计开始日期
//   - End: 统计结束日期
//   - Package: 包名称
//   - Downloads: 每日下载统计列表
type DownloadRangeStats struct {
	Start     string           `json:"start"`     // 统计开始日期
	End       string           `json:"end"`       // 统计结束日期
	Package   string           `json:"package"`   // 包名称
	Downloads []DailyDownloads `json:"downloads"` // 每日下载统计
}

// DailyDownloads 表示单日的下载统计
//
// 包含了特定日期的下载次数信息。
//
// 主要字段说明:
//   - Day: 日期 (YYYY-MM-DD)
//   - Downloads: 当日下载次数
type DailyDownloads struct {
	Day       string `json:"day"`       // 日期
	Downloads int    `json:"downloads"` // 下载次数
}

// ToJsonString 将 DownloadRangeStats 对象转换为 JSON 字符串
//
// 此方法将时间范围下载统计信息序列化为 JSON 格式的字符串。
//
// 返回值:
//   - string: JSON 格式的字符串表示
func (drs *DownloadRangeStats) ToJsonString() string {
	bytes, err := json.Marshal(drs)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}
