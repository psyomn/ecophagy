package common

import (
	"bufio"
	"bytes"
	"os"
)

// ReadSpaceConfig reads a text file, which each line has two
// continuous strings, delimited by one blank space
func ReadSpaceConfig(filename string) (map[string]string, error) {
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
