package matcher

import (
	"regexp"
)

//FilterMatcher matches when a text contains a givven pattern
type FilterMatcher struct {
	matcher []string
	text    string
}

//NewFilterMatcher create a new FilterMatcher
func NewFilterMatcher(text string, matcher []string) *FilterMatcher {
	return &FilterMatcher{
		text:    text,
		matcher: matcher,
	}
}

//Matches applies the filter
func (f *FilterMatcher) Matches() bool {
	for _, m := range f.matcher {
		matches, _ := regexp.MatchString(m, f.text)
		if matches {
			return true
		}
	}
	return false
}
