# Getting Started

## Install as Claude Code Plugin in 1 minute

```bash
claude plugin marketplace add scagogogo/npm-skills
claude plugin install npm@npm-skills
```

Then ask in natural language — the AI agent auto-invokes NPM Skills:

> *"Find info about the axios package"*
> *"Download the react tarball"*
> *"Get vue download stats for last month via the China mirror"*

The full call chain from natural language to the NPM Registry:

```mermaid
sequenceDiagram
    autonumber
    actor U as You
    participant AI as AI Agent<br/>Claude Code
    participant S as SKILL.md<br/>progressive disclosure
    participant CLI as npm-skills CLI
    participant R as Registry SDK
    participant N as NPM Registry / mirror

    U->>AI: "Get vue download stats for last month via China mirror"
    AI->>S: match skill description, load instructions
    S-->>AI: return command usage (expanded on demand)
    AI->>CLI: npm-skills download-stats vue -p last-month -m npm-mirror
    CLI->>R: build request (mirror/proxy/timeout)
    R->>N: HTTP GET /downloads/point/last-month/vue
    N-->>R: JSON response
    R-->>CLI: structured data
    CLI-->>AI: JSON (stdout)
    AI-->>U: summarize result in natural language
```

## Four Ways to Integrate

All four entry points share one Registry SDK core and converge on the NPM Registry or a mirror:

```mermaid
flowchart TB
    subgraph Users
        A1["AI Agents<br/>Claude / Cursor / Windsurf"]
        A2["Shell / CI scripts"]
        A3["Go applications"]
    end

    subgraph Layer["NPM Skills integration layer"]
        B1["Skill / Plugin<br/>SKILL.md"]
        B2["CLI<br/>npm-skills"]
        B3["MCP Server<br/>npm-mcp-server"]
        B4["Go SDK<br/>pkg/registry"]
    end

    C["Registry client core<br/>Options · HTTP pool · typed errors"]
    D["NPM Registry<br/>official / 8 mirrors / private"]

    A1 --> B1
    A1 --> B3
    A2 --> B2
    A3 --> B4
    B1 --> B2
    B2 --> C
    B3 --> C
    B4 --> C
    C -->|HTTP/HTTPS · proxy| D

    classDef core fill:#cb3837,stroke:#8b1a1a,color:#fff;
    class C core;
```

| Way | Use case | Entry |
|-----|----------|-------|
| **Skill / Plugin** | AI agents auto-invoke | `claude plugin install npm@npm-skills` |
| **CLI Tool** | Shell / scripts | `npm-skills <command>` |
| **Go SDK** | Go programs | `import "github.com/scagogogo/npm-skills/pkg/registry"` |
| **MCP Server** | MCP-compatible clients | `npm-mcp-server` |

## CLI Cheat Sheet (90% of cases)

```bash
npm-skills package-summary react            # Lightweight info (recommended)
npm-skills search "http client" -l 10
npm-skills versions react --latest
npm-skills dist-tags get react
npm-skills download-stats react -p last-month
npm-skills mirrors
npm-skills package react -m npm-mirror      # China mirror
```

## Mirrors & Proxy

```bash
npm-skills package react -m npm-mirror                       # China mirror
npm-skills package react --proxy http://127.0.0.1:7890       # HTTP proxy
npm-skills package my-lib --registry https://npm.co.com -t npm_x  # Private

export NPM_MIRROR=npm-mirror
export NPM_PROXY=http://127.0.0.1:7890
export NPM_TOKEN=npm_xxxxx
```

When the same setting is specified in multiple places, precedence is "CLI flag > env var > default":

```mermaid
flowchart LR
    Start(["Resolve one setting<br/>e.g. mirror / proxy / token"]) --> Q1{"CLI passed<br/>--mirror ?"}
    Q1 -->|yes| U1["use CLI flag"]
    Q1 -->|no| Q2{"env var<br/>NPM_MIRROR ?"}
    Q2 -->|yes| U2["use env var"]
    Q2 -->|no| U3["use built-in default<br/>official"]

    U1 --> Done(["effective config"])
    U2 --> Done
    U3 --> Done

    classDef win fill:#e6f4ea,stroke:#34a853,color:#1e4620;
    class U1,U2,U3 win;
```

## Next Steps

- [CLI Reference](/en/cli) — All 26 commands
- [Go SDK](/en/api/registry) — Programmatic access
- [MCP Server](/en/mcp-server) — AI tool chain integration
