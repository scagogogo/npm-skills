# 快速开始

## 一分钟安装为 Claude Code 插件

```bash
# 1. 添加 marketplace
claude plugin marketplace add scagogogo/npm-skills

# 2. 安装插件
claude plugin install npm@npm-skills
```

安装后直接用自然语言提问，AI 智能体会自动调用 NPM Skills：

> *"查找 axios 包的信息"*
> *"下载 react 的 tarball"*
> *"用国内镜像查看 vue 上个月下载量"*

从自然语言到 NPM Registry 的完整调用链路如下：

```mermaid
sequenceDiagram
    autonumber
    actor U as 你
    participant AI as AI 智能体<br/>Claude Code
    participant S as SKILL.md<br/>渐进式披露
    participant CLI as npm-skills CLI
    participant R as Registry SDK
    participant N as NPM Registry / 镜像

    U->>AI: "用国内镜像查看 vue 上个月下载量"
    AI->>S: 匹配技能描述，加载指令
    S-->>AI: 返回命令用法（按需展开）
    AI->>CLI: npm-skills download-stats vue -p last-month -m npm-mirror
    CLI->>R: 构造请求（镜像/代理/超时）
    R->>N: HTTP GET /downloads/point/last-month/vue
    N-->>R: JSON 响应
    R-->>CLI: 结构化数据
    CLI-->>AI: JSON（stdout）
    AI-->>U: 自然语言总结结果
```

## 四种接入方式

四种入口共享同一套 Registry SDK 核心，最终统一走向 NPM Registry 或镜像源：

```mermaid
flowchart TB
    subgraph 使用者
        A1["AI 智能体<br/>Claude / Cursor / Windsurf"]
        A2["命令行 / CI 脚本"]
        A3["Go 应用程序"]
    end

    subgraph 接入层["NPM Skills 接入层"]
        B1["Skill / Plugin<br/>SKILL.md"]
        B2["CLI<br/>npm-skills"]
        B3["MCP 服务器<br/>npm-mcp-server"]
        B4["Go SDK<br/>pkg/registry"]
    end

    C["Registry 客户端核心<br/>Options · HTTP 连接池 · 类型化错误"]
    D["NPM Registry<br/>官方 / 8 镜像源 / 私有仓库"]

    A1 --> B1
    A1 --> B3
    A2 --> B2
    A3 --> B4
    B1 --> B2
    B2 --> C
    B3 --> C
    B4 --> C
    C -->|HTTP/HTTPS · 代理| D

    classDef core fill:#cb3837,stroke:#8b1a1a,color:#fff;
    class C core;
```

| 方式 | 适用场景 | 入口 |
|------|---------|------|
| **Skill / Plugin** | AI 智能体自动调用 | `claude plugin install npm@npm-skills` |
| **CLI 工具** | 命令行 / 脚本 | `npm-skills <command>` |
| **Go SDK** | Go 程序集成 | `import "github.com/scagogogo/npm-skills/pkg/registry"` |
| **MCP 服务器** | MCP 兼容客户端 | `npm-mcp-server` |

## CLI 速查（90% 场景）

26 个命令按功能域划分，90% 场景只需记住这几类：

```mermaid
mindmap
  root((npm-skills CLI))
    查包
      package-summary
      package
      pkg-version
      versions --latest
    搜索
      search -l N
    版本标签
      dist-tags get
      versions
    统计
      download-stats -p last-week
      download-range
    镜像
      mirrors
      -m npm-mirror
    其他
      registry-info
      whoami --token
```

```bash
npm-skills package-summary react            # 轻量包信息（推荐）
npm-skills search "http client" -l 10       # 搜索包
npm-skills versions react --latest          # 最新版本
npm-skills dist-tags get react              # dist-tags
npm-skills download-stats react -p last-month  # 下载统计
npm-skills mirrors                          # 镜像源列表
npm-skills package react -m npm-mirror      # 国内镜像
```

## 镜像与代理

```bash
# 国内镜像（无需代理）
npm-skills package react -m npm-mirror

# HTTP 代理
npm-skills package react --proxy http://127.0.0.1:7890

# 私有仓库
npm-skills package my-lib --registry https://npm.my-company.com -t npm_xxxxx

# 环境变量（推荐写入 shell 配置）
export NPM_MIRROR=npm-mirror
export NPM_PROXY=http://127.0.0.1:7890
export NPM_TOKEN=npm_xxxxx
```

当同一项配置被多处指定时，遵循「命令行参数 > 环境变量 > 默认值」的优先级：

```mermaid
flowchart LR
    Start(["解析一项配置<br/>如 mirror / proxy / token"]) --> Q1{"命令行传了<br/>--mirror ?"}
    Q1 -->|是| U1["采用 CLI 参数"]
    Q1 -->|否| Q2{"环境变量<br/>NPM_MIRROR ?"}
    Q2 -->|是| U2["采用环境变量"]
    Q2 -->|否| U3["采用内置默认值<br/>official"]

    U1 --> Done(["生效配置"])
    U2 --> Done
    U3 --> Done

    classDef win fill:#e6f4ea,stroke:#34a853,color:#1e4620;
    class U1,U2,U3 win;
```

## 下一步

- [CLI 命令手册](/cli) — 完整 26 个命令
- [Go SDK](/api/registry) — 程序化访问
- [MCP 服务器](/mcp-server) — 接入 AI 工具链
