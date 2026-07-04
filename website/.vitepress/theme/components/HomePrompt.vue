<script setup lang="ts">
import { computed, ref } from 'vue'
import { useData } from 'vitepress'

const { lang } = useData()
const isZh = computed(() => lang.value === 'zh-CN')

// AI Agent 提示词：用户一键复制后粘贴给 AI Agent，Agent 顺着指引安装并使用 NPM Skills。
// 中英文版本内容一一对应，与 website/index.md、website/en/index.md 保持同步。
const promptZh = `---
title: 使用 NPM Skills 操作 NPM Registry
---

# 你的任务

你现在可以使用 **NPM Skills** —— 一个完整的 NPM Registry 客户端，支持包查询、版本管理、dist-tags、下载统计、发布、访问控制、安全审计、组织/团队、webhooks 等，内置 8 个镜像源（中国大陆免代理）与 HTTP/SOCKS5 代理。

## 第一步：安装（按你的环境任选一种，推荐方式 1）

### 方式 1：作为 Claude Code 插件（AI 原生，推荐）
\`\`\`bash
claude plugin marketplace add scagogogo/npm-skills
claude plugin install npm@npm-skills
\`\`\`
安装后你（AI）会自动发现 SKILL.md 并按需调用，无需手动 shell。

### 方式 2：下载预编译二进制（无 Claude Code 时）
\`\`\`bash
# Linux x86_64
curl -sL https://github.com/scagogogo/npm-skills/releases/latest/download/npm-skills_0.2.0_linux_x86_64.tar.gz | tar -xz
sudo mv npm-skills npm-mcp-server /usr/local/bin/
# macOS Apple Silicon：把 linux_x86_64 换成 aarch64；Windows 用 zip 包
\`\`\`

### 方式 3：从源码构建
\`\`\`bash
bash scripts/install.sh   # 编译到 ~/.local/bin/
\`\`\`

### 方式 4：go install
\`\`\`bash
go install github.com/scagogogo/npm-skills/cmd/npm-skills@latest
go install github.com/scagogogo/npm-skills/cmd/mcp-server@latest
\`\`\`

## 第二步：验证安装
\`\`\`bash
npm-skills --version
npm-skills mirrors          # 列出 8 个镜像源
\`\`\`

## 第三步：常用命令速查（90% 场景）

### 查包信息
\`\`\`bash
npm-skills package-summary <name>        # 轻量（推荐，KB 级响应）
npm-skills package <name>                # 完整元数据（可能 10MB+）
npm-skills versions <name> --latest      # 最新版本
npm-skills pkg-version <name> <ver>      # 特定版本详情
\`\`\`

### 搜索与统计
\`\`\`bash
npm-skills search "<query>" -l 10                  # 搜索包
npm-skills download-stats <name> -p last-month    # 下载量
npm-skills download-range <name> -p last-week     # 每日趋势
npm-skills download-stats-bulk react,vue -p last-month  # 批量对比
\`\`\`

### 镜像与代理（中国大陆用户重点）
\`\`\`bash
npm-skills package react -m npm-mirror                              # 国内镜像，免代理
npm-skills package react --proxy http://127.0.0.1:7890              # HTTP 代理
npm-skills package my-lib --registry https://npm.my-company.com     # 私有仓库
\`\`\`
环境变量：\`NPM_MIRROR=npm-mirror\`、\`NPM_PROXY=...\`、\`NPM_TOKEN=...\`

### 写操作（需要 token）
\`\`\`bash
npm-skills publish ./pkg.tgz --name my-pkg --version 1.0.0 -t npm_xxxxx
npm-skills deprecate my-pkg 1.0.0 -M "Use v2.0.0" -t npm_xxxxx
npm-skills dist-tags set my-pkg next --version 2.0.0-rc.1 -t npm_xxxxx
npm-skills access get my-pkg -t npm_xxxxx
npm-skills audit quick --deps "lodash=4.17.11,express=4.17.1"
\`\`\`

## 关键规则
1. **所有命令输出 JSON 到 stdout**，状态信息走 stderr —— 你可以用 \`jq\` 解析。
2. **读取操作免认证**；写入操作必须 \`--token\` 或 \`NPM_TOKEN\`。用 \`whoami\` 验证。
3. **优先 \`package-summary\`** 而非 \`package\`，响应小、快。
4. **下载统计始终查 api.npmjs.org**，与镜像设置无关。
5. 中国大陆用户默认 \`-m npm-mirror\`；受限网络叠 \`--proxy\`；私有仓库用 \`--registry\`。
6. 类型化错误：\`ErrNotFound\`/\`ErrUnauthorized\`/\`ErrRateLimited\` 等，可程序化处理。

## 当用户提出 NPM 相关需求时
直接选对应命令执行，例如：
- "查 axios 信息" → \`npm-skills package-summary axios\`
- "react 上个月下载量" → \`npm-skills download-stats react -p last-month\`
- "国内镜像查 vue" → \`npm-skills package vue -m npm-mirror\`
- "发布我的包" → \`npm-skills publish ... -t $NPM_TOKEN\`
- "审计依赖安全" → \`npm-skills audit quick --deps "..."\`

完整 26 命令见 \`npm-skills --help\`，SDK 文档见 https://scagogogo.github.io/npm-skills/api/`

const promptEn = `---
title: Use NPM Skills to operate on the NPM Registry
---

# Your Task

You now have access to **NPM Skills** — a complete NPM Registry client supporting package query, version management, dist-tags, download stats, publishing, access control, security audit, org/team, webhooks, etc. Includes 8 mirrors (China: no proxy needed) and HTTP/SOCKS5 proxy support.

## Step 1: Install (pick one based on your environment; option 1 recommended)

### Option 1: Claude Code plugin (AI-native, recommended)
\`\`\`bash
claude plugin marketplace add scagogogo/npm-skills
claude plugin install npm@npm-skills
\`\`\`
After install, you (AI) auto-discover SKILL.md and invoke on demand — no manual shell needed.

### Option 2: Download prebuilt binary (without Claude Code)
\`\`\`bash
# Linux x86_64
curl -sL https://github.com/scagogogo/npm-skills/releases/latest/download/npm-skills_0.2.0_linux_x86_64.tar.gz | tar -xz
sudo mv npm-skills npm-mcp-server /usr/local/bin/
# macOS Apple Silicon: replace linux_x86_64 with aarch64; Windows: use .zip
\`\`\`

### Option 3: Build from source
\`\`\`bash
bash scripts/install.sh   # Compiles to ~/.local/bin/
\`\`\`

### Option 4: go install
\`\`\`bash
go install github.com/scagogogo/npm-skills/cmd/npm-skills@latest
go install github.com/scagogogo/npm-skills/cmd/mcp-server@latest
\`\`\`

## Step 2: Verify installation
\`\`\`bash
npm-skills --version
npm-skills mirrors          # List 8 mirror sources
\`\`\`

## Step 3: Common commands cheat sheet (90% of cases)

### Get package info
\`\`\`bash
npm-skills package-summary <name>        # Lightweight (recommended, KB response)
npm-skills package <name>                # Full metadata (can be 10MB+)
npm-skills versions <name> --latest      # Latest version
npm-skills pkg-version <name> <ver>      # Specific version details
\`\`\`

### Search and stats
\`\`\`bash
npm-skills search "<query>" -l 10                  # Search packages
npm-skills download-stats <name> -p last-month    # Download count
npm-skills download-range <name> -p last-week     # Daily trend
npm-skills download-stats-bulk react,vue -p last-month  # Bulk compare
\`\`\`

### Mirrors and proxy (China users: key section)
\`\`\`bash
npm-skills package react -m npm-mirror                              # China mirror, no proxy
npm-skills package react --proxy http://127.0.0.1:7890              # HTTP proxy
npm-skills package my-lib --registry https://npm.my-company.com     # Private registry
\`\`\`
Env vars: \`NPM_MIRROR=npm-mirror\`, \`NPM_PROXY=...\`, \`NPM_TOKEN=...\`

### Write operations (require token)
\`\`\`bash
npm-skills publish ./pkg.tgz --name my-pkg --version 1.0.0 -t npm_xxxxx
npm-skills deprecate my-pkg 1.0.0 -M "Use v2.0.0" -t npm_xxxxx
npm-skills dist-tags set my-pkg next --version 2.0.0-rc.1 -t npm_xxxxx
npm-skills access get my-pkg -t npm_xxxxx
npm-skills audit quick --deps "lodash=4.17.11,express=4.17.1"
\`\`\`

## Key rules
1. **All commands output JSON to stdout**; status/info goes to stderr — parse with \`jq\`.
2. **Read ops = no auth**; write ops need \`--token\` or \`NPM_TOKEN\`. Verify with \`whoami\`.
3. **Prefer \`package-summary\`** over \`package\` — smaller, faster.
4. **Download stats always query api.npmjs.org** regardless of mirror.
5. China users: default to \`-m npm-mirror\`; restricted networks: add \`--proxy\`; private registries: \`--registry\`.
6. Typed errors: \`ErrNotFound\`/\`ErrUnauthorized\`/\`ErrRateLimited\` for programmatic handling.

## When the user asks for NPM-related tasks
Pick the matching command and run:
- "Check axios info" → \`npm-skills package-summary axios\`
- "React downloads last month" → \`npm-skills download-stats react -p last-month\`
- "Query vue via China mirror" → \`npm-skills package vue -m npm-mirror\`
- "Publish my package" → \`npm-skills publish ... -t $NPM_TOKEN\`
- "Audit dependencies" → \`npm-skills audit quick --deps "..."\`

Full 26 commands: \`npm-skills --help\`. SDK docs: https://scagogogo.github.io/npm-skills/en/api/`

const prompt = computed(() => (isZh.value ? promptZh : promptEn))

const copied = ref(false)
async function copyPrompt() {
  try {
    await navigator.clipboard.writeText(prompt.value)
    copied.value = true
    setTimeout(() => (copied.value = false), 2000)
  } catch {
    // 剪贴板 API 不可用时静默失败，用户仍可手动选择复制
  }
}
</script>

<template>
  <section class="home-prompt">
    <div class="home-prompt-container">
      <div class="home-prompt-header">
        <h2>
          <span class="home-prompt-emoji">📋</span>
          {{ isZh ? '一键复制给 AI Agent' : 'One-Click Copy for AI Agents' }}
        </h2>
        <button class="home-prompt-copy" :class="{ copied }" @click="copyPrompt">
          <span v-if="copied">✓ {{ isZh ? '已复制' : 'Copied' }}</span>
          <span v-else>{{ isZh ? '复制提示词' : 'Copy prompt' }}</span>
        </button>
      </div>
      <p class="home-prompt-desc">
        {{ isZh
          ? '把下面这段提示词复制粘贴给你的 AI Agent（Claude Code、Cursor、Windsurf、ChatGPT 等），Agent 会顺着指引自动安装并使用 NPM Skills，无需你逐条解释命令：'
          : 'Copy the prompt below and paste it into your AI agent (Claude Code, Cursor, Windsurf, ChatGPT, etc.). The agent will follow the instructions to install and use NPM Skills automatically — no need to explain commands line by line.' }}
      </p>
      <div class="home-prompt-code">
        <pre><code>{{ prompt }}</code></pre>
      </div>
      <p class="home-prompt-tip">
        {{ isZh
          ? '💡 Agent 会按提示词完成安装、验证并开始执行 NPM 相关任务。'
          : '💡 The agent will handle installation, verification, and NPM tasks for you.' }}
      </p>
    </div>
  </section>
</template>

<style scoped>
.home-prompt {
  padding: 28px 24px 8px;
}

.home-prompt-container {
  max-width: 1152px;
  margin: 0 auto;
  border: 1px solid var(--vp-c-divider);
  border-radius: 12px;
  background: var(--vp-c-bg-soft);
  padding: 24px 28px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.04);
}

.home-prompt-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 12px;
}

.home-prompt-header h2 {
  margin: 0;
  font-size: 22px;
  font-weight: 700;
  border-top: none;
  padding-top: 0;
  letter-spacing: -0.02em;
}

.home-prompt-emoji {
  margin-right: 4px;
}

.home-prompt-copy {
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 18px;
  font-size: 14px;
  font-weight: 600;
  color: var(--vp-button-brand-text);
  background: var(--vp-button-brand-bg);
  border: 1px solid var(--vp-button-brand-bg);
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.2s, border-color 0.2s, transform 0.1s;
}

.home-prompt-copy:hover {
  background: var(--vp-button-brand-hover-bg);
  border-color: var(--vp-button-brand-hover-bg);
}

.home-prompt-copy:active {
  transform: translateY(1px);
}

.home-prompt-copy.copied {
  background: #16a34a;
  border-color: #16a34a;
}

.home-prompt-desc {
  margin: 0 0 16px;
  color: var(--vp-c-text-2);
  font-size: 15px;
  line-height: 1.6;
}

.home-prompt-code {
  margin: 0;
  border-radius: 8px;
  overflow: hidden;
  background: var(--vp-code-block-bg);
  border: 1px solid var(--vp-c-divider);
}

.home-prompt-code pre {
  margin: 0;
  padding: 16px 18px;
  overflow-x: auto;
  font-size: 13px;
  line-height: 1.6;
  max-height: 480px;
}

.home-prompt-code code {
  font-family: var(--vp-font-family-mono);
  color: var(--vp-c-text-1);
  white-space: pre;
}

.home-prompt-tip {
  margin: 14px 0 0;
  color: var(--vp-c-text-2);
  font-size: 14px;
  line-height: 1.6;
}

@media (max-width: 768px) {
  .home-prompt {
    padding: 20px 16px 4px;
  }
  .home-prompt-container {
    padding: 18px 16px;
  }
  .home-prompt-header h2 {
    font-size: 19px;
  }
  .home-prompt-code pre {
    font-size: 12px;
    padding: 12px;
  }
}
</style>
