package install

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/scagogogo/npm-skills/pkg/models"
	"github.com/scagogogo/npm-skills/pkg/registry"
	"github.com/stretchr/testify/assert"
)

// buildTestTgz 构造 npm 风格 tarball，返回字节与 SHA1 hex。
func buildTestTgz(t *testing.T, pkgJSON string) ([]byte, string) {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	files := []struct{ name, body string }{
		{"package/package.json", pkgJSON},
		{"package/index.js", "module.exports = 1;"},
	}
	for _, f := range files {
		_ = tw.WriteHeader(&tar.Header{Name: f.name, Mode: 0o644, Size: int64(len(f.body)), Typeflag: tar.TypeReg})
		_, _ = tw.Write([]byte(f.body))
	}
	_ = tw.Close()
	_ = gz.Close()
	data := buf.Bytes()
	sum := sha1.Sum(data)
	return data, hex.EncodeToString(sum[:])
}

// installMockServer 建一个 mock registry，tarball 指向自身。
// deps 为可选的 testpkg dependencies JSON 片段（如 `{"dep-a":"^1.0.0"}`），
// 注入到 /testpkg/1.0.0 版本元数据；未传时默认 {"dep-a":"^1.0.0"}。
func installMockServer(t *testing.T, tarballBytes []byte, shasum string, deps ...string) *httptest.Server {
	t.Helper()
	pkgDeps := `{"dep-a":"^1.0.0"}`
	if len(deps) > 0 {
		pkgDeps = deps[0]
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		base := "http://" + r.Host
		switch r.URL.Path {
		case "/testpkg/1.0.0", "/testpkg/latest":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"name":"testpkg","version":"1.0.0","dependencies":` + pkgDeps + `,"dist":{"shasum":"` + shasum + `","tarball":"` + base + `/testpkg/-/testpkg-1.0.0.tgz"}}`))
		case "/testpkg":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"name":"testpkg","dist-tags":{"latest":"1.0.0"},"versions":{"1.0.0":{"name":"testpkg","version":"1.0.0"}}}`))
		case "/testpkg/-/testpkg-1.0.0.tgz":
			w.Write(tarballBytes)
		case "/dep-a/1.0.0", "/dep-a/latest":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"name":"dep-a","version":"1.0.0","dist":{"shasum":"` + shasum + `","tarball":"` + base + `/dep-a/-/dep-a-1.0.0.tgz"}}`))
		case "/dep-a":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"name":"dep-a","dist-tags":{"latest":"1.0.0"},"versions":{"1.0.0":{"name":"dep-a","version":"1.0.0"}}}`))
		case "/dep-a/-/dep-a-1.0.0.tgz":
			w.Write(tarballBytes)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestInstallExactVersion(t *testing.T) {
	tarBytes, shasum := buildTestTgz(t, `{"name":"testpkg","version":"1.0.0"}`)
	server := installMockServer(t, tarBytes, shasum)
	defer server.Close()

	reg := registry.NewRegistry(registry.NewOptions().SetRegistryURL(server.URL))
	inst := NewInstaller(reg)

	destDir := t.TempDir()
	res, err := inst.Install(context.Background(), "testpkg", "1.0.0", destDir, InstallOptions{MaxDepth: 0})
	assert.NoError(t, err)
	assert.Equal(t, "1.0.0", res.Version)
	assert.FileExists(t, filepath.Join(destDir, "node_modules", "testpkg", "index.js"))
	assert.FileExists(t, filepath.Join(destDir, "node_modules", "testpkg", "package.json"))
}

func TestInstallLatest(t *testing.T) {
	tarBytes, shasum := buildTestTgz(t, `{"name":"testpkg","version":"1.0.0"}`)
	server := installMockServer(t, tarBytes, shasum)
	defer server.Close()

	reg := registry.NewRegistry(registry.NewOptions().SetRegistryURL(server.URL))
	inst := NewInstaller(reg)

	destDir := t.TempDir()
	res, err := inst.Install(context.Background(), "testpkg", "latest", destDir, InstallOptions{MaxDepth: 0})
	assert.NoError(t, err)
	assert.Equal(t, "1.0.0", res.Version)
}

func TestInstallRangeResolvesToBest(t *testing.T) {
	tarBytes, shasum := buildTestTgz(t, `{"name":"testpkg","version":"1.0.0"}`)
	server := installMockServer(t, tarBytes, shasum)
	defer server.Close()

	reg := registry.NewRegistry(registry.NewOptions().SetRegistryURL(server.URL))
	inst := NewInstaller(reg)

	destDir := t.TempDir()
	res, err := inst.Install(context.Background(), "testpkg", "^1.0.0", destDir, InstallOptions{MaxDepth: 0})
	assert.NoError(t, err)
	assert.Equal(t, "1.0.0", res.Version)
}

func TestInstallAlreadyInstalledSkips(t *testing.T) {
	tarBytes, shasum := buildTestTgz(t, `{"name":"testpkg","version":"1.0.0"}`)
	server := installMockServer(t, tarBytes, shasum)
	defer server.Close()

	reg := registry.NewRegistry(registry.NewOptions().SetRegistryURL(server.URL))
	inst := NewInstaller(reg)

	destDir := t.TempDir()
	_, err := inst.Install(context.Background(), "testpkg", "1.0.0", destDir, InstallOptions{MaxDepth: 0})
	assert.NoError(t, err)

	res, err := inst.Install(context.Background(), "testpkg", "1.0.0", destDir, InstallOptions{MaxDepth: 0})
	assert.NoError(t, err)
	assert.NotEmpty(t, res.Skipped)
}

func TestInstallIntegrityMismatch(t *testing.T) {
	tarBytes, _ := buildTestTgz(t, `{"name":"testpkg","version":"1.0.0"}`)
	server := installMockServer(t, tarBytes, "deadbeef")
	defer server.Close()

	reg := registry.NewRegistry(registry.NewOptions().SetRegistryURL(server.URL))
	inst := NewInstaller(reg)

	destDir := t.TempDir()
	_, err := inst.Install(context.Background(), "testpkg", "1.0.0", destDir, InstallOptions{MaxDepth: 0})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "integrity")
}

func TestInstallSkipIntegrityCheck(t *testing.T) {
	tarBytes, _ := buildTestTgz(t, `{"name":"testpkg","version":"1.0.0"}`)
	server := installMockServer(t, tarBytes, "deadbeef")
	defer server.Close()

	reg := registry.NewRegistry(registry.NewOptions().SetRegistryURL(server.URL))
	inst := NewInstaller(reg)

	destDir := t.TempDir()
	res, err := inst.Install(context.Background(), "testpkg", "1.0.0", destDir, InstallOptions{MaxDepth: 0, SkipIntegrityCheck: true})
	assert.NoError(t, err)
	assert.Equal(t, "1.0.0", res.Version)
}

func TestInstallRecursiveDependencies(t *testing.T) {
	tarBytes, shasum := buildTestTgz(t, `{"name":"testpkg","version":"1.0.0","dependencies":{"dep-a":"^1.0.0"}}`)
	server := installMockServer(t, tarBytes, shasum)
	defer server.Close()

	reg := registry.NewRegistry(registry.NewOptions().SetRegistryURL(server.URL))
	inst := NewInstaller(reg)

	destDir := t.TempDir()
	res, err := inst.Install(context.Background(), "testpkg", "1.0.0", destDir, InstallOptions{MaxDepth: 3})
	assert.NoError(t, err)
	assert.Contains(t, res.Dependencies, "dep-a@1.0.0")
	assert.FileExists(t, filepath.Join(destDir, "node_modules", "testpkg", "node_modules", "dep-a", "package.json"))
}

func TestExtractTarballRejectsPathTraversal(t *testing.T) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	_ = tw.WriteHeader(&tar.Header{Name: "../evil.txt", Mode: 0o644, Size: 4, Typeflag: tar.TypeReg})
	_, _ = tw.Write([]byte("evil"))
	_ = tw.Close()
	_ = gz.Close()

	tmpPath := filepath.Join(t.TempDir(), "evil.tgz")
	assert.NoError(t, os.WriteFile(tmpPath, buf.Bytes(), 0o644))
	err := extractTarball(tmpPath, t.TempDir())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "traversal")
}

func TestExtractTarballRejectsSymlink(t *testing.T) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	_ = tw.WriteHeader(&tar.Header{Name: "link", Typeflag: tar.TypeSymlink, Linkname: "/etc/passwd", Mode: 0o777})
	_ = tw.Close()
	_ = gz.Close()

	tmpPath := filepath.Join(t.TempDir(), "sym.tgz")
	assert.NoError(t, os.WriteFile(tmpPath, buf.Bytes(), 0o644))
	err := extractTarball(tmpPath, t.TempDir())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "symlink")
}

func TestExtractTarballDirectoryEntry(t *testing.T) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	_ = tw.WriteHeader(&tar.Header{Name: "sub", Mode: 0o755, Typeflag: tar.TypeDir})
	_ = tw.WriteHeader(&tar.Header{Name: "sub/file.js", Mode: 0o644, Size: 1, Typeflag: tar.TypeReg})
	_, _ = tw.Write([]byte("x"))
	_ = tw.Close()
	_ = gz.Close()

	tmpPath := filepath.Join(t.TempDir(), "dir.tgz")
	assert.NoError(t, os.WriteFile(tmpPath, buf.Bytes(), 0o644))
	dest := t.TempDir()
	assert.NoError(t, extractTarball(tmpPath, dest))
	assert.DirExists(t, filepath.Join(dest, "sub"))
	assert.FileExists(t, filepath.Join(dest, "sub", "file.js"))
}

func TestParseIntegrityInvalid(t *testing.T) {
	_, _, ok := parseIntegrity("nodashhere")
	assert.False(t, ok)
	_, _, ok = parseIntegrity("-noalgo")
	assert.False(t, ok)
}

func TestVerifyIntegrityNilDist(t *testing.T) {
	f := filepath.Join(t.TempDir(), "x.tmp")
	assert.NoError(t, os.WriteFile(f, []byte("x"), 0o644))
	assert.NoError(t, verifyIntegrity(f, nil))
}

func TestResolveVersionRangeFallback(t *testing.T) {
	tarBytes, shasum := buildTestTgz(t, `{"name":"testpkg","version":"1.0.0"}`)
	server := installMockServer(t, tarBytes, shasum)
	defer server.Close()
	reg := registry.NewRegistry(registry.NewOptions().SetRegistryURL(server.URL))
	inst := NewInstaller(reg)
	_, _, err := inst.resolveVersion(context.Background(), "testpkg", "^9.9.9")
	assert.Error(t, err)
}

func TestInstallDependencyFailureSkipped(t *testing.T) {
	tarBytes, shasum := buildTestTgz(t, `{"name":"testpkg","version":"1.0.0","dependencies":{"notexist":"^1.0.0"}}`)
	server := installMockServer(t, tarBytes, shasum, `{"notexist":"^1.0.0"}`)
	defer server.Close()
	reg := registry.NewRegistry(registry.NewOptions().SetRegistryURL(server.URL))
	inst := NewInstaller(reg)
	destDir := t.TempDir()
	res, err := inst.Install(context.Background(), "testpkg", "1.0.0", destDir, InstallOptions{MaxDepth: 3})
	assert.NoError(t, err)
	assert.NotEmpty(t, res.Skipped)
	assert.Contains(t, res.Skipped[0], "notexist")
}

// === 覆盖率补测：verifyIntegrity 的 integrity 分支 ===

func TestVerifyIntegritySha512Match(t *testing.T) {
	tarBytes, _ := buildTestTgz(t, `{"name":"x","version":"1.0.0"}`)
	sum512 := sha512.Sum512(tarBytes)
	integrity := "sha512-" + hex.EncodeToString(sum512[:])
	f := filepath.Join(t.TempDir(), "x.tgz")
	assert.NoError(t, os.WriteFile(f, tarBytes, 0o644))
	dist := &models.Dist{Integrity: integrity}
	assert.NoError(t, verifyIntegrity(f, dist))
}

func TestVerifyIntegritySha512Mismatch(t *testing.T) {
	f := filepath.Join(t.TempDir(), "x.tgz")
	assert.NoError(t, os.WriteFile(f, []byte("data"), 0o644))
	dist := &models.Dist{Integrity: "sha512-" + strings.Repeat("0", 128)}
	err := verifyIntegrity(f, dist)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sha512")
}

func TestVerifyIntegritySha1ViaIntegrity(t *testing.T) {
	tarBytes, _ := buildTestTgz(t, `{"name":"x","version":"1.0.0"}`)
	sum1 := sha1.Sum(tarBytes)
	integrity := "sha1-" + hex.EncodeToString(sum1[:])
	f := filepath.Join(t.TempDir(), "x.tgz")
	assert.NoError(t, os.WriteFile(f, tarBytes, 0o644))
	dist := &models.Dist{Integrity: integrity}
	assert.NoError(t, verifyIntegrity(f, dist))
}

func TestVerifyIntegritySha1ViaIntegrityMismatch(t *testing.T) {
	f := filepath.Join(t.TempDir(), "x.tgz")
	assert.NoError(t, os.WriteFile(f, []byte("data"), 0o644))
	dist := &models.Dist{Integrity: "sha1-" + strings.Repeat("0", 40)}
	err := verifyIntegrity(f, dist)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sha1")
}

func TestVerifyIntegrityUnknownAlgoFallsBackToShasum(t *testing.T) {
	// 未知算法 → 回退到 shasum 路径
	tarBytes, _ := buildTestTgz(t, `{"name":"x","version":"1.0.0"}`)
	sum1 := sha1.Sum(tarBytes)
	f := filepath.Join(t.TempDir(), "x.tgz")
	assert.NoError(t, os.WriteFile(f, tarBytes, 0o644))
	// Integrity 用未知算法 + 不匹配 → 回退 shasum（匹配）→ nil
	dist := &models.Dist{Integrity: "md5-abcdef", Shasum: hex.EncodeToString(sum1[:])}
	assert.NoError(t, verifyIntegrity(f, dist))
}

func TestVerifyIntegrityReadFileError(t *testing.T) {
	// 文件不存在 → ReadFile 失败
	dist := &models.Dist{Shasum: "abc"}
	err := verifyIntegrity(filepath.Join(t.TempDir(), "no-such-file.tgz"), dist)
	assert.Error(t, err)
}

func TestVerifyIntegrityNoChecksums(t *testing.T) {
	// dist 无 shasum 无 integrity → nil
	f := filepath.Join(t.TempDir(), "x.tgz")
	assert.NoError(t, os.WriteFile(f, []byte("x"), 0o644))
	dist := &models.Dist{}
	assert.NoError(t, verifyIntegrity(f, dist))
}

func TestParseIntegrityValid(t *testing.T) {
	algo, val, ok := parseIntegrity("sha512-abcdef")
	assert.True(t, ok)
	assert.Equal(t, "sha512", algo)
	assert.Equal(t, "abcdef", val)
}

// === extractTarball 错误分支 ===

func TestExtractTarballOpenError(t *testing.T) {
	err := extractTarball(filepath.Join(t.TempDir(), "no-such.tgz"), t.TempDir())
	assert.Error(t, err)
}

func TestExtractTarballNotGzip(t *testing.T) {
	// 非 gzip 数据 → gzip.NewReader 失败
	tmpPath := filepath.Join(t.TempDir(), "plain.tgz")
	assert.NoError(t, os.WriteFile(tmpPath, []byte("not gzip data at all"), 0o644))
	err := extractTarball(tmpPath, t.TempDir())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "gzip")
}

func TestExtractTarballCorruptTar(t *testing.T) {
	// 有效 gzip 但 tar 数据损坏
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write([]byte("not a tar stream"))
	gz.Close()
	tmpPath := filepath.Join(t.TempDir(), "corrupt.tgz")
	assert.NoError(t, os.WriteFile(tmpPath, buf.Bytes(), 0o644))
	err := extractTarball(tmpPath, t.TempDir())
	assert.Error(t, err)
}

func TestExtractTarballHardlinkSkipped(t *testing.T) {
	// TypeLink 等其他类型 → default 跳过
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	_ = tw.WriteHeader(&tar.Header{Name: "hardlink", Typeflag: tar.TypeLink, Linkname: "target", Mode: 0o644})
	_ = tw.Close()
	_ = gz.Close()
	tmpPath := filepath.Join(t.TempDir(), "link.tgz")
	assert.NoError(t, os.WriteFile(tmpPath, buf.Bytes(), 0o644))
	// hardlink 被 default 跳过，不报错
	assert.NoError(t, extractTarball(tmpPath, t.TempDir()))
}

func TestExtractTarballAbsoluteEntryRejected(t *testing.T) {
	// 绝对路径条目 → 拒绝
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	_ = tw.WriteHeader(&tar.Header{Name: "/etc/passwd", Mode: 0o644, Size: 1, Typeflag: tar.TypeReg})
	_, _ = tw.Write([]byte("x"))
	_ = tw.Close()
	_ = gz.Close()
	tmpPath := filepath.Join(t.TempDir(), "abs.tgz")
	assert.NoError(t, os.WriteFile(tmpPath, buf.Bytes(), 0o644))
	err := extractTarball(tmpPath, t.TempDir())
	assert.Error(t, err)
}

// === resolveVersion 错误分支 ===

func TestResolveVersionExactVersionError(t *testing.T) {
	// 精确版本但 server 不可达 → GetPackageVersion 失败
	reg := registry.NewRegistry(registry.NewOptions().SetRegistryURL("http://localhost:1"))
	inst := NewInstaller(reg)
	_, _, err := inst.resolveVersion(context.Background(), "testpkg", "1.0.0")
	assert.Error(t, err)
}

func TestResolveVersionGetVersionsError(t *testing.T) {
	// 范围表达式但 GetPackageVersions 失败（server 不可达）
	reg := registry.NewRegistry(registry.NewOptions().SetRegistryURL("http://localhost:1"))
	inst := NewInstaller(reg)
	_, _, err := inst.resolveVersion(context.Background(), "testpkg", "^1.0.0")
	assert.Error(t, err)
}

func TestResolveVersionPickedFetchError(t *testing.T) {
	// picked 非空但 GetPackageVersion(picked) 失败
	// 构造 mock：/testpkg 返回 versions 含 1.0.0（^1.0.0 匹配 1.0.0），
	// 但 /testpkg/1.0.0 返回 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/testpkg" {
			w.Write([]byte(`{"name":"testpkg","dist-tags":{"latest":"1.0.0"},"versions":{"1.0.0":{"name":"testpkg","version":"1.0.0"}}}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()
	reg := registry.NewRegistry(registry.NewOptions().SetRegistryURL(server.URL))
	inst := NewInstaller(reg)
	_, _, err := inst.resolveVersion(context.Background(), "testpkg", "^1.0.0")
	assert.Error(t, err)
}

// === isExactVersion / ifEmpty 边界 ===

func TestIsExactVersionEmpty(t *testing.T) {
	assert.False(t, isExactVersion(""))
}

func TestIsExactVersionWithOperator(t *testing.T) {
	assert.False(t, isExactVersion("^1.0.0"))
	assert.False(t, isExactVersion("1.x"))
}

func TestIsExactVersionValid(t *testing.T) {
	assert.True(t, isExactVersion("1.2.3"))
}

func TestIfEmptyFallback(t *testing.T) {
	assert.Equal(t, "fallback", ifEmpty("", "fallback"))
	assert.Equal(t, "value", ifEmpty("value", "fallback"))
}

// === installDependency 循环与重复 ===

func TestInstallDependencyCycleSkipped(t *testing.T) {
	// pkg-a 依赖 pkg-b，pkg-b 依赖 pkg-a → 循环：pkg-b 安装时发现 pkg-a 已 installed，跳过
	tarBytes, shasum := buildTestTgz(t, `{"name":"pkg-a","version":"1.0.0","dependencies":{"pkg-b":"^1.0.0"}}`)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		base := "http://" + r.Host
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/pkg-a":
			// packument 端点：含 versions map 供 GetPackageVersions 解析
			w.Write([]byte(`{"name":"pkg-a","dist-tags":{"latest":"1.0.0"},"versions":{"1.0.0":{"name":"pkg-a","version":"1.0.0","dependencies":{"pkg-b":"^1.0.0"}}}}`))
		case "/pkg-a/1.0.0", "/pkg-a/latest":
			w.Write([]byte(`{"name":"pkg-a","version":"1.0.0","dependencies":{"pkg-b":"^1.0.0"},"dist":{"shasum":"` + shasum + `","tarball":"` + base + `/pkg-a/-/pkg-a-1.0.0.tgz"}}`))
		case "/pkg-a/-/pkg-a-1.0.0.tgz":
			w.Write(tarBytes)
		case "/pkg-b":
			// packument 端点：含 versions map
			w.Write([]byte(`{"name":"pkg-b","dist-tags":{"latest":"1.0.0"},"versions":{"1.0.0":{"name":"pkg-b","version":"1.0.0","dependencies":{"pkg-a":"^1.0.0"}}}}`))
		case "/pkg-b/1.0.0", "/pkg-b/latest":
			// pkg-b 依赖 pkg-a（循环）
			w.Write([]byte(`{"name":"pkg-b","version":"1.0.0","dependencies":{"pkg-a":"^1.0.0"},"dist":{"shasum":"` + shasum + `","tarball":"` + base + `/pkg-b/-/pkg-b-1.0.0.tgz"}}`))
		case "/pkg-b/-/pkg-b-1.0.0.tgz":
			w.Write(tarBytes)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	reg := registry.NewRegistry(registry.NewOptions().SetRegistryURL(server.URL))
	inst := NewInstaller(reg)
	destDir := t.TempDir()
	// MaxDepth 足够深以触发循环
	res, err := inst.Install(context.Background(), "pkg-a", "1.0.0", destDir, InstallOptions{MaxDepth: 5})
	assert.NoError(t, err)
	// pkg-b 被安装（作为 pkg-a 的依赖），pkg-b 再依赖 pkg-a 时因循环被跳过
	assert.Contains(t, res.Dependencies, "pkg-b@1.0.0")
}

// === Install 下载失败分支 ===

func TestInstallDownloadFailure(t *testing.T) {
	// 版本元数据 OK 但 tarball URL 不可达
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/testpkg/1.0.0" {
			w.Write([]byte(`{"name":"testpkg","version":"1.0.0","dist":{"shasum":"abc","tarball":"http://localhost:1/testpkg.tgz"}}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()
	reg := registry.NewRegistry(registry.NewOptions().SetRegistryURL(server.URL))
	inst := NewInstaller(reg)
	_, err := inst.Install(context.Background(), "testpkg", "1.0.0", t.TempDir(), InstallOptions{MaxDepth: 0, SkipIntegrityCheck: true})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "download")
}

// === Install 解压失败分支 ===

func TestInstallExtractFailure(t *testing.T) {
	// tarball 内容不是有效 gzip → extract 失败
	tarballBytes := []byte("not a valid tgz")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		base := "http://" + r.Host
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/testpkg/1.0.0":
			w.Write([]byte(`{"name":"testpkg","version":"1.0.0","dist":{"shasum":"abc","tarball":"` + base + `/testpkg/-/testpkg-1.0.0.tgz"}}`))
		case "/testpkg/-/testpkg-1.0.0.tgz":
			w.Write(tarballBytes)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	reg := registry.NewRegistry(registry.NewOptions().SetRegistryURL(server.URL))
	inst := NewInstaller(reg)
	_, err := inst.Install(context.Background(), "testpkg", "1.0.0", t.TempDir(), InstallOptions{MaxDepth: 0, SkipIntegrityCheck: true})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "extract")
}

// === extractTarball: package 目录条目分支 ===

func TestExtractTarballPackageRootEntry(t *testing.T) {
	// tar 含 "package" 目录条目（target == "package" → target = "."）
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	_ = tw.WriteHeader(&tar.Header{Name: "package", Mode: 0o755, Typeflag: tar.TypeDir})
	_ = tw.WriteHeader(&tar.Header{Name: "package/file.js", Mode: 0o644, Size: 1, Typeflag: tar.TypeReg})
	_, _ = tw.Write([]byte("x"))
	_ = tw.Close()
	_ = gz.Close()
	tmpPath := filepath.Join(t.TempDir(), "pkg.tgz")
	assert.NoError(t, os.WriteFile(tmpPath, buf.Bytes(), 0o644))
	dest := t.TempDir()
	assert.NoError(t, extractTarball(tmpPath, dest))
	assert.FileExists(t, filepath.Join(dest, "file.js"))
}

func TestInstallDependencyNilResultSkipped(t *testing.T) {
	// 依赖安装成功后 installDependency 返回 (*InstallResult, nil)，
	// 这里构造一个场景：安装到已存在目录但 .installed-version 不存在，
	// 结果 Name/Version 为空字符串 → depRes.Name == "" 分支
	tarBytes, shasum := buildTestTgz(t, `{"name":"testpkg","version":"1.0.0","dependencies":{"dep-a":"^1.0.0"}}`)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		base := "http://" + r.Host
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/testpkg":
			w.Write([]byte(`{"name":"testpkg","dist-tags":{"latest":"1.0.0"},"versions":{"1.0.0":{"name":"testpkg","version":"1.0.0","dependencies":{"dep-a":"^1.0.0"}}}}`))
		case "/testpkg/1.0.0", "/testpkg/latest":
			w.Write([]byte(`{"name":"testpkg","version":"1.0.0","dependencies":{"dep-a":"^1.0.0"},"dist":{"shasum":"` + shasum + `","tarball":"` + base + `/testpkg/-/testpkg-1.0.0.tgz"}}`))
		case "/testpkg/-/testpkg-1.0.0.tgz":
			w.Write(tarBytes)
		case "/dep-a":
			w.Write([]byte(`{"name":"dep-a","dist-tags":{"latest":"1.0.0"},"versions":{"1.0.0":{"name":"dep-a","version":"1.0.0"}}}`))
		case "/dep-a/1.0.0", "/dep-a/latest":
			w.Write([]byte(`{"name":"dep-a","version":"1.0.0","dist":{"shasum":"` + shasum + `","tarball":"` + base + `/dep-a/-/dep-a-1.0.0.tgz"}}`))
		case "/dep-a/-/dep-a-1.0.0.tgz":
			w.Write(tarBytes)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	reg := registry.NewRegistry(registry.NewOptions().SetRegistryURL(server.URL))
	inst := NewInstaller(reg)
	destDir := t.TempDir()
	res, err := inst.Install(context.Background(), "testpkg", "1.0.0", destDir, InstallOptions{MaxDepth: 3})
	assert.NoError(t, err)
	assert.Contains(t, res.Dependencies, "dep-a@1.0.0")
}

// === extractTarball: 文件系统错误分支 ===

func TestExtractTarballDirMkdirFailure(t *testing.T) {
	// destDir 是一个已存在的文件（非目录）→ MkdirAll 对 TypeDir 失败
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	_ = tw.WriteHeader(&tar.Header{Name: "sub", Mode: 0o755, Typeflag: tar.TypeDir})
	_ = tw.Close()
	_ = gz.Close()
	tmpPath := filepath.Join(t.TempDir(), "d.tgz")
	assert.NoError(t, os.WriteFile(tmpPath, buf.Bytes(), 0o644))
	// destDir 是一个已存在文件，MkdirAll 失败
	destFile := filepath.Join(t.TempDir(), "notadir")
	assert.NoError(t, os.WriteFile(destFile, []byte("x"), 0o644))
	err := extractTarball(tmpPath, destFile)
	assert.Error(t, err)
}

func TestExtractTarballRegMkdirFailure(t *testing.T) {
	// TypeReg 的父目录 MkdirAll 失败：destDir 是已存在文件
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	_ = tw.WriteHeader(&tar.Header{Name: "sub/file.js", Mode: 0o644, Size: 1, Typeflag: tar.TypeReg})
	_, _ = tw.Write([]byte("x"))
	_ = tw.Close()
	_ = gz.Close()
	tmpPath := filepath.Join(t.TempDir(), "r.tgz")
	assert.NoError(t, os.WriteFile(tmpPath, buf.Bytes(), 0o644))
	destFile := filepath.Join(t.TempDir(), "blockeddir")
	assert.NoError(t, os.WriteFile(destFile, []byte("x"), 0o644))
	err := extractTarball(tmpPath, destFile)
	assert.Error(t, err)
}

func TestExtractTarballOpenFileFailure(t *testing.T) {
	// OpenFile 失败：目标路径是一个已存在目录（无法以文件形式打开）
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	_ = tw.WriteHeader(&tar.Header{Name: "file.js", Mode: 0o644, Size: 1, Typeflag: tar.TypeReg})
	_, _ = tw.Write([]byte("x"))
	_ = tw.Close()
	_ = gz.Close()
	tmpPath := filepath.Join(t.TempDir(), "o.tgz")
	assert.NoError(t, os.WriteFile(tmpPath, buf.Bytes(), 0o644))
	// destDir 下 "file.js" 已是一个目录 → OpenFile(文件) 失败
	dest := t.TempDir()
	assert.NoError(t, os.MkdirAll(filepath.Join(dest, "file.js"), 0o755))
	err := extractTarball(tmpPath, dest)
	assert.Error(t, err)
}

func TestExtractTarballCopyFailure(t *testing.T) {
	// io.Copy 失败：tar header 声明 Size=10 但实际只写 1 字节 → tar reader 在 Copy 时 EOF 前 size 不足报错
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	// 声明 10 字节但写入 1 字节，tar.Next 读取时会因 size 不匹配触发 io.Copy 的 ErrUnexpectedEOF
	_ = tw.WriteHeader(&tar.Header{Name: "short.js", Mode: 0o644, Size: 10, Typeflag: tar.TypeReg})
	_, _ = tw.Write([]byte("x"))
	_ = tw.Close()
	_ = gz.Close()
	tmpPath := filepath.Join(t.TempDir(), "c.tgz")
	assert.NoError(t, os.WriteFile(tmpPath, buf.Bytes(), 0o644))
	err := extractTarball(tmpPath, t.TempDir())
	assert.Error(t, err)
}

// === Install: 顶层依赖去重 continue 分支（自循环依赖）===

func TestInstallSelfDependencySkipped(t *testing.T) {
	// testpkg 依赖自身 → installed[packageName] 命中 → continue（去重分支）
	tarBytes, shasum := buildTestTgz(t, `{"name":"testpkg","version":"1.0.0"}`)
	server := installMockServer(t, tarBytes, shasum, `{"testpkg":"^1.0.0"}`)
	defer server.Close()
	reg := registry.NewRegistry(registry.NewOptions().SetRegistryURL(server.URL))
	inst := NewInstaller(reg)
	// 注入自循环依赖到版本元数据需要 mock 返回 dependencies 含 testpkg
	// installMockServer 的 /testpkg/1.0.0 用 deps 注入
	res, err := inst.Install(context.Background(), "testpkg", "1.0.0", t.TempDir(), InstallOptions{MaxDepth: 3})
	assert.NoError(t, err)
	// 自循环依赖被跳过，Dependencies 不含 testpkg
	assert.NotContains(t, res.Dependencies, "testpkg@")
}

// === resolveVersion: picked==="" 后 GetPackageVersion 成功返回分支 ===

func TestResolveVersionRangeFallbackToExact(t *testing.T) {
	// ^9.9.9 无匹配 → picked==="" → GetPackageVersion(ctx, name, "^9.9.9")
	// 构造 mock：/testpkg 的 versions 空（PickBestVersion 返回 ""），
	// 但 /testpkg/^9.9.9 端点返回有效版本（模拟 registry 接受裸范围作 version）
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/testpkg":
			// versions 为空 → PickBestVersion 返回 ""
			w.Write([]byte(`{"name":"testpkg","dist-tags":{"latest":"1.0.0"},"versions":{}}`))
		case "/testpkg/%5E9.9.9", "/testpkg/^9.9.9":
			// 裸范围作 version → 返回有效版本
			w.Write([]byte(`{"name":"testpkg","version":"1.0.0","dist":{"shasum":"abc"}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	reg := registry.NewRegistry(registry.NewOptions().SetRegistryURL(server.URL))
	inst := NewInstaller(reg)
	v, vi, err := inst.resolveVersion(context.Background(), "testpkg", "^9.9.9")
	assert.NoError(t, err)
	assert.Equal(t, "1.0.0", v)
	assert.NotNil(t, vi)
}

