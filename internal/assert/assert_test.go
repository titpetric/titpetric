package assert_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/titpetric/titpetric/internal/assert"
)

// recordingTB stands in for a *testing.T so a test can call the failure
// reporters without failing itself. The reporters reach only Fail and
// FailNow; nothing touches the embedded nil TB.
type recordingTB struct {
	testing.TB

	failed bool
}

func (r *recordingTB) Fail() {
	r.failed = true
}

func (r *recordingTB) FailNow() {
	r.failed = true
}

// mustFail calls a failing assertion against a recordingTB, which records
// the failure and leaves the test that drives it passing. Every assertion
// prints its message to stdout when it fails, and a passing run prints
// one line per case. Those lines come from the reporter working; no test fails.
func mustFail(t *testing.T, name string, failing func(testing.TB) bool) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		tb := &recordingTB{}
		held := failing(tb)
		assert.True(t, tb.failed, name+" must fail the TB it is given")
		assert.False(t, held, name+" must report that it did not hold")
	})
}

// wrappedError is the type the test asks ErrorAs to find.
type wrappedError struct{ msg string }

func (e *wrappedError) Error() string { return e.msg }

// errTarget is the sentinel the test passes to ErrorIs, which does not
// find it.
var errTarget = errors.New("target")

// suiteFixture exercises TestSuite and Run.
type suiteFixture struct {
	assert.TestSuite

	setup    int
	tearDown int
	ran      int
}

func (s *suiteFixture) SetupTest() {
	s.setup++
}

func (s *suiteFixture) TearDownTest() {
	s.tearDown++
}

func (s *suiteFixture) TestPasses() {
	s.ran++
	assert.True(s.T(), true, "the suite hands its T to the test")
}

// bareSuite has neither SetupTest nor TearDownTest, which is the branch of
// Run that suiteFixture does not reach.
type bareSuite struct {
	assert.TestSuite

	ran int
}

func (s *bareSuite) TestBare() {
	s.ran++
}

// fixturePaths returns a regular file and a directory that both exist.
// The four existence assertions read them.
func fixturePaths(t *testing.T) (file, dir string) {
	t.Helper()

	dir = t.TempDir()
	file = filepath.Join(dir, "file")
	if err := os.WriteFile(file, []byte("x"), 0o600); err != nil {
		t.Fatalf("writing the fixture file: %v", err)
	}
	return file, dir
}

// TestAssert invokes every assertion on the path where it holds.
func TestAssert(t *testing.T) {
	t.Run("equality", assertEquality)
	t.Run("errors", assertErrors)
	t.Run("emptiness", assertEmptiness)
	t.Run("containment", assertContainment)
	t.Run("ordering", assertOrdering)
	t.Run("files", assertFiles)
	t.Run("json", assertJSON)
	t.Run("polling", assertPolling)
	t.Run("reporters", assertReporters)
	t.Run("suite", assertSuite)
}

// TestAssertFailure invokes every assertion on the path where it does not
// hold, which is the branch that raises the failure.
func TestAssertFailure(t *testing.T) {
	t.Run("equality", assertEqualityFailure)
	t.Run("errors", assertErrorsFailure)
	t.Run("emptiness", assertEmptinessFailure)
	t.Run("containment", assertContainmentFailure)
	t.Run("ordering", assertOrderingFailure)
	t.Run("files", assertFilesFailure)
	t.Run("json", assertJSONFailure)
	t.Run("polling", assertPollingFailure)
}

// TestAssertReturnsTrue checks the contract the guard form relies on: an
// assertion that holds reports true, and `if !assert.X(...) { return }`
// carries on through a passing test.
func TestAssertReturnsTrue(t *testing.T) {
	assert.True(t, assert.Equal(t, 1, 1), "Equal returns true when it holds")
	assert.True(t, assert.NotNil(t, &struct{}{}), "NotNil returns true when it holds")
	assert.True(t, assert.Len(t, "abc", 3), "Len returns true when it holds")
	assert.True(t, assert.NoError(t, nil), "NoError returns true when it holds")
	assert.True(t, assert.Greater(t, 2, 1), "Greater returns true when it holds")
	assert.True(t, assert.CheckEquals(t, 1, 1), "CheckEquals returns true when it holds")
	assert.True(t, assert.Assert(t, true, "holds"), "Assert returns true when it holds")
}

func assertEquality(t *testing.T) {
	assert.True(t, assert.ObjectsAreEqualValues(1, 1))
	assert.False(t, assert.ObjectsAreEqualValues(1, 2))

	assert.Equal(t, "a", "a", "strings match")
	assert.EqualValues(t, []int{1, 2}, []int{1, 2})
	assert.NotEqual(t, 1, 2)

	assert.CheckEquals(t, 42, 42, "integers match")
	assert.Assert(t, true, "condition holds for %d", 42)

	assert.True(t, true)
	assert.False(t, false)
}

func assertEqualityFailure(t *testing.T) {
	mustFail(t, "Equal", func(tb testing.TB) bool { return assert.Equal(tb, "a", "b") })
	mustFail(t, "EqualValues", func(tb testing.TB) bool { return assert.EqualValues(tb, []int{1}, []int{2}) })
	mustFail(t, "NotEqual", func(tb testing.TB) bool { return assert.NotEqual(tb, 1, 1) })

	mustFail(t, "CheckEquals", func(tb testing.TB) bool { return assert.CheckEquals(tb, 41, 42, "integers differ") })
	mustFail(t, "Assert", func(tb testing.TB) bool { return assert.Assert(tb, false, "condition fails for %d", 42) })

	mustFail(t, "True", func(tb testing.TB) bool { return assert.True(tb, false) })
	mustFail(t, "False", func(tb testing.TB) bool { return assert.False(tb, true) })
}

func assertErrors(t *testing.T) {
	assert.NoError(t, nil)
	assert.Error(t, errors.New("boom"))
	assert.Error(t, assert.AnError)

	assert.ErrorIs(t, fmt.Errorf("open: %w", os.ErrNotExist), os.ErrNotExist)
	assert.ErrorContains(t, errors.New("boom"), "oo")

	var target *wrappedError
	assert.ErrorAs(t, errors.Join(assert.AnError, &wrappedError{"boom"}), &target)
	assert.Equal(t, "boom", target.Error())
}

func assertErrorsFailure(t *testing.T) {
	mustFail(t, "NoError", func(tb testing.TB) bool { return assert.NoError(tb, errors.New("boom")) })
	mustFail(t, "Error", func(tb testing.TB) bool { return assert.Error(tb, nil) })

	mustFail(t, "ErrorIs", func(tb testing.TB) bool { return assert.ErrorIs(tb, errors.New("boom"), errTarget) })
	mustFail(t, "ErrorContainsNil", func(tb testing.TB) bool { return assert.ErrorContains(tb, nil, "boom") })
	mustFail(t, "ErrorContains", func(tb testing.TB) bool { return assert.ErrorContains(tb, errors.New("boom"), "bang") })

	var target *wrappedError
	mustFail(t, "ErrorAs", func(tb testing.TB) bool { return assert.ErrorAs(tb, assert.AnError, &target) })
}

func assertEmptiness(t *testing.T) {
	var typed *int

	assert.Nil(t, nil)
	assert.Nil(t, typed)
	assert.NotNil(t, &struct{}{})

	assert.Empty(t, nil)
	assert.Empty(t, "")
	assert.Empty(t, []int{})
	assert.Empty(t, 0)
	assert.NotEmpty(t, "x")

	assert.Zero(t, 0)
	assert.Zero(t, "")
	assert.Zero(t, nil)
	assert.NotZero(t, 1)

	assert.Positive(t, 1)
	assert.Positive(t, uint(1))
	assert.Positive(t, 0.5)
	assert.Positive(t, time.Second)

	assert.Len(t, "abc", 3)
	assert.Len(t, []int{1, 2}, 2)
	assert.Len(t, map[string]int{"a": 1}, 1)
}

func assertEmptinessFailure(t *testing.T) {
	var typed *int

	mustFail(t, "Nil", func(tb testing.TB) bool { return assert.Nil(tb, &struct{}{}) })
	mustFail(t, "NotNil", func(tb testing.TB) bool { return assert.NotNil(tb, typed) })

	mustFail(t, "Empty", func(tb testing.TB) bool { return assert.Empty(tb, "x") })
	mustFail(t, "NotEmpty", func(tb testing.TB) bool { return assert.NotEmpty(tb, "") })

	mustFail(t, "Zero", func(tb testing.TB) bool { return assert.Zero(tb, 1) })
	mustFail(t, "NotZero", func(tb testing.TB) bool { return assert.NotZero(tb, 0) })

	mustFail(t, "Positive", func(tb testing.TB) bool { return assert.Positive(tb, 0) })
	mustFail(t, "PositiveUint", func(tb testing.TB) bool { return assert.Positive(tb, uint(0)) })
	mustFail(t, "PositiveFloat", func(tb testing.TB) bool { return assert.Positive(tb, -0.5) })
	// A kind that is no number is never positive.
	mustFail(t, "PositiveKind", func(tb testing.TB) bool { return assert.Positive(tb, "1") })
	mustFail(t, "PositiveNil", func(tb testing.TB) bool { return assert.Positive(tb, nil) })

	mustFail(t, "Len", func(tb testing.TB) bool { return assert.Len(tb, "abc", 2) })
	mustFail(t, "LenNoLength", func(tb testing.TB) bool { return assert.Len(tb, 42, 1) })
}

func assertContainment(t *testing.T) {
	assert.Contains(t, "haystack", "stack")
	assert.Contains(t, []string{"a", "b"}, "b")
	assert.Contains(t, map[string]int{"a": 1}, "a")

	assert.NotContains(t, "haystack", "needle")
	assert.NotContains(t, []string{"a"}, "b")
	assert.NotContains(t, map[string]int{"a": 1}, "b")

	assert.ElementsMatch(t, []int{1, 2, 2}, []int{2, 1, 2})
	assert.ElementsMatch(t, []string(nil), []string{})
}

func assertContainmentFailure(t *testing.T) {
	mustFail(t, "Contains", func(tb testing.TB) bool { return assert.Contains(tb, "haystack", "needle") })
	// A haystack reflect cannot read at all takes the invalid branch.
	mustFail(t, "ContainsInvalid", func(tb testing.TB) bool { return assert.Contains(tb, nil, "needle") })
	mustFail(t, "NotContains", func(tb testing.TB) bool { return assert.NotContains(tb, "haystack", "stack") })

	mustFail(t, "ElementsMatch", func(tb testing.TB) bool { return assert.ElementsMatch(tb, []int{1, 2}, []int{1, 3}) })
	mustFail(t, "ElementsMatchCount", func(tb testing.TB) bool { return assert.ElementsMatch(tb, []int{1, 1}, []int{1}) })
	mustFail(t, "ElementsMatchKind", func(tb testing.TB) bool { return assert.ElementsMatch(tb, 1, []int{1}) })
}

func assertOrdering(t *testing.T) {
	assert.Greater(t, 2, 1)
	assert.Greater(t, uint(2), uint(1))
	assert.Greater(t, 2.5, 1.5)
	assert.Greater(t, "b", "a")

	assert.Less(t, 1, 2)
	assert.LessOrEqual(t, 2, 2)
	assert.GreaterOrEqual(t, 2, 2)
	assert.LessOrEqual(t, 1*time.Second, 2*time.Second)

	assert.IsIncreasing(t, []int{1, 2, 3})
	assert.IsIncreasing(t, []string{"a", "b"})
}

func assertOrderingFailure(t *testing.T) {
	mustFail(t, "Greater", func(tb testing.TB) bool { return assert.Greater(tb, 1, 2) })
	mustFail(t, "GreaterUint", func(tb testing.TB) bool { return assert.Greater(tb, uint(1), uint(2)) })
	mustFail(t, "GreaterFloat", func(tb testing.TB) bool { return assert.Greater(tb, 1.5, 2.5) })
	mustFail(t, "GreaterString", func(tb testing.TB) bool { return assert.Greater(tb, "a", "b") })
	// Kinds that do not compare fail; they never count as equal.
	mustFail(t, "GreaterMixedKinds", func(tb testing.TB) bool { return assert.Greater(tb, 1, "a") })

	mustFail(t, "Less", func(tb testing.TB) bool { return assert.Less(tb, 2, 1) })
	mustFail(t, "LessOrEqual", func(tb testing.TB) bool { return assert.LessOrEqual(tb, 3, 2) })
	mustFail(t, "GreaterOrEqual", func(tb testing.TB) bool { return assert.GreaterOrEqual(tb, 1, 2) })
	mustFail(t, "LessOrEqualMixedKinds", func(tb testing.TB) bool { return assert.LessOrEqual(tb, 1, "a") })
	mustFail(t, "GreaterOrEqualMixedKinds", func(tb testing.TB) bool { return assert.GreaterOrEqual(tb, 1, "a") })

	mustFail(t, "IsIncreasing", func(tb testing.TB) bool { return assert.IsIncreasing(tb, []int{1, 3, 2}) })
	mustFail(t, "IsIncreasingKind", func(tb testing.TB) bool { return assert.IsIncreasing(tb, 42) })
}

func assertFiles(t *testing.T) {
	file, dir := fixturePaths(t)

	assert.FileExists(t, file)
	assert.NoFileExists(t, filepath.Join(dir, "absent"))
	// A directory is no file, and NoFileExists holds on one.
	assert.NoFileExists(t, dir)
	assert.DirExists(t, dir)
	assert.NoDirExists(t, filepath.Join(dir, "absent"))
	// And a file is no directory.
	assert.NoDirExists(t, file)
}

func assertFilesFailure(t *testing.T) {
	file, dir := fixturePaths(t)
	absent := filepath.Join(dir, "absent")

	mustFail(t, "FileExists", func(tb testing.TB) bool { return assert.FileExists(tb, absent) })
	mustFail(t, "FileExistsDir", func(tb testing.TB) bool { return assert.FileExists(tb, dir) })
	mustFail(t, "NoFileExists", func(tb testing.TB) bool { return assert.NoFileExists(tb, file) })
	mustFail(t, "DirExists", func(tb testing.TB) bool { return assert.DirExists(tb, absent) })
	mustFail(t, "DirExistsFile", func(tb testing.TB) bool { return assert.DirExists(tb, file) })
	mustFail(t, "NoDirExists", func(tb testing.TB) bool { return assert.NoDirExists(tb, dir) })
}

func assertJSON(t *testing.T) {
	assert.JSONEq(t, `{"a":1,"b":2}`, `{"b":2,  "a":1}`)
}

func assertJSONFailure(t *testing.T) {
	mustFail(t, "JSONEq", func(tb testing.TB) bool { return assert.JSONEq(tb, `{"a":1}`, `{"a":2}`) })
	mustFail(t, "JSONEqExpected", func(tb testing.TB) bool { return assert.JSONEq(tb, `{`, `{}`) })
	mustFail(t, "JSONEqActual", func(tb testing.TB) bool { return assert.JSONEq(tb, `{}`, `{`) })
}

func assertPolling(t *testing.T) {
	// The condition is already true, and Eventually returns on the first
	// read without waiting for a tick.
	assert.Eventually(t, func() bool { return true }, time.Second, time.Millisecond)

	// And here it turns true on a later read, which is the polling path.
	reads := 0
	assert.Eventually(t, func() bool { reads++; return reads > 2 }, time.Second, time.Millisecond)
}

func assertPollingFailure(t *testing.T) {
	mustFail(t, "Eventually", func(tb testing.TB) bool {
		return assert.Eventually(tb, func() bool { return false }, 5*time.Millisecond, time.Millisecond)
	})
}

// assertReporters drives Errorf and Fail against a stand-in TB, which
// records the failure they raise and leaves this test passing. Both print
// the message to stdout, and those are the two lines in a passing run.
func assertReporters(t *testing.T) {
	errorf := &recordingTB{}
	assert.Errorf(errorf, "reporter check: %d", 42)
	assert.True(t, errorf.failed, "Errorf fails the TB it is given")

	fail := &recordingTB{}
	assert.Fail(fail, "reporter check", "fail")
	assert.True(t, fail.failed, "Fail fails the TB it is given")
}

func assertSuite(t *testing.T) {
	fixture := &suiteFixture{}
	fixture.SetT(t)
	assert.True(t, fixture.T() == t, "SetT and T agree")

	assert.Run(t, fixture)
	assert.Equal(t, 1, fixture.setup)
	assert.Equal(t, 1, fixture.ran)
	assert.Equal(t, 1, fixture.tearDown)

	// A suite with neither hook takes the other branch of Run.
	bare := &bareSuite{}
	assert.Run(t, bare)
	assert.Equal(t, 1, bare.ran)
}
