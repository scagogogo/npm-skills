# NPM Registry MCP Server

> NPM 注册表 MCP 服务器 — 为 AI 代理提供 NPM 注册表操作工具

## 什么是 MCP Server？

MCP (Model Context Protocol) 是一种标准协议，允许 AI 代理（如 Claude、GPT 等）通过 JSON-RPC 发现和调用工具。npm-mcp-server 将 NPM 注册表操作暴露为 MCP 工具，使 AI 代理能够直接查询包信息、搜索包、获取下载统计等，无需通过 CLI 命令行。

## 安装

```bash
# 编译 MCP 服务器
go build -o ~/.local/bin/npm-mcp-server ./cmd/mcp-server/

# 或使用安装脚本（同时构建 CLI 和 MCP 服务器）
bash scripts/install.sh
```

## 配置

### 命令行参数

| 参数 | 环境变量 | 默认值 | 说明 |
|------|---------|--------|------|
| `--registry` | `NPM_REGISTRY` | | 自定义注册表 URL（覆盖 --mirror） |
| `--mirror` | `NPM_MIRROR` | `official` | 镜像源名称 |
| `--proxy` | `NPM_PROXY` | | HTTP 代理 URL |
| `--token` | `NPM_TOKEN` | | NPM 认证令牌 |
| `--timeout` | `NPM_TIMEOUT` | `120` | 请求超时时间（秒） |

**优先级：** 命令行参数 > 环境变量 > 默认值

### Claude Code 集成

添加到你的 Claude Code 设置（`~/.claude/settings.json` 或项目 `.claude/settings.json`）：

```json
{
  "mcpServers": {
    "npm-registry": {
      "command": "npm-mcp-server",
      "args": ["--mirror", "npm-mirror"]
    }
  }
}
```

使用代理：

```json
{
  "mcpServers": {
    "npm-registry": {
      "command": "npm-mcp-server",
      "args": ["--mirror", "npm-mirror", "--proxy", "http://127.0.0.1:7890"]
    }
  }
}
```

使用自定义注册表：

```json
{
  "mcpServers": {
    "npm-registry": {
      "command": "npm-mcp-server",
      "args": ["--registry", "https://npm.my-company.com", "--token", "npm_xxxxx"]
    }
  }
}
```

## MCP 工具列表（12 个）

### 注册表信息工具

| 工具 | 说明 | 参数 |
|------|------|------|
| `npm_registry_info` | 获取 NPM 注册表状态和统计信息 | 无 |
| `npm_mirrors` | 列出所有可用的镜像源 | 无 |

### 包信息工具

| 工具 | 说明 | 参数 |
|------|------|------|
| `npm_package` | 获取完整的包元数据（响应可能很大，10MB+） | `name` (必填) |
| `npm_package_summary` | 获取精简的包元数据（推荐用于大多数查询） | `name` (必填) |

### 搜索工具

| 工具 | 说明 | 参数 |
|------|------|------|
| `npm_search` | 搜索 NPM 包，支持分页和权重调整 | `query` (必填), `limit`, `from`, `quality`, `popularity`, `maintenance` (可选) |

### 版本工具

| 工具 | 说明 | 参数 |
|------|------|------|
| `npm_version` | 获取特定版本的元数据 | `name` (必填), `version` (必填) |
| `npm_versions` | 列出所有已发布的版本号 | `name` (必填) |
| `npm_latest_version` | 获取最新版本号 | `name` (必填) |
| `npm_dist_tags` | 获取分发标签（latest, next, beta 等） | `name` (必填) |

### 下载统计工具

| 工具 | 说明 | 参数 |
|------|------|------|
| `npm_download_stats` | 获取下载统计 | `name` (必填), `period` (可选，默认 last-week) |
| `npm_download_range` | 获取每日下载趋势数据 | `name` (必填), `period` (可选，默认 last-week) |

### 认证工具

| 工具 | 说明 | 参数 |
|------|------|------|
| `npm_whoami` | 检查认证状态（需要 --token） | 无 |

## 响应格式

所有 MCP 工具返回 JSON 格式的文本内容。对于超过 50KB 的响应，会自动截断：

- 包 README 截断为前 2000 字符
- 包 Versions 映射替换为 version_keys 数组（仅版本号字符串）
- 硬限制：超过 100KB 的响应会被截断，附带尾部通知

## 使用示例

### 查看包信息

```json
// 工具调用
{"name": "npm_package_summary", "arguments": {"name": "react"}}
```

### 搜索包

```json
// 基本搜索
{"name": "npm_search", "arguments": {"query": "http client", "limit": 10}}

// 带权重搜索（侧重流行度）
{"name": "npm_search", "arguments": {"query": "http client", "popularity": 1.0, "quality": 0.0}}
```

### 查看最新版本

```json
{"name": "npm_latest_version", "arguments": {"name": "react"}}
```

### 下载统计

```json
{"name": "npm_download_stats", "arguments": {"name": "react", "period": "last-month"}}
```

## 架构说明

MCP 服务器通过 stdio 传输运行，使用 JSON-RPC 协议通信。它包装了现有的 Go SDK (`pkg/registry/`)，将每个 SDK 方法映射为一个 MCP 工具。

```
AI Agent → JSON-RPC → npm-mcp-server (stdio) → pkg/registry/ → NPM Registry API
```

关键设计决策：

1. **排除 DownloadTarball** — 写磁盘操作不适合 MCP 工具
2. **npm_package vs npm_package_summary** — 完整包信息可达 10MB+，建议大多数查询使用精简版
3. **下载统计始终查询 api.npmjs.org** — NPM 的下载统计 API 是独立的，不受 mirror/registry 设置影响
4. **认证令牌通过服务器级配置** — 使用 --token 或 NPM_TOKEN 环境变量，而非每个工具参数