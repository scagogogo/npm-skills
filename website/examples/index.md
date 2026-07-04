# 示例

NPM Skills Go SDK 的实用模式。

```mermaid
flowchart LR
    Start(["想做什么？"]) --> Q1{"查包 / 版本<br/>信息？"}
    Q1 -->|是| B["基础用法"]
    Q1 -->|否| Q2{"下载 .tgz<br/>到磁盘？"}
    Q2 -->|是| D["下载 Tarball"]
    Q2 -->|否| Q3{"国内访问慢<br/>/ 需代理？"}
    Q3 -->|是| M["镜像源配置"]
    Q3 -->|否| Q4{"认证 / 限流 /<br/>高级错误处理？"}
    Q4 -->|是| A["高级用法"]
    Q4 -->|否| B

    classDef pick fill:#e8f0fe,stroke:#4285f4,color:#174ea6;
    class B,D,M,A pick;
```

- [基础用法](/examples/basic) — 创建客户端、查询包、搜索
- [高级用法](/examples/advanced) — 自定义仓库、认证、超时、错误处理
- [镜像源配置](/examples/mirrors) — 国内镜像与代理
- [下载 Tarball](/examples/download) — 下载 .tgz 文件到磁盘
