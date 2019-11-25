package main

import (
	"bufio"
	"bytes"
	"math/rand"
	"os"
)

func sampleArray(arr []string) string {
	// tempting to use unsafe here, but I know that this will be
	// used on a windows machine, and not too sure how that will
	// behave there.

	index := rand.Intn(len(arr))
	return arr[index]
}

func readSpaceConfig(filename string) (map[string]string, error) {
	// space configuration is a key value file, whose key values are
	// separated by a space, on each line

	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	reader := bufio.NewReader(fd)
	var line []byte
	ret := map[string]string{}

	for {
		line, _, err = reader.ReadLine()
		if err != nil {
			break
		}

		parts := bytes.Split(line, []byte(" "))

		if len(parts) != 2 {
			panic(string(parts[0]))
		}

		ret[string(parts[0])] = string(parts[1])
	}

	return ret, nil
}
