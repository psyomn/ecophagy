/**
 * A more experimental parser for the tinystory project. I'm mostly
 * leveraging around generating json for now. Maybe we can scrap that
 * in the future for our very own parser...
 *
 * I want to rewrite all of this because I can do better. Wanted to
 * keep it simple but it kind of ended up being a mess.
 */
package tinystory

import (
	"fmt" // TODO remove me

	"bufio"
	"errors"
	"io"
	"os"
	"regexp"
	"strconv"
)

type TokenTypeEnum uint64

const (
	TokenKeyword TokenTypeEnum = iota
	TokenWord
	TokenNewline
	TokenNumber
	TokenSemicolon

	TokenError
)

var keywords = []string{
	"TITLE",
	"COMMENTS",
	"AUTHORS",
	"FRAGMENT",
	"ENDFRAGMENT",
	"GOTO",
}

var skip = []byte{
	' ',
	'\t',
	'\n',
}

var (
	isNumReg = regexp.MustCompile(`^\d+$`)

	// The stricter way of doing this is specifying ^ and $, however
	// we'd need to explicitly specify anything else which is glued to
	// a word.  Since this is a best effort, and a somewhat loose
	// parser, not specifying some of this stuff is OK.
	isAlphaNumReg = regexp.MustCompile(`[0-9a-zA-Z]+`)
)

type Token struct {
	Type       TokenTypeEnum
	Value      string
	LineNumber uint64
}

func TokenTypeEnumString(t TokenTypeEnum) string {
	switch t {
	case TokenKeyword:
		return "KEYWORD"
	case TokenWord:
		return "WORD"
	case TokenNewline:
		return "NEWLINE"
	case TokenNumber:
		return "NUMBER"
	case TokenSemicolon:
		return "SEMICOLON"

	case TokenError:
		fallthrough
	default:
		return "UNKNOWN_TOKEN_ERROR"
	}
}

func (s Token) String() string {
	if s.Type == TokenNewline {
		return fmt.Sprintf("<Token \\n %s>", TokenTypeEnumString(s.Type))
	}
	return fmt.Sprintf("<Token \"%s\" %s>", s.Value, TokenTypeEnumString(s.Type))
}

func classifyToken(value string) TokenTypeEnum {
	if stringArrIncludes(value, keywords) {
		return TokenKeyword
	}

	if isNumReg.Match([]byte(value)) {
		return TokenNumber
	}

	if value == "\n" {
		return TokenNewline
	}

	if isAlphaNumReg.Match([]byte(value)) {
		return TokenWord
	}

	return TokenError
}

func TokenFromString(value string) (*Token, error) {
	token := &Token{
		Type:       classifyToken(value),
		Value:      value,
		LineNumber: 0,
	}
	return token, nil
}

func stringArrIncludes(e string, arr []string) bool {
	for index := range arr {
		if arr[index] == e {
			return true
		}
	}
	return false
}

func byteArrIncludes(b byte, arr []byte) bool {
	for index := range arr {
		if arr[index] == b {
			return true
		}
	}
	return false
}

func tokenize(reader io.Reader) []Token {
	cursor := make([]byte, 1)
	breader := bufio.NewReader(reader)
	buffer := make([]byte, 0, 1024)
	tokens := make([]Token, 0, 1024)

	for {
		_, err := breader.Read(cursor)
		if errors.Is(err, io.EOF) {
			break
		}

		if cursor[0] == '\r' {
			// ignore windows like newlines
			continue
		}

		if cursor[0] == ' ' {
			continue
		}

		if cursor[0] == '\n' {
			newToken := Token{
				Type:       TokenNewline,
				Value:      "\n",
				LineNumber: 0,
			}

			tokens = append(tokens, newToken)

			continue
		}

		buffer = append(buffer, cursor[0])

		peekBytes, err := breader.Peek(1)

		if errors.Is(err, bufio.ErrBufferFull) || byteArrIncludes(peekBytes[0], skip) {
			strToken := string(buffer)
			buffer = nil

			token, err := TokenFromString(strToken)
			if err != nil {
				// TODO: dont panic
				panic(err)
			}

			tokens = append(tokens, *token)

			continue
		}
	}

	return tokens
}

func parseKeyword(tokens []Token, cursor *int, doc *Document) {
	maxIndex := len(tokens) - 1

	switch tokens[*cursor].Value {
	case "TITLE":
		// skip title
		*cursor++
		fmt.Println("parse title")

		for ; tokens[*cursor].Type != TokenKeyword && *cursor < maxIndex; *cursor++ {
			if tokens[*cursor].Type != TokenNewline {
				doc.Title += tokens[*cursor].Value + " "
			}
		}

	case "AUTHORS":
		// skip authors
		*cursor++
		fmt.Println("parse authors")

		currAuthor := ""
		for ; tokens[*cursor].Type != TokenKeyword && *cursor < maxIndex; *cursor++ {
			if tokens[*cursor].Type != TokenNewline {
				currAuthor += tokens[*cursor].Value + " "
				continue
			}

			doc.Authors = append(doc.Authors, currAuthor)
			currAuthor = ""
		}

	case "COMMENTS":
		// TODO
	case "FRAGMENT":
		*cursor++

		frag := StoryFragment{}

		fmt.Println("parse fragment")
		if tokens[*cursor].Type != TokenNumber {
			// error case here -- fragment should be followed by a
			// (non negative) number
		}

		{ /* parse inded */
			// no error check since the regex makes sure the input is
			// a number
			num, _ := strconv.Atoi(tokens[*cursor].Value)

			// not the best, but we can fix this in a rewrite, when we
			// don't care about json
			frag.Index = num
		}

		// skip extraneous newlines
		for ; tokens[*cursor].Type != TokenNewline && *cursor < maxIndex; *cursor++ {
		}

		for ; tokens[*cursor].Type != TokenKeyword && *cursor < maxIndex; *cursor++ {
			frag.Content += tokens[*cursor].Value + " "
		}

		// TODO: parse GOTO

		fmt.Println("end parse fragment")
		if tokens[*cursor].Value == "ENDFRAGMENT" {
			// extra check here
			*cursor++
		}

		doc.Fragments = append(doc.Fragments, frag)
	}
}

func parseTokens(tokens []Token) (*Document, error) {
	var doc Document

	for cursor := 0; cursor < len(tokens); {
		/* NOTE: parsing functions must be responsible on updating the
		   cursor position */
		fmt.Println(tokens[cursor])
		switch tokens[cursor].Type {
		case TokenKeyword:
			parseKeyword(tokens, &cursor, &doc)
		case TokenNewline:
			cursor++
		}

	}

	return &doc, nil
}

func ParseTinyStoryFormatFile(path string) (*Document, error) {
	fs, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	doc, err := ParseTinystoryFormat(fs)
	return doc, err
}

func ParseTinystoryFormat(reader io.ReadCloser) (*Document, error) {
	defer reader.Close()

	var doc *Document
	tokens := tokenize(reader)
	fmt.Println(tokens)

	doc, err := parseTokens(tokens)
	if err != nil {
		return nil, err
	}

	fmt.Println(*doc)

	return doc, nil
}
