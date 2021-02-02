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

	for index := range tcs {
		t.Run(tcs[index].input, func(t *testing.T) {
			parts, err := PartsOfURLSafe(tcs[index].input)
			if err != nil {
				t.Fatal("err should be nil:", err)
			}
			if !reflect.DeepEqual(parts, tcs[index].expected) {
				t.Fatal("expected", tcs[index].expected, "got:", parts)
			}
		})
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

	for index := range tcs {
		t.Run(tcs[index].input, func(t *testing.T) {
			val, err := PartsOfURLSafe(tcs[index].input)
			if err == nil && tcs[index].expectedErrNil {
				t.Fatal("err should not be nil:", err, "value:", val)
			}
			if val != nil && tcs[index].expectedErrNil {
				t.Fatal("value should be nil:", val)
			}
		})
	}
}
