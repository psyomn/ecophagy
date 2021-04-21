/**
 * A very simple parser that relies on stories that are written in
 * json
 */
package tinystory

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"

	"github.com/psyomn/ecophagy/common"
)

type Choice struct {
	Description string
	Index       int
}

func (s Choice) String() string {
	return fmt.Sprintf("<Choice index:%d description:%s>", s.Index, s.Description)
}

func (s *Choice) UnmarshalJSON(data []byte) error {
	elements := []interface{}{&s.Description, &s.Index}
	return json.Unmarshal(data, &elements)
}

type StoryFragment struct {
	Index   int
	Content string
	Choices []Choice
}

func (s StoryFragment) String() string {
	return fmt.Sprintf("<StoryFragment %d %s %s>", s.Index, s.Content, s.Choices)
}

func (s *StoryFragment) UnmarshalJSON(data []byte) error {
	elements := []interface{}{&s.Index, &s.Content, &s.Choices}
	return json.Unmarshal(data, &elements)
}

type Document struct {
	Title     string          `json:"title"`
	Comment   string          `json:"comment"`
	Authors   []string        `json:"authors"`
	Website   string          `json:"website"`
	Fragments []StoryFragment `json:"story"`
}

func (s Document) String() string {
	return fmt.Sprintf(`Title: %s

Comment: %s

Authors: %v
`, s.Title, s.Comment, s.Authors)
}

func Parse(bjson []byte) (*Document, error) {
	doc := &Document{}

	err := json.Unmarshal(bjson, &doc)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func ParseAllInDir(dirpath string) ([]Document, error) {
	docs := make([]Document, 0, 256)

	err := filepath.Walk(dirpath, func(currpath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path.Ext(currpath) != ".json" {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		data, err := common.FileToBytes(currpath)
		if err != nil {
			//nolint
			return nil
		}

		document, err := Parse(data)
		if err != nil {
			return err
		}

		docs = append(docs, *document)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return docs, nil
}
