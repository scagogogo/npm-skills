package models

import "encoding/json"

// ChangesOptions CouchDB Changes Feed 查询参数
type ChangesOptions struct {
	Since      string `json:"since,omitempty"`       // 起始更新序列号
	Limit      int    `json:"limit,omitempty"`        // 限制返回结果数量
	IncludeDocs bool   `json:"include_docs,omitempty"` // 是否包含完整文档
}

// ChangesResult CouchDB Changes Feed 结果
type ChangesResult struct {
	LastSeq string        `json:"last_seq"`
	Pending int           `json:"pending"`
	Results []ChangeEntry `json:"results"`
}

// ToJsonString 将 ChangesResult 对象转换为 JSON 字符串
func (x *ChangesResult) ToJsonString() string {
	bytes, _ := json.Marshal(x)
	return string(bytes)
}

// ChangeEntry CouchDB Changes Feed 中的单个变更条目
type ChangeEntry struct {
	Seq     string          `json:"seq"`
	ID      string          `json:"id"`
	Changes []ChangeVersion `json:"changes"`
	Deleted bool            `json:"deleted,omitempty"`
	Doc     json.RawMessage `json:"doc,omitempty"`
}

// ChangeVersion 变更版本信息
type ChangeVersion struct {
	Rev string `json:"rev"`
}

// AllDocsOptions CouchDB All Docs 查询参数
type AllDocsOptions struct {
	StartKey    string `json:"startkey,omitempty"`
	EndKey      string `json:"endkey,omitempty"`
	Limit       int    `json:"limit,omitempty"`
	Skip        int    `json:"skip,omitempty"`
	IncludeDocs bool   `json:"include_docs,omitempty"`
	Descending  bool   `json:"descending,omitempty"`
}

// AllDocsResult CouchDB All Docs 查询结果
type AllDocsResult struct {
	TotalRows int      `json:"total_rows"`
	Offset    int      `json:"offset"`
	Rows      []DocRow `json:"rows"`
}

// ToJsonString 将 AllDocsResult 对象转换为 JSON 字符串
func (x *AllDocsResult) ToJsonString() string {
	bytes, _ := json.Marshal(x)
	return string(bytes)
}

// DocRow CouchDB 文档行
type DocRow struct {
	ID    string          `json:"id"`
	Key   string          `json:"key"`
	Value DocRowValue     `json:"value"`
	Doc   json.RawMessage `json:"doc,omitempty"`
}

// DocRowValue 文档行值
type DocRowValue struct {
	Rev string `json:"rev"`
}

// ViewOptions CouchDB 视图查询参数
type ViewOptions struct {
	Key        string `json:"key,omitempty"`
	StartKey   string `json:"startkey,omitempty"`
	EndKey     string `json:"endkey,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	Skip       int    `json:"skip,omitempty"`
	Group      bool   `json:"group,omitempty"`
	GroupLevel int    `json:"group_level,omitempty"`
	Descending bool   `json:"descending,omitempty"`
}

// ViewResult CouchDB 视图查询结果
type ViewResult struct {
	TotalRows int       `json:"total_rows"`
	Offset    int       `json:"offset"`
	Rows      []ViewRow `json:"rows"`
}

// ToJsonString 将 ViewResult 对象转换为 JSON 字符串
func (x *ViewResult) ToJsonString() string {
	bytes, _ := json.Marshal(x)
	return string(bytes)
}

// ViewRow CouchDB 视图行
type ViewRow struct {
	ID    string          `json:"id"`
	Key   json.RawMessage `json:"key"`
	Value json.RawMessage `json:"value"`
}
