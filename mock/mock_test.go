package mock

import (
	"testing"
)

type SomeInterface interface {
	DoSomething(a string, b string) string
}

type SomeImplementingStruct struct {
	Mock
}

type TestStruct struct {
	dependency SomeInterface
}

func (t *TestStruct) Call() {
	t.dependency.DoSomething("Hello", "World")
}

func (s *SomeImplementingStruct) DoSomething(a string, b string) string {
	args := s.Called("DoSomething", a, b)
	return args[0].(string)
}

func Test_Mock_On_And_When_And_AssertCalled(t *testing.T) {
	s := new(SomeImplementingStruct)
	s.On("DoSomething", "Hello", "World").
		ThenReturn("Hello World!")

	sut := TestStruct{
		dependency: s,
	}

	sut.Call()

	s.AssertCalled(t, "DoSomething", "Hello", "World")
}

func Test_AssertNumberOfCalls(t *testing.T) {
	s := new(SomeImplementingStruct)
	s.On("DoSomething", "Hello", "World").ThenReturn("Hello World")

	sut := TestStruct{
		dependency: s,
	}

	sut.Call()
	sut.Call()
	sut.Call()

	s.AssertNumberOfCalls(t, "DoSomething", 3)
}

func Test_AssertNotCalled(t *testing.T) {
	s := new(SomeImplementingStruct)
	s.On("DoSomething", "Hello", "World").ThenReturn("Hello World")

	_ = TestStruct{
		dependency: s,
	}

	s.AssertNotCalled(t, "DoSomething")
}
