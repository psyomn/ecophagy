package tinystory

import (
	"testing"

	"github.com/psyomn/ecophagy/common"
)

const (
	base                    = "../stories"
	fixture                 = base + "/simple.json"
	tinystoryFixtureSimple  = base + "/simple.tinystory"
	tinystoryFixtureSmall01 = base + "/small.tinystory"
	tinystoryFixtureSmall02 = base + "/small-2.tinystory"
)

func TestTinyStoryParserTitle(t *testing.T) {
	data, err := common.FileToBytes(fixture)

	if err != nil {
		t.Fatalf("no such file: %s", fixture)
	}

	document, err := Parse(data)
	if err != nil {
		t.Fatalf("error: %s", err.Error())
	}

	if document == nil {
		t.Fatalf("document must not be nil")
	}
}

func TestTinyStoryFormat(t *testing.T) {
	type testCase struct {
		name string
		path string
	}

	tcs := []testCase{
		{"simple", tinystoryFixtureSimple},
		{"small", tinystoryFixtureSmall01},
		{"small-2", tinystoryFixtureSmall02},
	}

	for index := range tcs {
		t.Run(tcs[index].name, func(t *testing.T) {
			//nolint:scopelint // succint tests good
			doc, err := ParseTinyStoryFormatFile(tcs[index].path)

			if doc == nil {
				t.Error("docs should not be nil")
			}

			if err != nil {
				t.Errorf("error not nil: %s", err.Error())
			}
		})
	}
}
