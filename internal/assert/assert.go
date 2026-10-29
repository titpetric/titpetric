// Copyright 2022 Fortio Authors
// Licensed under the Apache License, Version 2.0.

// Package assert is a drop-in replacement for stretchr/testify/assert in
// one file with no dependency outside the standard library. Copy the
// directory into a module and the tests compile against it unchanged.
//
// It derives from fortio.org/assert and takes testing.TB in place of
// *testing.T, and a helper taking a TB asserts as a test does.
//
// There is no require half and no FailNow, except in Fail, which exists
// to stop a test. A failed assertion reports and the test carries on,
// which is a problem for a test that dereferences the value it just
// checked: the assertion fails, the test keeps running and panics on the
// next line. Every assertion returns whether it held, so that test reads:
//
//	if !assert.NotNil(t, got) {
//		return
//	}
//	assert.Equal(t, "x", got.Field)
package assert

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
)

// AnError is a sentinel for a test that needs an error and does not care
// which. It replaces testify's assert.AnError.
var AnError = errors.New("assert.AnError general error for testing")

// Errorf reports a failure at the line of the assertion that called it.
// It reads two frames up: one for the assertion, one for the test.
func Errorf(t testing.TB, format string, rest ...any) {
	_, file, line, _ := runtime.Caller(2)
	file = file[strings.LastIndex(file, "/")+1:]
	fmt.Printf("%s:%d %s\n", file, line, fmt.Sprintf(format, rest...))
	t.Fail()
}

// Fail stops the test, the one reporter here that calls FailNow.
func Fail(t testing.TB, msg string, args ...any) {
	_, file, line, _ := runtime.Caller(1)
	file = file[strings.LastIndex(file, "/")+1:]
	if len(args) > 0 {
		msg = fmt.Sprintf("%s: %v", msg, args)
	}
	fmt.Printf("%s:%d %s\n", file, line, msg)
	t.FailNow()
}

// ObjectsAreEqualValues reports whether a and b are deeply equal.
func ObjectsAreEqualValues(a, b any) bool {
	return reflect.DeepEqual(a, b)
}

// Equal checks that a and b are deeply equal.
func Equal(t testing.TB, a, b any, msg ...any) bool {
	if !ObjectsAreEqualValues(a, b) {
		Errorf(t, "%v unexpectedly not equal %v: %v", a, b, msg)
		return false
	}
	return true
}

// EqualValues checks that a and b are deeply equal. Equal under the name
// testify reserves for the looser comparison.
func EqualValues(t testing.TB, a, b any, msg ...any) bool {
	if !ObjectsAreEqualValues(a, b) {
		Errorf(t, "%v unexpectedly not equal %v: %v", a, b, msg)
		return false
	}
	return true
}

// NotEqual checks that a and b are not deeply equal.
func NotEqual(t testing.TB, a, b any, msg ...any) bool {
	if ObjectsAreEqualValues(a, b) {
		Errorf(t, "%v unexpectedly equal: %v", a, msg)
		return false
	}
	return true
}

// CheckEquals checks that actual equals expected with ==, and prints
// each value on its own line when they differ.
func CheckEquals(t testing.TB, actual, expected any, msg ...any) bool {
	if expected != actual {
		_, file, line, _ := runtime.Caller(1)
		file = file[strings.LastIndex(file, "/")+1:]
		fmt.Printf("%s:%d mismatch!\nactual:\n%+v\nexpected:\n%+v\nfor %+v\n", file, line, actual, expected, msg)
		t.Fail()
		return false
	}
	return true
}

// Assert is True with a printf message, under the name of the earlier
// fortio/stats tests.
func Assert(t testing.TB, cond bool, msg string, rest ...any) bool {
	if !cond {
		_, file, line, _ := runtime.Caller(1)
		file = file[strings.LastIndex(file, "/")+1:]
		fmt.Printf("%s:%d assert failure: %s\n", file, line, fmt.Sprintf(msg, rest...))
		t.Fail()
		return false
	}
	return true
}

// True checks that b is true.
func True(t testing.TB, b bool, msg ...any) bool {
	if !b {
		Errorf(t, "expecting true, didn't: %v", msg)
		return false
	}
	return true
}

// False checks that b is false.
func False(t testing.TB, b bool, msg ...any) bool {
	if b {
		Errorf(t, "expecting false, didn't: %v", msg)
		return false
	}
	return true
}

// NoError checks that err is nil.
func NoError(t testing.TB, err error, msg ...any) bool {
	if err != nil {
		Errorf(t, "expecting no error, got %v: %v", err, msg)
		return false
	}
	return true
}

// Error checks that err is not nil.
func Error(t testing.TB, err error, msg ...any) bool {
	if err == nil {
		Errorf(t, "expecting an error, didn't get it: %v", msg)
		return false
	}
	return true
}

// ErrorIs checks that err, or an error it wraps, matches target.
func ErrorIs(t testing.TB, err, target error, msg ...any) bool {
	if !errors.Is(err, target) {
		Errorf(t, "%v does not match %v: %v", err, target, msg)
		return false
	}
	return true
}

// ErrorAs checks that err, or an error it wraps, is assignable to target,
// a non-nil pointer.
func ErrorAs(t testing.TB, err error, target any, msg ...any) bool {
	if !errors.As(err, target) {
		Errorf(t, "%v does not match %T: %v", err, target, msg)
		return false
	}
	return true
}

// ErrorContains checks that err is not nil and its message contains want.
func ErrorContains(t testing.TB, err error, want string, msg ...any) bool {
	if err == nil {
		Errorf(t, "expecting an error containing %q, got nil: %v", want, msg)
		return false
	}
	if !strings.Contains(err.Error(), want) {
		Errorf(t, "%v does not contain %q: %v", err, want, msg)
		return false
	}
	return true
}

// Nil checks that value is nil, a typed nil included.
func Nil(t testing.TB, value any, msg ...any) bool {
	if !isNil(value) {
		Errorf(t, "expecting nil, got %v: %v", value, msg)
		return false
	}
	return true
}

// NotNil checks that value is not nil, a typed nil included.
func NotNil(t testing.TB, value any, msg ...any) bool {
	if isNil(value) {
		Errorf(t, "expecting non-nil: %v", msg)
		return false
	}
	return true
}

// Empty checks that value is nil, the zero value, or of length zero.
func Empty(t testing.TB, value any, msg ...any) bool {
	if !isEmpty(value) {
		Errorf(t, "expecting empty, got %v: %v", value, msg)
		return false
	}
	return true
}

// NotEmpty checks that value is neither nil nor empty.
func NotEmpty(t testing.TB, value any, msg ...any) bool {
	if isEmpty(value) {
		Errorf(t, "expecting non-empty: %v", msg)
		return false
	}
	return true
}

// Zero checks that value is the zero value for its type.
func Zero(t testing.TB, value any, msg ...any) bool {
	if !isZero(value) {
		Errorf(t, "%v is not the zero value: %v", value, msg)
		return false
	}
	return true
}

// NotZero checks that value is not the zero value for its type.
func NotZero(t testing.TB, value any, msg ...any) bool {
	if isZero(value) {
		Errorf(t, "expecting a non-zero value: %v", msg)
		return false
	}
	return true
}

// Positive checks that value holds a numeric kind greater than zero.
func Positive(t testing.TB, value any, msg ...any) bool {
	if !isPositive(value) {
		Errorf(t, "%v is not positive: %v", value, msg)
		return false
	}
	return true
}

// Len checks the length of a string, array, channel, map or slice.
func Len(t testing.TB, value any, expected int, msg ...any) bool {
	v := reflect.ValueOf(value)
	if !hasLength(v) {
		Errorf(t, "%T has no length: %v", value, msg)
		return false
	}
	if actual := v.Len(); actual != expected {
		Errorf(t, "length %d unexpectedly not equal %d: %v", actual, expected, msg)
		return false
	}
	return true
}

// Contains checks that needle is in haystack, which may be a string, an
// array, a slice or a map.
func Contains(t testing.TB, haystack, needle any, msg ...any) bool {
	if !containsElement(haystack, needle) {
		Errorf(t, "%v doesn't contain %v: %v", haystack, needle, msg)
		return false
	}
	return true
}

// NotContains checks that needle is not in haystack.
func NotContains(t testing.TB, haystack, needle any, msg ...any) bool {
	if containsElement(haystack, needle) {
		Errorf(t, "%v unexpectedly contains %v: %v", haystack, needle, msg)
		return false
	}
	return true
}

// ElementsMatch checks that two lists hold the same elements, as often
// each, in any order.
func ElementsMatch(t testing.TB, listA, listB any, msg ...any) bool {
	extraA, extraB, ok := diffElements(listA, listB)
	switch {
	case !ok:
		Errorf(t, "%T and %T are not both arrays or slices: %v", listA, listB, msg)
		return false
	case len(extraA) > 0 || len(extraB) > 0:
		Errorf(t, "elements differ, only in %v: %v, only in %v: %v: %v", listA, extraA, listB, extraB, msg)
		return false
	}
	return true
}

// IsIncreasing checks that every item is greater than the one before it.
func IsIncreasing(t testing.TB, value any, msg ...any) bool {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Array && v.Kind() != reflect.Slice {
		Errorf(t, "%T is not an array or slice: %v", value, msg)
		return false
	}
	for i := 1; i < v.Len(); i++ {
		if cmp, ok := ordered(v.Index(i).Interface(), v.Index(i-1).Interface()); !ok || cmp <= 0 {
			Errorf(t, "%v is not increasing at index %d: %v", value, i, msg)
			return false
		}
	}
	return true
}

// Greater checks that actual is greater than expected.
func Greater(t testing.TB, actual, expected any, msg ...any) bool {
	if cmp, ok := ordered(actual, expected); !ok || cmp <= 0 {
		Errorf(t, "%v is not greater than %v: %v", actual, expected, msg)
		return false
	}
	return true
}

// GreaterOrEqual checks that actual is greater than or equal to expected.
func GreaterOrEqual(t testing.TB, actual, expected any, msg ...any) bool {
	if cmp, ok := ordered(actual, expected); !ok || cmp < 0 {
		Errorf(t, "%v is not greater than or equal to %v: %v", actual, expected, msg)
		return false
	}
	return true
}

// Less checks that actual is less than expected.
func Less(t testing.TB, actual, expected any, msg ...any) bool {
	if cmp, ok := ordered(actual, expected); !ok || cmp >= 0 {
		Errorf(t, "%v is not less than %v: %v", actual, expected, msg)
		return false
	}
	return true
}

// LessOrEqual checks that actual is less than or equal to expected.
func LessOrEqual(t testing.TB, actual, expected any, msg ...any) bool {
	if cmp, ok := ordered(actual, expected); !ok || cmp > 0 {
		Errorf(t, "%v is not less than or equal to %v: %v", actual, expected, msg)
		return false
	}
	return true
}

// FileExists checks that path exists and is a regular file.
func FileExists(t testing.TB, path string, msg ...any) bool {
	info, err := os.Lstat(path)
	switch {
	case err != nil:
		Errorf(t, "%q does not exist: %v", path, msg)
		return false
	case info.IsDir():
		Errorf(t, "%q is a directory, not a file: %v", path, msg)
		return false
	}
	return true
}

// NoFileExists checks that path is no regular file. A directory passes,
// as in testify.
func NoFileExists(t testing.TB, path string, msg ...any) bool {
	info, err := os.Lstat(path)
	if err == nil && !info.IsDir() {
		Errorf(t, "%q unexpectedly exists: %v", path, msg)
		return false
	}
	return true
}

// DirExists checks that path exists and is a directory.
func DirExists(t testing.TB, path string, msg ...any) bool {
	info, err := os.Lstat(path)
	switch {
	case err != nil:
		Errorf(t, "%q does not exist: %v", path, msg)
		return false
	case !info.IsDir():
		Errorf(t, "%q is a file, not a directory: %v", path, msg)
		return false
	}
	return true
}

// NoDirExists checks that path is no directory. A regular file passes, as
// in testify.
func NoDirExists(t testing.TB, path string, msg ...any) bool {
	info, err := os.Lstat(path)
	if err == nil && info.IsDir() {
		Errorf(t, "%q unexpectedly exists: %v", path, msg)
		return false
	}
	return true
}

// JSONEq checks that two JSON documents carry the same value, ignoring
// key order and whitespace.
func JSONEq(t testing.TB, expected, actual string, msg ...any) bool {
	var want, got any
	if err := json.Unmarshal([]byte(expected), &want); err != nil {
		Errorf(t, "expected value is not valid JSON: %v: %v", err, msg)
		return false
	}
	if err := json.Unmarshal([]byte(actual), &got); err != nil {
		Errorf(t, "actual value is not valid JSON: %v: %v", err, msg)
		return false
	}
	if !reflect.DeepEqual(want, got) {
		Errorf(t, "%v unexpectedly not equal %v: %v", actual, expected, msg)
		return false
	}
	return true
}

// Eventually polls condition every tick until it returns true or waitFor
// elapses.
func Eventually(t testing.TB, condition func() bool, waitFor, tick time.Duration, msg ...any) bool {
	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	deadline := time.After(waitFor)
	for {
		if condition() {
			return true
		}
		select {
		case <-deadline:
			// One last read: the condition may have turned true between
			// the final tick and the deadline.
			if condition() {
				return true
			}
			Errorf(t, "condition was not met within %v: %v", waitFor, msg)
			return false
		case <-ticker.C:
		}
	}
}

// TestSuite is the base of a test suite, in place of
// github.com/stretchr/testify/suite.
type TestSuite struct {
	t *testing.T
}

// T returns the testing.T the suite runs under.
func (s *TestSuite) T() *testing.T {
	return s.t
}

// SetT sets the testing.T the suite runs under.
func (s *TestSuite) SetT(t *testing.T) {
	s.t = t
}

type hasT interface {
	T() *testing.T
	SetT(t *testing.T)
}

type hasSetupTest interface {
	SetupTest()
}

type hasTearDown interface {
	TearDownTest()
}

// Run runs each Test method of suite as a subtest, with SetupTest before
// and TearDownTest after when the suite has them.
func Run(t *testing.T, suite hasT) {
	suite.SetT(t)

	setup, _ := suite.(hasSetupTest)
	tearDown, _ := suite.(hasTearDown)

	for method := range reflect.TypeOf(suite).Methods() {
		if !strings.HasPrefix(method.Name, "Test") {
			continue
		}
		if setup != nil {
			setup.SetupTest()
		}
		t.Run(method.Name, func(_ *testing.T) {
			method.Func.Call([]reflect.Value{reflect.ValueOf(suite)})
		})
		if tearDown != nil {
			tearDown.TearDownTest()
		}
	}
}

func isNil(value any) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

func isEmpty(value any) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	if hasLength(v) {
		return v.Len() == 0
	}
	return v.IsZero()
}

func isZero(value any) bool {
	if value == nil {
		return true
	}
	return reflect.ValueOf(value).IsZero()
}

func isPositive(value any) bool {
	v := reflect.ValueOf(value)
	if !v.IsValid() {
		return false
	}
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() > 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() > 0
	case reflect.Float32, reflect.Float64:
		return v.Float() > 0
	default:
		return false
	}
}

func hasLength(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return true
	default:
		return false
	}
}

func isList(v reflect.Value) bool {
	return v.IsValid() && (v.Kind() == reflect.Array || v.Kind() == reflect.Slice)
}

func containsElement(haystack, needle any) bool {
	v := reflect.ValueOf(haystack)
	if !v.IsValid() {
		return false
	}
	switch v.Kind() {
	case reflect.String:
		s, ok := needle.(string)
		return ok && strings.Contains(v.String(), s)
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			if reflect.DeepEqual(v.Index(i).Interface(), needle) {
				return true
			}
		}
	case reflect.Map:
		key := reflect.ValueOf(needle)
		return key.IsValid() && key.Type().AssignableTo(v.Type().Key()) && v.MapIndex(key).IsValid()
	}
	return false
}

// diffElements reports the elements of listA missing from listB and the
// other way around. The final return is false when either value is no
// array and no slice.
func diffElements(listA, listB any) (extraA, extraB []any, ok bool) {
	a := reflect.ValueOf(listA)
	b := reflect.ValueOf(listB)
	if !isList(a) || !isList(b) {
		return nil, nil, false
	}

	matched := make([]bool, b.Len())
	for i := 0; i < a.Len(); i++ {
		element := a.Index(i).Interface()
		found := false
		for j := 0; j < b.Len(); j++ {
			if matched[j] || !ObjectsAreEqualValues(element, b.Index(j).Interface()) {
				continue
			}
			matched[j] = true
			found = true
			break
		}
		if !found {
			extraA = append(extraA, element)
		}
	}
	for j := 0; j < b.Len(); j++ {
		if !matched[j] {
			extraB = append(extraB, b.Index(j).Interface())
		}
	}
	return extraA, extraB, true
}

// ordered compares two values of the same broad kind. The second return
// is false when the two do not compare at all, which every caller treats
// as a failure: a bare 0 would read as "equal" and let LessOrEqual and
// GreaterOrEqual pass on an int against a string.
func ordered(actual, expected any) (int, bool) {
	a := reflect.ValueOf(actual)
	b := reflect.ValueOf(expected)
	if !a.IsValid() || !b.IsValid() {
		return 0, false
	}

	switch a.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if b.Kind() >= reflect.Int && b.Kind() <= reflect.Int64 {
			return compare(a.Int(), b.Int()), true
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if b.Kind() >= reflect.Uint && b.Kind() <= reflect.Uintptr {
			return compare(a.Uint(), b.Uint()), true
		}
	case reflect.Float32, reflect.Float64:
		if b.Kind() == reflect.Float32 || b.Kind() == reflect.Float64 {
			return compare(a.Float(), b.Float()), true
		}
	case reflect.String:
		if b.Kind() == reflect.String {
			return strings.Compare(a.String(), b.String()), true
		}
	}
	return 0, false
}

func compare[T ~int64 | ~uint64 | ~float64](a, b T) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}
