package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubPackageMatcher(t *testing.T) {
	assert.False(t, NewSubPackageMatcher("my-base", "os").Matches())
	assert.True(t, NewSubPackageMatcher("my-base", "my-base/asd").Matches())
	assert.True(t, NewSubPackageMatcher("my-base", "my-base/ddd/eee").Matches())
}
