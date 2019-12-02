package common

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

const OrgName = "ecophagy"

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

func HasHome() bool { return os.Getenv("HOME") != "" }

func DataPath() string {
	return path.Join(os.Getenv("HOME"), ".config", OrgName)
}

func PathExists(path string) bool {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return true
}

// FileToMd5Sum will calculate the hash of a particular file
// contents. The intended use is for quick checks rather than
// secure.
func FileToMd5Sum(path string) (string, error) {
	fh, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer fh.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, fh); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func FileList(dirpath string) ([]string, error) {
	var ret []string

	err := filepath.Walk(dirpath,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			ret = append(ret, path)
			return nil
		})

	if err != nil {
		return nil, err
	}

	return ret, nil
}
