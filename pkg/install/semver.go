package install

import (
	"fmt"
	"strconv"
	"strings"
)

// Version 表示一个语义化版本号（major.minor.patch）。
type Version struct {
	Major int
	Minor int
	Patch int
}

// String 返回 "major.minor.patch" 格式字符串。
func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// parseVersion 解析 "1.2.3" 形式的版本号，忽略 prerelease/build 后缀。
// 缺失的部分补 0（"1" → 1.0.0）。
func parseVersion(s string) (Version, error) {
	if i := strings.IndexAny(s, "-+"); i >= 0 {
		s = s[:i]
	}
	parts := strings.Split(s, ".")
	if len(parts) == 0 || parts[0] == "" {
		return Version{}, fmt.Errorf("empty version")
	}
	v := Version{}
	var err error
	if v.Major, err = strconv.Atoi(parts[0]); err != nil {
		return Version{}, fmt.Errorf("invalid major: %s", parts[0])
	}
	if len(parts) >= 2 {
		if v.Minor, err = strconv.Atoi(parts[1]); err != nil {
			return Version{}, fmt.Errorf("invalid minor: %s", parts[1])
		}
	}
	if len(parts) >= 3 {
		if v.Patch, err = strconv.Atoi(parts[2]); err != nil {
			return Version{}, fmt.Errorf("invalid patch: %s", parts[2])
		}
	}
	return v, nil
}

// Compare 比较 v 与 other：返回 -1/0/1。
func (v Version) Compare(other Version) int {
	if v.Major != other.Major {
		if v.Major < other.Major {
			return -1
		}
		return 1
	}
	if v.Minor != other.Minor {
		if v.Minor < other.Minor {
			return -1
		}
		return 1
	}
	if v.Patch != other.Patch {
		if v.Patch < other.Patch {
			return -1
		}
		return 1
	}
	return 0
}

// satisfiesRange 判断候选版本是否满足 npm 范围表达式。
// 支持：*（任意）、^1.2.3、~1.2.3、>=1.2.3、>1.2.3、<=1.2.3、<1.2.3、=1.2.3、1.2.3（精确）、
// 1.2.x / 1.x / x（x-range）。不支持 || 复合范围与 prerelease。
func satisfiesRange(candidate Version, rangeExpr string) bool {
	expr := strings.TrimSpace(rangeExpr)
	if expr == "" || expr == "*" || expr == "latest" || expr == "x" || expr == "X" {
		return true
	}
	if strings.HasPrefix(expr, "^") {
		base := strings.TrimPrefix(expr, "^")
		v, err := parseVersion(base)
		if err != nil {
			return false
		}
		return candidate.Compare(v) >= 0 && candidate.Major == v.Major
	}
	if strings.HasPrefix(expr, "~") {
		base := strings.TrimPrefix(expr, "~")
		v, err := parseVersion(base)
		if err != nil {
			return false
		}
		return candidate.Compare(v) >= 0 && candidate.Major == v.Major && candidate.Minor == v.Minor
	}
	for _, op := range []string{">=", "<=", ">", "<", "="} {
		if strings.HasPrefix(expr, op) {
			base := strings.TrimSpace(strings.TrimPrefix(expr, op))
			v, err := parseVersion(base)
			if err != nil {
				return false
			}
			c := candidate.Compare(v)
			switch op {
			case ">=":
				return c >= 0
			case "<=":
				return c <= 0
			case ">":
				return c > 0
			case "<":
				return c < 0
			case "=":
				return c == 0
			}
		}
	}
	if strings.Contains(expr, "x") || strings.Contains(expr, "X") {
		parts := strings.Split(expr, ".")
		v := Version{}
		xIndex := -1
		for i, p := range parts {
			if p == "x" || p == "X" {
				xIndex = i
				break
			}
			n, err := strconv.Atoi(p)
			if err != nil {
				return false
			}
			switch i {
			case 0:
				v.Major = n
			case 1:
				v.Minor = n
			}
		}
		if xIndex == 0 {
			return true
		}
		if candidate.Major != v.Major {
			return false
		}
		if xIndex == 1 {
			return true
		}
		return candidate.Minor == v.Minor
	}
	v, err := parseVersion(expr)
	if err != nil {
		return false
	}
	return candidate.Compare(v) == 0
}

// Satisfies 判断 candidate 版本字符串是否满足 rangeExpr。
func Satisfies(candidate, rangeExpr string) bool {
	c, err := parseVersion(candidate)
	if err != nil {
		return false
	}
	return satisfiesRange(c, rangeExpr)
}

// PickBestVersion 从 versions 列表中选出满足 rangeExpr 的最高版本。
// 空列表或无匹配返回空字符串。
func PickBestVersion(versions []string, rangeExpr string) string {
	var best *Version
	var bestStr string
	for _, vs := range versions {
		if !Satisfies(vs, rangeExpr) {
			continue
		}
		v, _ := parseVersion(vs) // Satisfies 已验证可解析，err 不可达
		if best == nil || v.Compare(*best) > 0 {
			b := v
			best = &b
			bestStr = vs
		}
	}
	return bestStr
}
