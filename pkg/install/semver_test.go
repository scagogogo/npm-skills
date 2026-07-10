package install

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseVersion(t *testing.T) {
	v, err := parseVersion("1.2.3")
	assert.NoError(t, err)
	assert.Equal(t, Version{1, 2, 3}, v)

	v, _ = parseVersion("2")
	assert.Equal(t, Version{2, 0, 0}, v)

	v, _ = parseVersion("1.2.3-beta.1")
	assert.Equal(t, Version{1, 2, 3}, v)

	v, _ = parseVersion("1.2.3+build.5")
	assert.Equal(t, Version{1, 2, 3}, v)

	_, err = parseVersion("abc")
	assert.Error(t, err)

	_, err = parseVersion("")
	assert.Error(t, err)
}

func TestVersionCompare(t *testing.T) {
	assert.Equal(t, 0, Version{1, 2, 3}.Compare(Version{1, 2, 3}))
	assert.Equal(t, -1, Version{1, 2, 3}.Compare(Version{2, 0, 0}))
	assert.Equal(t, -1, Version{1, 2, 3}.Compare(Version{1, 3, 0}))
	assert.Equal(t, -1, Version{1, 2, 3}.Compare(Version{1, 2, 4}))
	assert.Equal(t, 1, Version{1, 2, 4}.Compare(Version{1, 2, 3}))
}

func TestSatisfiesCaret(t *testing.T) {
	assert.True(t, Satisfies("1.5.0", "^1.2.3"))
	assert.True(t, Satisfies("1.2.3", "^1.2.3"))
	assert.False(t, Satisfies("2.0.0", "^1.2.3"))
	assert.False(t, Satisfies("1.2.2", "^1.2.3"))
}

func TestSatisfiesTilde(t *testing.T) {
	assert.True(t, Satisfies("1.2.5", "~1.2.3"))
	assert.True(t, Satisfies("1.2.3", "~1.2.3"))
	assert.False(t, Satisfies("1.3.0", "~1.2.3"))
	assert.False(t, Satisfies("1.2.2", "~1.2.3"))
}

func TestSatisfiesOperators(t *testing.T) {
	assert.True(t, Satisfies("1.2.3", ">=1.2.0"))
	assert.False(t, Satisfies("1.1.0", ">=1.2.0"))
	assert.True(t, Satisfies("1.2.0", ">1.1.0"))
	assert.False(t, Satisfies("1.1.0", ">1.1.0"))
	assert.True(t, Satisfies("1.2.0", "<=1.2.0"))
	assert.False(t, Satisfies("1.3.0", "<1.3.0"))
	assert.True(t, Satisfies("1.2.3", "=1.2.3"))
	assert.False(t, Satisfies("1.2.4", "=1.2.3"))
}

func TestSatisfiesWildcardAndXRange(t *testing.T) {
	assert.True(t, Satisfies("99.0.0", "*"))
	assert.True(t, Satisfies("1.0.0", ""))
	assert.True(t, Satisfies("1.0.0", "latest"))
	assert.True(t, Satisfies("1.5.0", "1.x"))
	assert.False(t, Satisfies("2.0.0", "1.x"))
	assert.True(t, Satisfies("1.2.9", "1.2.x"))
	assert.False(t, Satisfies("1.3.0", "1.2.x"))
	assert.True(t, Satisfies("5.0.0", "x"))
}

func TestSatisfiesExact(t *testing.T) {
	assert.True(t, Satisfies("1.2.3", "1.2.3"))
	assert.False(t, Satisfies("1.2.4", "1.2.3"))
}

func TestSatisfiesInvalidCandidate(t *testing.T) {
	assert.False(t, Satisfies("not-a-version", "^1.2.3"))
	assert.False(t, Satisfies("1.2.3", "garbage-range"))
}

func TestPickBestVersion(t *testing.T) {
	versions := []string{"1.0.0", "1.2.0", "1.2.5", "2.0.0", "1.1.0"}
	assert.Equal(t, "1.2.5", PickBestVersion(versions, "^1.2.0"))
	assert.Equal(t, "1.2.5", PickBestVersion(versions, "~1.2.0"))
	assert.Equal(t, "1.0.0", PickBestVersion(versions, "1.0.0"))
	assert.Equal(t, "", PickBestVersion(versions, "^3.0.0"))
	assert.Equal(t, "", PickBestVersion(nil, "^1.0.0"))
	assert.Equal(t, "2.0.0", PickBestVersion(versions, "*"))
}

func TestVersionString(t *testing.T) {
	assert.Equal(t, "1.2.3", Version{1, 2, 3}.String())
}

func TestParseVersionErrorBranches(t *testing.T) {
	_, err := parseVersion("1.abc.3")
	assert.Error(t, err)
	_, err = parseVersion("1.2.abc")
	assert.Error(t, err)
}

func TestSatisfiesOperatorErrorBranches(t *testing.T) {
	assert.False(t, Satisfies("1.2.3", ">=abc"))
	assert.False(t, Satisfies("1.2.3", "^abc"))
	assert.False(t, Satisfies("1.2.3", "~abc"))
}

func TestSatisfiesXRangeMajorWildcard(t *testing.T) {
	assert.True(t, Satisfies("3.0.0", "X"))
	assert.True(t, Satisfies("3.0.0", "x.2.3"))
	assert.False(t, Satisfies("2.0.0", "1.abc.x"))
}

func TestPickBestVersionSkipsUnparseable(t *testing.T) {
	versions := []string{"not-a-version", "1.2.5"}
	assert.Equal(t, "1.2.5", PickBestVersion(versions, "^1.0.0"))
}
