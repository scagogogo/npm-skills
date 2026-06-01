package registry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 创建模拟的 NPM Registry HTTP 服务器
func setupMockRegistryServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 提供 registry 根路径响应
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"db_name": "registry",
				"doc_count": 1000,
				"doc_del_count": 10,
				"update_seq": 5000,
				"compact_running": false,
				"sizes": {"file": 1000000, "active": 900000}
			}`))
			return
		}

		// 提供特定包的响应
		if r.URL.Path == "/react" || r.URL.Path == "/lodash" ||
			r.URL.Path == "/express" || r.URL.Path == "/vue" ||
			r.URL.Path == "/webpack" || r.URL.Path == "/typescript" ||
			r.URL.Path == "/react-dom" {
			packageName := r.URL.Path[1:] // 去掉开头的斜杠
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"_id": "` + packageName + `",
				"name": "` + packageName + `",
				"description": "Test package",
				"dist-tags": {"latest": "1.0.0"},
				"versions": {
					"1.0.0": {
						"name": "` + packageName + `", 
						"version": "1.0.0",
						"dist": {
							"shasum": "abc123", 
							"tarball": "https://registry.npmjs.org/` + packageName + `/-/` + packageName + `-1.0.0.tgz"
						}
					}
				}
			}`))
			return
		}

		// 默认返回 404
		w.WriteHeader(http.StatusNotFound)
	}))
}

func TestNewYarnRegistry(t *testing.T) {
	// 创建模拟服务器
	server := setupMockRegistryServer()
	defer server.Close()

	// 创建使用模拟服务器的 Registry
	r := NewRegistry(NewOptions().SetRegistryURL(server.URL))

	// 验证基本配置
	assert.NotNil(t, r)

	// 用超时上下文测试基本功能
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 测试获取 Registry 信息
	info, err := r.GetRegistryInformation(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, "registry", info.DbName)

	// 测试获取知名包信息
	pkg, err := r.GetPackageInformation(ctx, "react")
	assert.Nil(t, err)
	assert.NotNil(t, pkg)
	assert.Equal(t, "react", pkg.Name)
}

func TestAllRegistryImplementations(t *testing.T) {
	// 创建模拟服务器
	server := setupMockRegistryServer()
	defer server.Close()

	testCases := []struct {
		name        string
		createFn    func() *Registry
		expectedURL string
		packageName string
	}{
		{"YarnRegistry", NewYarnRegistry, RegistryUrlYarn, "react"},
		{"CnpmRegistry", NewCnpmRegistry, RegistryUrlCnpm, "lodash"},
		{"HuaWeiCloudRegistry", NewHuaWeiCloudRegistry, RegistryUrlHuaWeiCloud, "express"},
		{"NpmMirrorRegistry", NewNpmMirrorRegistry, RegistryUrlNpmMirror, "vue"},
		{"NpmjsComRegistry", NewNpmjsComRegistry, RegistryUrlNpmjsCom, "webpack"},
		{"TaoBaoRegistry", NewTaoBaoRegistry, RegistryUrlTaoBao, "typescript"},
		{"TencentRegistry", NewTencentRegistry, RegistryUrlTencent, "react-dom"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建 Registry 实例
			registry := tc.createFn()

			// 基础配置测试
			assert.NotNil(t, registry)
			assert.Equal(t, tc.expectedURL, registry.GetOptions().RegistryURL)

			// 使用本地模拟服务器覆盖 URL 进行接口测试
			registry = NewRegistry(NewOptions().SetRegistryURL(server.URL))

			// 使用超时上下文
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// 获取 Registry 信息
			info, err := registry.GetRegistryInformation(ctx)
			assert.Nil(t, err)
			assert.NotNil(t, info)
			assert.Equal(t, "registry", info.DbName)

			// 获取包信息
			pkg, err := registry.GetPackageInformation(ctx, tc.packageName)
			assert.Nil(t, err)
			assert.NotNil(t, pkg)
			assert.Equal(t, tc.packageName, pkg.Name)
		})
	}
}

func TestRegistryWithProxy(t *testing.T) {
	// 创建带有代理的 Registry
	proxyRegistry := NewRegistry(
		NewOptions().
			SetRegistryURL(RegistryUrlYarn).
			SetProxy("http://your-proxy-server:8080"),
	)

	assert.NotNil(t, proxyRegistry)
	assert.Equal(t, RegistryUrlYarn, proxyRegistry.GetOptions().RegistryURL)
	assert.Equal(t, "http://your-proxy-server:8080", proxyRegistry.GetOptions().Proxy)
}

func TestMirrorRegistryCreation(t *testing.T) {
	mirrors := []struct {
		name        string
		registry    *Registry
		expectedURL string
	}{
		{"Yarn", NewYarnRegistry(), RegistryUrlYarn},
		{"CNPM", NewCnpmRegistry(), RegistryUrlCnpm},
		{"HuaWeiCloud", NewHuaWeiCloudRegistry(), RegistryUrlHuaWeiCloud},
		{"NpmMirror", NewNpmMirrorRegistry(), RegistryUrlNpmMirror},
		{"NpmjsCom", NewNpmjsComRegistry(), RegistryUrlNpmjsCom},
		{"TaoBao", NewTaoBaoRegistry(), RegistryUrlTaoBao},
		{"Tencent", NewTencentRegistry(), RegistryUrlTencent},
	}

	for _, m := range mirrors {
		t.Run(m.name, func(t *testing.T) {
			assert.NotNil(t, m.registry)
			assert.Equal(t, m.expectedURL, m.registry.GetOptions().RegistryURL)
		})
	}
}

func TestListMirrors(t *testing.T) {
	mirrors := ListMirrors()

	// 测试返回的镜像数量
	assert.Equal(t, 8, len(mirrors), "应该返回 8 个镜像源")

	// 测试第一个镜像源是 official
	assert.Equal(t, "official", mirrors[0].Name)
	assert.Equal(t, DefaultRegistryURL, mirrors[0].URL)
	assert.Equal(t, "Global", mirrors[0].Region)

	// 测试每个镜像源都有必要字段
	for _, m := range mirrors {
		assert.NotEmpty(t, m.Name, "镜像源名称不能为空")
		assert.NotEmpty(t, m.URL, "镜像源 URL 不能为空")
		assert.NotEmpty(t, m.Region, "镜像源区域不能为空")
		assert.NotEmpty(t, m.Description, "镜像源描述不能为空")
	}

	// 测试镜像源 URL 与常量一致
	urlMap := map[string]string{
		"official":   DefaultRegistryURL,
		"taobao":     RegistryUrlTaoBao,
		"npm-mirror": RegistryUrlNpmMirror,
		"huawei":     RegistryUrlHuaWeiCloud,
		"tencent":    RegistryUrlTencent,
		"cnpm":       RegistryUrlCnpm,
		"yarn":       RegistryUrlYarn,
		"npmjscom":   RegistryUrlNpmjsCom,
	}
	for _, m := range mirrors {
		expectedURL, ok := urlMap[m.Name]
		assert.True(t, ok, "镜像源 %s 不在 URL 映射中", m.Name)
		assert.Equal(t, expectedURL, m.URL, "镜像源 %s 的 URL 与常量不一致", m.Name)
	}
}
