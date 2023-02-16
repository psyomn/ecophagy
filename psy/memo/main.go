/*
Package memo is the tool for storing labels, in a shitty way, on
files in filesystems.

I am a chronir user of ~/.local/bin, and sometimes I want to remember
why the hell I've installed ~/.local/bin/satan. This is a shitty way
to just add notes in a familiar and quick way, and maintain a key
value store my computer to recheck why I originally did such a thing.

more experimental than anything.

Copyright 2019 Simon Symeonidis (psyomn)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/
package memo

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/psyomn/ecophagy/common"
	psycommon "github.com/psyomn/ecophagy/psy/common"
)

func memoDirPath() string      { return path.Join(common.ConfigDir(), "memo") }
func memoDataFilePath() string { return path.Join(memoDirPath(), "data.gobbin") }

type memoRecord struct {
	ID   uint64
	Data string
}
type memoStore struct {
	Data map[string]*memoRecord
}

func memoStoreNew() *memoStore {
	var store memoStore
	store.Data = make(map[string]*memoRecord)
	return &store
}

func initialize() error {
	if _, err := os.Stat(memoDirPath()); os.IsNotExist(err) {
		if err := os.MkdirAll(memoDirPath(), os.ModePerm); err != nil {
			return err
		}
	}

	if _, err := os.Stat(memoDataFilePath()); os.IsNotExist(err) {
		initStore := memoStoreNew()
		return store(initStore)
	}

	return nil
}

func (s *memoStore) encode() (bytes.Buffer, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	err := enc.Encode(*s)
	if err != nil {
		log.Println("problem encoding memo store file: ", err)
	}

	return buffer, err
}

func (s *memoStore) Add(key, value string) {
	maxID := s.maxID() + 1
	s.Data[key] = &memoRecord{ID: maxID, Data: value}
}

func (s *memoStore) Get(key string) (*memoRecord, bool) {
	val, ok := s.Data[key]
	return val, ok
}

func (s *memoStore) maxID() uint64 {
	var maxID uint64
	for _, v := range s.Data {
		if maxID < v.ID {
			maxID = v.ID
		}
	}

	return maxID
}

func decode(cmdInFile string) *memoStore {
	var buff bytes.Buffer

	if !common.FileExists(cmdInFile) {
		return memoStoreNew()
	}

	dat, err := os.ReadFile(cmdInFile)
	if err != nil {
		log.Fatal("problem opening file:", cmdInFile, ":", err)
		os.Exit(1)
	}

	dec := gob.NewDecoder(&buff)
	var store memoStore
	buff.Write(dat)

	err = dec.Decode(&store)
	if err != nil {
		log.Println("problem decoding store: ", cmdInFile, ", ", err)
		os.Exit(1)
	}

	return &store
}

func store(memos *memoStore) error {
	bytes, err := memos.encode()
	if err != nil {
		return err
	}

	file, err := os.Create(memoDataFilePath())
	if err != nil {
		return fmt.Errorf("problem opening file for storing gob: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(bytes.Bytes()); err != nil {
		return fmt.Errorf("problem writing memo file: %w", err)
	}

	return nil
}

// Run the memo command
// TODO: this needs some cleaning up and a better argument parsing strategy
func Run(args psycommon.RunParams) psycommon.RunReturn {
	if err := initialize(); err != nil {
		return err
	}

	type memoFlags struct {
		fileName string
		list     bool
	}

	sess := memoFlags{}

	memoCmd := flag.NewFlagSet("memo", flag.ExitOnError)
	memoCmd.StringVar(&sess.fileName, "file", sess.fileName, "<message> - the filename to write a memo about")
	memoCmd.BoolVar(&sess.list, "list", sess.list, "list all current memos")
	if err := memoCmd.Parse(args); err != nil {
		return fmt.Errorf("problem parsing: %w", err)
	}

	if sess.list {
		theStore := decode(memoDataFilePath())

		writer := new(tabwriter.Writer)
		writer.Init(os.Stdout, 0, 8, 0, '\t', 0)

		for k, v := range theStore.Data {
			fmt.Fprintf(writer, "%v\t%v\t%v\n", v.ID, k, v.Data)
		}

		writer.Flush()

		return nil
	}

	if _, err := os.Stat(sess.fileName); os.IsNotExist(err) {
		return ErrNonExistFile
	}

	if sess.fileName == "" {
		return nil
	}

	absPath, err := filepath.Abs(sess.fileName)
	if err != nil {
		return fmt.Errorf("problem getting abs path: %w", err)
	}

	if len(memoCmd.Args()) > 0 {
		message := strings.Join(memoCmd.Args(), " ")
		theStore := decode(memoDataFilePath())
		theStore.Add(absPath, message)
		return store(theStore)
	}

	// read operations
	store := decode(memoDataFilePath())
	value, ok := store.Get(absPath)
	if !ok {
		return fmt.Errorf("%w: %v", ErrCantFindEntry, value)
	}

	fmt.Println(value)

	return nil
}
