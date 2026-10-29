# assert

A drop-in replacement for `stretchr/testify/assert` in one file, with no
dependency outside the standard library.

Copy the directory into a module and the tests compile against it:

```
cp -r titpetric/internal/assert <your-module>/internal/assert
```

The package is under `internal/`, and no module imports it from another.
Every module carries its own copy and a test dependency leaves the
go.mod. A copy diverges only where its module made
it diverge, and the file is small enough to read in full before trusting
it.

## Where it comes from

It derives from [fortio.org/assert](https://github.com/fortio/assert),
Apache 2.0, and is the superset of the copies that grew in these
repositories: atkins, task-ui, platform, platform-app, cli, vuego and
cdn. Every assertion any of them defined is here, and the file merges the
three the largest copy had split.

Two changes against the copy it came from:

- Assertions take `testing.TB` in place of `*testing.T`, and a helper
  taking a TB asserts as a test does.
- `ordered` returns whether the two operands compare at all. `Greater`
  read a bare 0 as "not greater", correct by accident; `LessOrEqual` and
  `GreaterOrEqual` read it as "equal" and passed an int against a string.

`EqualHTML` from the vuego copy is out. It compares two DOMs through
vuego's own differ, which is a dependency this file does not take.

## No require

There is no `require` half and no `FailNow`, except in `Fail`, which
exists to stop a test. A failed assertion reports and the test carries
on. That is a problem for a test that dereferences the value it just
checked: the assertion fails, the test keeps running and panics on the
next line, and the panic reads worse in the output than the assertion it
followed.

Every assertion returns whether it held, so that test reads:

```go
if !assert.NotNil(t, got) {
    return
}
assert.Equal(t, "x", got.Field)
```

The cdn copy returned a bool from `Len` alone. Here every assertion
does, and a caller may ignore it.

## The API

`Assert`, `CheckEquals`, `Contains`, `DirExists`, `ElementsMatch`,
`Empty`, `Equal`, `EqualValues`, `Error`, `ErrorAs`, `ErrorContains`,
`Errorf`, `ErrorIs`, `Eventually`, `Fail`, `False`, `FileExists`,
`Greater`, `GreaterOrEqual`, `IsIncreasing`, `JSONEq`, `Len`, `Less`,
`LessOrEqual`, `Nil`, `NoDirExists`, `NoError`, `NoFileExists`,
`NotContains`, `NotEmpty`, `NotEqual`, `NotNil`, `NotZero`,
`ObjectsAreEqualValues`, `Positive`, `True`, `Zero`.

`AnError` is a sentinel for a test that needs an error and does not care
which. `TestSuite` and `Run` replace `testify/suite`: embed `TestSuite`,
write `Test` methods, and `Run` calls each one with `SetupTest` before
and `TearDownTest` after when the suite defines them.
