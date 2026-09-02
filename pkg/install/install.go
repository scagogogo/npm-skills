package install

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/scagogogo/npm-skills/pkg/models"
	"github.com/scagogogo/npm-skills/pkg/registry"
)

// InstallOptions 控制安装器行为。
//
// MaxDepth 语义：
//   - -1：无限递归（默认）。
//   - 0：仅安装顶层包，不递归依赖。
//   - >0：递归安装 N 层依赖。
//
// 调用方必须显式传入所需值（CLI 的 --max-depth 默认 -1，MCP 的 max_depth 默认 -1）。
type InstallOptions struct {
	MaxDepth           int
	IncludeDev         bool
	SkipIntegrityCheck bool
}

// defaultInstallOptions 返回带默认值的 InstallOptions。
//
// 注意：此处不再把 MaxDepth==0 转换为 -1。MaxDepth 是显式语义：
// 0 表示"仅顶层"，-1 表示"无限"。调用方（CLI/MCP）始终传入显式值，
// 故 defaultInstallOptions 当前为透传 no-op，保留作为未来扩展点。
func defaultInstallOptions(opts InstallOptions) InstallOptions {
	return opts
}

// InstallResult 描述一次安装的结果。
type InstallResult struct {
	Name         string   `json:"name"`
	Version      string   `json:"version"`
	Path         string   `json:"path"`
	Dependencies []string `json:"dependencies,omitempty"`
	Skipped      []string `json:"skipped,omitempty"`
}

// Installer 负责将 NPM 包安装到本地 node_modules 目录。
type Installer struct {
	registry *registry.Registry
}

// NewInstaller 用给定 Registry 客户端创建安装器。
func NewInstaller(reg *registry.Registry) *Installer {
	return &Installer{registry: reg}
}

// Install 解析版本、下载、校验、解压 packageName 到 destDir/node_modules/<name>，
// 并按 InstallOptions 递归安装运行时依赖。
func (inst *Installer) Install(ctx context.Context, packageName, versionRange, destDir string, opts InstallOptions) (*InstallResult, error) {
	opts = defaultInstallOptions(opts)

	resolvedVersion, versionInfo, err := inst.resolveVersion(ctx, packageName, versionRange)
	if err != nil {
		return nil, err
	}

	pkgDir := filepath.Join(destDir, "node_modules", packageName)
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create package dir %s: %w", pkgDir, err)
	}

	if existing := readInstalledVersion(pkgDir); existing == resolvedVersion {
		return &InstallResult{
			Name: packageName, Version: resolvedVersion, Path: pkgDir,
			Skipped: []string{packageName + "@" + existing + " (already installed)"},
		}, nil
	}

	tmpFile, err := os.CreateTemp("", "npm-install-*.tgz")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	if err := inst.registry.DownloadTarball(ctx, packageName, resolvedVersion, tmpPath); err != nil {
		return nil, fmt.Errorf("failed to download %s@%s: %w", packageName, resolvedVersion, err)
	}

	if !opts.SkipIntegrityCheck {
		if err := verifyIntegrity(tmpPath, versionInfo.Dist); err != nil {
			return nil, fmt.Errorf("integrity check failed for %s@%s: %w", packageName, resolvedVersion, err)
		}
	}

	if err := extractTarball(tmpPath, pkgDir); err != nil {
		return nil, fmt.Errorf("failed to extract %s@%s: %w", packageName, resolvedVersion, err)
	}

	result := &InstallResult{
		Name:    packageName,
		Version: resolvedVersion,
		Path:    pkgDir,
	}

	writeInstalledVersion(pkgDir, resolvedVersion)

	if opts.MaxDepth != 0 && len(versionInfo.Dependencies) > 0 {
		installed := map[string]string{packageName: resolvedVersion}
		for depName, depRange := range versionInfo.Dependencies {
			if _, ok := installed[depName]; ok {
				continue
			}
			depOpts := opts
			if opts.MaxDepth > 0 {
				depOpts.MaxDepth = opts.MaxDepth - 1
			}
			depRes, err := inst.installDependency(ctx, depName, depRange, pkgDir, depOpts, installed)
			if err != nil {
				result.Skipped = append(result.Skipped, depName+" (install failed: "+err.Error()+")")
				continue
			}
			if depRes != nil {
				result.Skipped = append(result.Skipped, depRes.Skipped...)
				if depRes.Name != "" && depRes.Version != "" {
					result.Dependencies = append(result.Dependencies, depRes.Name+"@"+depRes.Version)
				}
			}
		}
	}

	return result, nil
}

// installDependency 递归安装一个依赖，记录到 installed 集合供下层去重。
//
// 注意：循环防护在 Install 顶层遍历依赖时完成（installed[depName] 命中即 continue），
// 故此处不再重复检查 name 是否在 installed 中——调用方已保证。
func (inst *Installer) installDependency(ctx context.Context, name, versionRange, parentDir string, opts InstallOptions, installed map[string]string) (*InstallResult, error) {
	res, err := inst.Install(ctx, name, versionRange, parentDir, opts)
	if err != nil {
		return nil, err
	}
	installed[name] = res.Version
	for _, d := range res.Dependencies {
		parts := strings.SplitN(d, "@", 2)
		if len(parts) == 2 {
			installed[parts[0]] = parts[1]
		}
	}
	return res, nil
}

// resolveVersion 把 versionRange 解析成具体版本号 + *models.Version 元数据。
func (inst *Installer) resolveVersion(ctx context.Context, name, versionRange string) (string, *models.Version, error) {
	if versionRange == "" || versionRange == "latest" || isExactVersion(versionRange) {
		v, err := inst.registry.GetPackageVersion(ctx, name, ifEmpty(versionRange, "latest"))
		if err != nil {
			return "", nil, err
		}
		return v.Version, v, nil
	}
	versions, err := inst.registry.GetPackageVersions(ctx, name)
	if err != nil {
		return "", nil, err
	}
	picked := PickBestVersion(versions, versionRange)
	if picked == "" {
		v, err := inst.registry.GetPackageVersion(ctx, name, versionRange)
		if err != nil {
			return "", nil, fmt.Errorf("no version of %s satisfies %s: %w", name, versionRange, err)
		}
		return v.Version, v, nil
	}
	v, err := inst.registry.GetPackageVersion(ctx, name, picked)
	if err != nil {
		return "", nil, err
	}
	return v.Version, v, nil
}

// isExactVersion 判断 s 是否是纯版本号（不含范围运算符）。
func isExactVersion(s string) bool {
	if s == "" {
		return false
	}
	if strings.ContainsAny(s, "^~><=*x") {
		return false
	}
	_, err := parseVersion(s)
	return err == nil
}

// ifEmpty 返回 s，若空则返回 fallback。
func ifEmpty(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

// verifyIntegrity 用 Dist 中的 shasum（SHA1）或 integrity（SHA512）校验文件。
func verifyIntegrity(path string, dist *models.Dist) error {
	if dist == nil {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if dist.Integrity != "" {
		algo, expected, ok := parseIntegrity(dist.Integrity)
		if ok {
			switch algo {
			case "sha512":
				sum := sha512.Sum512(data)
				if hex.EncodeToString(sum[:]) != expected {
					return fmt.Errorf("sha512 mismatch: expected %s", expected)
				}
				return nil
			case "sha1":
				sum := sha1.Sum(data)
				if hex.EncodeToString(sum[:]) != expected {
					return fmt.Errorf("sha1 mismatch: expected %s", expected)
				}
				return nil
			}
		}
	}
	if dist.Shasum != "" {
		sum := sha1.Sum(data)
		if hex.EncodeToString(sum[:]) != dist.Shasum {
			return fmt.Errorf("shasum mismatch: expected %s got %s", dist.Shasum, hex.EncodeToString(sum[:]))
		}
		return nil
	}
	return nil
}

// parseIntegrity 解析 "sha512-abcdef" 形式，返回算法名和 hex 值。
func parseIntegrity(s string) (algo, hexVal string, ok bool) {
	idx := strings.Index(s, "-")
	if idx <= 0 {
		return "", "", false
	}
	return s[:idx], s[idx+1:], true
}

// extractTarball 解压 .tgz（gzip+tar）到 destDir，带路径穿越防护。
func extractTarball(tgzPath, destDir string) error {
	f, err := os.Open(tgzPath)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("gzip: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := strings.TrimPrefix(hdr.Name, "./")
		if strings.HasPrefix(target, "package/") {
			target = strings.TrimPrefix(target, "package/")
		} else if target == "package" {
			target = "."
		}

		cleaned := filepath.Clean(target)
		if strings.HasPrefix(cleaned, "..") || filepath.IsAbs(cleaned) {
			return fmt.Errorf("refusing to extract path traversal entry: %s", hdr.Name)
		}
		absDest := filepath.Join(destDir, cleaned)
		rel, err := filepath.Rel(destDir, absDest)
		if err != nil || strings.HasPrefix(rel, "..") {
			return fmt.Errorf("refusing to extract outside dest: %s", hdr.Name)
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(absDest, 0o755); err != nil {
				return err
			}
		case tar.TypeReg, tar.TypeRegA:
			if err := os.MkdirAll(filepath.Dir(absDest), 0o755); err != nil {
				return err
			}
			out, err := os.OpenFile(absDest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(hdr.Mode)&0o755)
			if err != nil {
				return err
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return err
			}
			out.Close()
		case tar.TypeSymlink:
			return fmt.Errorf("refusing to extract symlink: %s -> %s", hdr.Name, hdr.Linkname)
		default:
			// 跳过其他类型（hardlink、device 等）
		}
	}
	return nil
}

// readInstalledVersion 从 pkgDir/.installed-version 读取已安装版本。
func readInstalledVersion(pkgDir string) string {
	b, err := os.ReadFile(filepath.Join(pkgDir, ".installed-version"))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

// writeInstalledVersion 写入版本标记文件。
func writeInstalledVersion(pkgDir, version string) {
	_ = os.WriteFile(filepath.Join(pkgDir, ".installed-version"), []byte(version), 0o644)
}
