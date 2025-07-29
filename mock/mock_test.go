package mock

import (
	"testing"
)

type SomeErrorInterface interface {
	ReturnError() error
}

type SomeInterface interface {
	DoSomething(a string, b string) string
}

type SomeImplementingStruct struct {
	Mock
}

type SomeImplementingErrorStruct struct {
	Mock
}

type TestStruct struct {
	dependency SomeInterface
}

type TestErrorStruct struct {
	dependency SomeErrorInterface
}

func (t *TestStruct) Call() {
	t.dependency.DoSomething("Hello", "World")
}

func (t *TestErrorStruct) Call() {
	t.dependency.ReturnError()
}

func (s *SomeImplementingStruct) DoSomething(a string, b string) string {
	args := s.Called("DoSomething", a, b)
	return args.Get(0).(string)
}

func (s *SomeImplementingErrorStruct) ReturnError() error {
	args := s.Called("ReturnError")
	return args.Error(0)
}

func Test_On_And_When_And_AssertCalled(t *testing.T) {
	s := new(SomeImplementingStruct)
	s.On("DoSomething", "Hello", "World").
		ThenReturn("Hello World!")

	sut := TestStruct{
		dependency: s,
	}

	sut.Call()

	s.AssertCalled(t, "DoSomething", "Hello", "World")
}

func Test_ErrorReturningMethods(t *testing.T) {
	s := new(SomeImplementingErrorStruct)
	s.On("ReturnError").ThenReturn(nil)

	sut := TestErrorStruct{
		dependency: s,
	}

	sut.Call()

	s.AssertCalled(t, "ReturnError")
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
