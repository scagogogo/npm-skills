# Data Models

This document describes all data structures returned by the NPM Skills API.

## Package Information

### PackageInformation
Complete package metadata including all versions and distribution tags.

```go
type PackageInformation struct {
    ID          string                     `json:"_id"`
    Name        string                     `json:"name"`
    Description string                     `json:"description"`
    DistTags    map[string]string         `json:"dist-tags"`
    Versions    map[string]PackageVersion `json:"versions"`
    Time        map[string]string         `json:"time"`
    Author      Author                    `json:"author"`
    Maintainers []Author                  `json:"maintainers"`
    Homepage    string                    `json:"homepage"`
    Keywords    []string                  `json:"keywords"`
    Repository  Repository                `json:"repository"`
    Bugs        Bugs                      `json:"bugs"`
    License     string                    `json:"license"`
    Readme      string                    `json:"readme"`
}
```

**Fields:**
- `ID` - Package identifier (usually same as name)
- `Name` - Package name
- `Description` - Package description
- `DistTags` - Distribution tags (e.g., "latest", "beta")
- `Versions` - Map of all package versions
- `Time` - Creation/modification timestamps
- `Author` - Package author information
- `Maintainers` - List of package maintainers
- `Homepage` - Package homepage URL
- `Keywords` - Package keywords for discovery
- `Repository` - Source repository information
- `Bugs` - Bug reporting information
- `License` - Package license
- `Readme` - Package README content

### PackageVersion
Information about a specific package version.

```go
type PackageVersion struct {
    Name                 string            `json:"name"`
    Version              string            `json:"version"`
    Description          string            `json:"description"`
    Main                 string            `json:"main"`
    Scripts              map[string]string `json:"scripts"`
    Dependencies         map[string]string `json:"dependencies"`
    DevDependencies      map[string]string `json:"devDependencies"`
    PeerDependencies     map[string]string `json:"peerDependencies"`
    OptionalDependencies map[string]string `json:"optionalDependencies"`
    BundleDependencies   []string          `json:"bundleDependencies"`
    Keywords             []string          `json:"keywords"`
    Author               Author            `json:"author"`
    License              string            `json:"license"`
    Repository           Repository        `json:"repository"`
    Bugs                 Bugs              `json:"bugs"`
    Homepage             string            `json:"homepage"`
    Dist                 Distribution      `json:"dist"`
    Deprecated           string            `json:"deprecated"`
}
```

**Fields:**
- `Name` - Package name
- `Version` - Specific version string
- `Description` - Version description
- `Main` - Main entry point file
- `Scripts` - NPM scripts defined in package.json
- `Dependencies` - Runtime dependencies
- `DevDependencies` - Development dependencies
- `PeerDependencies` - Peer dependencies
- `OptionalDependencies` - Optional dependencies
- `BundleDependencies` - Bundled dependencies
- `Keywords` - Package keywords
- `Author` - Package author
- `License` - Package license
- `Repository` - Repository information
- `Bugs` - Bug reporting information
- `Homepage` - Package homepage
- `Dist` - Distribution information
- `Deprecated` - Deprecation message (if deprecated)

## Supporting Structures

### Author
Author or maintainer information.

```go
type Author struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    URL   string `json:"url"`
}
```

### Repository
Source repository information.

```go
type Repository struct {
    Type      string `json:"type"`
    URL       string `json:"url"`
    Directory string `json:"directory"`
}
```

### Bugs
Bug reporting information.

```go
type Bugs struct {
    URL   string `json:"url"`
    Email string `json:"email"`
}
```

### Distribution
Package distribution information.

```go
type Distribution struct {
    Shasum     string `json:"shasum"`
    Tarball    string `json:"tarball"`
    FileCount  int    `json:"fileCount"`
    UnpackSize int    `json:"unpackedSize"`
}
```

## Search Results

### SearchResult
Results from package search operations.

```go
type SearchResult struct {
    Objects []SearchObject `json:"objects"`
    Total   int           `json:"total"`
    Time    string        `json:"time"`
}
```

**Fields:**
- `Objects` - Array of search result objects
- `Total` - Total number of matching packages
- `Time` - Search execution time

### SearchObject
Individual search result item.

```go
type SearchObject struct {
    Package     SearchPackage `json:"package"`
    Score       SearchScore   `json:"score"`
    SearchScore float64       `json:"searchScore"`
}
```

### SearchPackage
Package information in search results.

```go
type SearchPackage struct {
    Name        string            `json:"name"`
    Scope       string            `json:"scope"`
    Version     string            `json:"version"`
    Description string            `json:"description"`
    Keywords    []string          `json:"keywords"`
    Date        string            `json:"date"`
    Links       SearchLinks       `json:"links"`
    Author      Author            `json:"author"`
    Publisher   Author            `json:"publisher"`
    Maintainers []Author          `json:"maintainers"`
}
```

### SearchScore
Scoring information for search results.

```go
type SearchScore struct {
    Final   float64 `json:"final"`
    Detail  ScoreDetail `json:"detail"`
}

type ScoreDetail struct {
    Quality     float64 `json:"quality"`
    Popularity  float64 `json:"popularity"`
    Maintenance float64 `json:"maintenance"`
}
```

### SearchLinks
Links associated with search results.

```go
type SearchLinks struct {
    NPM        string `json:"npm"`
    Homepage   string `json:"homepage"`
    Repository string `json:"repository"`
    Bugs       string `json:"bugs"`
}
```

## Statistics

### DownloadStats
Package download statistics.

```go
type DownloadStats struct {
    Downloads int    `json:"downloads"`
    Start     string `json:"start"`
    End       string `json:"end"`
    Package   string `json:"package"`
}
```

**Fields:**
- `Downloads` - Number of downloads in the period
- `Start` - Period start date
- `End` - Period end date
- `Package` - Package name

## Registry Information

### RegistryInformation
Information about the registry itself.

```go
type RegistryInformation struct {
    DbName            string `json:"db_name"`
    DocCount          int    `json:"doc_count"`
    DocDelCount       int    `json:"doc_del_count"`
    UpdateSeq         int    `json:"update_seq"`
    PurgeSeq          int    `json:"purge_seq"`
    CompactRunning    bool   `json:"compact_running"`
    DiskSize          int    `json:"disk_size"`
    DataSize          int    `json:"data_size"`
    InstanceStartTime string `json:"instance_start_time"`
    DiskFormatVersion int    `json:"disk_format_version"`
    CommittedUpdateSeq int   `json:"committed_update_seq"`
}
```

**Fields:**
- `DbName` - Database name
- `DocCount` - Total number of documents (packages)
- `DocDelCount` - Number of deleted documents
- `UpdateSeq` - Update sequence number
- `PurgeSeq` - Purge sequence number
- `CompactRunning` - Whether compaction is running
- `DiskSize` - Total disk usage in bytes
- `DataSize` - Data size in bytes
- `InstanceStartTime` - Registry instance start time
- `DiskFormatVersion` - Disk format version
- `CommittedUpdateSeq` - Committed update sequence

## JSON Examples

### Package Information Example
```json
{
    "_id": "react",
    "name": "react",
    "description": "React is a JavaScript library for building user interfaces.",
    "dist-tags": {
        "latest": "18.2.0",
        "beta": "18.3.0-beta"
    },
    "versions": {
        "18.2.0": {
            "name": "react",
            "version": "18.2.0",
            "description": "React is a JavaScript library for building user interfaces.",
            "main": "index.js",
            "dependencies": {
                "loose-envify": "^1.1.0"
            },
            "license": "MIT"
        }
    },
    "author": {
        "name": "React Team",
        "email": "react@meta.com"
    },
    "license": "MIT",
    "homepage": "https://reactjs.org/"
}
```

### Search Result Example
```json
{
    "objects": [
        {
            "package": {
                "name": "react",
                "version": "18.2.0",
                "description": "React is a JavaScript library for building user interfaces.",
                "keywords": ["react", "javascript", "ui"],
                "links": {
                    "npm": "https://www.npmjs.com/package/react",
                    "homepage": "https://reactjs.org/"
                }
            },
            "score": {
                "final": 0.95,
                "detail": {
                    "quality": 0.98,
                    "popularity": 0.99,
                    "maintenance": 0.88
                }
            }
        }
    ],
    "total": 1,
    "time": "Wed Jan 01 2024 12:00:00 GMT+0000 (UTC)"
}
```

### Download Stats Example
```json
{
    "downloads": 18500000,
    "start": "2024-01-01",
    "end": "2024-01-31",
    "package": "react"
}
```

## Usage Examples

### Accessing Package Information
```go
pkg, err := client.GetPackageInformation(ctx, "react")
if err != nil {
    return err
}

// Access basic information
fmt.Printf("Name: %s\n", pkg.Name)
fmt.Printf("Latest version: %s\n", pkg.DistTags["latest"])
fmt.Printf("Description: %s\n", pkg.Description)

// Access author information
if pkg.Author.Name != "" {
    fmt.Printf("Author: %s <%s>\n", pkg.Author.Name, pkg.Author.Email)
}

// List all versions
for version := range pkg.Versions {
    fmt.Printf("Version: %s\n", version)
}
```

### Working with Search Results
```go
results, err := client.SearchPackages(ctx, "react ui", 5)
if err != nil {
    return err
}

fmt.Printf("Found %d packages\n", results.Total)

for _, obj := range results.Objects {
    pkg := obj.Package
    score := obj.Score
    
    fmt.Printf("Package: %s (score: %.2f)\n", pkg.Name, score.Final)
    fmt.Printf("  Description: %s\n", pkg.Description)
    fmt.Printf("  Quality: %.2f, Popularity: %.2f, Maintenance: %.2f\n",
        score.Detail.Quality, score.Detail.Popularity, score.Detail.Maintenance)
}
```

### Analyzing Dependencies
```go
version, err := client.GetPackageVersion(ctx, "react", "18.2.0")
if err != nil {
    return err
}

fmt.Printf("Dependencies for %s@%s:\n", version.Name, version.Version)

if len(version.Dependencies) > 0 {
    fmt.Println("Runtime dependencies:")
    for dep, ver := range version.Dependencies {
        fmt.Printf("  %s: %s\n", dep, ver)
    }
}

if len(version.DevDependencies) > 0 {
    fmt.Println("Dev dependencies:")
    for dep, ver := range version.DevDependencies {
        fmt.Printf("  %s: %s\n", dep, ver)
    }
}
```

## Next Steps

- Review [Registry API](registry.md) for method documentation
- Check [Configuration Options](configuration.md) for client setup
- Explore [Examples](../examples/) for practical usage patterns 