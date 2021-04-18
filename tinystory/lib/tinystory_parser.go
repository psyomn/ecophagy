/**
 * A more experimental parser for the tinystory project. I'm mostly
 * leveraging around generating json for now. Maybe we can scrap that
 * in the future for our very own parser...
 */
package tinystory

import (
	"errors"
	"io"
	"os"
)

type TokenTypeEnum uint64

const (
	TokenKeyword TokenTypeEnum = iota
	TokenWord
	TokenWhitespace
	TokenNewline
	TokenNumber
)

// TODO: this will eventually be used
// nolint
var terminals = []string{
	"TITLE",
	"COMMENTS",
	"AUTHORS",
	"CHOICE",
	"FRAGMENT",
	"ENDFRAGMENT",
	"GOTO",
}

type Token struct {
	Type       TokenTypeEnum
	Value      string
	LineNumber uint64
}

func ParseTinyStoryFormat(path string) (*Document, error) {
	fs, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	doc, err := ParseTinystoryFormat(fs)
	return doc, err
}

func ParseTinystoryFormat(reader io.ReadCloser) (*Document, error) {
	var b [1]byte

	for {
		_, err := reader.Read(b[:])
		if errors.Is(err, io.EOF) {
			break
		}
	}

	defer reader.Close()
	return nil, nil
}
