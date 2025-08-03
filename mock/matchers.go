package mock

type GenericArgMatcher[T any] func(arg T) bool

func MatchedBy[T any](matcher GenericArgMatcher[T]) ArgMatcher {
	return func(actual any) bool {
		val, ok := actual.(T)
		if !ok {
			return false
		}
		return matcher(val)
	}
}
