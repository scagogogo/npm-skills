package models

import "encoding/json"

// Advisory NPM安全公告
//
// 表示一个 NPM 安全公告，包含漏洞详情、影响范围和修复建议。
type Advisory struct {
	ID             int             `json:"id"`
	Created        string          `json:"created"`
	Updated        string          `json:"updated"`
	Title          string          `json:"title"`
	Severity       string          `json:"severity"` // "low", "moderate", "high", "critical"
	CVE            string          `json:"cve,omitempty"`
	CWE            string          `json:"cwe,omitempty"`
	ModuleName     string          `json:"module_name"`
	Vulnerable     string          `json:"vulnerable_versions"`
	Patched        string          `json:"patched_versions"`
	URL            string          `json:"url"`
	Overview       string          `json:"overview,omitempty"`
	Recommendation string          `json:"recommendation,omitempty"`
	References     json.RawMessage `json:"references,omitempty"` // 可以是字符串或字符串数组
	Access         string          `json:"access,omitempty"`
}

// ToJsonString 将 Advisory 对象转换为 JSON 字符串
func (x *Advisory) ToJsonString() string {
	bytes, _ := json.Marshal(x)
	return string(bytes)
}

// GetReferences returns the references as a string slice.
// Handles the fact that References can be a string, []string, or nil.
func (x *Advisory) GetReferences() []string {
	if x.References == nil {
		return nil
	}

	// Try as array first
	var arr []string
	if err := json.Unmarshal(x.References, &arr); err == nil {
		return arr
	}

	// Try as single string
	var s string
	if err := json.Unmarshal(x.References, &s); err == nil {
		return []string{s}
	}

	return nil
}

// AdvisoryListOptions 安全公告列表查询参数
type AdvisoryListOptions struct {
	Page            int    `json:"page,omitempty"`
	PerPage         int    `json:"per_page,omitempty"`
	AffectedPackage string `json:"affected_package,omitempty"`
}

// QuickAuditRequest 快速审计请求
//
// 提交依赖列表进行安全审计，key 为包名，value 为版本范围。
type QuickAuditRequest struct {
	Dependencies map[string]string `json:"dependencies"`
}

// ToJsonString 将 QuickAuditRequest 对象转换为 JSON 字符串
func (x *QuickAuditRequest) ToJsonString() string {
	bytes, _ := json.Marshal(x)
	return string(bytes)
}

// QuickAuditResult 快速审计结果
//
// 包含漏洞统计信息，按严重程度分类。
type QuickAuditResult struct {
	Metadata struct {
		Vulnerabilities struct {
			Low      int `json:"low"`
			Moderate int `json:"moderate"`
			High     int `json:"high"`
			Critical int `json:"critical"`
		} `json:"vulnerabilities"`
		Dependencies         int `json:"dependencies"`
		DevDependencies      int `json:"devDependencies"`
		OptionalDependencies int `json:"optionalDependencies"`
		TotalDependencies    int `json:"totalDependencies"`
	} `json:"metadata"`
	Error interface{} `json:"error,omitempty"`
}

// ToJsonString 将 QuickAuditResult 对象转换为 JSON 字符串
func (x *QuickAuditResult) ToJsonString() string {
	bytes, _ := json.Marshal(x)
	return string(bytes)
}
