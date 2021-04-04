package common

import (
	"reflect"
	"testing"
)

func TestPartsOfURLSafeSimple(t *testing.T) {
	type tc struct {
		input    string
		expected []string
	}

	tcs := []tc{
		{"/a/b/c", []string{"a", "b", "c"}},
		{"/a/b/c/", []string{"a", "b", "c"}},
	}

	// TODO: think about a better way to do this in the future
	// nolint
	makeTestFn := func(tc *tc, upt *testing.T) func(t *testing.T) {
		return func(t *testing.T) {
			parts, err := PartsOfURLSafe(tc.input)
			if err != nil {
				t.Fatal("err should be nil:", err)
			}
			if !reflect.DeepEqual(parts, tc.expected) {
				t.Fatal("expected:", tc.expected, "got:", parts)
			}
		}
	}

	for index := range tcs {
		t.Run(tcs[index].input, makeTestFn(&tcs[index], t))
	}
}

func TestPartsOfURLSafeDirectoryTraversalAttempt(t *testing.T) {
	type tc struct {
		input          string
		expectedErrNil bool
	}

	tcs := []tc{
		{"/path/to/place/../oh/no", true},
		{"/path/to/..", true},
		{"/path/../to", true},
		{"/../path/to", true},
		{"..", true},
		{"/a/b/c", false},
	}

	// TODO: think about a better way to do this in the future
	// nolint
	makeTestFn := func(tc *tc, upt *testing.T) func(t *testing.T) {
		return func(t *testing.T) {
			val, err := PartsOfURLSafe(tc.input)
			if err == nil && tc.expectedErrNil {
				t.Fatal("err should not be nil:", err, "value:", val)
			}
			if val != nil && tc.expectedErrNil {
				t.Fatal("value should be nil:", val)
			}
		}
	}

	for index := range tcs {
		t.Run(tcs[index].input, makeTestFn(&tcs[index], t))
	}
}
