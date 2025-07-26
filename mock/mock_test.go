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

func Test_Mock_On_And_When(t *testing.T) {
	s := new(SomeImplementingStruct)
	s.On("DoSomething", "Hello", "World").ThenReturn("Hello World!")

	sut := TestStruct{
		dependency: s,
	}
	sut.Call()

}
