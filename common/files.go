package common

import (
	"io/ioutil"
	"os"
)

func FileToBytes(filename string) ([]byte, error) {
	fs, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fs.Close()

	bytes, err := ioutil.ReadAll(fs)
	return bytes, err
}
