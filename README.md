# 🦆 Duck Assert

Small experimental mocking library to understand how [Testify](https://github.com/stretchr/testify)'s mocking capabilities work under the hood.

## 🧠 Motivation

Coming from a Java background where Mocking frameworks like [Mockito](https://github.com/mockito/mockito) are the standard
for mocking dependencies, I wondered:

> "If Go isn't object-oriented, how does mocking even work here?"

That curiosity led me to explore `testify/mock`. To _really understand_ it, I decided to rebuild a minimal version of it
from scratch and that's how **Duck Assert** was born.

## 🔬 Overview

Duck Assert focuses on three key aspects of mocking:

1. Record Expectations – via `On(...).ThenReturn(...)`

Define what a method should return when called with specific arguments.

2. Record Calls – via `Called(...)`

When your code calls the mock, it records the method name and arguments, then finds and returns the matching stub’s values.

3. Assertions – via `AssertCalled`, `AssertNotCalled`, `AssertNumberOfCalls`

You verify later that methods were called as expected.

## ⚙️ Core Data Structures

`Mock`

Holds everything (literally):
- `stubs` — predefined expectations (`On(...).ThenReturn(...)`)
- `calls` — actual called made (`Called(...)`)

``` go
type Mock struct {
	stubs map[string][]Stub // expectations for each method
	calls map[string][]Call // actual calls recorded
}
```
___

`Stub`

Defines one expectation:

``` go
type Stub struct {
	ArgMatchers []ArgMatcher // matchers for args
	Returns     []any        // what to return
}
```

So if you'd write in the code:

``` go
m.On("GetUser", 42).ThenReturn("Alice", nil)
```

It creates a `Stub` for method `"GetUser"` with:
- One matcher for arg `42`
- Two return values: `"Alice"` and `nil`

___

`Call`

Represent an _actual invocation_ of a mocked method.

``` go
type Call struct {
  args []any
}
```

Every time `Called(...)` is triggered, a Call is recorded.

___

`ArgMatcher`

A function type for flexible argument matching.

``` go
type ArgMatcher func(actual any) bool
```

By default, argument matching uses a deep equality check.
This design allows extending the library later to support matchers like mock.Anything or mock.MatchedBy.

## 🔁 Core Flow

### 1. Defining Behavior - `On(...).ThenReturn(...)

``` go
m.On("GetUser", 42).ThenReturn("Alice", nil)
```

`On` builds a `StubBuilder` holding the method name and argument matchers. `ThenReturn` adds the new `Stub` into `m.stubs[method]`

Afterward:

``` go
m.stubs["GetUser"] = []Stub {
  {
    ArgMathcers: [matcherFor(42)],
    Returns: ["Alice", nil],
  }
}
```

___

### 2. Calling the mock - `Called(method, args...)`

This is the part the _code under test_ uses:

``` go
type MockUserService struct {
  mock.Mock
}

func (s *MockUserService) GetUser(id int) (string, error) {
  return s.mock.Called("GetUser", id).Get(0).(string), s.mock.Called("GetUser", id).Error(1)
}
```

What is happening here?

1. The method name + args are recorded in `MockUserService.mock.calls`
2. It loops over `stubs[method]` looking for the _first_ matching stub (same number of args, and all matchers return true)
3. If found, returns its predefined return values
4. If not found, returns `nil`

This mirrors `testify/mock`'s stub matching logic.

___

### 3. Getting return values - `ReturnValues`

You can extract values conveniently using several methods available:

``` go
r := m.Called("GetUser", 42)
name := r.Get(0).(string)
err := r.Error(1)
```

This mirrors how `testify` let's you extract the stubbed return values of a method.

### 4. Assertions - verifying call history

`AssertCalled`

Checks if the method was called with specific arguments:

- Looks in `stubs` (not `calls`) for a stub that matches the provided arguments.
- Reports an error if not found

Conceptually `testify` checks `calls` instead of stubs, but it still works the same.

___

`AssertNumberOfCalls`

Ensures the number of invocations matches expectations:

``` go
len(m.calls[method]) == numCalls
```

___

`AssertNotCalled`

Ensured a method was never called at all

## 🧩 Argument Matching (`deepMatch`)

The heart of this libary is its argument comparison logic.

When your test defines:

For example:

``` go
m.On("GetUser", User{Id: 10, Name: "Alice"}).ThenReturn(...)
```

and your code calls:

``` go
m.Called("GetUser", User{Id: 10, Name: "Alice"})
```

the mock must decide:

> "Do these two `User` structs represent the same input?"

You can't just use `==` because this is invalid for slices, maps, or structs with uncomparable fields, and even `reflect.DeepEqual`
might not be flexible enough if we want to support _custom matchers_ like `mock.Anything`, or predicate based matches.

So `deepMatch` solves this by:

1. Allowing **custom matching logic** (via `ArgMatcher`)
2. Falling back to **reflect-based recursive equality** when not using a matcher
3. Support **struct field-by-field comparisons**
4. Skipping unexported fields safely

This makes matching robust, extensible, and a key learning point in how `testify` works internally

___

## 🧾 Notes

Duck Assert wasn't built for production — it's a learning project. It's goal is to demistify what's going on _underneath_
mocking framworks like `testify/mock`
