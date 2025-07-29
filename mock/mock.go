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
	if m.calls == nil {
		m.calls = make(map[string][]Call)
	}

	m.calls[method] = append(m.calls[method], Call{
		args: args,
	})

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
	if s.mock.stubs == nil {
		s.mock.stubs = make(map[string][]Stub)
	}
	s.mock.stubs[s.method] = append(s.mock.stubs[s.method], Stub{
		ArgMatchers: s.matchers,
		Returns:     vals,
	})
}

func (m *Mock) AssertCalled(t *testing.T, method string, args ...any) {
	if calls, exists := m.calls[method]; !exists {
		t.Errorf("No matching method calls for method [%s]", method)
	} else {
		for _, call := range calls {
			if len(args) != len(call.args) {
				continue
			}

			matched := true
			for i, match := range call.args {
				if match != args[i] {
					matched = false
					break
				}
			}

			if matched {
				return
			}
		}
		t.Errorf("No matching calls for method [%s] with arguments %v", method, args)
	}
}

func (m *Mock) AssertNumberOfCalls(t *testing.T, method string, numCalls int) {
	calls, exists := m.calls[method]
	if !exists {
		t.Errorf("No matching method %s", method)
	}

	if len(calls) != numCalls {
		t.Errorf(`
			Error with number of calls for method %s
			Expecting: %d
			Found: %d
			`, method, numCalls, len(calls))
	}
}

func (m *Mock) AssertNotCalled(t *testing.T, method string) {
	calls, exists := m.calls[method]
	if exists {
		t.Errorf(`
			Not expecting calls for method %s
			Found %d calls`, method, len(calls))
	}
}
