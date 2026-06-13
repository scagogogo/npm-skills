package models

import "encoding/json"

// PublishMetadata 发布包时需要的元数据
//
// 用于构造发布请求中的包文档信息，包含 package.json 中的标准字段。
// 通过 PublishPackageFromTarball 方法使用。
type PublishMetadata struct {
	Name                 string            `json:"name"`
	Description          string            `json:"description,omitempty"`
	Version              string            `json:"version"`
	Main                 string            `json:"main,omitempty"`
	Module               string            `json:"module,omitempty"`
	Types                string            `json:"types,omitempty"`
	Typings              string            `json:"typings,omitempty"`
	Scripts              map[string]string `json:"scripts,omitempty"`
	Dependencies         map[string]string `json:"dependencies,omitempty"`
	DevDependencies      map[string]string `json:"devDependencies,omitempty"`
	PeerDependencies     map[string]string `json:"peerDependencies,omitempty"`
	OptionalDependencies map[string]string `json:"optionalDependencies,omitempty"`
	BundledDependencies  []string          `json:"bundledDependencies,omitempty"`
	Keywords             []string          `json:"keywords,omitempty"`
	License              string            `json:"license,omitempty"`
	Author               *Author           `json:"author,omitempty"`
	Contributors         []Contributor     `json:"contributors,omitempty"`
	Repository           *Repository       `json:"repository,omitempty"`
	Bugs                 *Bugs             `json:"bugs,omitempty"`
	Homepage             string            `json:"homepage,omitempty"`
	Private              bool              `json:"private,omitempty"`
}

// ToJsonString 将 PublishMetadata 对象转换为 JSON 字符串
func (x *PublishMetadata) ToJsonString() string {
	bytes, _ := json.Marshal(x)
	return string(bytes)
}
