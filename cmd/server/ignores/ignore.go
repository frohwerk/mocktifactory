package ignores

import "strings"

type predicate func(string) bool

type set []predicate

func Equal(file string) predicate {
	return func(f string) bool { return f == file }
}

func HasSuffix(suffix string) predicate {
	return func(f string) bool { return strings.HasSuffix(f, suffix) }
}

func New(values ...predicate) *set {
	set := make(set, len(values))
	for i, value := range values {
		set[i] = value
	}
	return &set
}

func (predicates *set) Matches(v string) bool {
	for _, predicate := range *predicates {
		if predicate(v) {
			return true
		}
	}
	return false
}
