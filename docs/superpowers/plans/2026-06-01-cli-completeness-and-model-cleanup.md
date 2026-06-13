# CLI Completeness & Model Cleanup Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: `superpowers:subagent-driven-development`
> Steps use checkbox (`- [ ]`) syntax.

**Goal:** 补全 CLI 缺失能力（config 命令、镜像数据与 SDK 统一）并清理冗余模型（DistTags、Script），使 CLI 100% 封装 SDK 能力且数据源单一。

**Architecture:** SDK 层新增 `MirrorEntry` 结构体和 `ListMirrors()` 函数作为镜像元数据的唯一来源 → CLI 的 `mirrors` 命令和 `mirrorToURL()` 改为引用 SDK 数据 → 新增 `config` 命令暴露 `GetOptions()` → 清理 `DistTags` 冗余结构体 → `Script` 改为 `map[string]string` 以支持任意 npm scripts → 补 CLI 层单元测试。

**Tech Stack:** Go 1.20, Cobra (spf13/cobra), testify (stretchr/testify), net/http/httptest

**Risks:**
- Task 1 修改 `mirror.go`，已有 `mirror_test.go` 依赖现有构造函数签名 → 缓解：只新增函数和类型，不修改已有构造函数
- Task 3 修改 `Script` 从 struct 到 `map[string]string`，JSON 反序列化行为变化 → 缓解：`map[string]string` 天然兼容 JSON object，反序列化行为等价，且覆盖面更广
- Task 3 删除 `DistTags` 结构体，需确认无其他引用 → 缓解：调研确认 `Package.DistTags` 使用 `map[string]string`，`DistTags` 结构体未被引用

---

### Task 1: SDK Mirror Metadata Unification

**Depends on:** None
**Files:**
- Modify: `pkg/registry/mirror.go:1-167`（在文件末尾新增 MirrorEntry + ListMirrors）
- Modify: `pkg/registry/mirror_test.go`（新增 ListMirrors 测试）

- [ ] **Step 1: 新增 MirrorEntry 结构体和 ListMirrors 函数到 mirror.go — 统一镜像元数据的唯一来源**

在 `pkg/registry/mirror.go` 文件末尾（第 167 行之后）追加以下代码：

```go
// MirrorEntry 表示一个镜像源的元数据信息
//
// 包含镜像源的名称、URL、地理区域和描述信息。
// 此结构体是 CLI 和 SDK 共享镜像元数据的标准数据结构。
//
// 主要字段说明:
//   - Name: 镜像源标识名，用于 --mirror flag
//   - URL: 镜像源的 Registry URL
//   - Region: 地理区域（"Global" 或 "China"）
//   - Description: 镜像源的可读描述
type MirrorEntry struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Region      string `json:"region"`
	Description string `json:"description"`
}

// ListMirrors 返回所有支持的镜像源列表
//
// 返回值:
//   - []MirrorEntry: 包含所有镜像源元数据的切片
//
// 使用示例:
//
//	mirrors := registry.ListMirrors()
//	for _, m := range mirrors {
//	    fmt.Printf("%s: %s (%s)\n", m.Name, m.URL, m.Region)
//	}
func ListMirrors() []MirrorEntry {
	return []MirrorEntry{
		{"official", DefaultRegistryURL, "Global", "NPM Official Registry"},
		{"taobao", RegistryUrlTaoBao, "China", "Taobao NPM Mirror (Alibaba)"},
		{"npm-mirror", RegistryUrlNpmMirror, "China", "NPM Mirror (new Taobao domain)"},
		{"huawei", RegistryUrlHuaWeiCloud, "China", "Huawei Cloud Mirror"},
		{"tencent", RegistryUrlTencent, "China", "Tencent Cloud Mirror"},
		{"cnpm", RegistryUrlCnpm, "China", "CNPM Mirror"},
		{"yarn", RegistryUrlYarn, "Global", "Yarn Official Mirror"},
		{"npmjscom", RegistryUrlNpmjsCom, "Global", "NPM CouchDB Mirror"},
	}
}
```

- [ ] **Step 2: 新增 ListMirrors 单元测试**

在 `pkg/registry/mirror_test.go` 文件末尾（第 177 行之后）追加以下测试函数：

```go
func TestListMirrors(t *testing.T) {
	mirrors := ListMirrors()

	// 测试返回的镜像数量
	assert.Equal(t, 8, len(mirrors), "应该返回 8 个镜像源")

	// 测试第一个镜像源是 official
	assert.Equal(t, "official", mirrors[0].Name)
	assert.Equal(t, DefaultRegistryURL, mirrors[0].URL)
	assert.Equal(t, "Global", mirrors[0].Region)

	// 测试每个镜像源都有必要字段
	for _, m := range mirrors {
		assert.NotEmpty(t, m.Name, "镜像源名称不能为空")
		assert.NotEmpty(t, m.URL, "镜像源 URL 不能为空")
		assert.NotEmpty(t, m.Region, "镜像源区域不能为空")
		assert.NotEmpty(t, m.Description, "镜像源描述不能为空")
	}

	// 测试镜像源 URL 与常量一致
	urlMap := map[string]string{
		"official":   DefaultRegistryURL,
		"taobao":     RegistryUrlTaoBao,
		"npm-mirror": RegistryUrlNpmMirror,
		"huawei":     RegistryUrlHuaWeiCloud,
		"tencent":    RegistryUrlTencent,
		"cnpm":       RegistryUrlCnpm,
		"yarn":       RegistryUrlYarn,
		"npmjscom":   RegistryUrlNpmjsCom,
	}
	for _, m := range mirrors {
		expectedURL, ok := urlMap[m.Name]
		assert.True(t, ok, "镜像源 %s 不在 URL 映射中", m.Name)
		assert.Equal(t, expectedURL, m.URL, "镜像源 %s 的 URL 与常量不一致", m.Name)
	}
}
```

- [ ] **Step 3: 修改 CLI cmd_mirrors.go — 改为引用 SDK 的 ListMirrors() 数据源**

替换 `cmd/npm-skills/cmd_mirrors.go:10-37`（mirrorEntry 结构体定义和硬编码数据）：

```go
// 文件: cmd/npm-skills/cmd_mirrors.go
// 替换第 10-37 行（删除本地 mirrorEntry 结构体和硬编码 mirrors 变量）
// 改为使用 SDK 的 registry.MirrorEntry 和 registry.ListMirrors()

package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/scagogogo/npm-skills/pkg/registry"
	"github.com/spf13/cobra"
)

var mirrorsCmd = &cobra.Command{
	Use:   "mirrors",
	Short: "List available NPM mirror sources",
	Long: color.New(color.FgCyan).Sprintf("List available NPM mirror sources") + "\n\n" +
		"Shows all supported mirror sources with their URLs and descriptions.\n" +
		"Use the mirror name with --mirror flag in other commands.",
	Aliases: []string{"mirror", "ms"},
	Example: `  npm-skills mirrors`,
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		mirrors := registry.ListMirrors()

		// Print a nice table to stderr
		printHeader("Available NPM Mirror Sources")
		fmt.Fprintln(os.Stderr)
		for i, m := range mirrors {
			nameColor := color.New(color.FgGreen, color.Bold)
			urlColor := color.HiBlackString
			regionColor := color.New(color.FgYellow)
			descColor := color.HiWhiteString

			regionIcon := "🌍"
			if m.Region == "China" {
				regionIcon = "🇨🇳"
			}

			fmt.Fprintf(os.Stderr, "  %s %-12s  %s  %s %s  %s\n",
				color.HiBlackString("%d.", i+1),
				nameColor.Sprint(m.Name),
				urlColor(m.URL),
				regionIcon,
				regionColor.Sprint(m.Region),
				descColor(m.Description),
			)
		}
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, color.HiBlackString("  Usage: npm-skills <command> -m <name>"))
		fmt.Fprintln(os.Stderr, color.HiBlackString("  Or:    NPM_MIRROR=<name> npm-skills <command>"))
		fmt.Fprintln(os.Stderr)

		outputJSON(mirrors)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(mirrorsCmd)
}
```

- [ ] **Step 4: 修改 CLI helpers.go 的 mirrorToURL — 改为从 SDK ListMirrors 动态查找**

替换 `cmd/npm-skills/helpers.go:35-62`（mirrorToURL 函数）：

```go
// mirrorToURL converts a mirror name to its registry URL
// 使用 SDK 的 ListMirrors() 作为唯一数据源
func mirrorToURL(name string) string {
	// 先在 SDK 镜像列表中查找
	lowerName := strings.ToLower(name)
	for _, m := range registry.ListMirrors() {
		if strings.ToLower(m.Name) == lowerName {
			return m.URL
		}
	}

	// 特殊别名支持（SDK 不包含，CLI 独有）
	switch lowerName {
	case "npmmirror":
		return registry.RegistryUrlNpmMirror
	case "huaweicloud":
		return registry.RegistryUrlHuaWeiCloud
	case "tencentcloud":
		return registry.RegistryUrlTencent
	}

	// If it looks like a URL, use it directly
	if strings.HasPrefix(strings.ToLower(name), "http://") ||
		strings.HasPrefix(strings.ToLower(name), "https://") {
		return name
	}

	// Fallback to official
	return registry.DefaultRegistryURL
}
```

> 注意：CLI 独有别名（npmmirror、huaweicloud、tencentcloud）保留在此函数中，因为它们是 CLI 用户体验优化，不属于 SDK 层概念。

- [ ] **Step 5: 验证 Task 1**
Run: `cd /home/cc11001100/github/scagogogo/npm-skills && go build ./pkg/registry/ && go test ./pkg/registry/ -run TestListMirrors -v`
Expected:
  - Exit code: 0
  - Output contains: "PASS"

- [ ] **Step 6: 提交**
Run: `cd /home/cc11001100/github/scagogogo/npm-skills && git add pkg/registry/mirror.go pkg/registry/mirror_test.go cmd/npm-skills/cmd_mirrors.go cmd/npm-skills/helpers.go && git commit -m "feat(registry): add MirrorEntry and ListMirrors to unify mirror metadata source"`

---

### Task 2: CLI config Command

**Depends on:** Task 1
**Files:**
- Create: `cmd/npm-skills/cmd_config.go`

- [ ] **Step 1: 创建 config 命令 — 显示当前生效的 Registry 配置（暴露 GetOptions 能力到 CLI）**

```go
// cmd/npm-skills/cmd_config.go
package main

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current effective configuration",
	Long: color.New(color.FgCyan).Sprintf("Show current effective configuration") + "\n\n" +
		"Displays the currently active registry URL, mirror, proxy, and timeout settings.\n" +
		"Useful for debugging which registry/mirror/proxy is actually being used\n" +
		"after applying CLI flags and environment variables.",
	Aliases: []string{"cfg", "conf"},
	Example: `  npm-skills config
  npm-skills config -m npm-mirror
  npm-skills config --registry https://registry.npmmirror.com --proxy http://127.0.0.1:7890`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := resolveClient()
		opts := client.GetOptions()

		// Determine the effective mirror name (for display only)
		effectiveMirror := currentMirrorLabel()

		config := map[string]interface{}{
			"registry_url": opts.RegistryURL,
			"mirror":       effectiveMirror,
			"proxy":        opts.Proxy,
			"timeout":      globalTimeout,
		}
		if opts.Proxy == "" {
			config["proxy"] = "(none)"
		}

		printInfo("Current effective configuration:")
		outputJSON(config)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
```

- [ ] **Step 2: 验证 config 命令**
Run: `cd /home/cc11001100/github/scagogogo/npm-skills && go build ./cmd/npm-skills/ && ./npm-skills config`
Expected:
  - Exit code: 0
  - Output contains: "registry_url" and "mirror" and "official"

Run: `cd /home/cc11001100/github/scagogogo/npm-skills && ./npm-skills config -m npm-mirror`
Expected:
  - Exit code: 0
  - Output contains: "npmmirror.com"

- [ ] **Step 3: 提交**
Run: `cd /home/cc11001100/github/scagogogo/npm-skills && git add cmd/npm-skills/cmd_config.go && git commit -m "feat(cli): add config command to show effective registry configuration"`

---

### Task 3: Model Cleanup (DistTags + Script)

**Depends on:** None
**Files:**
- Modify: `pkg/models/dist_tags.go:1-7`（删除整个文件）
- Modify: `pkg/models/scripts.go:1-14`（替换 Script 结构体为 map 类型别名）
- Modify: `pkg/models/version.go:20`（更新 Scripts 字段类型）

- [ ] **Step 1: 删除冗余 DistTags 结构体 — Package.DistTags 实际使用 map[string]string，此结构体从未被引用**

删除文件 `pkg/models/dist_tags.go` 的全部内容。由于整个文件只有一个未被使用的结构体，直接删除文件：

Run: `cd /home/cc11001100/github/scagogogo/npm-skills && rm pkg/models/dist_tags.go`

- [ ] **Step 2: 修改 Script 模型 — 改为 map[string]string 以支持任意 npm scripts**

替换 `pkg/models/scripts.go:1-14` 的全部内容：

```go
// Script 表示 NPM 包的脚本命令定义
//
// 使用 map[string]string 类型以支持 npm 包中定义的任意脚本命令，
// 例如 "build"、"lint"、"dev" 等非标准脚本。
// 最常用的脚本命令包括:
//   - "test": 测试脚本命令
//   - "start": 启动项目脚本命令
//   - "build": 构建脚本命令
//   - "lint": 代码检查脚本命令
type Script map[string]string
```

- [ ] **Step 3: 更新 Version 结构体的 Scripts 字段类型 — 从 *Script 改为 Script**

替换 `pkg/models/version.go:20`（Scripts 字段定义）：

将:
```go
		Scripts     *Script     `json:"scripts"`     // 脚本命令定义
```

替换为:
```go
		Scripts     Script      `json:"scripts"`      // 脚本命令定义（支持任意 npm script key）
```

- [ ] **Step 4: 验证模型改动 — 编译和测试**
Run: `cd /home/cc11001100/github/scagogogo/npm-skills && go build ./... && go test ./pkg/models/ -v`
Expected:
  - Exit code: 0
  - Output contains: "PASS"

- [ ] **Step 5: 提交**
Run: `cd /home/cc11001100/github/scagogogo/npm-skills && git add -A pkg/models/ && git commit -m "refactor(models): remove unused DistTags struct, change Script to map[string]string for flexibility"`

---

### Task 4: CLI Unit Tests

**Depends on:** Task 1, Task 2
**Files:**
- Create: `cmd/npm-skills/helpers_test.go`

- [ ] **Step 1: 创建 CLI helpers 单元测试 — 覆盖 mirrorToURL、resolveClient、currentMirrorLabel**

```go
// cmd/npm-skills/helpers_test.go
package main

import (
	"testing"

	"github.com/scagogogo/npm-skills/pkg/registry"
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
```

- [ ] **Step 2: 验证 CLI 测试**
Run: `cd /home/cc11001100/github/scagogogo/npm-skills && go test ./cmd/npm-skills/ -v`
Expected:
  - Exit code: 0
  - Output contains: "PASS"

- [ ] **Step 3: 全量回归测试**
Run: `cd /home/cc11001100/github/scagogogo/npm-skills && go test ./... -count=1`
Expected:
  - Exit code: 0
  - Output does NOT contain: "FAIL"

- [ ] **Step 4: 提交**
Run: `cd /home/cc11001100/github/scagogogo/npm-skills && git add cmd/npm-skills/helpers_test.go && git commit -m "test(cli): add unit tests for mirrorToURL, resolveClient and currentMirrorLabel"`
