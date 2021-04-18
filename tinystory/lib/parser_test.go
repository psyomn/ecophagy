package tinystory

import (
	"testing"

	"github.com/psyomn/ecophagy/common"
)

const fixture = "../fixtures/simple.json"

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
