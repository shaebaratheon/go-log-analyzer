package filter

import "regexp"

type PatternFilter struct {
	include []*regexp.Regexp
	exclude []*regexp.Regexp
}

func NewPatternFilter(includePatterns, excludePatterns []string) (*PatternFilter, error) {
	inc := make([]*regexp.Regexp, 0, len(includePatterns))
	for _, p := range includePatterns {
		r, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		inc = append(inc, r)
	}
	exc := make([]*regexp.Regexp, 0, len(excludePatterns))
	for _, p := range excludePatterns {
		r, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		exc = append(exc, r)
	}
	return &PatternFilter{include: inc, exclude: exc}, nil
}

func (f *PatternFilter) Match(line string) bool {
	for _, ex := range f.exclude {
		if ex.MatchString(line) {
			return false
		}
	}
	if len(f.include) == 0 {
		return true
	}
	for _, in := range f.include {
		if in.MatchString(line) {
			return true
		}
	}
	return false
}
