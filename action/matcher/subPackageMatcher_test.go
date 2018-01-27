package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubPackageMatcher(t *testing.T) {
	assert.True(t, NewSubPackageMatcher("my-base", "os").Matches())
	assert.False(t, NewSubPackageMatcher("my-base", "my-base/asd").Matches())
	assert.False(t, NewSubPackageMatcher("my-base", "my-base/ddd/eee").Matches())
}
