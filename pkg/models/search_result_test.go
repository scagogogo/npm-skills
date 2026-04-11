package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchResult(t *testing.T) {
	// 创建一个测试搜索结果
	result := &SearchResult{
		Objects: []SearchObject{
			{
				Package: SearchPackage{
					Name:        "react",
					Version:     "18.2.0",
					Description: "React is a JavaScript library for building user interfaces.",
					Keywords:    []string{"react", "ui", "library"},
					Date:        "2022-06-14T18:58:43.231Z",
					Links: Links{
						NPM:        "https://www.npmjs.com/package/react",
						Homepage:   "https://reactjs.org/",
						Repository: "https://github.com/facebook/react",
						Bugs:       "https://github.com/facebook/react/issues",
					},
					Author: &User{
						Name:  "facebook",
						Email: "react@fb.com",
					},
					Publisher: &User{
						Name:  "react-team",
						Email: "react@fb.com",
					},
					Maintainers: []*User{
						{
							Name:  "gaearon",
							Email: "dan.abramov@gmail.com",
						},
					},
				},
				Score: Score{
					Final: 0.95,
					Detail: ScoreDetail{
						Quality:     0.98,
						Popularity:  0.99,
						Maintenance: 0.95,
					},
				},
				SearchScore: 0.95,
			},
		},
		Total: 1,
		Time:  "5ms",
	}

	// 基本属性测试
	assert.Equal(t, 1, len(result.Objects))
	assert.Equal(t, 1, result.Total)
	assert.Equal(t, "5ms", result.Time)

	// 测试第一个包对象
	pkg := result.Objects[0]
	assert.Equal(t, "react", pkg.Package.Name)
	assert.Equal(t, "18.2.0", pkg.Package.Version)
	assert.Contains(t, pkg.Package.Keywords, "react")
	assert.Equal(t, 0.95, pkg.SearchScore)

	// 测试评分信息
	assert.Equal(t, 0.95, pkg.Score.Final)
	assert.Equal(t, 0.98, pkg.Score.Detail.Quality)
	assert.Equal(t, 0.99, pkg.Score.Detail.Popularity)
	assert.Equal(t, 0.95, pkg.Score.Detail.Maintenance)

	// 测试链接信息
	assert.Equal(t, "https://www.npmjs.com/package/react", pkg.Package.Links.NPM)
	assert.Equal(t, "https://reactjs.org/", pkg.Package.Links.Homepage)

	// 测试作者信息
	assert.NotNil(t, pkg.Package.Author)
	assert.Equal(t, "facebook", pkg.Package.Author.Name)

	// 测试维护者信息
	assert.Equal(t, 1, len(pkg.Package.Maintainers))
	assert.Equal(t, "gaearon", pkg.Package.Maintainers[0].Name)
}

func TestSearchResultToJsonString(t *testing.T) {
	result := &SearchResult{
		Objects: []SearchObject{
			{
				Package: SearchPackage{
					Name:        "lodash",
					Version:     "4.17.21",
					Description: "Lodash modular utilities.",
					Keywords:    []string{"lodash", "utility"},
				},
				Score: Score{
					Final: 0.85,
					Detail: ScoreDetail{
						Quality:     0.80,
						Popularity:  0.90,
						Maintenance: 0.85,
					},
				},
				SearchScore: 0.85,
			},
		},
		Total: 1,
		Time:  "3ms",
	}

	jsonStr := result.ToJsonString()
	assert.NotEmpty(t, jsonStr)
	assert.Contains(t, jsonStr, "lodash")
	assert.Contains(t, jsonStr, "4.17.21")

	// 测试反序列化
	var parsed SearchResult
	err := json.Unmarshal([]byte(jsonStr), &parsed)
	assert.NoError(t, err)
	assert.Equal(t, result.Total, parsed.Total)
	assert.Equal(t, result.Objects[0].Package.Name, parsed.Objects[0].Package.Name)
}

func TestSearchPackage(t *testing.T) {
	pkg := SearchPackage{
		Name:        "express",
		Scope:       "",
		Version:     "4.18.2",
		Description: "Fast, unopinionated, minimalist web framework",
		Keywords:    []string{"express", "framework", "web"},
		Date:        "2022-10-08T17:48:22.853Z",
		Links: Links{
			NPM:        "https://www.npmjs.com/package/express",
			Homepage:   "http://expressjs.com/",
			Repository: "https://github.com/expressjs/express",
			Bugs:       "https://github.com/expressjs/express/issues",
		},
		Author: &User{
			Name:  "tjholowaychuk",
			Email: "tj@vision-media.ca",
		},
	}

	assert.Equal(t, "express", pkg.Name)
	assert.Equal(t, "4.18.2", pkg.Version)
	assert.Contains(t, pkg.Keywords, "framework")
	assert.Equal(t, "tjholowaychuk", pkg.Author.Name)
}

func TestScore(t *testing.T) {
	score := Score{
		Final: 0.89,
		Detail: ScoreDetail{
			Quality:     0.85,
			Popularity:  0.95,
			Maintenance: 0.87,
		},
	}

	assert.Equal(t, 0.89, score.Final)
	assert.Equal(t, 0.85, score.Detail.Quality)
	assert.Equal(t, 0.95, score.Detail.Popularity)
	assert.Equal(t, 0.87, score.Detail.Maintenance)
}

func TestLinks(t *testing.T) {
	links := Links{
		NPM:        "https://www.npmjs.com/package/test",
		Homepage:   "https://example.com",
		Repository: "https://github.com/user/test",
		Bugs:       "https://github.com/user/test/issues",
	}

	assert.Equal(t, "https://www.npmjs.com/package/test", links.NPM)
	assert.Equal(t, "https://example.com", links.Homepage)
	assert.Equal(t, "https://github.com/user/test", links.Repository)
	assert.Equal(t, "https://github.com/user/test/issues", links.Bugs)
}

func TestUser(t *testing.T) {
	user := &User{
		Name:  "testuser",
		Email: "test@example.com",
	}

	assert.Equal(t, "testuser", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestSearchResultToJsonStringErrorBranch(t *testing.T) {
	// 测试 SearchResult.ToJsonString 的 error 分支
	// 当 TestMarshalErr 字段非 nil 时，json.Marshal 会返回错误
	result := &SearchResult{
		Objects: []SearchObject{
			{
				Package: SearchPackage{
					Name:        "test",
					Version:     "1.0.0",
					Description: "test package",
				},
				Score: Score{
					Final: 0.5,
					Detail: ScoreDetail{
						Quality:    0.5,
						Popularity: 0.5,
					},
				},
				SearchScore: 0.5,
			},
		},
		Total:          1,
		Time:           "1ms",
		TestMarshalErr: &testMarshalFailType{},
	}

	jsonStr := result.ToJsonString()
	assert.Contains(t, jsonStr, "forced marshal error for testing")
}
