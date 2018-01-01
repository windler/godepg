package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoPackagesMatcher(t *testing.T) {
	assert.False(t, NewGoPackagesMatcher("my-pkg").Matches())
	assert.True(t, NewGoPackagesMatcher("os").Matches())
	assert.True(t, NewGoPackagesMatcher("fmt").Matches())
}
