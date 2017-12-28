package matcher

import "regexp"

type SubPackageMatcher struct {
	base string
	text string
}

var _ Matcher = &FilterMatcher{}

func NewSubPackageMatcher(base, text string) *SubPackageMatcher {
	return &SubPackageMatcher{
		text: text,
		base: base,
	}
}

func (f *SubPackageMatcher) Matches() bool {
	matches, _ := regexp.MatchString(f.base+".*", f.text)
	return matches
}
