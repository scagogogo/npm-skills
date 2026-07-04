# Installation

Five installation methods exist — pick by use case rather than trying all. Use the decision diagram below to find the one that fits:

```mermaid
flowchart TD
    Start(["Want NPM Skills"]) --> Q1{"Where will you<br/>mainly use it?"}
    Q1 -->|"Claude Code / AI agent"| M1["Option 1: Plugin<br/>claude plugin install"]
    Q1 -->|"CLI / CI scripts"| Q2{"Need cross-platform<br/>prebuilt binaries?"}
    Q1 -->|"Inside a Go program"| Q3{"Just the SDK library,<br/>or CLI too?"}
    Q2 -->|"Yes (fastest start)"| M2["Option 2: Prebuilt binary<br/>34 platforms"]
    Q2 -->|"No, Go is available"| M4["Option 4: go install"]
    Q3 -->|"Library only"| M5["Option 5: go get"]
    Q3 -->|"Library + CLI"| M3["Option 3: Build from source<br/>bash scripts/install.sh"]

    M1 --> V["Verify: npm-skills --version"]
    M2 --> V
    M3 --> V
    M4 --> V
    M5 --> V2["No CLI to verify<br/>just import"]

    classDef pick fill:#e8f0fe,stroke:#4285f4,color:#174ea6;
    classDef verify fill:#e6f4ea,stroke:#34a853,color:#1e4620;
    class M1,M2,M3,M4,M5 pick;
    class V,V2 verify;
```

## Option 1: Claude Code Plugin (recommended for AI agents)

```bash
claude plugin marketplace add scagogogo/npm-skills
claude plugin install npm@npm-skills
```

## Option 2: Pre-built Binary (recommended for CLI)

![Release Pipeline](/release-pipeline.svg)

Download the archive for your platform from [GitHub Releases](https://github.com/scagogogo/npm-skills/releases/latest):

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

### Supported Platforms (34 combinations)

| OS | Architectures |
|----|---------------|
| Linux | amd64, arm64, 386, arm, loong64, mips, mips64, mips64le, mipsle, ppc64, ppc64le, riscv64, s390x |
| macOS | amd64, arm64 |
| Windows | amd64, 386 |
| FreeBSD | amd64, arm64, 386, arm, mips, mipsle, ppc64, ppc64le |
| OpenBSD | amd64, arm64, 386, arm, mips, mipsle, ppc64, ppc64le |
| NetBSD | amd64, arm64, 386, arm, mips, mipsle, ppc64, ppc64le |
| Illumos | amd64 |
| Solaris | amd64 |

## Option 3: Build from Source

```bash
git clone https://github.com/scagogogo/npm-skills.git
cd npm-skills
bash scripts/install.sh   # builds to ~/.local/bin/
```

## Option 4: go install

```bash
go install github.com/scagogogo/npm-skills/cmd/npm-skills@latest
go install github.com/scagogogo/npm-skills/cmd/mcp-server@latest
```

## Option 5: Go Module

```bash
go get github.com/scagogogo/npm-skills
```

## Verify

```bash
npm-skills --version
npm-skills mirrors
```

## Next Steps

- Read [Getting Started](/en/getting-started) for the four entry points
- Check the [CLI Reference](/en/cli) for all 26 commands
- Browse [API docs](/en/api/) for the underlying SDK methods
