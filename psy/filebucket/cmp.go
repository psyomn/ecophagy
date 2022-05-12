/*
Copyright 2022 Simon Symeonidis (psyomn)

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

package filebucket

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/psyomn/ecophagy/psy/common"
)

type pathWithHash struct {
	path string
	hash string
}

type workerPool struct {
	done       bool
	numWorkers int
	workerWg   sync.WaitGroup
	bufferSize uint
	input      chan pathWithHash
	output     chan pathWithHash
}

func workerPoolNew() *workerPool {
	const bufsz = 500
	return &workerPool{
		done:       false,
		numWorkers: runtime.GOMAXPROCS(0),
		workerWg:   sync.WaitGroup{},
		bufferSize: bufsz,
		input:      make(chan pathWithHash, bufsz),
		output:     make(chan pathWithHash, bufsz),
	}
}

func (s *workerPool) Run() {
	if s.done {
		panic("create a new workerpool instead of reusing old")
	}

	workerFn := func(wg *sync.WaitGroup) {
		defer wg.Done()

		for {
			el, more := <-s.input
			if !more {
				return
			}

			hash, _ := hashFile(el.path)
			el.hash = hash
			s.output <- el
		}
	}

	for i := 0; i < s.numWorkers; i++ {
		s.workerWg.Add(1)
		go workerFn(&s.workerWg)
	}
}

func (s *workerPool) Finish() {
	s.done = true

	close(s.input)
	s.workerWg.Wait()

	close(s.output)
}

func (s *workerPool) Process(entry *pathWithHash) {
	// PERF: passing pointers through channels might be slightly more
	// performant here.
	s.input <- *entry
}

type session struct {
	dirPath string
}

func usage(fs *flag.FlagSet) error {
	fs.Usage()
	return ErrWrongUsage
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}

	// PERF: returning actual byte arrays here might be overall better
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

type fileBucket map[string][]string

func listOffensiveTxt(bucket fileBucket) {
	for k, v := range bucket {
		if len(v) > 1 {
			fmt.Println("# === bucket ", k, "=======")
			for _, el := range v {
				fmt.Printf("# \"%s\"\n", el)
			}
		}
	}
}

func walkAndBucket(path string) fileBucket {
	bucket := make(fileBucket)

	workpl := workerPoolNew()
	workpl.Run()

	var hashWg sync.WaitGroup
	hashWg.Add(1)
	go func(hw *sync.WaitGroup) {
		defer hw.Done()
		for {
			fhash, more := <-workpl.output

			if !more {
				return
			}

			hash := fhash.hash
			currentFile := fhash.path
			bucket[hash] = append(bucket[hash], currentFile)
		}
	}(&hashWg)

	fszMap := make(map[int64][]string)

	err := filepath.WalkDir(path, func(currentFile string, info fs.DirEntry, err error) error {
		if info.IsDir() {
			return nil
		}

		maybeFsz, infoErr := info.Info()
		if infoErr != nil {
			fmt.Fprintf(os.Stderr, "err: could not read: %v\n", infoErr)
			return nil
		}

		fsz := maybeFsz.Size()
		fszMap[fsz] = append(fszMap[fsz], currentFile)

		return nil
	})
	if err != nil {
		fmt.Printf("%#v\n", err)
	}

	for _, v := range fszMap {
		if len(v) < 2 {
			continue
		}

		for _, fileName := range v {
			workpl.Process(&pathWithHash{fileName, ""})
		}
	}

	// signal we're done -- this closes the channel, and should end the
	// goroutine above
	workpl.Finish()

	hashWg.Wait()

	return bucket
}

func Run(args common.RunParams) common.RunReturn {
	sess := session{}

	cmpCmd := flag.NewFlagSet("compare-dirs", flag.ExitOnError)
	cmpCmd.StringVar(&sess.dirPath, "path", sess.dirPath, "the dir to check for duplicates")

	if err := cmpCmd.Parse(args); err != nil {
		return err
	}

	if sess.dirPath == "" {
		goto printUsageAndExit
	}

	listOffensiveTxt(walkAndBucket(sess.dirPath))

	return nil

printUsageAndExit:
	return usage(cmpCmd)
}
