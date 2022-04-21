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

	"github.com/psyomn/ecophagy/psy/common"
)

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

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

type fileBucket map[string][]string

func listOffensiveTxt(bucket fileBucket) {
	for k, v := range bucket {
		if len(v) > 1 {
			fmt.Println("# === bucket ", k, "=======")
			for _, el := range v {
				fmt.Println("# ", el)
			}
		}
	}
}

func WalkAndBucket(path string) fileBucket {
	bucket := make(fileBucket)

	filepath.Walk(path, func(currentFile string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		hash, _ := hashFile(currentFile)

		if _, ok := bucket[hash]; ok {
			bucket[hash] = append(bucket[hash], currentFile)
		} else {
			bucket[hash] = []string{currentFile}
		}

		return nil
	})

	return bucket
}

func Run(args common.RunParams) common.RunReturn {
	fmt.Println("compare dirs")

	sess := session{}

	cmpCmd := flag.NewFlagSet("compare-dirs", flag.ExitOnError)
	cmpCmd.StringVar(&sess.dirPath, "path", sess.dirPath, "the dir to check for duplicates")

	if err := cmpCmd.Parse(args); err != nil {
		return err
	}

	if sess.dirPath == "" {
		goto printUsageAndExit
	}

	listOffensiveTxt(WalkAndBucket(sess.dirPath))

	return nil

printUsageAndExit:
	return usage(cmpCmd)
}
