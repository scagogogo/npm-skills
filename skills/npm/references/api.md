# NPM Crawler — Complete API Reference

> This reference is loaded on demand when deeper SDK details are needed.
> For quick CLI usage, see the main SKILL.md.

## CLI Reference

### Global Flags

All commands support these persistent flags:

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--token` | `-t` | | NPM authentication token (env: NPM_TOKEN) |
| `--mirror` | `-m` | `official` | Mirror source: official\|taobao\|npm-mirror\|huawei\|tencent\|cnpm\|yarn\|npmjscom, or any URL |
| `--registry` | | | Custom registry URL (overrides --mirror) |
| `--proxy` | | | HTTP proxy URL (e.g. http://127.0.0.1:7890) |
| `--timeout` | | `120` | Request timeout in seconds |
| `--no-color` | | `false` | Disable colored output |

### Environment Variables

| Variable | Equivalent Flag | Example |
|----------|----------------|---------|
| `NPM_MIRROR` | `--mirror` | `NPM_MIRROR=npm-mirror` |
| `NPM_REGISTRY` | `--registry` | `NPM_REGISTRY=https://npm.company.com` |
| `NPM_PROXY` | `--proxy` | `NPM_PROXY=http://127.0.0.1:7890` |
| `NPM_TOKEN` | `--token` | `NPM_TOKEN=npm_xxxxx` |

Priority: CLI flag > Environment variable > Default

---

### Command: `registry-info`

Get NPM registry status and statistics.

```bash
npm-skills registry-info
npm-skills info -m npm-mirror
npm-skills info --registry https://registry.npmmirror.com
```

**Arguments:** None

**Output fields:**
- `db_name` — Database name (usually "registry")
- `doc_count` — Total number of packages
- `doc_del_count` — Deleted documents count
- `disk_size` — Disk usage in bytes
- `data_size` — Data size in bytes
- `update_seq` — Update sequence number
- `instance_start_time` — Registry instance start timestamp

---

### Command: `package`

Get complete package metadata.

```bash
npm-skills package <name>
npm-skills pkg <name> -m taobao
npm-skills package <name> --proxy http://127.0.0.1:7890
```

**Arguments:**
- `name` (required): Package name, e.g. "react", "@nestjs/core"

**Output fields:**
- `name` — Package name
- `description` — Package description
- `dist-tags` — Map of tags to versions (e.g. `{"latest": "18.2.0"}`)
- `versions` — Map of version strings to Version objects
- `maintainers` — Array of maintainer objects `{name, email}`
- `time` — Map of timestamps (created, modified, per-version)
- `repository` — Repository info `{type, url}`
- `readme` — Full README content
- `homepage` — Project homepage URL
- `license` — License identifier
- `deprecated` — Deprecation notice (string or boolean)
- `keywords` — Array of keyword strings
- `author` — Author object `{name, email}`

> **Tip:** For most queries, use `package-summary` instead — much smaller response (KB vs MB).

---

### Command: `package-summary`

Get lightweight package metadata.

```bash
npm-skills package-summary <name>
npm-skills ps <name> -m taobao
npm-skills pkgsum <name> --proxy http://127.0.0.1:7890
```

**Arguments:**
- `name` (required): Package name

**Output fields:** Same as `package` but may lack `readme`, `maintainers`, and some other fields. Response is typically 10-100x smaller.

---

### Command: `search`

Search NPM registry by keyword.

```bash
npm-skills search <query>
npm-skills s <query> -l 10 -m npm-mirror
npm-skills s <query> --from 20 -l 10
npm-skills s <query> --popularity 1.0 --quality 0.0
```

**Arguments:**
- `query` (required): Search keywords (quote multi-word queries)

**Flags:**
| Flag | Default | Description |
|------|---------|-------------|
| `-l, --limit` | `20` | Max results |
| `--from` | `0` | Pagination offset (0-based) |
| `--quality` | `0` | Quality weight (0.0-1.0) |
| `--popularity` | `0` | Popularity weight (0.0-1.0) |
| `--maintenance` | `0` | Maintenance weight (0.0-1.0) |

**Output fields:**
- `objects` — Array of search results
  - `package.name` — Package name
  - `package.version` — Latest version
  - `package.description` — Description
  - `package.keywords` — Keywords array
  - `package.links` — Links `{npm, homepage, repository, bugs}`
  - `score.detail` — Quality, popularity, maintenance scores
  - `searchScore` — Overall search relevance score
- `total` — Total matching packages

---

### Command: `pkg-version`

Get metadata for a specific package version.

```bash
npm-skills pkg-version <name> <version>
npm-skills ver <name> latest -m npm-mirror
```

**Arguments:**
- `name` (required): Package name
- `version` (required): Version string or "latest"

**Output fields:**
- `name`, `version`, `description`, `main`
- `scripts` — npm scripts map
- `dependencies`, `devDependencies` — Dependency maps
- `dist` — Distribution info:
  - `tarball` — Download URL
  - `shasum` — SHA-1 hash
  - `integrity` — SHA-512 integrity hash
- `repository`, `license`, `homepage`, `keywords`

---

### Command: `versions`

List all published versions of a package.

```bash
npm-skills versions <name>
npm-skills vs <name> --latest
npm-skills vers <name> -m npm-mirror
```

**Arguments:**
- `name` (required): Package name

**Flags:**
| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--latest` | `-L` | `false` | Show only the latest version |

**Output fields:**
- `package` — Package name
- `version_count` — Total number of published versions
- `versions` — Array of version strings (sorted)
- (with `--latest`) `latest` — Latest version string

---

### Command: `dist-tags`

Get distribution tags for a package.

```bash
npm-skills dist-tags <name>
npm-skills tags <name> --abbreviated
npm-skills dt <name> -m npm-mirror
```

**Arguments:**
- `name` (required): Package name

**Flags:**
| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--abbreviated` | `-a` | `false` | Use lightweight dist-tags endpoint (faster) |

**Output fields:**
- Map of tag names to version strings, e.g. `{"latest": "18.2.0", "next": "19.0.0-rc.1"}`

---

### Command: `download-stats`

Get download statistics for a package.

```bash
npm-skills download-stats <name>
npm-skills stats <name> -p last-month
```

**Arguments:**
- `name` (required): Package name

**Flags:**
- `-p, --period string` — `last-day`|`last-week`|`last-month` (default: `last-week`)

> **Note:** Download stats always query api.npmjs.org regardless of --mirror/--registry. Proxy settings are still applied.

**Output fields:**
- `downloads` — Number of downloads in the period
- `start` — Period start date (YYYY-MM-DD)
- `end` — Period end date (YYYY-MM-DD)
- `package` — Package name

---

### Command: `download-stats-date`

Get download statistics for a custom date range.

```bash
npm-skills download-stats-date <name> --start 2024-01-01 --end 2024-01-31
npm-skills dsd <name> --start 2024-06-01 --end 2024-06-30
npm-skills stats-date <name> --start 2024-01-01 --end 2024-12-31
```

**Arguments:**
- `name` (required): Package name

**Flags:**
- `--start` (required): Start date (YYYY-MM-DD)
- `--end` (required): End date (YYYY-MM-DD)

> **Note:** Download stats always query api.npmjs.org regardless of --mirror/--registry.

**Output fields:**
- `downloads` — Total downloads in the date range
- `start` — Start date
- `end` — End date
- `package` — Package name

---

### Command: `download-stats-bulk`

Get download statistics for multiple packages at once (up to 128).

```bash
npm-skills download-stats-bulk react,vue,angular
npm-skills dsb react vue angular -p last-month
npm-skills stats-bulk react,vue --range
```

**Arguments:**
- `names...` (required): Package names (comma-separated or space-separated)

**Flags:**
| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `-p, --period` | | `last-week` | last-day, last-week, last-month |
| `--range` | | `false` | Get daily trends instead of totals |

> **Note:** Download stats always query api.npmjs.org regardless of --mirror/--registry.

**Output fields (without --range):**
- Map of package name to `DownloadStats` (downloads, start, end, package)

**Output fields (with --range):**
- Map of package name to `DownloadRangeStats` (daily array, start, end, package)

**Arguments:**
- `name` (required): Package name

**Flags:**
- `-p, --period string` — `last-day`|`last-week`|`last-month` (default: `last-week`)

> **Note:** Download stats always query api.npmjs.org regardless of --mirror/--registry. Proxy settings are still applied.

**Output fields:**
- `downloads` — Number of downloads in the period
- `start` — Period start date (YYYY-MM-DD)
- `end` — Period end date (YYYY-MM-DD)
- `package` — Package name

---

### Command: `download-range`

Get daily download trends for a package.

```bash
npm-skills download-range <name>
npm-skills dr <name> -p last-month
npm-skills range <name> -p last-day
```

**Arguments:**
- `name` (required): Package name

**Flags:**
- `-p, --period string` — `last-day`|`last-week`|`last-month` (default: `last-week`)

> **Note:** Download stats always query api.npmjs.org regardless of --mirror/--registry. Proxy settings are still applied.

**Output fields:**
- `downloads` — Array of daily download counts:
  - `day` — Date (YYYY-MM-DD)
  - `downloads` — Download count for that day
- `start` — Period start date
- `end` — Period end date
- `package` — Package name

---

### Command: `download`

Download a package tarball (.tgz) file.

```bash
npm-skills download <name> <version> <dest>
npm-skills dl <name> <version> <dest> -m npm-mirror
npm-skills download <name> <version> <dest> --proxy http://127.0.0.1:7890
```

**Arguments:**
- `name` (required): Package name
- `version` (required): Version string or "latest"
- `dest` (required): Local file path to save the tarball

**Output fields:**
- `package` — Package name
- `version` — Version string
- `path` — Saved file path
- `source` — Mirror/registry used
- `status` — "downloaded"

---

### Command: `mirrors`

List all available mirror sources.

```bash
npm-skills mirrors
```

**Output:** Array of mirror objects with fields:
- `name` — Mirror identifier for use with --mirror
- `url` — Registry URL
- `region` — Geographic region
- `description` — Human-readable description

---

### Command: `whoami`

Check current NPM authentication status.

```bash
npm-skills whoami --token npm_xxxxx
NPM_TOKEN=npm_xxxxx npm-skills whoami
npm-skills me --token npm_xxxxx -m npm-mirror
```

**Arguments:** None

**Flags:**
| Flag | Default | Description |
|------|---------|-------------|
| `--token` | | NPM authentication token (env: NPM_TOKEN) |

**Output fields:**
- `username` — Authenticated username
- `status` — "authenticated"

---

### Command: `config`

Show current effective configuration.

```bash
npm-skills config
npm-skills cfg -m npm-mirror
npm-skills conf --registry https://registry.npmmirror.com --proxy http://127.0.0.1:7890
```

**Arguments:** None

**Output fields:**
- `registry_url` — Effective registry URL
- `mirror` — Mirror name or URL
- `proxy` — Proxy URL or "(none)"
- `timeout` — Request timeout in seconds

---

## Write Operations

> All write operations require authentication via `--token`/`-t` or the `NPM_TOKEN` environment variable, unless otherwise noted.

---

### Command: `publish`

Publish an NPM package to the registry.

```bash
npm-skills publish <tarball-path> --name <name> --version <version> --token npm_xxxxx
npm-skills pub ./my-package-1.0.0.tgz --name my-package --version 1.0.0 -t npm_xxxxx
npm-skills publish ./pkg.tgz --name @myorg/pkg --version 2.0.0 --description "My package" -t npm_xxxxx
```

**Arguments:**
- `tarball-path` (required): Path to the .tgz file to publish

**Flags:**
| Flag | Default | Description |
|------|---------|-------------|
| `--name` | | Package name (required) |
| `--version` | | Package version (required) |
| `--description` | | Package description |

**Output fields:**
- `package` — Published package name
- `version` — Published version
- `status` — "published"

> **Note:** Requires authentication. This is a destructive operation — make sure you have permission to publish.

---

### Command: `unpublish`

Unpublish an NPM package or a specific version.

```bash
npm-skills unpublish <package> --version <version> --token npm_xxxxx
npm-skills unpublish <package> --force --token npm_xxxxx
npm-skills unpub my-pkg --version 1.0.0 -t npm_xxxxx
```

**Arguments:**
- `package` (required): Package name to unpublish

**Flags:**
| Flag | Default | Description |
|------|---------|-------------|
| `--version` | | Version to unpublish (unpublish specific version) |
| `--force` | `false` | Unpublish entire package (dangerous) |

**Output fields:**
- `package` — Package name
- `version` — Version string (when unpublishing a specific version)
- `status` — "unpublished"

> **WARNING:** Destructive and potentially irreversible. Most registries only allow unpublish within 72 hours. You must specify either `--version` or `--force`.

---

### Command: `deprecate`

Deprecate a specific version of an NPM package.

```bash
npm-skills deprecate <package> <version> --message "Use v2.0.0 instead" --token npm_xxxxx
npm-skills dep my-pkg 1.0.0 -M "Security vulnerability, upgrade to 1.0.1" -t npm_xxxxx
```

**Arguments:**
- `package` (required): Package name
- `version` (required): Version string to deprecate

**Flags:**
| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--message` | `-M` | | Deprecation message (required) |

**Output fields:**
- `package` — Package name
- `version` — Deprecated version
- `message` — Deprecation message
- `status` — "deprecated"

> **Note:** Prefer deprecation over unpublish. It is safer and does not remove the version.

---

### Command: `dist-tags set`

Set a distribution tag to point to a specific version.

```bash
npm-skills dist-tags set <package> <tag> --version <version> --token npm_xxxxx
npm-skills dist-tags add react next --version 2.0.0-rc.1 -t npm_xxxxx
```

**Arguments:**
- `package` (required): Package name
- `tag` (required): Tag name (e.g. next, beta, stable)

**Flags:**
| Flag | Default | Description |
|------|---------|-------------|
| `--version` | | Version to point the tag to (required) |

**Output fields:**
- `package` — Package name
- `tag` — Tag name
- `version` — Version the tag now points to
- `status` — "updated"

> **Note:** If the tag already exists, it will be updated. Requires authentication.

---

### Command: `dist-tags delete`

Delete a distribution tag from a package.

```bash
npm-skills dist-tags delete <package> <tag> --token npm_xxxxx
npm-skills dist-tags rm my-pkg next -t npm_xxxxx
```

**Arguments:**
- `package` (required): Package name
- `tag` (required): Tag name to delete

**Output fields:**
- `package` — Package name
- `tag` — Deleted tag name
- `status` — "deleted"

> **WARNING:** Deleting the `latest` tag may cause issues. Requires authentication.

---

### Command: `access`

Package access and collaborator management. Subcommands: `get`, `set`, `collaborators`, `grant`, `revoke`.

#### Subcommand: `access get`

Get package access settings.

```bash
npm-skills access get <package> --token npm_xxxxx
```

**Arguments:**
- `package` (required): Package name

Requires authentication.

#### Subcommand: `access set`

Set package access level (public/restricted).

```bash
npm-skills access set <package> --visibility public --token npm_xxxxx
npm-skills access set <package> --visibility restricted -t npm_xxxxx
```

**Arguments:**
- `package` (required): Package name

**Flags:**
| Flag | Default | Description |
|------|---------|-------------|
| `--visibility` | | Access level: `public` or `restricted` (required) |

> **Note:** Changing from public to restricted will make the package inaccessible to unauthorized users.

#### Subcommand: `access collaborators`

List package collaborators.

```bash
npm-skills access collaborators <package> --token npm_xxxxx
npm-skills access collabs <package> -t npm_xxxxx
```

**Arguments:**
- `package` (required): Package name

#### Subcommand: `access grant`

Grant user access to a package.

```bash
npm-skills access grant <package> <user> --permission write --token npm_xxxxx
npm-skills access grant <package> myorg:devteam --permission read -t npm_xxxxx
```

**Arguments:**
- `package` (required): Package name
- `user` (required): Username or `org:team` format

**Flags:**
| Flag | Default | Description |
|------|---------|-------------|
| `--permission` | `read` | Permission level: `read` or `write` |

#### Subcommand: `access revoke`

Revoke user access from a package.

```bash
npm-skills access revoke <package> <user> --token npm_xxxxx
```

**Arguments:**
- `package` (required): Package name
- `user` (required): Username to revoke

---

### Command: `star`

Star/unstar NPM packages. Subcommands: `add`, `remove`, `list`, `stargazers`.

#### Subcommand: `star add`

Star (favorite) an NPM package.

```bash
npm-skills star add <package> --token npm_xxxxx
```

**Arguments:**
- `package` (required): Package name

**Output fields:** `package`, `status: "starred"`. Requires authentication.

#### Subcommand: `star remove`

Unstar (unfavorite) an NPM package.

```bash
npm-skills star remove <package> --token npm_xxxxx
npm-skills star unstar <package> -t npm_xxxxx
```

**Arguments:**
- `package` (required): Package name

**Output fields:** `package`, `status: "unstarred"`. Requires authentication.

#### Subcommand: `star list`

List packages starred by a user.

```bash
npm-skills star list <username>
npm-skills star ls <username>
```

**Arguments:**
- `username` (required): NPM username

No authentication required.

#### Subcommand: `star stargazers`

List users who starred a package.

```bash
npm-skills star stargazers <package>
npm-skills star sg <package>
```

**Arguments:**
- `package` (required): Package name

No authentication required.

---

### Command: `token`

Manage NPM access tokens. All subcommands require authentication.

#### Subcommand: `token list`

List all NPM access tokens.

```bash
npm-skills token list --token npm_xxxxx
npm-skills token ls -t npm_xxxxx
```

#### Subcommand: `token get`

Get details of a specific token.

```bash
npm-skills token get <token-id> --token npm_xxxxx
```

**Arguments:**
- `token-id` (required): Token ID

#### Subcommand: `token create`

Create a new access token.

```bash
npm-skills token create --password <pass> --token npm_xxxxx
npm-skills token create --password <pass> --readonly -t npm_xxxxx
npm-skills token create --password <pass> --cidr 192.168.1.0/24 -t npm_xxxxx
npm-skills token new --password <pass> -t npm_xxxxx
```

**Flags:**
| Flag | Default | Description |
|------|---------|-------------|
| `--password` | | Current user password (required) |
| `--readonly` | `false` | Create a read-only token |
| `--cidr` | | IP whitelist CIDR ranges (can be specified multiple times) |

#### Subcommand: `token delete`

Delete (revoke) an access token.

```bash
npm-skills token delete <token-id> --token npm_xxxxx
npm-skills token rm <token-id> -t npm_xxxxx
npm-skills token revoke <token-id> -t npm_xxxxx
```

**Arguments:**
- `token-id` (required): Token ID to delete

> **Note:** You cannot delete the token you are currently using.

---

### Command: `audit`

Security audit for NPM packages. Subcommands: `quick`, `bulk`, `advisory`, `advisories`.

#### Subcommand: `audit quick`

Quick security audit of dependencies.

```bash
npm-skills audit quick --deps "lodash=4.17.11,express=4.17.1"
```

**Flags:**
| Flag | Default | Description |
|------|---------|-------------|
| `--deps` | | Dependencies to audit, format: `name=version,name2=version2` (required) |

**Output fields:**
- `metadata.vulnerabilities` — Counts by severity: `low`, `moderate`, `high`, `critical`
- `metadata.dependencies` — Total dependencies count
- `metadata.totalDependencies` — Total including dev/optional

No authentication required.

#### Subcommand: `audit bulk`

Bulk security audit with detailed advisories.

```bash
npm-skills audit bulk --advisories "lodash=<4.17.12,express=<4.17.3"
```

**Flags:**
| Flag | Default | Description |
|------|---------|-------------|
| `--advisories` | | Packages to check, format: `name=<version,name2=<version2` (required) |

**Output fields:** Map of package name to array of `Advisory` objects.

No authentication required.

#### Subcommand: `audit advisory`

Get details of a specific security advisory.

```bash
npm-skills audit advisory <id>
```

**Arguments:**
- `id` (required): Advisory ID (numeric)

**Output fields:** Full advisory with `id`, `title`, `severity`, `cve`, `module_name`, `vulnerable_versions`, `patched_versions`, `overview`, `recommendation`.

No authentication required.

#### Subcommand: `audit advisories`

List security advisories with optional filtering.

```bash
npm-skills audit advisories
npm-skills audit advisories --package lodash --per-page 10
```

**Flags:**
| Flag | Default | Description |
|------|---------|-------------|
| `--page` | `0` | Page number |
| `--per-page` | `20` | Results per page |
| `--package` | | Filter by affected package name |

No authentication required.

---

### Command: `org`

Organization and team management. All subcommands require authentication.

#### Subcommand: `org get`

Get organization details.

```bash
npm-skills org get <org> --token npm_xxxxx
```

**Arguments:** `org` (required)

#### Subcommand: `org create`

Create a new organization.

```bash
npm-skills org create <org> --token npm_xxxxx
```

**Arguments:** `org` (required)

#### Subcommand: `org delete`

Delete an organization (irreversible).

```bash
npm-skills org delete <org> --token npm_xxxxx
npm-skills org rm <org> -t npm_xxxxx
```

**Arguments:** `org` (required)

#### Subcommand: `org members`

List organization members.

```bash
npm-skills org members <org> --token npm_xxxxx
npm-skills org member <org> -t npm_xxxxx
```

**Arguments:** `org` (required)

#### Subcommand: `org member-add`

Add a member to an organization.

```bash
npm-skills org member-add <org> <username> --token npm_xxxxx
```

**Arguments:** `org` (required), `username` (required)

#### Subcommand: `org member-remove`

Remove a member from an organization.

```bash
npm-skills org member-remove <org> <username> --token npm_xxxxx
npm-skills org member-rm <org> <username> -t npm_xxxxx
```

**Arguments:** `org` (required), `username` (required)

#### Subcommand: `org packages`

List packages owned by an organization.

```bash
npm-skills org packages <org> --token npm_xxxxx
npm-skills org pkgs <org> -t npm_xxxxx
```

**Arguments:** `org` (required)

#### Subcommand: `org team-list`

List teams in an organization.

```bash
npm-skills org team-list <org> --token npm_xxxxx
npm-skills org teams <org> -t npm_xxxxx
```

**Arguments:** `org` (required)

#### Subcommand: `org team-create`

Create a team in an organization.

```bash
npm-skills org team-create <org> <team> --token npm_xxxxx
```

**Arguments:** `org` (required), `team` (required)

#### Subcommand: `org team-delete`

Delete a team from an organization (irreversible).

```bash
npm-skills org team-delete <org> <team> --token npm_xxxxx
npm-skills org team-rm <org> <team> -t npm_xxxxx
```

**Arguments:** `org` (required), `team` (required)

#### Subcommand: `org team-members`

List members of a team.

```bash
npm-skills org team-members <org> <team> --token npm_xxxxx
```

**Arguments:** `org` (required), `team` (required)

#### Subcommand: `org team-member-add`

Add a member to a team. The user must already be an org member.

```bash
npm-skills org team-member-add <org> <team> <username> --token npm_xxxxx
```

**Arguments:** `org` (required), `team` (required), `username` (required)

#### Subcommand: `org team-member-remove`

Remove a member from a team.

```bash
npm-skills org team-member-remove <org> <team> <username> --token npm_xxxxx
npm-skills org team-member-rm <org> <team> <username> -t npm_xxxxx
```

**Arguments:** `org` (required), `team` (required), `username` (required)

#### Subcommand: `org team-packages`

List packages a team has access to.

```bash
npm-skills org team-packages <org> <team> --token npm_xxxxx
```

**Arguments:** `org` (required), `team` (required)

---

### Command: `hook`

Manage NPM webhooks. All subcommands require authentication.

#### Subcommand: `hook list`

List all webhooks, optionally filtered by package.

```bash
npm-skills hook list --token npm_xxxxx
npm-skills hook ls --package react -t npm_xxxxx
```

**Flags:**
| Flag | Default | Description |
|------|---------|-------------|
| `--package` | | Filter by package name |
| `--page` | `0` | Page number |
| `--per-page` | `20` | Results per page |

#### Subcommand: `hook get`

Get details of a specific webhook.

```bash
npm-skills hook get <hook-id> --token npm_xxxxx
```

**Arguments:** `hook-id` (required)

#### Subcommand: `hook create`

Create a new webhook.

```bash
npm-skills hook create --name my-hook --endpoint https://example.com/webhook --package react --token npm_xxxxx
npm-skills hook create --name ci-hook --endpoint https://ci.example.com/hook --secret mysecret --package my-pkg -t npm_xxxxx
```

**Flags:**
| Flag | Default | Description |
|------|---------|-------------|
| `--name` | | Hook name (required) |
| `--endpoint` | | Webhook endpoint URL (required) |
| `--secret` | | Secret for signature verification |
| `--package` | | Package to monitor (empty = all packages) |
| `--events` | | Event types (default: all, can specify multiple) |

#### Subcommand: `hook update`

Update a webhook. Only specified fields are changed.

```bash
npm-skills hook update <hook-id> --endpoint https://new.example.com/webhook --token npm_xxxxx
npm-skills hook update <hook-id> --secret newsecret -t npm_xxxxx
npm-skills hook edit <hook-id> --active --set-active -t npm_xxxxx
```

**Arguments:** `hook-id` (required)

**Flags:**
| Flag | Default | Description |
|------|---------|-------------|
| `--endpoint` | | New endpoint URL |
| `--secret` | | New secret |
| `--events` | | New event types |
| `--active` | `false` | Set hook active/inactive |
| `--set-active` | `false` | Explicitly set the active flag |

#### Subcommand: `hook delete`

Delete a webhook permanently.

```bash
npm-skills hook delete <hook-id> --token npm_xxxxx
npm-skills hook rm <hook-id> -t npm_xxxxx
```

**Arguments:** `hook-id` (required)

---

### Command: `user`

User operations. Subcommands: `login`, `signup`, `get`.

#### Subcommand: `user login`

Login to NPM registry and obtain an authentication token.

```bash
npm-skills user login --username myuser --password mypass
npm-skills user login -u myuser -p mypass -m npm-mirror
```

**Flags:**
| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--username` | `-u` | | Username (required) |
| `--password` | `-p` | | Password (required) |

**Output fields:** `id`, `rev`, `token`, `ok`

No prior authentication required — this command obtains a token.

#### Subcommand: `user signup`

Create a new NPM user account.

```bash
npm-skills user signup --username myuser --password mypass --email me@example.com
npm-skills user signup -u myuser -p mypass --email me@example.com
```

**Flags:**
| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--username` | `-u` | | Username (required) |
| `--password` | `-p` | | Password (required) |
| `--email` | | | Email address (required) |

No prior authentication required.

#### Subcommand: `user get`

Get user profile information.

```bash
npm-skills user get <username> --token npm_xxxxx
npm-skills user info <username> -t npm_xxxxx
```

**Arguments:** `username` (required)

Requires authentication.

---

### Command: `couchdb`

CouchDB advanced query operations. Low-level API access for mirroring and incremental sync. Most users do not need these commands.

#### Subcommand: `couchdb changes`

Get the registry changes feed.

```bash
npm-skills couchdb changes
npm-skills couchdb changes --limit 10
npm-skills couchdb changes --since 12345 --limit 50
npm-skills couchdb changes --include-docs
```

**Flags:**
| Flag | Default | Description |
|------|---------|-------------|
| `--since` | | Start sequence (resume from `last_seq`) |
| `--limit` | `0` | Limit number of results (0 = no limit) |
| `--include-docs` | `false` | Include full document content |

**Output fields:** `last_seq`, `pending`, `results` (array of `ChangeEntry` with `seq`, `id`, `changes`, `deleted`, `doc`)

No authentication required.

#### Subcommand: `couchdb all-docs`

List all document IDs in the registry.

```bash
npm-skills couchdb all-docs --limit 20
npm-skills couchdb all-docs --start-key "@nestjs" --end-key "@nestt" --limit 50
npm-skills couchdb all-docs --include-docs --limit 5
```

**Flags:**
| Flag | Default | Description |
|------|---------|-------------|
| `--start-key` | | Start key for range query |
| `--end-key` | | End key for range query |
| `--limit` | `0` | Limit number of results |
| `--skip` | `0` | Skip number of results |
| `--include-docs` | `false` | Include full document content (may be very large) |
| `--descending` | `false` | Reverse order |

**Output fields:** `total_rows`, `offset`, `rows` (array of `DocRow` with `id`, `key`, `value.rev`, `doc`)

No authentication required.

#### Subcommand: `couchdb view`

Query a CouchDB view.

```bash
npm-skills couchdb view starredByUser --key '"username"'
npm-skills couchdb view starredByPackage --key '"react"'
npm-skills couchdb view byKeyword --key '"http"' --limit 20
```

**Arguments:** `view-name` (required) — Common views: `starredByUser`, `starredByPackage`, `byKeyword`, `byUser`

**Flags:**
| Flag | Default | Description |
|------|---------|-------------|
| `--key` | | Exact key to query |
| `--start-key` | | Start key for range query |
| `--end-key` | | End key for range query |
| `--limit` | `0` | Limit number of results |
| `--skip` | `0` | Skip number of results |
| `--group` | `false` | Group results |
| `--group-level` | `0` | Group level for hierarchical grouping |
| `--descending` | `false` | Reverse order |

**Output fields:** `total_rows`, `offset`, `rows` (array of `ViewRow` with `id`, `key`, `value`)

> **Note:** Not all registries support all views. No authentication required.

---

## Go SDK Reference

### Package: `github.com/scagogogo/npm-skills/pkg/registry`

#### Creating a Client

```go
import "github.com/scagogogo/npm-skills/pkg/registry"

// Default (official npmjs.org)
client := registry.NewRegistry()

// Pre-configured mirrors
client = registry.NewTaoBaoRegistry()
client = registry.NewNpmMirrorRegistry()
client = registry.NewHuaWeiCloudRegistry()
client = registry.NewTencentRegistry()
client = registry.NewCnpmRegistry()
client = registry.NewYarnRegistry()
client = registry.NewNpmjsComRegistry()

// Custom registry with proxy and token
options := registry.NewOptions().
    SetRegistryURL("https://npm.my-company.com").
    SetProxy("http://proxy:8080").
    SetToken("npm_xxxxx")
client := registry.NewRegistry(options)
```

#### Method: `GetRegistryInformation`

```go
info, err := client.GetRegistryInformation(ctx) (*models.RegistryInformation, error)
```

Returns registry status: `DbName`, `DocCount`, `DiskSize`, `DataSize`, `UpdateSeq`, `InstanceStartTime`.

#### Method: `GetPackageInformation`

```go
pkg, err := client.GetPackageInformation(ctx, "react") (*models.Package, error)
```

Returns: `Name`, `Description`, `DistTags`, `Versions`, `Maintainers`, `Time`, `Repository`, `ReadMe`, `Homepage`, `License`, `Deprecated`, `Keywords`.

#### Method: `GetAbbreviatedPackageInformation`

```go
pkg, err := client.GetAbbreviatedPackageInformation(ctx, "react") (*models.Package, error)
```

Returns lightweight package metadata using `application/vnd.npm.install-v1+json` Accept header. Much smaller response than `GetPackageInformation` (KB vs MB). Useful when only version list and dist-tags are needed.

#### Method: `GetPackageVersion`

```go
version, err := client.GetPackageVersion(ctx, "react", "18.2.0") (*models.Version, error)
```

Returns: `Name`, `Version`, `Description`, `Dependencies`, `DevDependencies`, `Dist`, `Scripts`, `License`.

#### Method: `GetPackageVersions`

```go
versions, err := client.GetPackageVersions(ctx, "react") ([]string, error)
```

Returns a sorted list of all published version numbers. Uses the abbreviated API internally for efficiency.

#### Method: `GetPackageVersionCount`

```go
count, err := client.GetPackageVersionCount(ctx, "react") (int, error)
```

Returns the number of published versions. More lightweight than `GetPackageVersions` when only the count is needed.

#### Method: `GetPackageLatestVersion`

```go
latest, err := client.GetPackageLatestVersion(ctx, "react") (string, error)
```

Returns the version string of the `latest` dist-tag. Uses the lightweight dist-tags endpoint.

#### Method: `SearchPackages`

```go
result, err := client.SearchPackages(ctx, "react framework", 20) (*models.SearchResult, error)
```

Parameters: `query` string, `limit` int (default 20). Returns `SearchResult` with `Objects` array.

#### Method: `SearchPackagesWithOptions`

```go
result, err := client.SearchPackagesWithOptions(ctx, "http client", registry.SearchOptions{
    From:        20,
    Size:        10,
    Quality:     0.5,
    Popularity:  1.0,
    Maintenance: 0.0,
}) (*models.SearchResult, error)
```

Advanced search with pagination (`From`) and score weighting (`Quality`, `Popularity`, `Maintenance`). Weights are 0.0-1.0.

#### Method: `GetDistTags`

```go
tags, err := client.GetDistTags(ctx, "react") (map[string]string, error)
```

Returns all dist-tags (e.g. `{"latest": "18.2.0", "next": "19.0.0-rc.1"}`). Uses full package metadata internally.

#### Method: `GetDistTagsAbbreviated`

```go
tags, err := client.GetDistTagsAbbreviated(ctx, "react") (map[string]string, error)
```

Returns dist-tags using the lightweight `/-/package/{name}/dist-tags` endpoint. Faster than `GetDistTags`.

#### Method: `GetDownloadStats`

```go
stats, err := client.GetDownloadStats(ctx, "react", "last-week") (*models.DownloadStats, error)
```

Period options: `last-day`, `last-week`, `last-month`. Returns: `Downloads`, `Start`, `End`, `Package`.

#### Method: `GetDownloadRangeStats`

```go
stats, err := client.GetDownloadRangeStats(ctx, "react", "last-week") (*models.DownloadRangeStats, error)
```

Returns daily download counts for the period. Useful for trend visualization. Each entry has `Day` (date) and `Downloads` (count).

#### Method: `GetDownloadStatsByDateRange`

```go
stats, err := client.GetDownloadStatsByDateRange(ctx, "react", "2024-01-01", "2024-01-31") (*models.DownloadStats, error)
```

Custom date range query. Date format: `YYYY-MM-DD`.

#### Method: `GetBulkDownloadStats`

```go
stats, err := client.GetBulkDownloadStats(ctx, []string{"react", "vue", "angular"}, "last-week") (map[string]*models.DownloadStats, error)
```

Batch download stats for up to 128 packages in a single request.

#### Method: `GetBulkDownloadRangeStats`

```go
stats, err := client.GetBulkDownloadRangeStats(ctx, []string{"react", "vue"}, "last-week") (map[string]*models.DownloadRangeStats, error)
```

Batch daily download stats for up to 128 packages.

#### Method: `DownloadTarball`

```go
err := client.DownloadTarball(ctx, "react", "18.2.0", "./react.tgz") error
```

Downloads the .tgz file to the specified local path.

#### Method: `WhoAmI`

```go
username, err := client.WhoAmI(ctx) (string, error)
```

Checks authentication status using `/-/whoami` endpoint. Requires token to be set via `options.SetToken()`.

#### Method: `PublishPackage`

```go
err := client.PublishPackage(ctx, pkg) error
```

Publishes a full `Package` object to the registry. Requires token via `options.SetToken()`.

#### Method: `PublishPackageFromTarball`

```go
err := client.PublishPackageFromTarball(ctx, "my-pkg", "1.0.0", tarballBytes, metadata) error
```

Publishes a tarball with metadata. `metadata` is a `*PublishMetadata` with `Name`, `Version`, `Description` and other package.json fields. Requires token.

#### Method: `UnpublishPackage`

```go
err := client.UnpublishPackage(ctx, "my-pkg") error
```

Removes an entire package from the registry. Irreversible. Requires token.

#### Method: `UnpublishPackageVersion`

```go
err := client.UnpublishPackageVersion(ctx, "my-pkg", "1.0.0") error
```

Removes a specific version. Requires token.

#### Method: `DeprecateVersion`

```go
err := client.DeprecateVersion(ctx, "my-pkg", "1.0.0", "Use v2.0.0 instead") error
```

Marks a version as deprecated. Users see a warning when installing. Requires token.

#### Method: `SetDistTag`

```go
err := client.SetDistTag(ctx, "react", "next", "2.0.0-rc.1") error
```

Sets a dist-tag to point to a specific version. Requires token.

#### Method: `SetDistTags`

```go
err := client.SetDistTags(ctx, "react", map[string]string{"latest": "18.2.0", "next": "19.0.0"}) error
```

Sets multiple dist-tags at once. Requires token.

#### Method: `DeleteDistTag`

```go
err := client.DeleteDistTag(ctx, "react", "beta") error
```

Removes a dist-tag. Requires token.

#### Method: `GetPackageAccess`

```go
access, err := client.GetPackageAccess(ctx, "my-pkg") (*models.PackageAccess, error)
```

Returns package access settings. Requires token.

#### Method: `SetPackageAccess`

```go
err := client.SetPackageAccess(ctx, "my-pkg", &models.PackageAccessUpdate{Access: "public"}) error
```

Sets package visibility to `public` or `restricted`. Requires token.

#### Method: `ListCollaborators`

```go
collabs, err := client.ListCollaborators(ctx, "my-pkg") ([]models.Collaborator, error)
```

Returns collaborators with `Name`, `Email`, `Permissions`. Requires token.

#### Method: `GrantAccess`

```go
err := client.GrantAccess(ctx, "my-pkg", "username", models.PermissionWrite) error
```

Grants `read` or `write` access to a user or `org:team`. Requires token.

#### Method: `RevokeAccess`

```go
err := client.RevokeAccess(ctx, "my-pkg", "username") error
```

Removes a user's access to a package. Requires token.

#### Method: `StarPackage`

```go
err := client.StarPackage(ctx, "react") error
```

Adds a package to your starred list. Requires token.

#### Method: `UnstarPackage`

```go
err := client.UnstarPackage(ctx, "react") error
```

Removes a package from your starred list. Requires token.

#### Method: `GetStarredByUser`

```go
packages, err := client.GetStarredByUser(ctx, "username") ([]string, error)
```

Returns package names starred by the user. No authentication required.

#### Method: `GetStarredByPackage`

```go
users, err := client.GetStarredByPackage(ctx, "react") ([]string, error)
```

Returns usernames who starred the package. No authentication required.

#### Method: `ListTokens`

```go
tokens, err := client.ListTokens(ctx) ([]models.Token, error)
```

Returns all access tokens. Requires token.

#### Method: `GetToken`

```go
token, err := client.GetToken(ctx, "token-id") (*models.Token, error)
```

Returns a specific token's details. Requires token.

#### Method: `CreateToken`

```go
token, err := client.CreateToken(ctx, &models.TokenCreation{
    Password: "mypass",
    Readonly: true,
    CIDR:     []string{"192.168.1.0/24"},
}) (*models.Token, error)
```

Creates a new access token. Requires token + current password.

#### Method: `DeleteToken`

```go
err := client.DeleteToken(ctx, "token-id") error
```

Revokes an access token. Requires token.

#### Method: `Login`

```go
result, err := client.Login(ctx, "username", "password") (*models.LoginResult, error)
```

Authenticates and returns a `LoginResult` with `ID`, `Rev`, `Token`, `Ok`. No prior token required.

#### Method: `CreateUser`

```go
result, err := client.CreateUser(ctx, &models.UserCreation{
    Name: "username", Password: "pass", Email: "me@example.com",
}) (*models.LoginResult, error)
```

Registers a new user. Returns `LoginResult` with token. No prior token required.

#### Method: `GetUser`

```go
profile, err := client.GetUser(ctx, "username") (*models.UserProfile, error)
```

Returns `UserProfile` with `Name`, `Email`. Requires token.

#### Method: `QuickAudit`

```go
result, err := client.QuickAudit(ctx, &models.QuickAuditRequest{
    Dependencies: map[string]string{"lodash": "4.17.11"},
}) (*models.QuickAuditResult, error)
```

Returns vulnerability counts by severity. No authentication required.

#### Method: `BulkAudit`

```go
result, err := client.BulkAudit(ctx, map[string][]string{
    "lodash": {"<4.17.12"}, "express": {"<4.17.3"},
}) (map[string][]models.Advisory, error)
```

Returns detailed advisories per package. No authentication required.

#### Method: `GetAdvisory`

```go
advisory, err := client.GetAdvisory(ctx, 1234) (*models.Advisory, error)
```

Returns a single advisory with full details. No authentication required.

#### Method: `ListAdvisories`

```go
advisories, err := client.ListAdvisories(ctx, models.AdvisoryListOptions{
    Page: 0, PerPage: 20, AffectedPackage: "lodash",
}) ([]models.Advisory, error)
```

Returns paginated advisory list. No authentication required.

#### Method: `GetOrg`

```go
org, err := client.GetOrg(ctx, "myorg") (*models.Organization, error)
```

Returns org `Name` and `Scope`. Requires token.

#### Method: `CreateOrg`

```go
org, err := client.CreateOrg(ctx, "myorg") (*models.Organization, error)
```

Creates a new organization. Requires token.

#### Method: `DeleteOrg`

```go
err := client.DeleteOrg(ctx, "myorg") error
```

Deletes an organization. Irreversible. Requires token.

#### Method: `ListOrgMembers`

```go
members, err := client.ListOrgMembers(ctx, "myorg") ([]string, error)
```

Returns member usernames. Requires token.

#### Method: `AddOrgMember`

```go
err := client.AddOrgMember(ctx, "myorg", "username") error
```

Adds a user to the organization. Requires token.

#### Method: `RemoveOrgMember`

```go
err := client.RemoveOrgMember(ctx, "myorg", "username") error
```

Removes a user from the organization. Requires token.

#### Method: `ListOrgPackages`

```go
packages, err := client.ListOrgPackages(ctx, "myorg") ([]string, error)
```

Returns package names owned by the org. Requires token.

#### Method: `ListTeams`

```go
teams, err := client.ListTeams(ctx, "myorg") ([]models.Team, error)
```

Returns teams with `ID`, `Name`, `DisplayName`. Requires token.

#### Method: `CreateTeam`

```go
team, err := client.CreateTeam(ctx, "myorg", "devteam") (*models.Team, error)
```

Creates a team within an org. Requires token.

#### Method: `DeleteTeam`

```go
err := client.DeleteTeam(ctx, "myorg", "devteam") error
```

Deletes a team. Irreversible. Requires token.

#### Method: `ListTeamMembers`

```go
members, err := client.ListTeamMembers(ctx, "myorg", "devteam") ([]string, error)
```

Returns team member usernames. Requires token.

#### Method: `AddTeamMember`

```go
err := client.AddTeamMember(ctx, "myorg", "devteam", "username") error
```

Adds a user to a team. User must already be an org member. Requires token.

#### Method: `RemoveTeamMember`

```go
err := client.RemoveTeamMember(ctx, "myorg", "devteam", "username") error
```

Removes a user from a team. Requires token.

#### Method: `ListTeamPackages`

```go
packages, err := client.ListTeamPackages(ctx, "myorg", "devteam") ([]string, error)
```

Returns packages the team has access to. Requires token.

#### Method: `ListHooks`

```go
hooks, err := client.ListHooks(ctx, models.HookListOptions{
    Package: "react", Page: 0, PerPage: 20,
}) ([]models.Hook, error)
```

Returns webhook list. Requires token.

#### Method: `GetHook`

```go
hook, err := client.GetHook(ctx, "hook-id") (*models.Hook, error)
```

Returns hook details. Requires token.

#### Method: `CreateHook`

```go
hook, err := client.CreateHook(ctx, &models.HookCreation{
    Name: "my-hook", Endpoint: "https://example.com/webhook",
    Package: "react", Active: true,
}) (*models.Hook, error)
```

Creates a webhook. Requires token.

#### Method: `UpdateHook`

```go
hook, err := client.UpdateHook(ctx, "hook-id", &models.HookUpdate{
    Endpoint: "https://new.example.com/webhook",
}) (*models.Hook, error)
```

Updates specified fields only. Requires token.

#### Method: `DeleteHook`

```go
err := client.DeleteHook(ctx, "hook-id") error
```

Permanently removes a webhook. Requires token.

#### Method: `GetChanges`

```go
result, err := client.GetChanges(ctx, models.ChangesOptions{
    Since: "12345", Limit: 50, IncludeDocs: true,
}) (*models.ChangesResult, error)
```

Returns CouchDB changes feed. No authentication required.

#### Method: `GetAllDocs`

```go
result, err := client.GetAllDocs(ctx, models.AllDocsOptions{
    StartKey: "@nestjs", EndKey: "@nestt", Limit: 50,
}) (*models.AllDocsResult, error)
```

Returns all document IDs. No authentication required.

#### Method: `GetView`

```go
result, err := client.GetView(ctx, "starredByUser", models.ViewOptions{
    Key: '"username"', Limit: 20,
}) (*models.ViewResult, error)
```

Queries a CouchDB view. No authentication required.

#### Method: `GetOptions`

```go
opts := client.GetOptions() *Options
```

Returns current configuration options.

#### Utility Functions

```go
mirrors := registry.ListMirrors() []MirrorEntry
```

Returns all available mirror sources with `Name`, `URL`, `Region`, `Description` fields.

### Package: `github.com/scagogogo/npm-skills/pkg/models`

#### Key Types

| Type | Key Fields | Description |
|------|-----------|-------------|
| `Package` | Name, Description, DistTags, Versions, Maintainers, ReadMe, License | Full package metadata |
| `Version` | Version, Dependencies, DevDependencies, Dist, Scripts, License | Version-specific info |
| `Dist` | Tarball, Shasum, Integrity | Distribution/download info |
| `SearchResult` | Objects, Total | Search results container |
| `SearchOptions` | From, Size, Quality, Popularity, Maintenance | Search parameters |
| `DownloadStats` | Downloads, Start, End, Package | Download statistics |
| `DownloadRangeStats` | Downloads (daily array), Start, End, Package | Daily download trends |
| `DailyDownloads` | Day, Downloads | Single day download data |
| `RegistryInformation` | DbName, DocCount, DiskSize, DataSize | Registry status |
| `Maintainer` | Name, Email | Package maintainer |
| `Repository` | Type, Url | Repository info |
| `MirrorEntry` | Name, URL, Region, Description | Mirror source metadata |
| `PublishMetadata` | Name, Version, Description, Dependencies, Keywords, License | Package publish metadata |
| `Token` | ID, Token, Key, Created, Updated, Readonly, CIDR | Access token |
| `TokenCreation` | Password, Readonly, CIDR | Token creation request |
| `Hook` | ID, Name, Endpoint, Secret, Events, Package, Active | Webhook |
| `HookCreation` | Name, Endpoint, Secret, Events, Package, Active | Webhook creation request |
| `HookUpdate` | Endpoint, Secret, Events, Active | Webhook update request |
| `HookListOptions` | Package, Page, PerPage | Hook list query params |
| `Organization` | Name, Scope | NPM organization |
| `Team` | ID, Name, DisplayName, Description | NPM team |
| `PackageAccess` | Package, Access | Package access settings |
| `PackageAccessUpdate` | Access | Access update (public/restricted) |
| `Collaborator` | Name, Email, Permissions | Package collaborator |
| `Permission` | (const: read, write) | Permission level |
| `Advisory` | ID, Title, Severity, CVE, ModuleName, Vulnerable, Patched | Security advisory |
| `AdvisoryListOptions` | Page, PerPage, AffectedPackage | Advisory list query params |
| `QuickAuditRequest` | Dependencies | Audit request (name->version map) |
| `QuickAuditResult` | Metadata (vulnerabilities counts) | Audit result |
| `LoginResult` | ID, Rev, Token, Ok | Login/signup response |
| `UserCreation` | Name, Password, Email | User registration request |
| `UserProfile` | ID, Rev, Name, Email, Type | User profile info |
| `ChangesOptions` | Since, Limit, IncludeDocs | CouchDB changes query params |
| `ChangesResult` | LastSeq, Pending, Results | CouchDB changes feed result |
| `AllDocsOptions` | StartKey, EndKey, Limit, Skip, IncludeDocs, Descending | CouchDB all-docs query params |
| `AllDocsResult` | TotalRows, Offset, Rows | CouchDB all-docs result |
| `ViewOptions` | Key, StartKey, EndKey, Limit, Skip, Group, GroupLevel, Descending | CouchDB view query params |
| `ViewResult` | TotalRows, Offset, Rows | CouchDB view result |

All model types implement `ToJsonString() string` for easy debugging.

---

## MCP Server Reference

The project includes an MCP (Model Context Protocol) server at `cmd/mcp-server/` that exposes NPM registry operations as tools for AI agents.

### Package: `github.com/scagogogo/npm-skills/pkg/mcp`

#### Creating an MCP Server

```go
import npmMcp "github.com/scagogogo/npm-skills/pkg/mcp"
import "github.com/scagogogo/npm-skills/pkg/registry"

cfg := npmMcp.Config{
    RegistryOptions: registry.NewOptions().SetRegistryURL("https://registry.npmjs.org"),
    Timeout:         120 * time.Second,
}
mcpServer := npmMcp.NewServer(cfg)
```

### MCP Tools (12 total)

| Tool Name | SDK Method | Required Params | Optional Params | Description |
|-----------|-----------|-----------------|-----------------|-------------|
| `npm_registry_info` | `GetRegistryInformation` | (none) | (none) | Registry status and statistics |
| `npm_mirrors` | `ListMirrors` | (none) | (none) | Available mirror sources with URLs, regions, descriptions |
| `npm_package` | `GetPackageInformation` | `name` | (none) | Full package metadata (WARNING: can be 10MB+) |
| `npm_package_summary` | `GetAbbreviatedPackageInformation` | `name` | (none) | Lightweight package metadata (recommended for most queries) |
| `npm_search` | `SearchPackagesWithOptions` | `query` | `limit`, `from`, `quality`, `popularity`, `maintenance` | Search packages with pagination and score weighting |
| `npm_version` | `GetPackageVersion` | `name`, `version` | (none) | Specific version metadata |
| `npm_versions` | `GetPackageVersions` | `name` | (none) | All published version numbers (sorted) |
| `npm_latest_version` | `GetPackageLatestVersion` | `name` | (none) | Latest version number only |
| `npm_dist_tags` | `GetDistTagsAbbreviated` | `name` | (none) | Distribution tags (latest, next, beta) |
| `npm_download_stats` | `GetDownloadStats` | `name` | `period` | Download count for a period |
| `npm_download_range` | `GetDownloadRangeStats` | `name` | `period` | Daily download trend data |
| `npm_whoami` | `WhoAmI` | (none) | (none) | Check auth status (requires server --token) |

### Response Format

All MCP tools return JSON-formatted text content. Responses exceeding 50KB are truncated:
- Package README is truncated to 2000 characters
- Package Versions map is replaced with a version_keys array
- Hard limit: responses exceeding 100KB are truncated with a trailing notice

### Running the MCP Server

```bash
# Build
go build -o npm-mcp-server ./cmd/mcp-server/

# Run with defaults (stdio transport)
npm-mcp-server

# Run with China mirror
npm-mcp-server --mirror npm-mirror

# Run with proxy and token
npm-mcp-server --proxy http://127.0.0.1:7890 --token npm_xxxxx

# Using environment variables
NPM_MIRROR=npm-mirror npm-mcp-server
```

### Claude Code Integration

Add to your settings (`~/.claude/settings.json`):

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
