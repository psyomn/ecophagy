package img

import (
	"io/ioutil"
	"testing"
)

const FixtureFilePath = "../fixtures/test.jpg"

func fileIntoBytes(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return data
}

func TestGetExifComment(t *testing.T) {
	bytes := fileIntoBytes(FixtureFilePath)

	comment, err := GetExifComment(bytes)
	if err != nil {
		t.Errorf("%v", comment)
	}

	if comment != "expected" {
		t.Errorf("%s", "did not get the expected comment")
	}
}

func TestSetExifComment(t *testing.T) {
	bytes := fileIntoBytes(FixtureFilePath)
	withExif := SetExifComment(bytes, "hello there")
	check, err := GetExifComment(withExif)
	if err != nil {
		t.Errorf("%v", err)
	}
	if check != "hello there" {
		t.Errorf("%s", "did not get expected exif comment")
	}
}

func TestSetExifRuneComment(t *testing.T) {
	bytes := fileIntoBytes(FixtureFilePath)
	withExif := SetExifComment(bytes, "χαιρετε!")
	check, err := GetExifComment(withExif)
	if err != nil {
		t.Errorf("%v", err)
	}
	if check != "χαιρετε!" {
		t.Errorf("did not get expected exif comment")
	}
}
