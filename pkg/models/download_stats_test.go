package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloadStats(t *testing.T) {
	// 创建一个测试下载统计
	stats := &DownloadStats{
		Downloads: 1234567,
		Start:     "2023-01-01",
		End:       "2023-01-07",
		Package:   "react",
	}

	// 基本属性测试
	assert.Equal(t, 1234567, stats.Downloads)
	assert.Equal(t, "2023-01-01", stats.Start)
	assert.Equal(t, "2023-01-07", stats.End)
	assert.Equal(t, "react", stats.Package)
}

func TestDownloadStatsToJsonString(t *testing.T) {
	stats := &DownloadStats{
		Downloads: 98765,
		Start:     "2023-06-01",
		End:       "2023-06-30",
		Package:   "lodash",
	}

	jsonStr := stats.ToJsonString()
	assert.NotEmpty(t, jsonStr)
	assert.Contains(t, jsonStr, "lodash")
	assert.Contains(t, jsonStr, "98765")
	assert.Contains(t, jsonStr, "2023-06-01")

	// 测试反序列化
	var parsed DownloadStats
	err := json.Unmarshal([]byte(jsonStr), &parsed)
	assert.NoError(t, err)
	assert.Equal(t, stats.Downloads, parsed.Downloads)
	assert.Equal(t, stats.Package, parsed.Package)
	assert.Equal(t, stats.Start, parsed.Start)
	assert.Equal(t, stats.End, parsed.End)
}

func TestDownloadRangeStats(t *testing.T) {
	// 创建测试时间范围下载统计
	rangeStats := &DownloadRangeStats{
		Start:   "2023-01-01",
		End:     "2023-01-03",
		Package: "express",
		Downloads: []DailyDownloads{
			{
				Day:       "2023-01-01",
				Downloads: 10000,
			},
			{
				Day:       "2023-01-02",
				Downloads: 12000,
			},
			{
				Day:       "2023-01-03",
				Downloads: 11500,
			},
		},
	}

	// 基本属性测试
	assert.Equal(t, "2023-01-01", rangeStats.Start)
	assert.Equal(t, "2023-01-03", rangeStats.End)
	assert.Equal(t, "express", rangeStats.Package)
	assert.Equal(t, 3, len(rangeStats.Downloads))

	// 测试每日下载数据
	firstDay := rangeStats.Downloads[0]
	assert.Equal(t, "2023-01-01", firstDay.Day)
	assert.Equal(t, 10000, firstDay.Downloads)

	secondDay := rangeStats.Downloads[1]
	assert.Equal(t, "2023-01-02", secondDay.Day)
	assert.Equal(t, 12000, secondDay.Downloads)

	thirdDay := rangeStats.Downloads[2]
	assert.Equal(t, "2023-01-03", thirdDay.Day)
	assert.Equal(t, 11500, thirdDay.Downloads)
}

func TestDownloadRangeStatsToJsonString(t *testing.T) {
	rangeStats := &DownloadRangeStats{
		Start:   "2023-05-01",
		End:     "2023-05-02",
		Package: "vue",
		Downloads: []DailyDownloads{
			{
				Day:       "2023-05-01",
				Downloads: 8000,
			},
			{
				Day:       "2023-05-02",
				Downloads: 9000,
			},
		},
	}

	jsonStr := rangeStats.ToJsonString()
	assert.NotEmpty(t, jsonStr)
	assert.Contains(t, jsonStr, "vue")
	assert.Contains(t, jsonStr, "2023-05-01")
	assert.Contains(t, jsonStr, "8000")

	// 测试反序列化
	var parsed DownloadRangeStats
	err := json.Unmarshal([]byte(jsonStr), &parsed)
	assert.NoError(t, err)
	assert.Equal(t, rangeStats.Start, parsed.Start)
	assert.Equal(t, rangeStats.End, parsed.End)
	assert.Equal(t, rangeStats.Package, parsed.Package)
	assert.Equal(t, len(rangeStats.Downloads), len(parsed.Downloads))
	assert.Equal(t, rangeStats.Downloads[0].Downloads, parsed.Downloads[0].Downloads)
}

func TestDailyDownloads(t *testing.T) {
	daily := DailyDownloads{
		Day:       "2023-03-15",
		Downloads: 15000,
	}

	assert.Equal(t, "2023-03-15", daily.Day)
	assert.Equal(t, 15000, daily.Downloads)
}

func TestDownloadStatsEdgeCases(t *testing.T) {
	// 测试零下载量
	zeroStats := &DownloadStats{
		Downloads: 0,
		Start:     "2023-01-01",
		End:       "2023-01-01",
		Package:   "new-package",
	}

	jsonStr := zeroStats.ToJsonString()
	assert.Contains(t, jsonStr, "new-package")
	assert.Contains(t, jsonStr, "0")

	// 测试空包名
	emptyPackage := &DownloadStats{
		Downloads: 100,
		Start:     "2023-01-01",
		End:       "2023-01-01",
		Package:   "",
	}

	jsonStr = emptyPackage.ToJsonString()
	assert.NotEmpty(t, jsonStr)

	// 测试空的下载范围统计
	emptyRange := &DownloadRangeStats{
		Start:     "2023-01-01",
		End:       "2023-01-01",
		Package:   "test",
		Downloads: []DailyDownloads{},
	}

	jsonStr = emptyRange.ToJsonString()
	assert.Contains(t, jsonStr, "test")
}

func TestDownloadStatsToJsonStringErrorBranch(t *testing.T) {
	// 测试 ToJsonString 的 error 分支
	// 当 TestMarshalErr 字段非 nil 时，json.Marshal 会返回错误
	stats := &DownloadStats{
		Downloads:    123,
		Start:        "2023-01-01",
		End:          "2023-01-07",
		Package:      "test-pkg",
		TestMarshalErr: &testMarshalFailType{},
	}

	jsonStr := stats.ToJsonString()
	assert.Contains(t, jsonStr, "forced marshal error for testing")
}

func TestDownloadRangeStatsToJsonStringErrorBranch(t *testing.T) {
	// 测试 DownloadRangeStats.ToJsonString 的 error 分支
	rangeStats := &DownloadRangeStats{
		Start:   "2023-01-01",
		End:     "2023-01-03",
		Package: "express",
		Downloads: []DailyDownloads{
			{Day: "2023-01-01", Downloads: 10000},
		},
		TestMarshalErr: &testMarshalFailType{},
	}

	jsonStr := rangeStats.ToJsonString()
	assert.Contains(t, jsonStr, "forced marshal error for testing")
}
