# 数据模型

NPM Skills 定义了完整的 Go 数据结构来表示 NPM 包的各种元数据信息。

## Package 模型

表示 NPM 包的完整信息：

```go
type Package struct {
    ID             string                 `json:"_id"`             // 包 ID
    Name           string                 `json:"name"`            // 包名称
    Description    string                 `json:"description"`     // 包描述
    DistTags       map[string]string      `json:"dist-tags"`       // 分发标签
    Versions       map[string]Version     `json:"versions"`        // 版本信息
    Maintainers    []Maintainer           `json:"maintainers"`     // 维护者
    Time           map[string]string      `json:"time"`            // 时间信息
    Repository     Repository             `json:"repository"`      // 仓库信息
    Homepage       string                 `json:"homepage"`        // 主页
    License        string                 `json:"license"`         // 许可证
    Keywords       []string               `json:"keywords"`        // 关键词
    Author         Author                 `json:"author"`          // 作者
    // ... 其他字段
}
```

### 常用操作

```go
// 获取最新版本
latestVersion := pkg.DistTags["latest"]

// 列出所有可用版本
for version := range pkg.Versions {
    fmt.Printf("可用版本: %s\n", version)
}

// 获取特定版本详情
if versionInfo, exists := pkg.Versions["1.0.0"]; exists {
    fmt.Printf("依赖: %+v\n", versionInfo.Dependencies)
    fmt.Printf("开发依赖: %+v\n", versionInfo.DevDependencies)
}

// 访问作者信息
fmt.Printf("作者: %s <%s>\n", pkg.Author.Name, pkg.Author.Email)

// 访问仓库信息
if pkg.Repository.Type == "git" {
    fmt.Printf("Git 仓库: %s\n", pkg.Repository.URL)
}
```

## Version 模型

表示包的特定版本信息：

```go
type Version struct {
    Name            string               `json:"name"`            // 包名称
    Version         string               `json:"version"`         // 版本号
    Description     string               `json:"description"`     // 描述
    Main            string               `json:"main"`            // 入口点
    Scripts         *Script              `json:"scripts"`         // NPM 脚本
    Dependencies    map[string]string    `json:"dependencies"`    // 运行时依赖
    DevDependencies map[string]string    `json:"devDependencies"` // 开发依赖
    Repository      *Repository          `json:"repository"`      // 仓库
    License         string               `json:"license"`         // 许可证
    Dist            *Dist                `json:"dist"`            // 分发信息
    // ... 其他字段
}
```

### 使用示例

```go
// 检查版本依赖
if len(version.Dependencies) > 0 {
    fmt.Println("运行时依赖:")
    for dep, ver := range version.Dependencies {
        fmt.Printf("  %s: %s\n", dep, ver)
    }
}

// 检查开发依赖
if len(version.DevDependencies) > 0 {
    fmt.Println("开发依赖:")
    for dep, ver := range version.DevDependencies {
        fmt.Printf("  %s: %s\n", dep, ver)
    }
}

// 访问分发信息
if version.Dist != nil {
    fmt.Printf("包大小: %d 字节\n", version.Dist.UnpackedSize)
    fmt.Printf("Tarball: %s\n", version.Dist.Tarball)
    fmt.Printf("SHA-1: %s\n", version.Dist.Shasum)
}
```

## Author 模型

表示作者信息：

```go
type Author struct {
    Name  string `json:"name"`  // 作者姓名
    Email string `json:"email"` // 邮箱地址
    URL   string `json:"url"`   // 个人网站
}
```

## Maintainer 模型

表示维护者信息：

```go
type Maintainer struct {
    Name  string `json:"name"`  // 维护者姓名
    Email string `json:"email"` // 邮箱地址
}
```

## Repository 模型

表示代码仓库信息：

```go
type Repository struct {
    Type string `json:"type"` // 仓库类型 (git, svn 等)
    URL  string `json:"url"`  // 仓库 URL
}
```

### 使用示例

```go
// 检查是否为 Git 仓库
if pkg.Repository.Type == "git" {
    fmt.Printf("Git 仓库地址: %s\n", pkg.Repository.URL)
    
    // 从 GitHub URL 提取信息
    if strings.Contains(pkg.Repository.URL, "github.com") {
        fmt.Println("这是一个 GitHub 项目")
    }
}
```

## Dist 模型

表示包的分发信息：

```go
type Dist struct {
    Integrity    string `json:"integrity"`    // 完整性校验
    Shasum       string `json:"shasum"`       // SHA-1 校验和
    Tarball      string `json:"tarball"`      // 包下载地址
    FileCount    int    `json:"fileCount"`    // 文件数量
    UnpackedSize int64  `json:"unpackedSize"` // 解压后大小
}
```

### 使用示例

```go
if version.Dist != nil {
    fmt.Printf("包大小: %.2f KB\n", float64(version.Dist.UnpackedSize)/1024)
    fmt.Printf("文件数量: %d\n", version.Dist.FileCount)
    fmt.Printf("下载地址: %s\n", version.Dist.Tarball)
    
    // 验证校验和
    fmt.Printf("SHA-1: %s\n", version.Dist.Shasum)
    if version.Dist.Integrity != "" {
        fmt.Printf("完整性: %s\n", version.Dist.Integrity)
    }
}
```

## SearchResult 模型

表示搜索结果：

```go
type SearchResult struct {
    Objects []SearchObject `json:"objects"` // 搜索结果对象
    Total   int            `json:"total"`   // 总匹配数量
    Time    string         `json:"time"`    // 搜索耗时
}

type SearchObject struct {
    Package     SearchPackage `json:"package"`     // 包信息
    Score       Score         `json:"score"`       // 评分
    SearchScore float64       `json:"searchScore"` // 搜索得分
}

type SearchPackage struct {
    Name        string            `json:"name"`        // 包名称
    Version     string            `json:"version"`     // 版本
    Description string            `json:"description"` // 描述
    Keywords    []string          `json:"keywords"`    // 关键词
    Author      Author            `json:"author"`      // 作者
    Maintainers []Maintainer      `json:"maintainers"` // 维护者
    Links       map[string]string `json:"links"`       // 相关链接
}

type Score struct {
    Final   float64 `json:"final"`   // 最终得分
    Detail  Detail  `json:"detail"`  // 详细得分
}

type Detail struct {
    Quality     float64 `json:"quality"`     // 质量得分
    Popularity  float64 `json:"popularity"`  // 流行度得分
    Maintenance float64 `json:"maintenance"` // 维护得分
}
```

### 使用示例

```go
result, err := client.SearchPackages(ctx, "react", 10)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("搜索耗时: %s\n", result.Time)
fmt.Printf("找到 %d 个结果\n", result.Total)

for i, obj := range result.Objects {
    pkg := obj.Package
    score := obj.Score
    
    fmt.Printf("%d. %s@%s\n", i+1, pkg.Name, pkg.Version)
    fmt.Printf("   描述: %s\n", pkg.Description)
    fmt.Printf("   得分: %.2f (质量: %.2f, 流行度: %.2f, 维护: %.2f)\n", 
        score.Final, score.Detail.Quality, score.Detail.Popularity, score.Detail.Maintenance)
    
    if len(pkg.Keywords) > 0 {
        fmt.Printf("   关键词: %s\n", strings.Join(pkg.Keywords, ", "))
    }
    
    fmt.Println()
}
```

## DownloadStats 模型

表示下载统计信息：

```go
type DownloadStats struct {
    Downloads int    `json:"downloads"` // 下载次数
    Start     string `json:"start"`     // 统计开始日期
    End       string `json:"end"`       // 统计结束日期
    Package   string `json:"package"`   // 包名称
}
```

### 使用示例

```go
stats, err := client.GetDownloadStats(ctx, "react", "last-week")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("包: %s\n", stats.Package)
fmt.Printf("下载次数: %s\n", formatNumber(stats.Downloads))
fmt.Printf("统计周期: %s 到 %s\n", stats.Start, stats.End)

// 格式化数字显示
func formatNumber(n int) string {
    if n >= 1000000 {
        return fmt.Sprintf("%.1fM", float64(n)/1000000)
    }
    if n >= 1000 {
        return fmt.Sprintf("%.1fK", float64(n)/1000)
    }
    return fmt.Sprintf("%d", n)
}
```

## RegistryInformation 模型

表示注册表信息：

```go
type RegistryInformation struct {
    DbName            string `json:"db_name"`              // 数据库名称
    DocCount          int    `json:"doc_count"`            // 包总数
    DocDelCount       int    `json:"doc_del_count"`        // 已删除包数
    UpdateSeq         int    `json:"update_seq"`           // 更新序列
    CompactRunning    bool   `json:"compact_running"`      // 压缩状态
    DiskSize          int64  `json:"disk_size"`            // 磁盘使用
    DataSize          int64  `json:"data_size"`            // 数据大小
    InstanceStartTime string `json:"instance_start_time"`  // 启动时间
    // ... 其他字段
}
```

### 使用示例

```go
info, err := client.GetRegistryInformation(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("注册表: %s\n", info.DbName)
fmt.Printf("包总数: %s\n", formatNumber(info.DocCount))
fmt.Printf("已删除包: %s\n", formatNumber(info.DocDelCount))
fmt.Printf("活跃包: %s\n", formatNumber(info.DocCount-info.DocDelCount))

fmt.Printf("磁盘使用: %s\n", formatBytes(info.DiskSize))
fmt.Printf("数据大小: %s\n", formatBytes(info.DataSize))

if info.CompactRunning {
    fmt.Println("状态: 正在压缩")
} else {
    fmt.Println("状态: 正常")
}

// 格式化字节大小
func formatBytes(bytes int64) string {
    const unit = 1024
    if bytes < unit {
        return fmt.Sprintf("%d B", bytes)
    }
    div, exp := int64(unit), 0
    for n := bytes / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
```

## Script 模型

表示 NPM 脚本：

```go
type Script struct {
    Test    string `json:"test"`    // 测试脚本
    Build   string `json:"build"`   // 构建脚本
    Start   string `json:"start"`   // 启动脚本
    Dev     string `json:"dev"`     // 开发脚本
    Lint    string `json:"lint"`    // 代码检查脚本
    // ... 其他自定义脚本
}
```

## 数据验证

### 包名验证

```go
func isValidPackageName(name string) bool {
    // NPM 包名规则
    if len(name) == 0 || len(name) > 214 {
        return false
    }
    
    // 不能包含大写字母
    if strings.ToLower(name) != name {
        return false
    }
    
    // 不能以点或下划线开头
    if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") {
        return false
    }
    
    // 只能包含 URL 安全字符
    matched, _ := regexp.MatchString(`^[a-z0-9._~-]+$`, name)
    return matched
}
```

### 版本号验证

```go
func isValidSemVer(version string) bool {
    // 简单的语义化版本验证
    pattern := `^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`
    matched, _ := regexp.MatchString(pattern, version)
    return matched
}
```

## 类型转换工具

### JSON 序列化

```go
import "encoding/json"

// 将包信息转换为 JSON
func packageToJSON(pkg *models.Package) ([]byte, error) {
    return json.MarshalIndent(pkg, "", "  ")
}

// 从 JSON 解析包信息
func packageFromJSON(data []byte) (*models.Package, error) {
    var pkg models.Package
    err := json.Unmarshal(data, &pkg)
    return &pkg, err
}
```

### 数据提取

```go
// 提取包的所有依赖（包括开发依赖）
func getAllDependencies(version *models.Version) map[string]string {
    deps := make(map[string]string)
    
    // 添加运行时依赖
    for name, ver := range version.Dependencies {
        deps[name] = ver
    }
    
    // 添加开发依赖
    for name, ver := range version.DevDependencies {
        deps[name+"@dev"] = ver
    }
    
    return deps
}

// 获取包的所有版本号，按语义化版本排序
func getSortedVersions(pkg *models.Package) []string {
    versions := make([]string, 0, len(pkg.Versions))
    for version := range pkg.Versions {
        versions = append(versions, version)
    }
    
    // 这里可以使用语义化版本排序库
    // sort.Slice(versions, func(i, j int) bool {
    //     return semver.Compare(versions[i], versions[j]) < 0
    // })
    
    return versions
}
``` 