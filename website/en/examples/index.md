# Examples

Practical usage patterns for the NPM Skills Go SDK.

```mermaid
flowchart LR
    Start(["What do you want?"]) --> Q1{"Query package /<br/>version info?"}
    Q1 -->|yes| B["Basic Usage"]
    Q1 -->|no| Q2{"Download .tgz<br/>to disk?"}
    Q2 -->|yes| D["Download Tarball"]
    Q2 -->|no| Q3{"Slow in China /<br/>need a proxy?"}
    Q3 -->|yes| M["Mirror Configuration"]
    Q3 -->|no| Q4{"Auth / rate limit /<br/>advanced errors?"}
    Q4 -->|yes| A["Advanced Usage"]
    Q4 -->|no| B

    classDef pick fill:#e8f0fe,stroke:#4285f4,color:#174ea6;
    class B,D,M,A pick;
```

- [Basic Usage](/en/examples/basic) — Create a client, query packages, search
- [Advanced Usage](/en/examples/advanced) — Custom registries, auth, timeout, error handling
- [Mirror Configuration](/en/examples/mirrors) — China mirrors and proxy setup
- [Download Tarball](/en/examples/download) — Download .tgz files to disk
