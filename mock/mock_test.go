package mock

import (
	"fmt"
	"testing"
)

type ExampleStruct struct {
	Mock
}

// Function that returns a concatenated string (not really)
// Simplest example for a single-returning value function
func (s *ExampleStruct) Example1(a string, b string) string {
	args := s.Called("Example1", a, b)
	return args.Get(0).(string)
}

// Simple function but this one returns an error, we are testing the
// .Error(index int) function
func (s *ExampleStruct) Example2() error {
	args := s.Called("Example2")
	return args.Error(0)
}

func (s *ExampleStruct) Example3(a string, b string) (string, error) {
	args := s.Called("Example3", a, b)
	return args.Get(0).(string), args.Error(1)
}

type MyStruct struct {
	Output string
}

func (s *ExampleStruct) Example4(a string, b string) (MyStruct, error) {
	args := s.Called("Example4", a, b)
	return args.Get(0).(MyStruct), args.Error(1)
}

func Test_builderMatcher_tableTests(t *testing.T) {
	tests := map[string]struct {
		input    any
		expected any
		want     bool
	}{
		"int match": {
			input:    1,
			expected: 1,
			want:     true,
		},
		"int no match": {
			input:    1,
			expected: 2,
			want:     false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			matcher := buildMatcher(test.input)
			got := matcher(test.expected)
			if got != test.want {
				t.Errorf("Did not get the expected output. Want: %t. Got: %t", test.want, got)
			}
		})
	}
}

func Test_Example1(t *testing.T) {
	m := ExampleStruct{Mock: Mock{}}
	m.On("Example1", "Hello", "World").
		ThenReturn("Hello World!")

	m.Example1("Hello", "World")

	m.assertStubs(t, "Example1", []Stub{
		{
			ArgMatchers: []ArgMatcher{
				buildMatcher("Hello"),
				buildMatcher("World"),
			},
			Returns: []any{
				"Hello World!",
			},
		},
	},
	)

	m.assertCalls(t, "Example1", []Call{
		{
			args: []any{
				"Hello",
				"World",
			},
		},
	})

}

func Test_Example2(t *testing.T) {
	m := ExampleStruct{Mock: Mock{}}
	m.On("Example2").ThenReturn(nil)

	m.Example2()

	m.assertCalls(t, "Example2", []Call{
		{
			args: []any{},
		},
	})

	m.assertStubs(t, "Example2", []Stub{
		{
			ArgMatchers: []ArgMatcher{},
			Returns:     []any{},
		},
	})
}

func Test_Example3_ReturnSuccessfully(t *testing.T) {
	m := ExampleStruct{Mock: Mock{}}
	m.On("Example3", "Hello", "World").ThenReturn("Hello World", nil)

	value, err := m.Example3("Hello", "World")
	if err != nil {
		t.Errorf("Not expecting error: %v", err)
	}

	if value != "Hello World" {
		t.Error("Value is incorrect")
	}

	m.assertStubs(t, "Example3", []Stub{
		{
			ArgMatchers: []ArgMatcher{
				buildMatcher("Hello"),
				buildMatcher("World"),
			},
			Returns: []any{
				"Hello World",
				nil,
			},
		},
	})

	m.assertCalls(t, "Example3", []Call{
		{
			args: []any{
				"Hello", "World",
			},
		},
	})
}

func Test_Example3_ReturnError(t *testing.T) {
	m := ExampleStruct{Mock: Mock{}}
	error := fmt.Errorf("Error concatenating strings")
	m.On("Example3", "Hello", "World").ThenReturn("", error)

	_, err := m.Example3("Hello", "World")
	if err == nil {
		t.Errorf("Expecting errors, found none")
	}

	m.assertStubs(t, "Example3", []Stub{
		{
			ArgMatchers: []ArgMatcher{
				buildMatcher("Hello"),
				buildMatcher("World"),
			},
			Returns: []any{
				"",
				error,
			},
		},
	})

	m.assertCalls(t, "Example3", []Call{
		{
			args: []any{
				"Hello", "World",
			},
		},
	})
}

func Test_Example4_ReturnStruct(t *testing.T) {
	m := ExampleStruct{Mock: Mock{}}
	st := MyStruct{Output: "Hello World"}
	m.On("Example4", "Hello", "World").ThenReturn(st, nil)

	result, err := m.Example4("Hello", "World")
	if err != nil {
		t.Errorf("Not expecting errors, got: %v", err)
	}

	if result.Output != "Hello World" {
		t.Errorf("Wrong result")
	}

	m.assertStubs(t, "Example4", []Stub{
		{
			ArgMatchers: []ArgMatcher{
				buildMatcher("Hello"),
				buildMatcher("World"),
			},
			Returns: []any{
				st, nil,
			},
		},
	})

	m.assertCalls(t, "Example4", []Call{
		{
			args: []any{
				"Hello", "World",
			},
		},
	})
}

func Test_Example4_ReturnError(t *testing.T) {
	m := ExampleStruct{Mock: Mock{}}
	st := MyStruct{}
	m.On("Example4", "Hello", "World").ThenReturn(st, fmt.Errorf("Error getting struct"))

	_, err := m.Example4("Hello", "World")
	if err == nil {
		t.Errorf("Expecting error, found none")
	}

	m.assertStubs(t, "Example4", []Stub{
		{
			ArgMatchers: []ArgMatcher{
				buildMatcher("Hello"),
				buildMatcher("World"),
			},
		},
	})

	m.assertCalls(t, "Example4", []Call{
		{
			args: []any{
				"Hello", "World",
			},
		},
	})
}

func Test_AssertCalled(t *testing.T) {
	m := ExampleStruct{Mock: Mock{}}
	m.On("Example1", "d", "3").ThenReturn("d 3!")

	m.Example1("d", "3")

	m.AssertCalled(t, "Example1", "d", "3")
}

func Test_AssertNotCalled(t *testing.T) {
	m := ExampleStruct{Mock: Mock{}}
	m.On("Example1", "123", "456").ThenReturn("123 456!")

	m.AssertNotCalled(t, "Example1")
}

func Test_AssertNumberCalls(t *testing.T) {
	m := ExampleStruct{Mock: Mock{}}
	m.On("Example1", "123", "456").ThenReturn("123 456!")
	m.On("Example1", "456", "").ThenReturn("456", "")

	m.Example1("123", "456")
	m.Example1("456", "")

	m.AssertNumberOfCalls(t, "Example1", 2)
}

func (m *Mock) assertCalls(t *testing.T, method string, expectedCalls []Call) {
	for i := 0; i < len(m.calls); i++ {
		curr := m.calls[method][i]
		if len(curr.args) != len(expectedCalls[i].args) {
			t.Errorf("Length of args do not match")
		}

		matched := true
		for j, arg := range curr.args {
			if arg != expectedCalls[i].args[j] {
				matched = false
				break
			}
		}

		if !matched {
			t.Errorf("Call %s did not match the expected calls", m.calls["Example1"][i])
		}
	}
}

func (m *Mock) assertStubs(t *testing.T, method string, expectedStubs []Stub) {
	for i := 0; i < len(m.stubs); i++ {
		curr := m.stubs[method][i]
		if len(curr.ArgMatchers) != len(expectedStubs[i].ArgMatchers) {
			t.Errorf("Length of matchers do not match")
		}

		matched := true
		for j, matcher := range curr.ArgMatchers {
			currArg := m.calls[method][i].args[j]
			if matcher(currArg) != expectedStubs[i].ArgMatchers[j](currArg) {
				matched = false
				break
			}
		}

		if !matched {
			t.Errorf("Call %s did not match the expected calls", m.calls["Example1"][i])
		}
	}
}
