package models

import (
	"encoding/json"
)

// SearchResult 表示 NPM 包搜索结果
//
// 此结构包含了从 NPM Registry 搜索 API 返回的结果，包含匹配的包列表、
// 搜索得分、总数等信息。
//
// 主要字段说明:
//   - Objects: 匹配的包对象列表
//   - Total: 总匹配数量
//   - Time: 搜索耗时（毫秒）
type SearchResult struct {
	Objects []SearchObject `json:"objects"` // 搜索结果对象列表
	Total   int            `json:"total"`   // 总匹配数量
	Time    string         `json:"time"`    // 搜索耗时
}

// SearchObject 表示搜索结果中的单个包对象
//
// 包含了搜索匹配的包的基本信息、评分等。
//
// 主要字段说明:
//   - Package: 包的基本信息
//   - Score: 搜索评分信息
//   - SearchScore: 总体搜索得分
type SearchObject struct {
	Package     SearchPackage `json:"package"`     // 包基本信息
	Score       Score         `json:"score"`       // 评分详情
	SearchScore float64       `json:"searchScore"` // 搜索得分
}

// SearchPackage 表示搜索结果中的包信息
//
// 包含了包的基本元数据信息，是完整包信息的简化版本。
//
// 主要字段说明:
//   - Name: 包名称
//   - Scope: 包作用域（如果有）
//   - Version: 最新版本号
//   - Description: 包描述
//   - Keywords: 关键词列表
//   - Date: 发布日期
//   - Author: 作者信息
//   - Publisher: 发布者信息
//   - Maintainers: 维护者列表
type SearchPackage struct {
	Name        string   `json:"name"`        // 包名称
	Scope       string   `json:"scope"`       // 包作用域
	Version     string   `json:"version"`     // 版本号
	Description string   `json:"description"` // 包描述
	Keywords    []string `json:"keywords"`    // 关键词
	Date        string   `json:"date"`        // 发布日期
	Links       Links    `json:"links"`       // 相关链接
	Author      *User    `json:"author"`      // 作者
	Publisher   *User    `json:"publisher"`   // 发布者
	Maintainers []*User  `json:"maintainers"` // 维护者列表
	ExactName   string   `json:"exactName"`   // 精确匹配的包名（用于判断是否为精确搜索匹配）
}

// Links 表示包的相关链接信息
//
// 包含了包的各种相关链接，如 NPM 页面、仓库、主页等。
type Links struct {
	NPM        string `json:"npm"`        // NPM 页面链接
	Homepage   string `json:"homepage"`   // 主页链接
	Repository string `json:"repository"` // 仓库链接
	Bugs       string `json:"bugs"`       // 问题跟踪链接
}

// Score 表示搜索结果的评分信息
//
// 包含了搜索算法计算的各种评分维度。
//
// 主要字段说明:
//   - Final: 最终得分
//   - Detail: 详细评分信息
type Score struct {
	Final  float64     `json:"final"`  // 最终得分
	Detail ScoreDetail `json:"detail"` // 详细评分
}

// ScoreDetail 表示详细的评分信息
//
// 包含了搜索算法各个维度的评分详情。
//
// 主要字段说明:
//   - Quality: 质量得分
//   - Popularity: 流行度得分
//   - Maintenance: 维护状态得分
type ScoreDetail struct {
	Quality     float64 `json:"quality"`     // 质量得分
	Popularity  float64 `json:"popularity"`  // 流行度得分
	Maintenance float64 `json:"maintenance"` // 维护状态得分
}

// ToJsonString 将 SearchResult 对象转换为 JSON 字符串
//
// 此方法将搜索结果序列化为 JSON 格式的字符串，方便存储或传输。
// 如果序列化过程中发生错误，将返回包含错误信息的字符串。
//
// 返回值:
//   - string: JSON 格式的字符串表示
//
// 示例:
//
//	result := &SearchResult{...}
//	jsonStr := result.ToJsonString()
//	fmt.Println(jsonStr)
func (sr *SearchResult) ToJsonString() string {
	bytes, err := json.Marshal(sr)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}
