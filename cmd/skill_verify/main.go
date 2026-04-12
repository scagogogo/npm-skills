package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/scagogogo/npm-crawler/pkg/registry"
)

func runTests() {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	fmt.Println("========================================")
	fmt.Println("NPM Crawler Skill 完整能力测试")
	fmt.Println("========================================")

	tests := []struct {
		name string
		fn   func(context.Context) error
	}{
		// 核心能力测试
		{"P1-01. Registry信息查询(NPMMirror)", testRegistryInfo},
		{"P1-02. 包信息查询 - axios", testPackageInfo},
		{"P1-03. 包搜索 - axios", testSearchPackages},
		{"P1-04. 版本查询 - axios@1.0.0", testPackageVersion},
		{"P1-05. 下载统计 - axios", testDownloadStats},
		{"P1-06. Tarball下载", testDownloadTarball},

		// 镜像源测试
		{"P2-01. 淘宝镜像源", testTaoBaoMirror},
		{"P2-02. NPMMirror镜像源", testNpmMirror},
		{"P2-03. 华为云镜像源", testHuaWeiCloudMirror},
		{"P2-04. 腾讯云镜像源", testTencentMirror},
		{"P2-05. CNPM镜像源", testCnpmMirror},
		{"P2-06. Yarn镜像源", testYarnMirror},
		{"P2-07. NPMjsCom镜像源", testNpmjsComMirror},

		// 配置测试
		{"P3-01. 自定义Registry URL", testCustomRegistryURL},
		{"P3-02. 代理配置(空代理)", testProxyConfig},
	}

	passed := 0
	failed := 0

	for _, tt := range tests {
		fmt.Printf("\n[%s]\n", tt.name)
		fmt.Println("----------------------------------------")
		if err := tt.fn(ctx); err != nil {
			fmt.Printf("❌ FAILED: %v\n", err)
			failed++
		} else {
			fmt.Printf("✅ PASSED\n")
			passed++
		}
	}

	fmt.Println("\n========================================")
	fmt.Printf("测试结果: %d 通过, %d 失败\n", passed, failed)
	fmt.Println("========================================")

	if failed > 0 {
		os.Exit(1)
	}
}

func main() {
	runTests()
}

// ============ P1: 核心能力测试 ============

func testRegistryInfo(ctx context.Context) error {
	client := registry.NewNpmMirrorRegistry()
	info, err := client.GetRegistryInformation(ctx)
	if err != nil {
		return fmt.Errorf("GetRegistryInformation failed: %w", err)
	}
	if info.DocCount <= 0 {
		return fmt.Errorf("unexpected docCount: %d", info.DocCount)
	}
	fmt.Printf("   NPMMirror Registry包数量: %d\n", info.DocCount)
	return nil
}

func testPackageInfo(ctx context.Context) error {
	client := registry.NewRegistry()
	pkg, err := client.GetPackageInformation(ctx, "axios")
	if err != nil {
		return fmt.Errorf("GetPackageInformation failed: %w", err)
	}
	if pkg.Name != "axios" {
		return fmt.Errorf("unexpected name: %s", pkg.Name)
	}
	fmt.Printf("   包名: %s, 最新版本: %s\n", pkg.Name, pkg.DistTags["latest"])
	return nil
}

func testSearchPackages(ctx context.Context) error {
	client := registry.NewRegistry()
	result, err := client.SearchPackages(ctx, "axios", 5)
	if err != nil {
		return fmt.Errorf("SearchPackages failed: %w", err)
	}
	if len(result.Objects) == 0 {
		return fmt.Errorf("no results found")
	}
	if !strings.Contains(strings.ToLower(result.Objects[0].Package.Name), "axios") {
		return fmt.Errorf("unexpected search result: %s", result.Objects[0].Package.Name)
	}
	fmt.Printf("   找到 %d 个包, 第一个: %s\n", len(result.Objects), result.Objects[0].Package.Name)
	return nil
}

func testPackageVersion(ctx context.Context) error {
	client := registry.NewRegistry()
	version, err := client.GetPackageVersion(ctx, "axios", "1.0.0")
	if err != nil {
		return fmt.Errorf("GetPackageVersion failed: %w", err)
	}
	if version.Version != "1.0.0" {
		return fmt.Errorf("unexpected version: %s", version.Version)
	}
	fmt.Printf("   版本: %s, 描述: %s\n", version.Version, version.Description)
	return nil
}

func testDownloadStats(ctx context.Context) error {
	client := registry.NewRegistry()
	stats, err := client.GetDownloadStats(ctx, "axios", "last-week")
	if err != nil {
		return fmt.Errorf("GetDownloadStats failed: %w", err)
	}
	if stats.Downloads < 0 {
		return fmt.Errorf("unexpected downloads: %d", stats.Downloads)
	}
	fmt.Printf("   包: %s, 下载次数: %d\n", stats.Package, stats.Downloads)
	return nil
}

func testDownloadTarball(ctx context.Context) error {
	client := registry.NewRegistry()
	tmpfile := "/tmp/axios-1.0.0-skill-test.tgz"
	defer os.Remove(tmpfile)

	err := client.DownloadTarball(ctx, "axios", "1.0.0", tmpfile)
	if err != nil {
		return fmt.Errorf("DownloadTarball failed: %w", err)
	}

	info, err := os.Stat(tmpfile)
	if err != nil {
		return fmt.Errorf("file not created: %w", err)
	}
	if info.Size() == 0 {
		return fmt.Errorf("file is empty")
	}
	fmt.Printf("   下载成功, 文件大小: %d bytes\n", info.Size())
	return nil
}

// ============ P2: 镜像源测试 ============

func testTaoBaoMirror(ctx context.Context) error {
	client := registry.NewTaoBaoRegistry()
	info, err := client.GetRegistryInformation(ctx)
	if err != nil {
		return fmt.Errorf("TaoBao mirror failed: %w", err)
	}
	fmt.Printf("   淘宝镜像 Registry包数量: %d\n", info.DocCount)
	return nil
}

func testNpmMirror(ctx context.Context) error {
	client := registry.NewNpmMirrorRegistry()
	info, err := client.GetRegistryInformation(ctx)
	if err != nil {
		return fmt.Errorf("NPMMirror failed: %w", err)
	}
	fmt.Printf("   NPMMirror Registry包数量: %d\n", info.DocCount)
	return nil
}

func testHuaWeiCloudMirror(ctx context.Context) error {
	client := registry.NewHuaWeiCloudRegistry()
	_, err := client.GetRegistryInformation(ctx)
	if err != nil {
		return fmt.Errorf("HuaWeiCloud mirror failed: %w", err)
	}
	fmt.Printf("   华为云镜像 API正常响应\n")
	return nil
}

func testTencentMirror(ctx context.Context) error {
	client := registry.NewTencentRegistry()
	_, err := client.GetRegistryInformation(ctx)
	if err != nil {
		return fmt.Errorf("Tencent mirror failed: %w", err)
	}
	fmt.Printf("   腾讯云镜像 API正常响应\n")
	return nil
}

func testCnpmMirror(ctx context.Context) error {
	client := registry.NewCnpmRegistry()
	info, err := client.GetRegistryInformation(ctx)
	if err != nil {
		return fmt.Errorf("CNPM mirror failed: %w", err)
	}
	fmt.Printf("   CNPM Registry包数量: %d\n", info.DocCount)
	return nil
}

func testYarnMirror(ctx context.Context) error {
	client := registry.NewYarnRegistry()
	_, err := client.GetRegistryInformation(ctx)
	if err != nil {
		return fmt.Errorf("Yarn mirror failed: %w", err)
	}
	fmt.Printf("   Yarn镜像 API正常响应\n")
	return nil
}

func testNpmjsComMirror(ctx context.Context) error {
	client := registry.NewNpmjsComRegistry()
	info, err := client.GetRegistryInformation(ctx)
	if err != nil {
		return fmt.Errorf("NPMjsCom mirror failed: %w", err)
	}
	fmt.Printf("   NPMjsCom Registry包数量: %d\n", info.DocCount)
	return nil
}

// ============ P3: 配置测试 ============

func testCustomRegistryURL(ctx context.Context) error {
	options := registry.NewOptions().
		SetRegistryURL(registry.RegistryUrlNpmjsCom)
	client := registry.NewRegistry(options)
	info, err := client.GetRegistryInformation(ctx)
	if err != nil {
		return fmt.Errorf("Custom Registry URL failed: %w", err)
	}
	fmt.Printf("   NPMJS.com Registry包数量: %d\n", info.DocCount)
	return nil
}

func testProxyConfig(ctx context.Context) error {
	// 测试空代理配置（应该正常工作）
	options := registry.NewOptions().
		SetRegistryURL(registry.RegistryUrlNpmMirror).
		SetProxy("")
	client := registry.NewRegistry(options)
	info, err := client.GetRegistryInformation(ctx)
	if err != nil {
		return fmt.Errorf("Proxy config failed: %w", err)
	}
	fmt.Printf("   空代理配置测试通过, Registry包数量: %d\n", info.DocCount)
	return nil
}
