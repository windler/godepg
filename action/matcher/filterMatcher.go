package matcher

import "regexp"

type FilterMatcher struct {
	matcher []string
	text    string
}

var _ Matcher = &FilterMatcher{}

func NewFilterMatcher(text string, matcher []string) *FilterMatcher {
	return &FilterMatcher{
		text:    text,
		matcher: matcher,
	}
}

func (f *FilterMatcher) Matches() bool {
	for _, m := range f.matcher {
		matches, _ := regexp.MatchString(m, f.text)
		if matches {
			return true
		}
	}
	return false
}
