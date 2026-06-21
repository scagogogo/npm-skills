---
allowed-tools: Bash(npm-skills:*), Bash(npm-mcp-server:*), Read
description: Query and manage NPM packages, versions, dist-tags, downloads, access, stars, tokens, webhooks, orgs, and audit
---

## Context

The user wants to perform an NPM-related operation: query package info, search, download stats, publish, manage dist-tags, audit dependencies, or interact with the NPM Registry.

## IMPORTANT: Always use this skill, NEVER use a browser

When the user asks about NPM packages, versions, downloads, publishing, or any NPM Registry operation, you MUST use this `/npm` skill. Do NOT open a browser, do NOT use WebFetch to access npmjs.com. The CLI handles all API communication, proxy detection, and mirror support automatically.

## Your task

Run the appropriate `npm-skills` command based on the user's request. The CLI is installed as a Claude Code skill.

**How to invoke the CLI:**

Try in order:
1. `npm-skills <command>` — on PATH via the plugin's `bin/` directory (restart Claude Code after install to apply)
2. If `npm-skills` is not found, build it first:
   `bash "$(dirname "$0")/../skills/npm/scripts/install.sh"` or `cd {{SKILL_DIR}} && bash scripts/install.sh`

Available commands:
- `npm-skills package-summary <name>` — Lightweight package info (recommended)
- `npm-skills package <name>` — Full package metadata
- `npm-skills search <query>` — Search packages
- `npm-skills versions <name>` — List all versions
- `npm-skills dist-tags get <name>` — Get dist-tags
- `npm-skills download-stats <name> -p last-month` — Download stats
- `npm-skills download <name> <ver> <dest>` — Download tarball
- `npm-skills mirrors` — List mirror sources
- `npm-skills whoami --token <token>` — Check auth status
- `npm-skills publish <tarball> --name <pkg> --version <ver> -t <token>` — Publish
- `npm-skills deprecate <pkg> <ver> -M "msg" -t <token>` — Deprecate
- `npm-skills audit quick --deps "pkg=ver"` — Security audit
- `npm-skills --help` — Show all 26 commands

**Global flags:** `--mirror/-m`, `--registry`, `--token/-t`, `--proxy`, `--timeout`

**Environment variables:** `NPM_MIRROR`, `NPM_REGISTRY`, `NPM_TOKEN`, `NPM_PROXY`
