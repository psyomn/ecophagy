package main

import (
	"testing"
)

func TestMinifyGenreJson(t *testing.T) {
	line := `[{"id": 123, "name": "potato"}, {"id": 321, "name": "patata"}]`
	actual := minifyGenreJson(line)
	expected := `potato,patata`

	if actual != expected {
		t.Log("expected: ", expected)
		t.Log("actual: ", actual)
		t.Fatal()
	}
}
