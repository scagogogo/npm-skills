package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/scagogogo/npm-skills/pkg/models"
)

// GetOrg 获取组织详情
//
// 需要认证Token。返回指定组织的信息。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - orgName: 组织名称
//
// 返回值:
//   - *models.Organization: 组织信息
//   - error: 如果请求失败则返回错误
func (x *Registry) GetOrg(ctx context.Context, orgName string) (*models.Organization, error) {
	if err := x.requireToken(); err != nil {
		return nil, err
	}

	targetUrl := fmt.Sprintf("%s/-/org/%s", x.options.RegistryURL, orgName)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get org '%s': %w", orgName, err)
	}
	return unmarshalJson[*models.Organization](bytes)
}

// CreateOrg 创建组织
//
// 需要认证Token。创建一个新的 NPM 组织。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - orgName: 组织名称
//
// 返回值:
//   - *models.Organization: 创建的组织信息
//   - error: 如果请求失败则返回错误
func (x *Registry) CreateOrg(ctx context.Context, orgName string) (*models.Organization, error) {
	if err := x.requireToken(); err != nil {
		return nil, err
	}

	targetUrl := fmt.Sprintf("%s/-/org/%s", x.options.RegistryURL, orgName)
	payload := &models.OrgCreation{Name: orgName}
	bytes, err := x.sendJSON(ctx, http.MethodPut, targetUrl, payload, http.StatusOK, http.StatusCreated)
	if err != nil {
		return nil, fmt.Errorf("failed to create org '%s': %w", orgName, err)
	}
	return unmarshalJson[*models.Organization](bytes)
}

// DeleteOrg 删除组织
//
// 需要认证Token。删除指定的 NPM 组织。这是一个不可逆操作。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - orgName: 组织名称
//
// 返回值:
//   - error: 如果请求失败则返回错误
func (x *Registry) DeleteOrg(ctx context.Context, orgName string) error {
	if err := x.requireToken(); err != nil {
		return err
	}

	targetUrl := fmt.Sprintf("%s/-/org/%s", x.options.RegistryURL, orgName)
	_, err := x.sendRequest(ctx, http.MethodDelete, targetUrl, nil, http.StatusOK, http.StatusNoContent)
	if err != nil {
		return fmt.Errorf("failed to delete org '%s': %w", orgName, err)
	}
	return nil
}

// ListOrgMembers 列出组织成员
//
// 需要认证Token。返回指定组织的所有成员用户名列表。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - orgName: 组织名称
//
// 返回值:
//   - []string: 成员用户名列表
//   - error: 如果请求失败则返回错误
func (x *Registry) ListOrgMembers(ctx context.Context, orgName string) ([]string, error) {
	if err := x.requireToken(); err != nil {
		return nil, err
	}

	targetUrl := fmt.Sprintf("%s/-/org/%s/member", x.options.RegistryURL, orgName)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to list members of org '%s': %w", orgName, err)
	}
	return unmarshalJson[[]string](bytes)
}

// AddOrgMember 添加组织成员
//
// 需要认证Token。将指定用户添加到组织中。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - orgName: 组织名称
//   - username: 要添加的用户名
//
// 返回值:
//   - error: 如果请求失败则返回错误
func (x *Registry) AddOrgMember(ctx context.Context, orgName, username string) error {
	if err := x.requireToken(); err != nil {
		return err
	}

	targetUrl := fmt.Sprintf("%s/-/org/%s/member/%s", x.options.RegistryURL, orgName, username)
	_, err := x.sendRequest(ctx, http.MethodPut, targetUrl, nil, http.StatusOK, http.StatusCreated)
	if err != nil {
		return fmt.Errorf("failed to add member '%s' to org '%s': %w", username, orgName, err)
	}
	return nil
}

// RemoveOrgMember 移除组织成员
//
// 需要认证Token。将指定用户从组织中移除。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - orgName: 组织名称
//   - username: 要移除的用户名
//
// 返回值:
//   - error: 如果请求失败则返回错误
func (x *Registry) RemoveOrgMember(ctx context.Context, orgName, username string) error {
	if err := x.requireToken(); err != nil {
		return err
	}

	targetUrl := fmt.Sprintf("%s/-/org/%s/member/%s", x.options.RegistryURL, orgName, username)
	_, err := x.sendRequest(ctx, http.MethodDelete, targetUrl, nil, http.StatusOK, http.StatusNoContent)
	if err != nil {
		return fmt.Errorf("failed to remove member '%s' from org '%s': %w", username, orgName, err)
	}
	return nil
}

// ListOrgPackages 列出组织拥有的包
//
// 需要认证Token。返回指定组织拥有的所有包名称列表。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - orgName: 组织名称
//
// 返回值:
//   - []string: 包名称列表
//   - error: 如果请求失败则返回错误
func (x *Registry) ListOrgPackages(ctx context.Context, orgName string) ([]string, error) {
	if err := x.requireToken(); err != nil {
		return nil, err
	}

	targetUrl := fmt.Sprintf("%s/-/org/%s/package", x.options.RegistryURL, orgName)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to list packages of org '%s': %w", orgName, err)
	}
	return unmarshalJson[[]string](bytes)
}

// ListTeams 列出组织中的团队
//
// 需要认证Token。返回指定组织中的所有团队列表。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - orgName: 组织名称
//
// 返回值:
//   - []models.Team: 团队列表
//   - error: 如果请求失败则返回错误
func (x *Registry) ListTeams(ctx context.Context, orgName string) ([]models.Team, error) {
	if err := x.requireToken(); err != nil {
		return nil, err
	}

	targetUrl := fmt.Sprintf("%s/-/org/%s/team", x.options.RegistryURL, orgName)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to list teams of org '%s': %w", orgName, err)
	}

	var result struct {
		Objects []models.Team `json:"objects"`
		Total   int           `json:"total"`
	}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse teams: %w", err)
	}
	return result.Objects, nil
}

// CreateTeam 创建团队
//
// 需要认证Token。在指定组织中创建一个新的团队。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - orgName: 组织名称
//   - teamName: 团队名称
//
// 返回值:
//   - *models.Team: 创建的团队信息
//   - error: 如果请求失败则返回错误
func (x *Registry) CreateTeam(ctx context.Context, orgName, teamName string) (*models.Team, error) {
	if err := x.requireToken(); err != nil {
		return nil, err
	}

	targetUrl := fmt.Sprintf("%s/-/org/%s/team/%s", x.options.RegistryURL, orgName, teamName)
	payload := &models.TeamCreation{Name: teamName}
	bytes, err := x.sendJSON(ctx, http.MethodPut, targetUrl, payload, http.StatusOK, http.StatusCreated)
	if err != nil {
		return nil, fmt.Errorf("failed to create team '%s' in org '%s': %w", teamName, orgName, err)
	}
	return unmarshalJson[*models.Team](bytes)
}

// DeleteTeam 删除团队
//
// 需要认证Token。删除指定组织中的团队。这是一个不可逆操作。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - orgName: 组织名称
//   - teamName: 团队名称
//
// 返回值:
//   - error: 如果请求失败则返回错误
func (x *Registry) DeleteTeam(ctx context.Context, orgName, teamName string) error {
	if err := x.requireToken(); err != nil {
		return err
	}

	targetUrl := fmt.Sprintf("%s/-/org/%s/team/%s", x.options.RegistryURL, orgName, teamName)
	_, err := x.sendRequest(ctx, http.MethodDelete, targetUrl, nil, http.StatusOK, http.StatusNoContent)
	if err != nil {
		return fmt.Errorf("failed to delete team '%s' in org '%s': %w", teamName, orgName, err)
	}
	return nil
}

// ListTeamMembers 列出团队成员
//
// 需要认证Token。返回指定团队的所有成员用户名列表。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - orgName: 组织名称
//   - teamName: 团队名称
//
// 返回值:
//   - []string: 成员用户名列表
//   - error: 如果请求失败则返回错误
func (x *Registry) ListTeamMembers(ctx context.Context, orgName, teamName string) ([]string, error) {
	if err := x.requireToken(); err != nil {
		return nil, err
	}

	targetUrl := fmt.Sprintf("%s/-/org/%s/team/%s/member", x.options.RegistryURL, orgName, teamName)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to list members of team '%s/%s': %w", orgName, teamName, err)
	}
	return unmarshalJson[[]string](bytes)
}

// AddTeamMember 添加团队成员
//
// 需要认证Token。将指定用户添加到团队中。
// 用户必须已经是组织的成员才能被添加到团队。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - orgName: 组织名称
//   - teamName: 团队名称
//   - username: 要添加的用户名
//
// 返回值:
//   - error: 如果请求失败则返回错误
func (x *Registry) AddTeamMember(ctx context.Context, orgName, teamName, username string) error {
	if err := x.requireToken(); err != nil {
		return err
	}

	targetUrl := fmt.Sprintf("%s/-/org/%s/team/%s/member/%s", x.options.RegistryURL, orgName, teamName, username)
	_, err := x.sendRequest(ctx, http.MethodPut, targetUrl, nil, http.StatusOK, http.StatusCreated)
	if err != nil {
		return fmt.Errorf("failed to add member '%s' to team '%s/%s': %w", username, orgName, teamName, err)
	}
	return nil
}

// RemoveTeamMember 移除团队成员
//
// 需要认证Token。将指定用户从团队中移除。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - orgName: 组织名称
//   - teamName: 团队名称
//   - username: 要移除的用户名
//
// 返回值:
//   - error: 如果请求失败则返回错误
func (x *Registry) RemoveTeamMember(ctx context.Context, orgName, teamName, username string) error {
	if err := x.requireToken(); err != nil {
		return err
	}

	targetUrl := fmt.Sprintf("%s/-/org/%s/team/%s/member/%s", x.options.RegistryURL, orgName, teamName, username)
	_, err := x.sendRequest(ctx, http.MethodDelete, targetUrl, nil, http.StatusOK, http.StatusNoContent)
	if err != nil {
		return fmt.Errorf("failed to remove member '%s' from team '%s/%s': %w", username, orgName, teamName, err)
	}
	return nil
}

// ListTeamPackages 列出团队有权限的包
//
// 需要认证Token。返回指定团队有权访问的所有包名称列表。
//
// 参数:
//   - ctx: 上下文，可用于取消请求或设置超时
//   - orgName: 组织名称
//   - teamName: 团队名称
//
// 返回值:
//   - []string: 包名称列表
//   - error: 如果请求失败则返回错误
func (x *Registry) ListTeamPackages(ctx context.Context, orgName, teamName string) ([]string, error) {
	if err := x.requireToken(); err != nil {
		return nil, err
	}

	targetUrl := fmt.Sprintf("%s/-/org/%s/team/%s/package", x.options.RegistryURL, orgName, teamName)
	bytes, err := x.getBytes(ctx, targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to list packages of team '%s/%s': %w", orgName, teamName, err)
	}
	return unmarshalJson[[]string](bytes)
}
