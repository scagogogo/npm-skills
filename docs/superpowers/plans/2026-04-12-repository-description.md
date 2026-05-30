# Improve Repository Description

> **For agentic workers:** REQUIRED SUB-SKILL: `superpowers:subagent-driven-development`
> Steps use checkbox (`- [ ]`) syntax.

**Goal:** Update the GitHub repository description from Chinese "npm仓库爬虫" to a descriptive English version that accurately reflects the library's functionality.

**Architecture:** N/A (metadata-only task using gh CLI)

**Tech Stack:** gh CLI, GitHub REST API

**Risks:**
- Task 1 is a read-only metadata update → no risk of breaking anything

---

### Task 1: Update GitHub Repository Description

**Depends on:** None
**Files:**
- None (remote metadata update via gh CLI)

- [ ] **Step 1: Verify current repository description**
Run: `gh repo view scagogogo/npm-crawler --json description`
Expected:
  - Output shows current description: "npm仓库爬虫"

- [ ] **Step 2: Update repository description to English**
Run: `gh repo edit scagogogo/npm-crawler --description "High-performance Go client library for NPM Registry with multi-mirror support and proxy configuration"`
Expected:
  - Exit code: 0
  - No output (silent success)

- [ ] **Step 3: Verify the update**
Run: `gh repo view scagogogo/npm-crawler --json description`
Expected:
  - Output shows new description: "High-performance Go client library for NPM Registry with multi-mirror support and proxy configuration"

- [ ] **Step 4: Review and optionally update repository topics**
Run: `gh repo view scagogogo/npm-crawler --json repositoryTopics`
Expected:
  - Current topics: crawler, npm, sca

  If topics need adjustment (optional), run:
  `gh repo edit scagogogo/npm-crawler --add-topic go --remove-topic sca`
  (Note: only run if user confirms topic changes are desired)
