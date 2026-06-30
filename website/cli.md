# CLI 命令手册

`npm-skills` CLI 共 26 个命令，所有命令输出 JSON 到 stdout（便于 AI 解析），状态信息走 stderr。

## 全局参数

| 参数 | 简写 | 默认值 | 说明 |
|------|------|--------|------|
| `--mirror` | `-m` | `official` | 镜像源名 |
| `--registry` | | | 自定义仓库 URL（覆盖 --mirror） |
| `--token` | `-t` | | NPM 认证 token（写操作必需，env: `NPM_TOKEN`） |
| `--proxy` | | | HTTP 代理 URL（env: `NPM_PROXY`） |
| `--timeout` | | `120` | 请求超时秒数 |

**优先级**：CLI 参数 > 环境变量 > 默认值

## 读取操作

### 包信息

```bash
npm-skills package-summary <name>     # 轻量包信息（推荐）
npm-skills package <name>             # 完整元数据（可能 10MB+）
npm-skills pkg-version <name> <ver>   # 特定版本
npm-skills versions <name>            # 所有版本
npm-skills versions <name> --latest   # 仅最新版本
```

> **提示**：优先用 `package-summary`，响应小得多、快得多。

### 搜索

```bash
npm-skills search <query>                  # 基础搜索
npm-skills search <query> -l 10            # 限制结果数
npm-skills search <query> --from 20 -l 10  # 分页
npm-skills search <query> --popularity 1.0 # 按流行度加权
```

| 参数 | 简写 | 默认 | 说明 |
|------|------|------|------|
| `--limit` | `-l` | 20 | 最大结果数 |
| `--from` | | 0 | 分页偏移 |
| `--quality` | | 0 | 质量权重 (0-1) |
| `--popularity` | | 0 | 流行度权重 (0-1) |
| `--maintenance` | | 0 | 维护度权重 (0-1) |

### Dist-Tags（读取）

```bash
npm-skills dist-tags get <name>
```

### 下载统计

```bash
npm-skills download-stats <name> -p last-month          # 单包
npm-skills download-range <name> -p last-week           # 每日趋势
npm-skills download-stats-date <name> --start 2024-01-01 --end 2024-06-30  # 自定义区间
npm-skills download-stats-bulk react,vue,angular -p last-month  # 批量（≤128）
```

> 下载统计始终查询 api.npmjs.org，与镜像/仓库设置无关。

### 其他读取

```bash
npm-skills registry-info                 # 仓库健康信息
npm-skills mirrors                       # 镜像源列表
npm-skills config                        # 当前配置
npm-skills whoami --token <token>        # 认证状态
npm-skills download <name> <ver> <dest>  # 下载 tarball
```

## 写入操作（需要 --token）

所有写操作都需要认证。用 `--token` 或设置 `NPM_TOKEN`。

### 发布 / 取消发布 / 弃用

```bash
npm-skills publish ./pkg.tgz --name my-pkg --version 1.0.0 -t <token>
npm-skills deprecate my-pkg 1.0.0 -M "Use v2.0.0" -t <token>
npm-skills unpublish my-pkg --version 1.0.0 -t <token>   # 危险
npm-skills unpublish my-pkg --force -t <token>           # 极危险
```

### Dist-Tags 管理

```bash
npm-skills dist-tags set <name> <tag> --version <ver> -t <token>
npm-skills dist-tags delete <name> <tag> -t <token>
```

### 访问控制与协作者

```bash
npm-skills access get <name> -t <token>
npm-skills access set <name> --visibility public -t <token>
npm-skills access collaborators <name> -t <token>
npm-skills access grant <name> <user> --permission read -t <token>
npm-skills access revoke <name> <user> -t <token>
```

### Stars

```bash
npm-skills star add <name> -t <token>
npm-skills star remove <name> -t <token>
npm-skills star list <username>
npm-skills star stargazers <name>
```

### Token 管理

```bash
npm-skills token list -t <token>
npm-skills token get <id> -t <token>
npm-skills token create --password <pass> -t <token>
npm-skills token delete <id> -t <token>
```

### 安全审计

```bash
npm-skills audit quick --deps "lodash=4.17.11,express=4.17.1"
npm-skills audit bulk --advisories "lodash=<4.17.12"
npm-skills audit advisory 123
npm-skills audit advisories --package lodash
```

### 组织与团队

```bash
npm-skills org get <org> -t <token>
npm-skills org members <org> -t <token>
npm-skills org packages <org> -t <token>
npm-skills org team-list <org> -t <token>
npm-skills org team-members <org> <team> -t <token>
# ... 完整列表见 npm-skills --help
```

### Webhooks

```bash
npm-skills hook list -t <token>
npm-skills hook get <id> -t <token>
npm-skills hook create --name my-hook --endpoint https://... -t <token>
npm-skills hook update <id> --endpoint https://new... -t <token>
npm-skills hook delete <id> -t <token>
```

## 镜像源

| 镜像 | 名称 | 地域 |
|------|------|------|
| `https://registry.npmjs.org` | `official` | 全球 |
| `https://registry.npmmirror.com` | `npm-mirror` | 中国（推荐） |
| `https://registry.npm.taobao.org` | `taobao` | 中国 |
| `https://mirrors.huaweicloud.com/repository/npm` | `huawei` | 中国 |
| `http://mirrors.cloud.tencent.com/npm` | `tencent` | 中国 |
| `http://r.cnpmjs.org` | `cnpm` | 中国 |
| `https://registry.yarnpkg.com` | `yarn` | 全球 |
| `https://skimdb.npmjs.com` | `npmjscom` | 全球 |

可直接传 URL：`--mirror https://your-registry.com`
