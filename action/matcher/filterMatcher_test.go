package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterMatcher(t *testing.T) {
	assert.True(t, NewFilterMatcher("my text", []string{"text"}).Matches())
	assert.False(t, NewFilterMatcher("my text", []string{"textt"}).Matches())
	assert.True(t, NewFilterMatcher("my text", []string{"t", "text"}).Matches())
}
