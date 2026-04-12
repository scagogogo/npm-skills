package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/scagogogo/npm-crawler/pkg/registry"
)

func main() {
	// 使用 CNPM 镜像
	options := registry.NewOptions().SetRegistryURL(registry.RegistryUrlCnpm)
	r := registry.NewRegistry(options)

	ctx := context.Background()
	packageName := "axios"
	version := "1.0.0"
	destDir := "/tmp/npm-tarballs"

	// 确保目标目录存在
	if err := os.MkdirAll(destDir, 0755); err != nil {
		panic(fmt.Sprintf("创建目录失败: %v", err))
	}
	destPath := filepath.Join(destDir, fmt.Sprintf("%s-%s.tgz", packageName, version))

	fmt.Printf("开始下载 %s@%s...\n", packageName, version)
	fmt.Printf("目标路径: %s\n", destPath)

	// 删除已存在的文件（如果有）
	os.Remove(destPath)

	// 下载 tarball
	if err := r.DownloadTarball(ctx, packageName, version, destPath); err != nil {
		panic(fmt.Sprintf("下载失败: %v", err))
	}

	// 验证文件
	info, err := os.Stat(destPath)
	if err != nil {
		panic(fmt.Sprintf("文件状态检查失败: %v", err))
	}

	fmt.Printf("下载成功！\n")
	fmt.Printf("文件大小: %d bytes\n", info.Size())
	fmt.Printf("文件路径: %s\n", destPath)
}
