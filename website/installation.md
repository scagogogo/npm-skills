# 安装指南

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
