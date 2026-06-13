package registry

import (
	"context"
	"fmt"
	"net/http"

	"github.com/scagogogo/npm-skills/pkg/models"
)

// Login 登录NPM Registry，返回认证Token
//
// 通过用户名和密码向 Registry 进行身份认证。
// 成功后返回包含 Token 的 LoginResult，可后续用于 SetToken() 设置认证。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - name: 用户名
//   - password: 密码
//
// 返回值:
//   - *models.LoginResult: 登录结果，包含认证Token
//   - error: 如果登录失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	result, err := registry.Login(ctx, "myuser", "mypassword")
//	if err != nil {
//	    // 处理错误
//	}
//	// 使用返回的Token创建认证客户端
//	authRegistry := NewRegistry(NewOptions().SetToken(result.Token))
func (x *Registry) Login(ctx context.Context, name, password string) (*models.LoginResult, error) {
	targetUrl := fmt.Sprintf("%s/-/user/org.couchdb.user:%s", x.options.RegistryURL, name)

	payload := map[string]string{
		"name":     name,
		"password": password,
	}

	bytes, err := x.sendJSON(ctx, http.MethodPut, targetUrl, payload, http.StatusOK, http.StatusCreated)
	if err != nil {
		return nil, fmt.Errorf("login failed for user '%s': %w", name, err)
	}

	return unmarshalJson[*models.LoginResult](bytes)
}

// CreateUser 创建新用户（注册）
//
// 向 Registry 注册一个新用户账号。成功后返回包含 Token 的 LoginResult。
// 创建用户不需要预先设置 Token。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - user: 用户创建信息，包含用户名、密码、邮箱等
//
// 返回值:
//   - *models.LoginResult: 注册结果，包含认证Token
//   - error: 如果注册失败则返回错误
//
// 使用示例:
//
//	registry := NewRegistry()
//	ctx := context.Background()
//	result, err := registry.CreateUser(ctx, &models.UserCreation{
//	    Name:     "myuser",
//	    Password: "mypassword",
//	    Email:    "user@example.com",
//	})
//	if err != nil {
//	    // 处理错误
//	}
//	fmt.Println("Token:", result.Token)
func (x *Registry) CreateUser(ctx context.Context, user *models.UserCreation) (*models.LoginResult, error) {
	// 设置 CouchDB 用户文档的固定字段
	user.ID = fmt.Sprintf("org.couchdb.user:%s", user.Name)
	user.Type = "user"

	targetUrl := fmt.Sprintf("%s/-/user/%s", x.options.RegistryURL, user.ID)

	bytes, err := x.sendJSON(ctx, http.MethodPut, targetUrl, user, http.StatusOK, http.StatusCreated)
	if err != nil {
		return nil, fmt.Errorf("failed to create user '%s': %w", user.Name, err)
	}

	return unmarshalJson[*models.LoginResult](bytes)
}

// GetUser 获取用户信息
//
// 需要认证Token（部分 Registry 可能允许未认证访问）。
// 返回指定用户的公开信息。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - name: 用户名
//
// 返回值:
//   - *models.UserProfile: 用户信息
//   - error: 如果请求失败则返回错误
//
// 使用示例:
//
//	options := NewOptions().SetToken("npm_xxxxx")
//	registry := NewRegistry(options)
//	ctx := context.Background()
//	profile, err := registry.GetUser(ctx, "myuser")
//	if err != nil {
//	    // 处理错误
//	}
//	fmt.Println("Name:", profile.Name)
//	fmt.Println("Email:", profile.Email)
func (x *Registry) GetUser(ctx context.Context, name string) (*models.UserProfile, error) {
	targetUrl := fmt.Sprintf("%s/-/user/org.couchdb.user:%s", x.options.RegistryURL, name)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get user '%s': %w", name, err)
	}
	return unmarshalJson[*models.UserProfile](bytes)
}
