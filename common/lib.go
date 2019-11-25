package common

import (
	"bufio"
	"io"
	"net/http"
	"os"
)

func FileExists(filename string) bool {
	fd, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	defer fd.Close()
	return err != nil
}

func DownloadIfNotExist(filename, url string) error {
	if FileExists(filename) {
		return nil
	}

	return DownloadFile(filename, url)
}

func DownloadFile(filename, url string) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func FileToLines(filename string) ([]string, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	reader := bufio.NewReader(fd)
	var line []byte
	var ret []string

	for {
		line, _, err = reader.ReadLine()
		if err != nil {
			break
		}

		ret = append(ret, string(line))
	}

	return ret, nil
}
