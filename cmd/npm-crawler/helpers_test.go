package main

import (
	"testing"

	"github.com/scagogogo/npm-crawler/pkg/registry"
	"github.com/stretchr/testify/assert"
)

func TestMirrorToURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"official mirror", "official", registry.DefaultRegistryURL},
		{"taobao mirror", "taobao", registry.RegistryUrlTaoBao},
		{"npm-mirror mirror", "npm-mirror", registry.RegistryUrlNpmMirror},
		{"huawei mirror", "huawei", registry.RegistryUrlHuaWeiCloud},
		{"tencent mirror", "tencent", registry.RegistryUrlTencent},
		{"cnpm mirror", "cnpm", registry.RegistryUrlCnpm},
		{"yarn mirror", "yarn", registry.RegistryUrlYarn},
		{"npmjscom mirror", "npmjscom", registry.RegistryUrlNpmjsCom},
		{"npmmirror alias", "npmmirror", registry.RegistryUrlNpmMirror},
		{"huaweicloud alias", "huaweicloud", registry.RegistryUrlHuaWeiCloud},
		{"tencentcloud alias", "tencentcloud", registry.RegistryUrlTencent},
		{"case insensitive", "TAOBAO", registry.RegistryUrlTaoBao},
		{"custom http URL", "http://my-registry.local:8080", "http://my-registry.local:8080"},
		{"custom https URL", "https://npm.company.com", "https://npm.company.com"},
		{"unknown name falls back to official", "unknown-mirror", registry.DefaultRegistryURL},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mirrorToURL(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMirrorToURLMatchesSDKListMirrors(t *testing.T) {
	// 验证 mirrorToURL 对每个 SDK 镜像名称都能返回与 SDK 常量一致的 URL
	for _, m := range registry.ListMirrors() {
		t.Run(m.Name, func(t *testing.T) {
			result := mirrorToURL(m.Name)
			assert.Equal(t, m.URL, result, "mirrorToURL(%q) should match SDK ListMirrors URL", m.Name)
		})
	}
}

func TestCurrentMirrorLabel(t *testing.T) {
	// 测试 --registry 优先时显示 registry URL
	globalRegistry = "https://npm.company.com"
	globalMirror = "official"
	assert.Equal(t, "https://npm.company.com", currentMirrorLabel())

	// 测试无 --registry 时显示 mirror 名称
	globalRegistry = ""
	globalMirror = "npm-mirror"
	assert.Equal(t, "npm-mirror", currentMirrorLabel())

	// 重置全局变量
	globalRegistry = ""
	globalMirror = "official"
}
