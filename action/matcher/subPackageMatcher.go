package matcher

import "regexp"

//SubPackageMatcher checks wether the given text does not begin with given base
type SubPackageMatcher struct {
	base string
	text string
}

//NewSubPackageMatcher creates a new SubPackageMatcher
func NewSubPackageMatcher(base, text string) *SubPackageMatcher {
	return &SubPackageMatcher{
		text: text,
		base: base,
	}
}

//Matches applies the filter
func (f *SubPackageMatcher) Matches() bool {
	matches, _ := regexp.MatchString(f.base+".*", f.text)
	return !matches
}
