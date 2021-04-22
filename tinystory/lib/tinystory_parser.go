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
	"strings"
)

type TokenTypeEnum uint64

const (
	TokenKeyword TokenTypeEnum = iota
	TokenKeywordTitle
	TokenKeywordComment
	TokenKeywordAuthors
	TokenKeywordEndAuthors
	TokenKeywordFragment
	TokenKeywordEndFragment
	TokenKeywordGoto

	TokenWord
	TokenNewline
	TokenNumber
	TokenSemicolon

	TokenError
)

var skip = []byte{
	' ',
	'\t',
	'\n',
	';',
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

type Parser struct {
	tokens []Token
	cursor int
	doc    *Document
}

func NewParser(tokens []Token, doc *Document) *Parser {
	return &Parser{
		tokens: tokens,
		cursor: 0,
		doc:    doc,
	}
}

func (s *Parser) Current() *Token {
	return &s.tokens[s.cursor]
}

func (s *Parser) Execute() {
	for s.cursor < len(s.tokens) {
		/* top level keywords go here */
		switch s.Current().Type {
		case TokenKeywordTitle:
			s.ParseTitle()
		case TokenKeywordAuthors:
			s.ParseAuthors()
		case TokenKeywordComment:
			s.ParseComment()
		case TokenKeywordFragment:
			s.ParseFragment()
		case TokenError:
			fallthrough
		case TokenNewline:
			// in the case of stray newlines, we don't care and we
			// skip to the next token
			s.cursor++
		default:
			fmt.Println(s.Current())
			return
		}
	}
}

func (s *Parser) CheckEndStatement() bool {
	// this is extraneous and I should recheck the grammar
	if s.Current().Type != TokenSemicolon {
		panic(
			fmt.Sprintf(
				"line %d: expected semicolon but got: %s",
				s.Current().LineNumber,
				s.Current().Value,
			),
		)
	}

	return true
}

func (s *Parser) ParseTitle() {
	// cursor on "TITLE" move
	s.cursor++

	s.CheckEndStatement()
	s.cursor++

	s.SkipNewlines()

	vals := s.TakeValuesUntilToken(TokenSemicolon)
	s.CheckEndStatement()
	s.cursor++

	s.doc.Title = strings.Join(vals, " ")
}

func (s *Parser) ParseAuthors() {
	// cursor on "AUTHORS" move
	s.cursor++

	s.CheckEndStatement()
	s.cursor++

	s.SkipNewlines()

	for s.Current().Type != TokenKeywordEndAuthors {
		tks := s.TakeValuesUntilToken(TokenSemicolon)

		s.CheckEndStatement()
		s.cursor++

		s.SkipNewlines()

		s.doc.Authors = append(s.doc.Authors, strings.Join(tks, " "))
	}

	// Skip over ENDAUTHORS
	s.cursor++

	// Make sure we're ending the fragment here
	s.CheckEndStatement()
	s.cursor++
}

func (s *Parser) ParseComment() {
	// move over COMMENT keyword
	s.cursor++

	s.CheckEndStatement()
	s.cursor++

	s.SkipNewlines()

	tks := s.TakeValuesUntilToken(TokenSemicolon)
	s.doc.Comment = strings.Join(tks, " ")

	s.CheckEndStatement()
	s.cursor++
}

func (s *Parser) ParseFragment() {
	// move over the FRAGMENT keyword
	s.cursor++

	fragIndex := 0
	{
		value, _ := strconv.Atoi(s.Current().Value)
		fragIndex = value
	}

	s.cursor++
	s.CheckEndStatement()

	s.cursor++
	vals := s.TakeValuesUntilToken(TokenSemicolon)

	s.CheckEndStatement()
	s.cursor++

	s.SkipNewlines()

	choices := s.ParseChoices()

	if s.Current().Type != TokenKeywordEndFragment {
		panic("expected ENDFRAGMENT got " + s.Current().String())
	}

	s.cursor++
	s.CheckEndStatement()
	s.cursor++

	s.doc.Fragments = append(s.doc.Fragments, StoryFragment{
		Index:   fragIndex,
		Content: strings.Join(vals, " "),
		Choices: choices,
	})
}

func (s *Parser) ParseChoices() []Choice {
	var choices []Choice

	fmt.Println("====", s.Current())

	for s.cursor < len(s.tokens) && s.Current().Type != TokenKeywordEndFragment {
		if s.Current().Type != TokenKeywordGoto {
			panic("expected GOTO keyword, got: " + s.Current().String())
		}
		s.cursor++

		if s.Current().Type != TokenNumber {
			panic("expected index/integer, got: " + s.Current().String())
		}

		gIndex := 0
		{
			// regex makes sure that it's indeed a number, don't need
			// to check for errors
			val, _ := strconv.Atoi(s.Current().Value)
			gIndex = val
		}
		s.cursor++

		desc := strings.Join(s.TakeValuesUntilToken(TokenSemicolon), " ")

		choices = append(choices, Choice{
			Index:       gIndex,
			Description: desc,
		})

		s.CheckEndStatement()
		s.cursor++

		s.SkipNewlines()
	}

	return choices
}

func (s *Parser) TakeValuesUntilToken(toktype TokenTypeEnum) []string {
	var tks []string

	maxIndex := len(s.tokens)
	for ; s.Current().Type != toktype && s.cursor < maxIndex; s.cursor++ {
		tks = append(tks, s.Current().Value)
	}

	return tks
}

func (s *Parser) SkipNewlines() {
	for s.cursor < len(s.tokens) && s.Current().Type == TokenNewline {
		s.cursor++
	}
}

func (s *Parser) ParseStringStatement() string {
	// this will parse tokens into a string. The cursor will be place
	// after the semicolon token
	vals := s.TakeValuesUntilToken(TokenSemicolon)
	s.cursor++

	return strings.Join(vals, " ")
}

func TokenTypeEnumString(t TokenTypeEnum) string {
	switch t {
	case TokenKeyword:
		return "KW"
	case TokenKeywordTitle:
		return "KW_TITLE"
	case TokenKeywordComment:
		return "KW_COMMENT"
	case TokenKeywordAuthors:
		return "KW_AUTHORS"
	case TokenKeywordFragment:
		return "KW_FRAGMENT"
	case TokenKeywordEndFragment:
		return "KW_ENDFRAGMENT"
	case TokenKeywordGoto:
		return "KW_GOTO"
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
	switch value {
	case "TITLE":
		return TokenKeywordTitle
	case "AUTHORS":
		return TokenKeywordAuthors
	case "ENDAUTHORS":
		return TokenKeywordEndAuthors
	case "COMMENT":
		return TokenKeywordComment
	case "FRAGMENT":
		return TokenKeywordFragment
	case "ENDFRAGMENT":
		return TokenKeywordEndFragment
	case "GOTO":
		return TokenKeywordGoto
	}

	if value == "\n" {
		return TokenNewline
	}

	if isNumReg.Match([]byte(value)) {
		return TokenNumber
	}

	if isAlphaNumReg.Match([]byte(value)) {
		return TokenWord
	}

	return TokenError
}

func TokenFromString(value string) *Token {
	token := &Token{
		Type:       classifyToken(value),
		Value:      value,
		LineNumber: 0,
	}
	return token
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
	var numLines uint64 = 1

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
				LineNumber: numLines,
			}

			numLines++

			tokens = append(tokens, newToken)
			continue
		}

		if cursor[0] == ';' {
			newToken := Token{
				Type:       TokenSemicolon,
				Value:      ";",
				LineNumber: numLines,
			}

			tokens = append(tokens, newToken)
			continue
		}

		buffer = append(buffer, cursor[0])

		peekBytes, err := breader.Peek(1)

		if errors.Is(err, bufio.ErrBufferFull) || byteArrIncludes(peekBytes[0], skip) {
			strToken := string(buffer)
			buffer = nil

			token := TokenFromString(strToken)
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

	doc := &Document{}
	tokens := tokenize(reader)
	parser := NewParser(tokens, doc)

	fmt.Println("===== TOKENIZE =====")
	fmt.Println(tokens)
	fmt.Println("===== PARSE =====")
	parser.Execute()

	// doc, err := parseTokens(tokens)
	// if err != nil {
	// 	return nil, err
	// }

	return doc, nil
}
