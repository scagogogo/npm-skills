# 安装指南

五种安装方式按使用场景选择，不必全试——先按下面的决策图定位最适合自己的那一种：

```mermaid
flowchart TD
    Start(["想用 NPM Skills"]) --> Q1{"主要在哪个环境<br/>使用？"}
    Q1 -->|"Claude Code / AI 智能体"| M1["方式一：插件<br/>claude plugin install"]
    Q1 -->|"命令行 / CI 脚本"| Q2{"需要跨平台预编译<br/>二进制？"}
    Q1 -->|"Go 程序里调用"| Q3{"只要 SDK 库<br/>还是也要 CLI？"}
    Q2 -->|"是（最快上手）"| M2["方式二：预编译二进制<br/>34 平台任选"]
    Q2 -->|"否，能跑 go"| M4["方式四：go install"]
    Q3 -->|"只要库"| M5["方式五：go get"]
    Q3 -->|"库 + CLI 都要"| M3["方式三：源码构建<br/>bash scripts/install.sh"]

    M1 --> V["验证：npm-skills --version"]
    M2 --> V
    M3 --> V
    M4 --> V
    M5 --> V2["无需验证 CLI<br/>import 即可"]

    classDef pick fill:#e8f0fe,stroke:#4285f4,color:#174ea6;
    classDef verify fill:#e6f4ea,stroke:#34a853,color:#1e4620;
    class M1,M2,M3,M4,M5 pick;
    class V,V2 verify;
```

## 方式一：Claude Code 插件（AI 智能体推荐）

```bash
claude plugin marketplace add scagogogo/npm-skills
claude plugin install npm@npm-skills
```

## 方式二：预编译二进制（CLI 推荐）

![发布流水线](/release-pipeline.svg)

从 [GitHub Releases](https://github.com/scagogogo/npm-skills/releases/latest) 下载对应平台的压缩包：

```bash
# Linux x86_64
curl -sL https://github.com/scagogogo/npm-skills/releases/latest/download/npm-skills_0.2.0_linux_x86_64.tar.gz | tar -xz
sudo mv npm-skills npm-mcp-server /usr/local/bin/

# macOS Apple Silicon
curl -sL https://github.com/scagogogo/npm-skills/releases/latest/download/npm-skills_0.2.0_aarch64.tar.gz | tar -xz
sudo mv npm-skills npm-mcp-server /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/scagogogo/npm-skills/releases/latest/download/npm-skills_0.2.0_windows_x86_64.zip" -OutFile "npm-skills.zip"
Expand-Archive npm-skills.zip
```

### 支持的平台（34 个组合）

| OS | 架构 |
|----|------|
| Linux | amd64, arm64, 386, arm, loong64, mips, mips64, mips64le, mipsle, ppc64, ppc64le, riscv64, s390x |
| macOS | amd64, arm64 |
| Windows | amd64, 386 |
| FreeBSD | amd64, arm64, 386, arm, mips, mipsle, ppc64, ppc64le |
| OpenBSD | amd64, arm64, 386, arm, mips, mipsle, ppc64, ppc64le |
| NetBSD | amd64, arm64, 386, arm, mips, mipsle, ppc64, ppc64le |
| Illumos | amd64 |
| Solaris | amd64 |

## 方式三：从源码构建

```bash
git clone https://github.com/scagogogo/npm-skills.git
cd npm-skills
bash scripts/install.sh   # 编译到 ~/.local/bin/
```

## 方式四：go install

```bash
go install github.com/scagogogo/npm-skills/cmd/npm-skills@latest
go install github.com/scagogogo/npm-skills/cmd/mcp-server@latest
```

## 方式五：Go Module

```bash
go get github.com/scagogogo/npm-skills
```

## 验证安装

```bash
npm-skills --version
npm-skills mirrors
```

## 下一步

- 阅读 [快速开始](/getting-started) 了解四种接入方式
- 查阅 [CLI 命令手册](/cli) 掌握全部 26 个命令
- 浏览 [API 文档](/api/) 了解底层 SDK 方法
