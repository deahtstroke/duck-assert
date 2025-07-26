package mock

import (
	"testing"
)

type Mock struct {
	stubs map[string][]Stub
	calls map[string][]Call
}

type Call struct {
	args []any
}

type Stub struct {
	ArgMatchers []ArgMatcher
	Returns     []any
}

type ArgMatcher func(arg any) bool

type StubBuilder struct {
	method   string
	matchers []ArgMatcher
	mock     *Mock
}

// Records a method call to the mock struct with given arguments
// This should return the return values recorded by the stub that first
// matches the given arugments
func (m *Mock) Called(method string, args ...any) []any {
	for _, stub := range m.stubs[method] {
		if len(args) != len(stub.ArgMatchers) {
			continue
		}

		matched := true
		for i, matcher := range stub.ArgMatchers {
			if !matcher(args[i]) {
				matched = false
				break
			}
		}

		if matched {
			return stub.Returns
		}
	}
	return nil
}

// Argument matcher that matches a value exactly
// TODO: Support structs and slices
func MatchExact(this any) ArgMatcher {
	return func(arg any) bool {
		return this == arg
	}
}

func (m *Mock) On(method string, args ...any) *StubBuilder {
	matchers := make([]ArgMatcher, len(args))
	for i, arg := range args {
		matchers[i] = MatchExact(arg)
	}
	return &StubBuilder{
		method:   method,
		matchers: matchers,
		mock:     m,
	}
}

func (s *StubBuilder) ThenReturn(vals ...any) {
	s.mock.stubs[s.method] = append(s.mock.stubs[s.method], Stub{
		ArgMatchers: s.matchers,
		Returns:     vals,
	})
}

func (m *Mock) AssertCalled(t *testing.T, method string, args ...any) {

}
